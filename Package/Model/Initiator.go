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
	DeleteTask(TaskID int64, WaitGroup *sync.WaitGroup) (bool, error)
	ListTask(WaitGroup *sync.WaitGroup) ([]TaskStoreResponse, error)
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
INSERT INTO Bank_Transfers (
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
UPDATE Bank_Transfers 
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
UPDATE Bank_Transfers 
SET Task_Status = false ,Edited_On = CURRENT_TIMESTAMP()
WHERE ID = ? 
;
`


func (Model *ModelStruct) DeleteTask(TaskID int64, WaitGroup *sync.WaitGroup) (bool, error) {
	defer WaitGroup.Done()

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	db, err := Model.Config.SqlDBConn.BeginTx(ctx, &Model.TxOption)

	if err != nil {
		return false, err
	}

	resp, err := db.ExecContext(ctx, DeleteTaskQuery, TaskID)

	
	if err != nil {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return false, errors.New(nerr)
		} else {
			return false, err
		}
	}

	numRowAffected, err := resp.RowsAffected()

	if err != nil {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return false, errors.New(nerr)
		} else {
			return false, err
		}
	}

	if numRowAffected > 1 || numRowAffected <= 0 {
		nerr := db.Rollback().Error()
		if len(nerr) >= 1 {
			return false, errors.New(nerr)
		} else {
			return false, errors.New("Invalid Query!")
		}
	}

	errMessage := db.Commit().Error()

	if len(errMessage) >= 1 {
		return false, errors.New(errMessage)
	}

	return true , nil 

}
func (Model *ModelStruct) ListTask(WaitGroup *sync.WaitGroup) ([]TaskStoreResponse, error) {
	defer WaitGroup.Done()
}
