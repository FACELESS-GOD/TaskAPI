package Model

import (
	"TaskManager/Helper/Startup"
	"TaskManager/Package/Configurator"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SuiteStruct struct {
	suite.Suite
	Model ModelStruct
}

func (Suite *SuiteStruct) SetupSuite() {
	config := Configurator.NewConfigurator()
	config.LoadConfig(Startup.DebugMode)
	config.LoadDBInstance()

	model := NewModel(*config)
	Suite.Model = model
}

func (Suite *SuiteStruct) TearDownSuite() {
	Suite.Model.Config.SqlDBConn.Close()
}

func (Suite *SuiteStruct) TestAddTask() {

}

func (Suite *SuiteStruct) TestEditTask() {

}

func (Suite *SuiteStruct) TestDeleteTask() {

}

func (Suite *SuiteStruct) TestListTask() {

}

func TestSuite(Testor *testing.T) {
	suite.Run(Testor, new(SuiteStruct))
}
