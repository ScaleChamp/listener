package repositories

import (
	"database/sql"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
)

type taskRepository struct {
	db *sql.DB
}

func (t *taskRepository) Get(taskId uuid.UUID) (*models.Task, error) {
	task := new(models.Task)
	task.Data = make(map[string]models.Stage)
	task.Metadata = make(map[string]json.RawMessage)
	if err := t.db.QueryRow("SELECT id, kind, action, state, metadata, data, instance_id FROM tasks WHERE id = $1", taskId).
		Scan(
			&task.Id,
			&task.Kind,
			&task.Action,
			&task.State,
			&task.Metadata,
			&task.Data,
			&task.InstanceId,
		); err != nil {
		return nil, err
	}
	return task, nil
}

func (t *taskRepository) Update(task *models.Task) error {
	return t.db.
		QueryRow("UPDATE tasks SET state = $2, metadata = $3, data = $4 WHERE id = $1 RETURNING id", task.Id, task.State, task.Metadata, task.Data).
		Scan(&task.Id)
}

func (t *taskRepository) Finish(task *models.Task) error {
	return t.db.
		QueryRow("UPDATE tasks SET state = $2, metadata = $3, data = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING id", task.Id, task.State, task.Metadata, task.Data).
		Scan(&task.Id)
}

func NewTaskRepository(db *sql.DB) components.TaskRepository {
	return &taskRepository{db: db}
}
