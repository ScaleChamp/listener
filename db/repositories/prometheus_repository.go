package repositories

import (
	"database/sql"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
)

type prometheusRepository struct {
	db *sql.DB
}

func (c *prometheusRepository) Size() (int, error) {
	var v int
	return v, c.db.QueryRow("select count(*) from prometheus;").Scan(&v)
}

func (c *prometheusRepository) Get() ([]*models.Prometheus, error) {
	rows, err := c.db.Query("select id, labels, targets from prometheus;")
	if err != nil {
		return nil, err
	}
	prometheus := make([]*models.Prometheus, 0)
	for rows.Next() {
		p := new(models.Prometheus)
		if err := rows.Scan(&p.Id, &p.Labels, pq.Array(p.Targets)); err != nil {
			return nil, err
		}
		prometheus = append(prometheus, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return prometheus, nil
}

func (c *prometheusRepository) Insert(p *models.Prometheus) error {
	return c.db.QueryRow("INSERT INTO prometheus(labels, targets, node_id) VALUES($1, $2, $3) RETURNING id", &p.Labels, pq.Array(&p.Targets), &p.NodeId).Scan(&p.Id)
}

func (c *prometheusRepository) Delete(id uuid.UUID) error {
	return c.db.QueryRow("DELETE FROM prometheus WHERE node_id = $1 RETURNING node_id", id).Scan(&id)
}

func NewPrometheusRepository(db *sql.DB) components.PrometheusRepository {
	return &prometheusRepository{db}
}
