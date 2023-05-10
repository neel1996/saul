package configuration

type Configuration struct {
	Port           string   `json:"port" validate:"required"`
	Cors           Cors     `json:"cors" validate:"required"`
	TrustedProxies []string `json:"trustedProxies" validate:"required"`
	CorsIgnoreUrls []string `json:"corsIgnoreUrls" validate:"required"`
	AuthIgnoreUrls []string `json:"authIgnoreUrls" validate:"required"`
}

type Cors struct {
	AllowedOrigins []string `json:"allowedOrigins" validate:"required"`
	AllowedHeaders []string `json:"allowedHeaders" validate:"required"`
	AllowedMethods []string `json:"allowedMethods" validate:"required"`
}
