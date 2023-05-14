package initializers

import (
	"github.com/go-resty/resty/v2"
	"github.com/neel1996/saul/clients"
	"github.com/neel1996/saul/clients/hugging_face"
	"github.com/neel1996/saul/configuration"
)

func InitializeClients(config configuration.Configuration) {
	httpClient := clients.NewHttpClient(resty.New())

	hugging_face.NewDocumentQAClient(config, httpClient)
}
