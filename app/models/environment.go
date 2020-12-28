package models

import (
	"time"
)

type Environment struct {
	RabbitURL  string `env:"RABBIT_URL,required"`
	TasksQueue string `env:"TASKS_QUEUE" envDefault:"tasks"`
	// db options
	DBConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"1m"`
	DBMaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"10"`
	DBMaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"1"`
	DatabaseUrl       string        `env:"DATABASE_URL,required"`
	// secrets
	SecretKeyPath     string `env:"SECRET_KEY_PATH,required"`
	SecretKeyPassword string `env:"SECRET_KEY_PASS"`
	// publish from
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"0.0.0.0:8011"`
	// softlayer
	SoftLayerApiKey string `env:"SL_API_KEY,required"`
	// alibabacloud
	AlibabaAccessKeyId string `env:"ALI_ACCESS_KEY_ID"`
	AlibabaSecretKey   string `env:"ALI_SECRET_KEY"`
	// AmazonWebServices
	AWSAccessKeyId     string `env:"AWS_ACCESS_KEY_ID,required"`
	AWSSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY,required"`
	// pusher-url
	PusherUrl string `env:"PUSHER_URL,required"`
	// cloudflare
	CloudflareToken  string `env:"CF_TOKEN,required"`
	CloudflareZone   string `env:"CF_ZONE,required"`
	CloudflareDomain string `env:"CF_DOMAIN,required"`
	// digitalocean
	DigitalOceanToken string `env:"DIGITAL_OCEAN_TOKEN,required"`
	// scaleway
	SCWAccessKey    string `env:"SCW_ACCESS_KEY,required"`
	SCWSecretKey    string `env:"SCW_SECRET_KEY,required"`
	SCWOrganisation string `env:"SCW_ORG,required"`
	// hetzner
	HetznerToken string `env:"HETZNER_TOKEN,required"`
	// linode
	LinodeAccessKey string `env:"LINODE_ACCESS_KEY,required"`
	// upcloud access
	UpcloudPassword string `env:"UPCLOUD_PASSWORD,required"`
	UpcloudUser     string `env:"UPCLOUD_USER,required"`
	// azure
	AzureSubscriptionId string `env:"AZURE_SUBSCRIPTION_ID,required"`
	AzureTenantId       string `env:"AZURE_TENANT_ID,required"`
	AzureClientId       string `env:"AZURE_CLIENT_ID,required"`
	AzureClientSecret   string `env:"AZURE_CLIENT_SECRET,required"`
	// vultr token
	VultrToken string `env:"VULTR_TOKEN,required"`
	// exoscale
	ExoscaleApiKey    string `env:"EXOSCALE_API_KEY"`
	ExoscaleApiSecret string `env:"EXOSCALE_API_SECRET"`
}
