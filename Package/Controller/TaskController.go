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
}

type ControllerStruct struct {
	Model  Model.ModelInterface
	router *gin.Engine
}

func NewController(Mdl Model.ModelInterface) ControllerStruct {
	ctrl := ControllerStruct{}
	router := gin.Default()

	router.POST(Route.PostURL, ctrl.AddData)

	ctrl.Model = Mdl
	ctrl.router = router

	return ctrl
}

func (Ctrl *ControllerStruct) StartServer(Address string) error {
	return Ctrl.router.Run(Address)
}
