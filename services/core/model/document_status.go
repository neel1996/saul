package model

type DocumentStatus struct {
	Status   string `json:"status" validate:"required,oneof=success error"`
	Checksum string `json:"checksum" validate:"required"`
}

const (
	DocumentStatusSuccess = "success"
	DocumentStatusError   = "error"
)
