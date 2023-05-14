package request

type DocumentQARequest struct {
	Inputs Inputs `json:"inputs" validate:"required"`
}

type Inputs struct {
	Question string `json:"question" validate:"required"`
	Image    string `json:"image" validate:"required"`
}
