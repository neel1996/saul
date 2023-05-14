package response

type DocumentQAResponse struct {
	Score  float64 `json:"score" validate:"required"`
	Answer string  `json:"answer" validate:"required"`
	Start  int     `json:"start"`
	End    int     `json:"end"`
}
