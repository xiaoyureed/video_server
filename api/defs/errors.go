package defs

type Error struct {
	ErrorMsg  string `json:"error_msg"`
	ErrorCode string `json:"error_code"`
}

type ErrorResponse struct {
	HttpCode int
	Error    Error
}

var (
	ErrorRequestBodyParseFailed = ErrorResponse{400, Error{"Request body is wrong", "001"}}
	ErrorNotAuthUser            = ErrorResponse{401, Error{"User authentication failed", "002"}}
	ErrorDBOperationFailed      = ErrorResponse{500, Error{"DB operation failed", "003"}}
	ErrorInternalError          = ErrorResponse{500, Error{"go sdk internal error", "004"}}
)
