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
	AddTask(Task TaskStoreRequest, Wg *sync.WaitGroup, ResultChannel chan<- TaskStoreResponse, ErrorChannel chan<- error)
	EditTask(Task UpdateTaskStoreRequest, Wg *sync.WaitGroup, ResultChannel chan<- TaskStoreResponse, ErrorChannel chan<- error)
	DeleteTask(Task DeleteTaskStoreRequest, Wg *sync.WaitGroup, ResultChannel chan<- DeleteTaskStoreResponse, ErrorChannel chan<- error)
	ListTask(ListTask ListTaskStore, WaitGroup *sync.WaitGroup)
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

func (Model *ModelStruct) AddTask(Task TaskStoreRequest, Wg *sync.WaitGroup, ResultChannel chan<- TaskStoreResponse, ErrorChannel chan<- error) {

	defer Wg.Done()

	isValid, errorMessage := Model.ValidateParamAddTask(Task)

	if isValid == true {
		errorObj := errors.New(errorMessage)
		ErrorChannel <- errorObj
		return
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Model.Config.SqlDBConn.BeginTx(ctx, &Model.TxOption)

	if err != nil {
		ErrorChannel <- err
		return
	}

	res, err := db.ExecContext(ctx, AddTaskQuery, Task.Title, Task.Task_Description)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			ErrorChannel <- nerr
			return
		} else {
			ErrorChannel <- err
			return
		}
	}

	taskID, err := res.LastInsertId()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			ErrorChannel <- nerr
			return
		} else {
			ErrorChannel <- err
			return
		}
	}

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			ErrorChannel <- nerr
			return
		} else {
			ErrorChannel <- err
			return
		}
	}

	resp := TaskStoreResponse{
		ID:   taskID,
		Task: Task,
	}

	ResultChannel <- resp
	return

}

func (Model *ModelStruct) ValidateParamAddTask(Task TaskStoreRequest) (bool, string) {
	errorMessages := []string{}
	var isValid bool = false
	if Task.Task_Status != true {
		isValid = true
		errorMessages = append(errorMessages, "Invalid Task Status")

	}
	if len(Task.Title) <= 0 {

		isValid = true
		errorMessages = append(errorMessages, "Invalid Title .")

	}

	if len(Task.Task_Description) <= 0 {

		isValid = true
		errorMessages = append(errorMessages, "Invalid Description")
	}

	errorMessage := ""

	for _, message := range errorMessages {
		errorMessage = errorMessage + message + " , "
	}

	return isValid, errorMessage

}

const EditTaskQuery string = `
UPDATE TaskStore 
SET Title = ? , Task_Description = ? ,Edited_On = CURRENT_TIMESTAMP()
WHERE ID = ? 
;
`

func (Model *ModelStruct) ValidateParamEditTask(Task UpdateTaskStoreRequest) (bool, string) {
	var IsValid bool = false
	errMessages := []string{}
	errorMessage := ""

	if Task.ID < 1 {
		IsValid = true
		errMessages = append(errMessages, "Invalid ID!")
	}

	validity, message := Model.ValidateParamAddTask(Task.Task)

	if validity == true {
		IsValid = validity
		errMessages = append(errMessages, message)
	}

	for _, message := range errMessages {
		errorMessage = errorMessage + message + " , "
	}

	return IsValid, errorMessage
}

func (Model *ModelStruct) EditTask(Task UpdateTaskStoreRequest, Wg *sync.WaitGroup, ResultChannel chan<- TaskStoreResponse, ErrorChannel chan<- error) {
	defer Wg.Done()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*100)
	defer cancelFunc()

	isValid, message := Model.ValidateParamEditTask(Task)

	if isValid == true {
		errObj := errors.New(message)
		ErrorChannel <- errObj
	}

	db, err := Model.Config.SqlDBConn.BeginTx(ctx, &Model.TxOption)

	if err != nil {
		ErrorChannel <- err
		return
	}

	resp, err := db.ExecContext(ctx, EditTaskQuery, Task.Task.Title, Task.Task.Task_Description, Task.ID)

	if err != nil {
		rollBackErr := db.Rollback()
		if rollBackErr != nil {
			ErrorChannel <- rollBackErr
			return
		} else {
			ErrorChannel <- err
			return
		}
	}

	numRowAffected, err := resp.RowsAffected()

	if err != nil {
		rollBackErr := db.Rollback()
		if rollBackErr != nil {
			ErrorChannel <- rollBackErr
			return
		} else {
			ErrorChannel <- err
			return
		}
	}

	if numRowAffected > 1 || numRowAffected <= 0 {
		rollBackErr := db.Rollback()
		if rollBackErr != nil {
			ErrorChannel <- rollBackErr
			return
		} else {
			ErrorChannel <- err
			return
		}
	}

	err = db.Commit()

	if err != nil {
		ErrorChannel <- err
		return
	}

	reslt := TaskStoreResponse{
		ID:   Task.ID,
		Task: Task.Task,
	}

	ResultChannel <- reslt
	return

}

const DeleteTaskQuery string = `
UPDATE TaskStore 
SET Task_Status = false ,Edited_On = CURRENT_TIMESTAMP()
WHERE ID = ? 
;
`

func (Model *ModelStruct) DeleteTask(Task DeleteTaskStoreRequest, Wg *sync.WaitGroup, ResultChannel chan<- DeleteTaskStoreResponse, ErrorChannel chan<- error) {

	defer Wg.Done()

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	db, err := Model.Config.SqlDBConn.BeginTx(ctx, &Model.TxOption)

	if err != nil {

		ErrorChannel <- err
		return
	}

	resp, err := db.ExecContext(ctx, DeleteTaskQuery, Task.ID)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			ErrorChannel <- nerr
			return
		} else {
			ErrorChannel <- err
			return
		}
	}

	numRowAffected, err := resp.RowsAffected()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			ErrorChannel <- nerr
			return
		} else {
			ErrorChannel <- err
			return
		}
	}

	if numRowAffected > 1 || numRowAffected <= 0 {
		nerr := db.Rollback()
		if nerr != nil {
			ErrorChannel <- nerr
			return
		} else {
			ErrorChannel <- err
			return
		}
	}

	errMessage := db.Commit()

	if errMessage != nil {
		ErrorChannel <- errMessage
		return
	}

	resl := DeleteTaskStoreResponse{
		ID:     Task.ID,
		Status: true,
		Task:   Task.Task,
	}

	ResultChannel <- resl

	return

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
