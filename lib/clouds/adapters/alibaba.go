package adapters

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/pkg/errors"
	"gitlab.com/scalablespace/listener/app/models"
	"sync"
)

type Alibaba struct {
	sync.RWMutex
	regionToClient map[string]*ecs.Client
	env            models.Environment
}

func (ali *Alibaba) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	client, err := ali.client(node.Region)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get client for region")
	}
	request := ecs.CreateCreateInstanceRequest()
	request.ImageId = plan.AlibabaCloud.ImageId
	request.ImageFamily = ""
	//request.ZoneId = ""
	request.InstanceType = plan.AlibabaCloud.InstanceType
	request.RegionId = plan.Region
	request.InternetChargeType = plan.AlibabaCloud.InternetChargeType
	request.KeyPairName = "main"
	request.HostName = fmt.Sprintf("%s-%s", plan.Kind, node.Id)
	request.InstanceName = request.HostName
	response, err := client.CreateInstance(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create instance")
	}
	req := ecs.CreateDescribeInstancesRequest()
	instanceIds, err := json.Marshal([]string{response.InstanceId})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal")
	}
	req.InstanceIds = string(instanceIds)
	resp, err := client.DescribeInstances(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to describe instances")
	}
	node.Metadata.IPv4 = resp.Instances.Instance[0].PublicIpAddress.IpAddress[0]
	if len(resp.Instances.Instance[0].PublicIpAddress.IpAddress) > 1 {
		node.Metadata.IPv6 = resp.Instances.Instance[0].PublicIpAddress.IpAddress[1]
	}
	node.Metadata.RecordId = response.InstanceId
	return node, nil
}

func (ali *Alibaba) DeleteNode(node *models.Node) error {
	client, err := ali.client(node.Region)
	if err != nil {
		return errors.Wrap(err, "failed to get client for region")
	}
	request := ecs.CreateDeleteInstanceRequest()
	request.InstanceId = node.Metadata.RecordId

	if _, err := client.DeleteInstance(request); err != nil {
		return errors.Wrap(err, "failed to delete instance")
	}
	return nil
}

func (ali *Alibaba) client(r string) (*ecs.Client, error) {
	ali.RLock()
	created := ali.regionToClient[r]
	ali.RUnlock()
	if created != nil {
		return created, nil
	}
	ali.Lock()
	defer ali.Unlock()
	c, err := ecs.NewClientWithAccessKey(r, ali.env.AlibabaAccessKeyId, ali.env.AlibabaSecretKey)
	if err != nil {
		return nil, err
	}
	ali.regionToClient[r] = c
	return c, nil

}

func NewAlibaba(env models.Environment) *Alibaba {
	return &Alibaba{env: env, regionToClient: make(map[string]*ecs.Client, 20)}
}
