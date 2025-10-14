package Controller

import (
	"TaskManager/Package/Model"
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

func (Ctr *ControllerStruct) AddData(GinCtx *gin.Context) {
	var req AddTaskStruct

	err := GinCtx.ShouldBind(&req)
	if err != nil {
		GinCtx.JSON(http.StatusBadRequest, ErrorObjInitiator(err))
		return
	}

	dbPayload := Model.TaskStoreRequest{}

	dbPayload.Task_Description = req.Task_Description
	dbPayload.Title = req.Title
	dbPayload.Task_Status, err = strconv.ParseBool(req.Task_Status)

	if err != nil {
		GinCtx.JSON(http.StatusBadRequest, ErrorObjInitiator(err))
		return
	}

	errChannel := make(chan error, 1)
	defer close(errChannel)
	resChannel := make(chan Model.TaskStoreResponse, 1)
	defer close(resChannel)
	wg := sync.WaitGroup{}
	//var IsexecutionFinished bool = false

	wg.Add(1)

	go Ctr.Model.AddTask(dbPayload, &wg, resChannel, errChannel)

	go func() {

		for err := range errChannel {
			GinCtx.JSON(http.StatusBadRequest, ErrorObjInitiator(err))
			break
		}
	}()

	go func() {

		for resl := range resChannel {
			GinCtx.JSON(http.StatusOK, resl)
			break
		}
	}()

	wg.Wait()
	return

}

func (Ctr *ControllerStruct) GetData(GinCtx *gin.Context) {
	var req GetTaskStruct

	err := GinCtx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		GinCtx.JSON(http.StatusBadRequest, ErrorObjInitiator(err))
		return
	}

	dbPayload := Model.GetTask{}

	dbPayload.ID = req.ID

	errChannel := make(chan error, 1)
	defer close(errChannel)
	resChannel := make(chan Model.TaskStoreResponse, 1)
	defer close(resChannel)
	wg := sync.WaitGroup{}
	//var IsexecutionFinished bool = false

	wg.Add(1)

	go Ctr.Model.GetTask(dbPayload, &wg, resChannel, errChannel)

	go func() {
		for err := range errChannel {
			GinCtx.JSON(http.StatusBadRequest, ErrorObjInitiator(err))
			return
		}
	}()

	go func() {

		for resl := range resChannel {
			if resl.ID >= 1 {
				GinCtx.JSON(http.StatusOK, resl)
			} else {
				GinCtx.JSON(http.StatusNotFound, ErrorObjInitiator(errors.New("Data Not Found")))
			}
			return
		}
	}()

	wg.Wait()
	return

}

func (Ctr *ControllerStruct) EditData(GinCtx *gin.Context) {
	var req UpdateTaskStruct

	err := GinCtx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		GinCtx.JSON(http.StatusBadRequest, ErrorObjInitiator(err))
		return
	}
	updatedTask := Model.TaskStoreRequest{
		Title:            req.Title,
		Task_Description: req.Task_Description,
	}

	updatedTask.Task_Status, err = strconv.ParseBool(req.Task_Status)
	if err != nil {
		GinCtx.JSON(http.StatusBadRequest, ErrorObjInitiator(err))
		return
	}

	dbPayload := Model.UpdateTaskStoreRequest{}

	dbPayload.ID = req.ID
	dbPayload.Task = updatedTask

	errChannel := make(chan error, 1)
	defer close(errChannel)
	resChannel := make(chan Model.TaskStoreResponse, 1)
	defer close(resChannel)
	wg := sync.WaitGroup{}
	//var IsexecutionFinished bool = false

	wg.Add(1)

	go Ctr.Model.EditTask(dbPayload, &wg, resChannel, errChannel)

	go func() {

		for err := range errChannel {
			GinCtx.JSON(http.StatusBadRequest, ErrorObjInitiator(err))
			return
		}
	}()

	go func() {

		for resl := range resChannel {
			if resl.ID >= 1 {
				GinCtx.JSON(http.StatusOK, resl)
			}
			return
		}
	}()

	wg.Wait()
	return
}

func (Ctr *ControllerStruct) DeleteData(GinCtx *gin.Context) {
	var req DeleteTaskStruct

	err := GinCtx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		GinCtx.JSON(http.StatusBadRequest, ErrorObjInitiator(err))
		return
	}

	dbPayload := Model.DeleteTaskStoreRequest{}

	dbPayload.ID = req.ID
	dbPayload.Task = Model.TaskStoreRequest{}

	errChannel := make(chan error, 1)
	defer close(errChannel)
	resChannel := make(chan Model.DeleteTaskStoreResponse, 1)
	defer close(resChannel)
	wg := sync.WaitGroup{}

	wg.Add(1)

	go Ctr.Model.DeleteTask(dbPayload, &wg, resChannel, errChannel)

	go func() {

		for err := range errChannel {
			GinCtx.JSON(http.StatusBadRequest, ErrorObjInitiator(err))
			return
		}
	}()

	go func() {

		for resl := range resChannel {
			if resl.ID >= 1 {
				GinCtx.JSON(http.StatusOK, resl)
			}
			return
		}
	}()

	wg.Wait()
	return
}

func (Ctr *ControllerStruct) ListData(GinCtx *gin.Context) {
	var req ListTaskStruct
	taskList := []Model.TaskStoreResponse{}

	err := GinCtx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		GinCtx.JSON(http.StatusBadRequest, ErrorObjInitiator(err))
		return
	}

	dbPayload := Model.ListTaskStore{}

	dbPayload.Limit = req.Limit
	dbPayload.Offset = req.Offset
	dbPayload.Page = req.Page

	taskList, err = Ctr.Model.ListTask(dbPayload)

	if err != nil {
		GinCtx.JSON(http.StatusBadRequest, ErrorObjInitiator(err))
		return
	}

	GinCtx.JSON(http.StatusOK, taskList)
	return
}

func ErrorObjInitiator(Err error) *gin.H {
	return &gin.H{
		"error": Err.Error(),
	}
}
