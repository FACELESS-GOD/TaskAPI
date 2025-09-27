package Model

import (
	"TaskManager/Package/Configurator"
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"
)

type ModelInterface interface {
	AddTask(Task TaskStoreRequest, WaitGroup *sync.WaitGroup) (TaskStoreResponse, error)
	EditTask(Task UpdateTaskStoreRequest, WaitGroup *sync.WaitGroup) (TaskStoreResponse, error)
	DeleteTask(Task DeleteTaskStoreRequest, WaitGroup *sync.WaitGroup) (DeleteTaskStoreResponse, error)
	ListTask(ListTask ListTaskStore, WaitGroup *sync.WaitGroup) ([]TaskStoreResponse, error)
}

type ModelStruct struct {
	Config   Configurator.ConfiguratorStruct
	TxOption sql.TxOptions
}

func NewModel(Configuration Configurator.ConfiguratorStruct) *ModelStruct {
	txOption := sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}
	return &ModelStruct{
		Config:   Configuration,
		TxOption: txOption,
	}
}

const AddTaskQuery string = `
INSERT INTO TaskStore (
  Title, Task_Description
) VALUES (
  ? , ? , ?, ?
)
;
`

func (Model *ModelStruct) AddTask(Task TaskStoreRequest, WaitGroup *sync.WaitGroup) (TaskStoreResponse, error) {
	defer WaitGroup.Done()

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)

	db, err := Model.Config.SqlDBConn.BeginTx(Model.Config.DbCtx, &Model.TxOption)

	if err != nil {
		return TaskStoreResponse{}, err
	}

	res, err := db.ExecContext(ctx, AddTaskQuery, Task.Title, Task.Task_Description)

	if err != nil {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return TaskStoreResponse{}, errors.New(nerr)
		} else {
			return TaskStoreResponse{}, err
		}
	}

	taskID, err := res.LastInsertId()

	if err != nil {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return TaskStoreResponse{}, errors.New(nerr)
		} else {
			return TaskStoreResponse{}, err
		}
	}

	errMessage := db.Commit().Error()

	if len(errMessage) >= 1 {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return TaskStoreResponse{}, errors.New(nerr)
		} else {
			return TaskStoreResponse{}, errors.New(errMessage)
		}
	}

	resp := TaskStoreResponse{
		ID:   taskID,
		Task: Task,
	}

	defer cancelFunc()

	return resp, nil

}

const EditTaskQuery string = `
UPDATE TaskStore 
SET Title = ? , Task_Description = ? ,Task_Status = ? ,Edited_On = CURRENT_TIMESTAMP()
WHERE ID = ? 
;
`

func (Model *ModelStruct) EditTask(Task UpdateTaskStoreRequest, WaitGroup *sync.WaitGroup) (TaskStoreResponse, error) {
	defer WaitGroup.Done()

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	db, err := Model.Config.SqlDBConn.BeginTx(ctx, &Model.TxOption)

	if err != nil {
		return TaskStoreResponse{}, err
	}

	resp, err := db.ExecContext(ctx, EditTaskQuery, Task.Task.Title, Task.Task.Task_Description, Task.Task.Task_Status, Task.ID)

	if err != nil {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return TaskStoreResponse{}, errors.New(nerr)
		} else {
			return TaskStoreResponse{}, err
		}
	}

	numRowAffected, err := resp.RowsAffected()

	if err != nil {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return TaskStoreResponse{}, errors.New(nerr)
		} else {
			return TaskStoreResponse{}, err
		}
	}

	if numRowAffected > 1 || numRowAffected <= 0 {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return TaskStoreResponse{}, errors.New(nerr)
		} else {
			return TaskStoreResponse{}, errors.New("Invalid Query!")
		}
	}

	errMessage := db.Commit().Error()

	if len(errMessage) >= 1 {
		return TaskStoreResponse{}, errors.New(errMessage)
	}

	return TaskStoreResponse{
		ID:   Task.ID,
		Task: Task.Task,
	}, nil

}

const DeleteTaskQuery string = `
UPDATE TaskStore 
SET Task_Status = false ,Edited_On = CURRENT_TIMESTAMP()
WHERE ID = ? 
;
`

func (Model *ModelStruct) DeleteTask(Task DeleteTaskStoreRequest, WaitGroup *sync.WaitGroup) (DeleteTaskStoreResponse, error) {
	defer WaitGroup.Done()

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	db, err := Model.Config.SqlDBConn.BeginTx(ctx, &Model.TxOption)

	if err != nil {
		return DeleteTaskStoreResponse{}, err
	}

	resp, err := db.ExecContext(ctx, DeleteTaskQuery, Task.ID)

	if err != nil {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return DeleteTaskStoreResponse{}, errors.New(nerr)
		} else {
			return DeleteTaskStoreResponse{}, err
		}
	}

	numRowAffected, err := resp.RowsAffected()

	if err != nil {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return DeleteTaskStoreResponse{}, errors.New(nerr)
		} else {
			return DeleteTaskStoreResponse{}, err
		}
	}

	if numRowAffected > 1 || numRowAffected <= 0 {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return DeleteTaskStoreResponse{}, errors.New(nerr)
		} else {
			return DeleteTaskStoreResponse{}, errors.New("Invalid Query!")
		}
	}

	errMessage := db.Commit().Error()

	if len(errMessage) >= 1 {
		return DeleteTaskStoreResponse{}, errors.New(errMessage)
	}

	return DeleteTaskStoreResponse{
		ID:     Task.ID,
		Status: true,
		Task:   Task.Task,
	}, nil

}

const ListTaskQuery string = `
SELECT * FROM TaskStore
LIMIT ?, ?
;
`

func (Model *ModelStruct) ListTask(Task ListTaskStore, WaitGroup *sync.WaitGroup) ([]TaskStoreResponse, error) {
	defer WaitGroup.Done()

	respList := []TaskStoreResponse{}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	db, err := Model.Config.SqlDBConn.BeginTx(ctx, &Model.TxOption)

	if err != nil {
		return respList, err
	}

	resp, err := db.QueryContext(ctx, ListTaskQuery, Task.Offset, Task.Limit)

	if err != nil {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return respList, errors.New(nerr)
		} else {
			return respList, err
		}
	}

	for resp.Next() {
		var taskResp TaskStoreResponse
		var task TaskStoreRequest

		if err := resp.Scan(
			&taskResp.ID,
			&task.Title,
			&task.Task_Description,
			&task.Task_Status,
		); err != nil {
			return nil, err
		}
		taskResp.Task = task
		respList = append(respList, taskResp)
	}

	errMessage := db.Commit().Error()

	if len(errMessage) >= 1 {
		return nil, errors.New(errMessage)
	}

	return respList, nil

}
