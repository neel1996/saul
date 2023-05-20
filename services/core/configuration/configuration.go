package configuration

type Configuration struct {
	Port           string      `json:"port" validate:"required"`
	Cors           Cors        `json:"cors" validate:"required,dive"`
	TrustedProxies []string    `json:"trustedProxies" validate:"required"`
	CorsIgnoreUrls []string    `json:"corsIgnoreUrls" validate:"required"`
	AuthIgnoreUrls []string    `json:"authIgnoreUrls" validate:"required"`
	HuggingFace    HuggingFace `json:"huggingFace" validate:"required,dive"`
	Kafka          Kafka       `json:"kafka" required:"required,dive"`
	Minio          Minio       `json:"minio" validate:"required,dive"`
}

type Cors struct {
	AllowedOrigins []string `json:"allowedOrigins" validate:"required"`
	AllowedHeaders []string `json:"allowedHeaders" validate:"required"`
	AllowedMethods []string `json:"allowedMethods" validate:"required"`
}

type HuggingFace struct {
	DocumentQA DocumentQA `json:"documentQA" validate:"required"`
}

type DocumentQA struct {
	Endpoint string `json:"endpoint" validate:"required"`
}

type Kafka struct {
	BrokerURL string `json:"brokerURL"`
	Topics    Topics `json:"topics"`
}

type Topics struct {
	ProcessDocument       TopicDetails `json:"processDocument"`
	ProcessDocumentStatus TopicDetails `json:"processDocumentStatus"`
}

type TopicDetails struct {
	Name    string `json:"name"`
	GroupId string `json:"groupId"`
}

type Minio struct {
	EndPoint string `json:"endPoint" required:"required"`
	Bucket   string `json:"bucket" required:"required"`
}
