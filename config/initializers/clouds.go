package initializers

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/digitalocean/godo"
	"github.com/exoscale/egoscale"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/linode/linodego"
	"github.com/scaleway/scaleway-sdk-go/scw"
	sl_session "github.com/softlayer/softlayer-go/session"
	"github.com/vultr/govultr"
	"gitlab.com/scalablespace/listener/app/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gcp "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"net/http"
	"time"
)

func NewAzureClient(env models.Environment) (compute.VirtualMachinesClient, network.InterfacesClient, network.PublicIPAddressesClient, compute.DisksClient, error) {
	computeClient := compute.NewVirtualMachinesClient(env.AzureSubscriptionId)
	nicClient := network.NewInterfacesClient(env.AzureSubscriptionId)
	publicIpsClient := network.NewPublicIPAddressesClient(env.AzureSubscriptionId)
	disksClient := compute.NewDisksClient(env.AzureSubscriptionId)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return computeClient, nicClient, publicIpsClient, disksClient, err
	}
	computeClient.Authorizer = authorizer
	nicClient.Authorizer = authorizer
	publicIpsClient.Authorizer = authorizer
	disksClient.Authorizer = authorizer
	return computeClient, nicClient, publicIpsClient, disksClient, nil
}

func NewVultrClient(env models.Environment, client *http.Client) *govultr.Client {
	return govultr.NewClient(client, env.VultrToken)
}

func NewHetznerClient(env models.Environment) *hcloud.Client {
	return hcloud.NewClient(hcloud.WithToken(env.HetznerToken))
}

func NewScalewayClient(env models.Environment) (*scw.Client, error) {
	opts := scw.WithAuth(env.SCWAccessKey, env.SCWSecretKey)
	return scw.NewClient(opts)
}

func NewUpcloud(env models.Environment) *service.Service {
	c := client.New(env.UpcloudUser, env.UpcloudPassword)
	return service.New(c)
}

func NewLinodeGo(env models.Environment) linodego.Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: env.LinodeAccessKey})
	oauth2Client := oauth2.NewClient(context.Background(), tokenSource)
	return linodego.NewClient(oauth2Client)
}

func NewDOClient(env models.Environment) *godo.Client {
	oauth := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: env.DigitalOceanToken})
	c := oauth2.NewClient(context.Background(), oauth)
	return godo.NewClient(c)
}

func NewIBMClient(env models.Environment) *sl_session.Session {
	return sl_session.New("apikey", env.SoftLayerApiKey)
}

func NewGCPClient() (*gcp.InstancesService, *gcp.ZoneOperationsService, error) {
	ui, _ := google.CredentialsFromJSON(context.TODO(), []byte(os.Getenv("GCP_KEY")), "https://www.googleapis.com/auth/compute")
	xc := oauth2.NewClient(context.TODO(), ui.TokenSource)
	srv, err := gcp.NewService(context.TODO(), option.WithHTTPClient(xc))
	if err != nil {
		return nil, nil, err
	}
	return gcp.NewInstancesService(srv), gcp.NewZoneOperationsService(srv), nil
}

func NewExoscale(env models.Environment) *egoscale.Client {
	c := egoscale.NewClient("https://api.exoscale.com/compute", env.ExoscaleApiKey, env.ExoscaleApiSecret)
	c.Timeout = 3 * time.Minute
	return c
}
