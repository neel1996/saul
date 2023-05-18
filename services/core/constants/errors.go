package constants

import "encoding/json"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var (
	ExternalApiError             = Error{Code: "EXTERNAL_API_ERROR", Message: "Failed to invoke external API"}
	DocumentQANoAnswerFoundError = Error{Code: "DOCUMENT_QA_NO_ANSWER_FOUND", Message: "No answer found"}
	UserNotFoundError            = Error{Code: "USER_NOT_FOUND", Message: "User not found"}
	UserLoginError               = Error{Code: "USER_LOGIN_ERROR", Message: "Error occurred while logging in user"}
	RequestValidationError       = Error{Code: "REQUEST_VALIDATION_ERROR", Message: "Invalid request"}
)

func (e Error) Error() string {
	return e.Message
}

func (e Error) GetCode() int {
	errorCodeToHttpCodeMap := map[string]int{
		"EXTERNAL_API_ERROR":          500,
		"DOCUMENT_QA_NO_ANSWER_FOUND": 404,
		"USER_NOT_FOUND":              404,
		"USER_LOGIN_ERROR":            500,
		"REQUEST_VALIDATION_ERROR":    400,
	}

	return errorCodeToHttpCodeMap[e.Code]
}

func (e Error) String() string {
	b, err := json.Marshal(e)
	if err != nil {
		return ""
	}

	return string(b)
}

func (e Error) GetGinResponse() (int, Error) {
	return e.GetCode(), e
}
