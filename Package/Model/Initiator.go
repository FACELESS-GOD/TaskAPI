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

func NewModel(Configuration Configurator.ConfiguratorStruct) ModelStruct {
	txOption := sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}
	return ModelStruct{
		Config:   Configuration,
		TxOption: txOption,
	}
}

const AddTaskQuery string = `
INSERT INTO TaskStore (
  Title, Task_Description
) VALUES (
  ? , ? 
)
;
`

func (Model *ModelStruct) AddTask(Task TaskStoreRequest) (TaskStoreResponse, error) {

	if Task.Task_Status != true {
		return TaskStoreResponse{}, errors.New("invalid data.")
	} else if len(Task.Title) <= 0 || len(Task.Task_Description) <= 0 {
		return TaskStoreResponse{}, errors.New("invalid data.")
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Model.Config.SqlDBConn.BeginTx(ctx, &Model.TxOption)

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

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return TaskStoreResponse{}, errors.New(nerr)
		} else {
			return TaskStoreResponse{}, err
		}
	}

	resp := TaskStoreResponse{
		ID:   taskID,
		Task: Task,
	}
	return resp, nil

}

const EditTaskQuery string = `
UPDATE TaskStore 
SET Title = ? , Task_Description = ? ,Edited_On = CURRENT_TIMESTAMP()
WHERE ID = ? 
;
`

func (Model *ModelStruct) EditTask(Task UpdateTaskStoreRequest) (TaskStoreResponse, error) {

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*100)
	defer cancelFunc()

	db, err := Model.Config.SqlDBConn.BeginTx(ctx, &Model.TxOption)

	if err != nil {
		return TaskStoreResponse{}, err
	}

	resp, err := db.ExecContext(ctx, EditTaskQuery, Task.Task.Title, Task.Task.Task_Description, Task.ID)

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
		err := db.Rollback()
		if err != nil {
			return TaskStoreResponse{}, err
		}
	}

	err = db.Commit()

	if err != nil {
		return TaskStoreResponse{}, err
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

func (Model *ModelStruct) DeleteTask(Task DeleteTaskStoreRequest) (DeleteTaskStoreResponse, error) {

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	db, err := Model.Config.SqlDBConn.BeginTx(ctx, &Model.TxOption)

	if err != nil {
		return DeleteTaskStoreResponse{}, err
	}

	resp, err := db.ExecContext(ctx, DeleteTaskQuery, Task.ID)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			return DeleteTaskStoreResponse{}, nerr
		}
	}

	numRowAffected, err := resp.RowsAffected()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			return DeleteTaskStoreResponse{}, nerr
		}
	}

	if numRowAffected > 1 || numRowAffected <= 0 {
		err := db.Rollback()
		if err != nil {
			return DeleteTaskStoreResponse{}, err
		}
	}

	errMessage := db.Commit()

	if errMessage != nil {
		return DeleteTaskStoreResponse{}, errMessage
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

func (Model *ModelStruct) ListTask(Task ListTaskStore) ([]TaskStoreResponse, error) {

	respList := []TaskStoreResponse{}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	db, err := Model.Config.SqlDBConn.BeginTx(ctx, &Model.TxOption)

	if err != nil {
		return respList, err
	}

	if Task.Limit < 1 || Task.Page < 1 {
		Task.Limit = 10
		Task.Page = 1
	}

	Task.Offset = (Task.Page - 1) * Task.Limit

	resp, err := db.QueryContext(ctx, ListTaskQuery, Task.Offset, Task.Limit)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			return respList, nerr
		}
	}

	for resp.Next() {
		var taskResp TaskStoreResponse
		var task TaskStoreRequest
		t := ""
		e := ""
		if err := resp.Scan(
			&taskResp.ID,
			&task.Title,
			&task.Task_Description,
			&task.Task_Status,
			&t,
			&e,
		); err != nil {
			return nil, err
		}
		taskResp.Task = task
		respList = append(respList, taskResp)
	}

	errMessage := db.Commit()

	if errMessage != nil {
		return nil, errMessage
	}

	return respList, nil

}
