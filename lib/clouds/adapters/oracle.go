package adapters

import (
	"gitlab.com/scalablespace/listener/app/models"
)

type Oracle struct {
}

func (o *Oracle) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	return node, nil
}

func (o *Oracle) DeleteNode(node *models.Node) error {
	return nil
}

func NewOracle() *Oracle {
	return &Oracle{}
}
