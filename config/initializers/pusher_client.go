package initializers

import (
	"github.com/pusher/pusher-http-go"
	"gitlab.com/scalablespace/listener/app/models"
	"net/http"
)

func NewPusherClient(httpClient *http.Client, env models.Environment) (*pusher.Client, error) {
	client, err := pusher.ClientFromURL(env.PusherUrl)
	if err != nil {
		return nil, err
	}
	client.HTTPClient = httpClient
	return client, nil
}
