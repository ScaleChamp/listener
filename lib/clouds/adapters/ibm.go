package adapters

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"gitlab.com/scalablespace/listener/app/models"
	"time"
)

// curl -X GET --user apikey:ExuZCe37vTnztoCktKsnWCeflGdHegYuSVqwrlW9UG7B "https://api.softlayer.com/rest/v3/SoftLayer_Location/getDatacenters.json" | jq
// curl -X GET --user apikey:ExuZCe37vTnztoCktKsnWCeflGdHegYuSVqwrlW9UG7B "https://api.softlayer.com/rest/v3.1/SoftLayer_Virtual_Guest/getCreateObjectOptions.json?objectMask=mask\[datacenter\]" | jq
type IBM struct {
	session *session.Session
	env     models.Environment
}

func (ibm *IBM) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	vGuestTemplate := datatypes.Virtual_Guest{
		Hostname:  sl.String(fmt.Sprintf("%.8s%.8s", node.InstanceId, node.Id)),
		MaxMemory: sl.Int(plan.IBMCloud.Memory),
		Domain:    sl.String(ibm.env.CloudflareDomain),
		StartCpus: sl.Int(plan.IBMCloud.Cpus),
		Datacenter: &datatypes.Location{
			Name: sl.String(plan.Region),
		},
		//OperatingSystemReferenceCode: sl.String(plan.IBMCloud.OS), // "UBUNTU_LATEST_64"
		LocalDiskFlag:     sl.Bool(true),
		HourlyBillingFlag: sl.Bool(true),
		SshKeys: []datatypes.Security_Ssh_Key{
			{
				Id: sl.Int(plan.IBMCloud.KeyId),
			},
		},
		BlockDeviceTemplateGroup: &datatypes.Virtual_Guest_Block_Device_Template_Group{
			GlobalIdentifier: sl.String(plan.IBMCloud.OS),
		},
		PrimaryNetworkComponent: &datatypes.Virtual_Guest_Network_Component{
			SecurityGroupBindings: []datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
				{
					SecurityGroupId: sl.Int(2786416),
				},
				{
					SecurityGroupId: sl.Int(2786418),
				},
			},
		},
	}

	service := services.GetVirtualGuestService(ibm.session)
	vGuest, err := service.Mask("id;domain").CreateObject(&vGuestTemplate)
	if err != nil {
		return nil, err
	}
	id := aws.IntValue(vGuest.Id)
	service = service.Id(id)

	time.Sleep(11 * time.Second)
	for transactions, err := service.GetActiveTransactions(); len(transactions) > 0 || err != nil; {
		time.Sleep(11 * time.Second)
		transactions, err = service.GetActiveTransactions()
		if err != nil {
			return nil, err
		}
	}

	node.Metadata.ServerId = id
	node.Metadata.IPv4, err = service.GetPrimaryIpAddress()
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (ibm *IBM) DeleteNode(node *models.Node) error {
	service := services.GetVirtualGuestService(ibm.session)
	vm := service.Id(node.Metadata.ServerId)
	success, err := vm.DeleteObject()
	if err != nil {
		return err
	}
	if !success {
		return errors.New("not success ")
	}
	return nil
}

func NewIBM(session *session.Session, env models.Environment) *IBM {
	return &IBM{session, env}
}
