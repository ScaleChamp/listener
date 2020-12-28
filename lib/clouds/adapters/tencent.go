package adapters

import (
	"gitlab.com/scalablespace/listener/app/models"
)

type TencentCloud struct {
}

func (dc *TencentCloud) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	return node, nil
}

func (dc *TencentCloud) DeleteNode(node *models.Node) error {
	return nil
}

func NewTencentCloud() *TencentCloud {
	return &TencentCloud{}
}
