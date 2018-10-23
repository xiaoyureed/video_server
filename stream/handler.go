package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"video_server/common"

	"github.com/julienschmidt/httprouter"
)

func streamHandler(writer http.ResponseWriter, req *http.Request, params httprouter.Params) {
	videoId := params.ByName("id")
	videoPath := common.VIDEO_DIR + videoId

	videoFile, err := os.Open(videoPath)
	if err != nil {
		sendErrorResponse(writer, http.StatusInternalServerError, "internal error")
		return
	}

	writer.Header().Set("Content-type", "video/mp4")
	http.ServeContent(writer, req, "", time.Now(), videoFile)

	defer videoFile.Close()

}

func uploadHandler(writer http.ResponseWriter, req *http.Request, params httprouter.Params) {
	req.Body = http.MaxBytesReader(writer, req.Body, common.MAX_UPLOAD_SIZE)
	if err := req.ParseMultipartForm(common.MAX_UPLOAD_SIZE); err != nil {
		sendErrorResponse(writer, http.StatusBadRequest, "File is too big.")
		return
	}

	file, _, err := req.FormFile("video") // form name == "file"
	if err != nil {
		log.Printf("Error while geting file: %v", err)
		sendErrorResponse(writer, http.StatusInternalServerError, "Internal error")
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Error while reading file: %v", err)
		sendErrorResponse(writer, http.StatusInternalServerError, "Internal error")
		return
	}

	videoName := params.ByName("id")
	err = ioutil.WriteFile(common.VIDEO_DIR+videoName, data, 0666) // 尽量不要用 0777
	if err != nil {
		log.Printf("Error while writing file : %v", err)
		sendErrorResponse(writer, http.StatusInternalServerError, "Internal error")
		return
	}

	writer.WriteHeader(http.StatusCreated)
	io.WriteString(writer, "Upload success")
}

func uploadPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles(common.VIDEO_DIR + "upload.html")
	t.Execute(w, nil)
}
