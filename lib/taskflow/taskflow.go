package taskflow

import (
	"bytes"
	"encoding/json"
	"errors"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
	"log"
	"reflect"
)

type Taskflow struct {
	taskRepository components.TaskRepository

	task *models.Task

	undoStack []*models.Step
}

var ErrAlreadyFailed = errors.New("already failed to undo")

var ErrUndo = errors.New("please undo")

func (tf *Taskflow) Do(cmd *models.Step) error {
	tf.undoStack = append(tf.undoStack, cmd)

	stage, ok := tf.task.Data[cmd.Name]
	if ok {
		if stage.Done {
			return nil
		}
		if stage.Failed {
			return ErrUndo
		}
	}

	f := reflect.ValueOf(cmd.Cmd).MethodByName("Do")

	functor := f.Type()
	var args []reflect.Value
	for i := 0; i < functor.NumIn(); i += 1 {
		inV := functor.In(i)
		f := reflect.New(inV)
		if i == functor.NumIn()-1 {
			if functor.IsVariadic() {
				var wants [][]byte
				for d := i; d < len(cmd.Wants); d += 1 {
					wants = append(wants, tf.task.Metadata[cmd.Wants[i]])
				}
				directWants := append([]byte("["), bytes.Join(wants, []byte(","))...)
				directWants = append(directWants, ']')
				if err := json.Unmarshal(directWants, f.Interface()); err != nil {
					panic(err)
				}
			} else {
				if err := json.Unmarshal(tf.task.Metadata[cmd.Wants[i]], f.Interface()); err != nil {
					panic(err)
				}
			}
		} else {
			if err := json.Unmarshal(tf.task.Metadata[cmd.Wants[i]], f.Interface()); err != nil {
				panic(err)
			}
		}

		args = append(args, f.Elem())
	}

	var results []reflect.Value
	if functor.IsVariadic() {
		results = f.CallSlice(args)
	} else {
		results = f.Call(args)
	}

	// set results anyway
	for i := 0; i < len(cmd.Provides); i += 1 {
		data, err := json.Marshal(results[i].Interface())
		if err != nil {
			panic(err)
		}
		tf.task.Metadata[cmd.Provides[i]] = data
	}

	if err := results[len(results)-1]; !err.IsNil() {
		log.Printf("taskflow err: %s\n", err)
		tf.task.Data[cmd.Name] = models.Stage{
			Failed: true,
		}

		if err := tf.taskRepository.Update(tf.task); err != nil {
			return err
		}
		return ErrUndo
	}
	tf.task.Data[cmd.Name] = models.Stage{
		Done: true,
	}
	return tf.taskRepository.Update(tf.task)
}

func (tf *Taskflow) Undo() error {
	for i := len(tf.undoStack) - 1; i >= 0; i -= 1 {
		cmd := tf.undoStack[i]

		stage := tf.task.Data[cmd.Name]

		if stage.Undone {
			continue
		} else if stage.UndoneFailed {
			return ErrAlreadyFailed
		}

		undo := reflect.ValueOf(cmd.Cmd).MethodByName("Undo")
		// allow to skip undo function in method
		//if undo.IsNil() {
		//	stage.Undone = true
		//	tf.task.Data[cmd.Name] = stage
		//	if err := tf.taskRepository.Update(tf.task); err != nil {
		//		return err
		//	}
		//	continue
		//}
		undoType := undo.Type()
		do := reflect.ValueOf(cmd.Cmd).MethodByName("Do").Type()
		var args []reflect.Value

		for i := 0; i < do.NumOut()-1; i += 1 {
			inV := do.Out(i)
			f := reflect.New(inV)
			_ = json.Unmarshal(tf.task.Metadata[cmd.Provides[i]], f.Interface())
			args = append(args, f.Elem())
		}

		for i := 0; i < do.NumIn(); i += 1 {
			inV := do.In(i)

			f := reflect.New(inV)

			if i == do.NumIn()-1 {
				if do.IsVariadic() {
					var wants [][]byte
					for d := i; d < len(cmd.Wants); d += 1 {
						wants = append(wants, tf.task.Metadata[cmd.Wants[i]])
					}
					directWants := append([]byte("["), bytes.Join(wants, []byte(","))...)
					directWants = append(directWants, ']')
					if err := json.Unmarshal(directWants, f.Interface()); err != nil {
						panic(err)
					}
				} else if err := json.Unmarshal(tf.task.Metadata[cmd.Wants[i]], f.Interface()); err != nil {
					panic(err)
				}
			} else {
				if err := json.Unmarshal(tf.task.Metadata[cmd.Wants[i]], f.Interface()); err != nil {
					panic(err)
				}
			}

			args = append(args, f.Elem())
		}

		var results []reflect.Value
		if undoType.IsVariadic() {
			results = undo.CallSlice(args)
		} else {
			results = undo.Call(args)
		}

		err := results[len(results)-1]
		if !err.IsNil() {
			stage.UndoneFailed = true
			tf.task.Data[cmd.Name] = stage
			return tf.taskRepository.Update(tf.task)
		}
		stage.Undone = true
		tf.task.Data[cmd.Name] = stage
		if err := tf.taskRepository.Update(tf.task); err != nil {
			return err
		}
	}
	return nil
}

func (tf *Taskflow) Execute(commands []*models.Step) error {
	for _, cmd := range commands {
		if err := tf.Do(cmd); err == ErrUndo {
			if err := tf.Undo(); err != nil {
				log.Println("taskflow undo failed:", err)
				tf.task.State = 2
				return tf.taskRepository.Finish(tf.task)
			}
			tf.task.State = 2
			return tf.taskRepository.Finish(tf.task)
		} else if err != nil {
			tf.task.State = 2
			return tf.taskRepository.Finish(tf.task)
		}
	}
	tf.task.State = 1
	return tf.taskRepository.Finish(tf.task)
}

func newTaskflow(task *models.Task, taskRepository components.TaskRepository) *Taskflow {
	return &Taskflow{
		taskRepository: taskRepository,
		task:           task,
	}
}

type Engine struct {
	taskRepository components.TaskRepository
}

func (c *Engine) NewTaskFlow(task *models.Task) *Taskflow {
	return newTaskflow(task, c.taskRepository)
}

func NewEngine(taskRepository components.TaskRepository) *Engine {
	return &Engine{taskRepository}
}
