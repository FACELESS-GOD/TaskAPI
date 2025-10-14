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

type GetTaskStruct struct {
	ID int64 `json:"ID" binding:"required,min=1"`
}

type UpdateTaskStruct struct {
	ID               int64  `json:"ID" binding:"required,min=1"`
	Title            string `json:"Title" binding:"required"`
	Task_Description string `json:"Task_Description" binding:"required"`
	Task_Status      string `json:"Task_Status" binding:"required,oneof=true false"`
}

type DeleteTaskStruct struct {
	ID int64 `json:"ID" binding:"required,min=1"`
}

type ListTaskStruct struct {
	Limit  int64 `json:"Limit" binding:"required"`
	Page   int64 `json:"Page" binding:"required"`
	Offset int64 `json:"Offset"`
}

func NewController(Mdl Model.ModelInterface) ControllerStruct {
	ctrl := ControllerStruct{}
	router := gin.Default()

	router.POST(Route.PostURL, ctrl.AddData)
	router.GET(Route.GetURL, ctrl.GetData)
	router.PUT(Route.EditURL, ctrl.EditData)
	router.DELETE(Route.DeleteURL, ctrl.DeleteData)
	router.GET(Route.ListPaginationURL, ctrl.ListData)

	ctrl.Model = Mdl
	ctrl.router = router

	return ctrl
}

func (Ctrl *ControllerStruct) StartServer(Address string) error {
	return Ctrl.router.Run(Address)
}
