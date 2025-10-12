package Controller

import (
	"TaskManager/Helper/Route"
	"TaskManager/Package/Model"

	"github.com/gin-gonic/gin"
)

type ControllerInterface interface {
	AddData(GinCtx *gin.Context)
	EditData(GinCtx *gin.Context)
	DeleteData(GinCtx *gin.Context)
	ListData(GinCtx *gin.Context)
	GetData(GinCtx *gin.Context)
}

type ControllerStruct struct {
	Model  Model.ModelInterface
	router *gin.Engine
}

type AddTaskStruct struct {
	Title            string `json:"Title" binding:"required"`
	Task_Description string `json:"Task_Description" binding:"required"`
	Task_Status      string `json:"Task_Status" binding:"required,oneof=true false"`
}

type GetTask struct {
	ID int64 `json:"ID" binding:"required,min=1"`
}

func NewController(Mdl Model.ModelInterface) ControllerStruct {
	ctrl := ControllerStruct{}
	router := gin.Default()

	router.POST(Route.PostURL, ctrl.AddData)
	router.GET(Route.GetURL, ctrl.GetData)

	ctrl.Model = Mdl
	ctrl.router = router

	return ctrl
}

func (Ctrl *ControllerStruct) StartServer(Address string) error {
	return Ctrl.router.Run(Address)
}
