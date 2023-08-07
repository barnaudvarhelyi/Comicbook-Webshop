package dtos

type ResponseDto struct {
	Message string `json:"message"`
}
type ErrorResponseDto struct {
	Error string `json:"error"`
}
