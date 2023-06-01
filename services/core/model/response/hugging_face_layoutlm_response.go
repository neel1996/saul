package response

type LayoutLMAnswer struct {
	Score  float64 `json:"score" validate:"required"`
	Answer string  `json:"answer" validate:"required"`
	Err    error   `json:"err,omitempty"`
}
