package datastructure

type BaseResponse struct {
	Success bool        `json:"success" default:"true"`
	Message string      `json:"msg" default:""`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success" default:"false"`
	Message string `json:"msg"`
}

type ErrorResponseWithCode struct {
	ErrorResponse
	Details map[string]string `json:"details"`
}
