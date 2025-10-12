package Controller

import (
	"TaskManager/Package/Model"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type AddTaskStruct struct {
	Title            string `json:"Title" binding:"required"`
	Task_Description string `json:"Task_Description" binding:"required"`
	Task_Status      string `json:"Task_Status" binding:"required,oneof=true false"`
}

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
			return
		}
	}()

	go func() {
		for resl := range resChannel {
			GinCtx.JSON(http.StatusOK, resl)
			return
		}
	}()

	wg.Wait()
	return

}

func (Ctr *ControllerStruct) EditData(GinCtx *gin.Context) {}

func (Ctr *ControllerStruct) DeleteData(GinCtx *gin.Context) {}

func (Ctr *ControllerStruct) ListData(GinCtx *gin.Context) {}

func ErrorObjInitiator(Err error) *gin.H {
	return &gin.H{
		"error": Err.Error(),
	}
}
