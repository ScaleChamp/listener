package repositories

import (
	"database/sql"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
)

type nodeRepository struct {
	db *sql.DB
}

const updateNodeState = `UPDATE nodes SET state = $2 WHERE id = $1`

func (r *nodeRepository) UpdateState(id uuid.UUID, to int) error {
	_, err := r.db.Exec(updateNodeState, id, to)
	return err
}

const getNddeState = `SELECT state FROM nodes WHERE id = $1`

func (r *nodeRepository) GetState(id uuid.UUID) (state int, err error) {
	return state, r.db.QueryRow(getNddeState, id).Scan(&state)
}

func (r *nodeRepository) FindFirstByInstanceId(id uuid.UUID) (*models.Node, error) {
	const findInstanceById = `
SELECT id, cloud, region, metadata, state, instance_id, whitelist FROM nodes WHERE id = $1 LIMIT 1
`
	node := new(models.Node)
	err := r.db.QueryRow(findInstanceById, id).Scan(&node.Id, &node.Cloud, &node.Region, &node.Metadata, &node.State, &node.InstanceId, pq.Array(&node.Whitelist))
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (r *nodeRepository) FindById(id uuid.UUID) (*models.Node, error) {
	const findInstanceById = `SELECT id, cloud, region, metadata, state, instance_id, whitelist FROM nodes WHERE id = $1`
	node := new(models.Node)
	err := r.db.QueryRow(findInstanceById, id).
		Scan(&node.Id, &node.Cloud, &node.Region, &node.Metadata, &node.State, &node.InstanceId, pq.Array(&node.Whitelist))
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (r *nodeRepository) FindByInstanceId(id uuid.UUID) ([]*models.Node, error) {
	const findInstanceById = `SELECT id, cloud, region, state, instance_id, metadata, whitelist FROM nodes WHERE instance_id = $1`
	rows, err := r.db.Query(findInstanceById, id)
	if err != nil {
		return nil, err
	}
	var nodes []*models.Node
	for rows.Next() {
		node := new(models.Node)
		if err := rows.Scan(&node.Id, &node.Cloud, &node.Region, &node.State, &node.InstanceId, &node.Metadata, pq.Array(&node.Whitelist)); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (r *nodeRepository) Insert(s *models.Node) error {
	const insert = `
INSERT INTO nodes(instance_id, cloud, region, metadata, state)
  VALUES($1, $2, $3, $4, $5)
  RETURNING id;
`
	return r.db.QueryRow(insert,
		s.InstanceId,
		s.Cloud,
		s.Region,
		s.Metadata,
		s.State,
	).Scan(&s.Id)
}

func (r *nodeRepository) Delete(id uuid.UUID) error {
	const delete = `DELETE FROM nodes WHERE id = $1 RETURNING id`
	return r.db.QueryRow(delete, id).Scan(&id)
}

func (r *nodeRepository) Orphan(id uuid.UUID) error {
	const update = `UPDATE nodes SET state = 3 WHERE id = $1 RETURNING id`
	return r.db.QueryRow(update, id).Scan(&id)
}

func (r *nodeRepository) UpdateMetadata(s *models.Node) error {
	_, err := r.db.Exec(updateNodeMetadata, s.Id, s.Metadata)
	return err
}

const updateNodeMetadata = `UPDATE nodes SET metadata = $2 WHERE id = $1`
const updateNodeWhitelist = `UPDATE nodes SET whitelist = $2 WHERE id = $1`

func (r *nodeRepository) UpdateWhitelist(s *models.Node) error {
	_, err := r.db.Exec(updateNodeWhitelist, s.Id, pq.Array(&s.Whitelist))
	return err
}
func NewNodeRepository(db *sql.DB) components.NodeRepository {
	return &nodeRepository{db}
}
