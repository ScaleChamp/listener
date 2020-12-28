package adapters

import (
	"gitlab.com/scalablespace/listener/app/models"
)

type Selectel struct {
}

func (dc *Selectel) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	return node, nil
}

func (dc *Selectel) DeleteNode(node *models.Node) error {
	return nil
}

func NewSelectel() *Selectel {
	return &Selectel{}
}
