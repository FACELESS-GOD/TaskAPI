package Controller

import "TaskManager/Package/Model"

type ControllerInterface interface {
	AddData()
	EditData()
	DeleteData()
	ListData()
}

type ControllerStruct struct {
	Model Model.ModelInterface
}

func NewModel() ControllerStruct {
	return ControllerStruct{}
}

func (Ctr *ControllerStruct) AddData() {}

func (Ctr *ControllerStruct) EditData() {}

func (Ctr *ControllerStruct) DeleteData() {}

func (Ctr *ControllerStruct) ListData() {}
