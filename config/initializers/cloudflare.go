package initializers

import (
	"github.com/cloudflare/cloudflare-go"
	"gitlab.com/scalablespace/listener/app/models"
)

func NewCloudflare(env models.Environment) (*cloudflare.API, error) {
	return cloudflare.NewWithAPIToken(env.CloudflareToken)
}
