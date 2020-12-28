package adapters

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"gitlab.com/scalablespace/listener/app/models"
	"sync"
	"time"
)

type AWS struct {
	sync.RWMutex
	regionToEC2 map[string]*ec2.EC2
	config      *aws.Config
}

func (dc *AWS) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	reservations, err := dc.clientByRegion(node.Region).RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(plan.AmazonWebServices.ImageId),
		InstanceType: aws.String(plan.AmazonWebServices.InstanceType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		Placement: &ec2.Placement{
			AvailabilityZone: aws.String(plan.AmazonWebServices.Zone),
		},
		NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
			{
				DeviceIndex: aws.Int64(0),
				Groups:      aws.StringSlice([]string{plan.AmazonWebServices.GroupId}),
				SubnetId:    aws.String(plan.AmazonWebServices.SubnetId),
			},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to run instance")
	}
	time.Sleep(time.Minute)
	for i := 0; i < 5; i += 1 {
		time.Sleep(5 * time.Second)
		output, err := dc.clientByRegion(node.Region).DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{reservations.Instances[0].InstanceId},
		})
		if err != nil {
			continue
		}
		node.Metadata.IPv4 = *output.Reservations[0].Instances[0].PublicIpAddress
		node.Metadata.StringId = *output.Reservations[0].Instances[0].InstanceId
		return node, nil
	}
	return nil, errors.New("failed to due to timeout")
}

func (dc *AWS) DeleteNode(node *models.Node) error {
	_, err := dc.clientByRegion(node.Region).TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(node.Metadata.StringId)},
	})
	if err != nil {
		return err
	}
	return nil
}

func (dc *AWS) clientByRegion(region string) *ec2.EC2 {
	dc.RLock()
	created := dc.regionToEC2[region]
	dc.RUnlock()
	if created != nil {
		return created
	}
	dc.Lock()
	defer dc.Unlock()
	s := session.Must(session.NewSession(dc.config.Copy().WithRegion(region)))
	c := ec2.New(s)
	dc.regionToEC2[region] = c
	return c

}

func NewAWS() *AWS {
	config := aws.NewConfig()
	config.WithCredentials(credentials.NewEnvCredentials())
	return &AWS{
		config:      config,
		regionToEC2: make(map[string]*ec2.EC2, 20),
	}
}
