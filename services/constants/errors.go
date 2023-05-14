package constants

type Error struct {
	Code    string
	Message string
}

var (
	ExternalApiError        = Error{Code: "EXTERNAL_API_ERROR", Message: "Failed to invoke external API"}
	DocumentQANoAnswerFound = Error{Code: "DOCUMENT_QA_NO_ANSWER_FOUND", Message: "No answer found"}
)

func (e Error) Error() string {
	return e.Message
}

func (e Error) GetCode() int {
	errorCodeToHttpCodeMap := map[string]int{
		"EXTERNAL_API_ERROR":          500,
		"DOCUMENT_QA_NO_ANSWER_FOUND": 404,
	}

	return errorCodeToHttpCodeMap[e.Code]
}
