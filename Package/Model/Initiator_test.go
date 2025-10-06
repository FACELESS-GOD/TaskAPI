package Model

import (
	"TaskManager/Helper/Startup"
	"TaskManager/Package/Configurator"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SuiteStruct struct {
	suite.Suite
	Model     ModelStruct
	RespStore []TaskStoreResponse
}

func (Suite *SuiteStruct) SetupSuite() {
	config := Configurator.NewConfigurator()
	config.LoadConfig(Startup.DebugMode)
	config.LoadDBInstance()
	//fmt.Println(config.DbConnString)
	//fmt.Println(config.DbDriver)

	model := NewModel(*config)
	Suite.Model = model
	Suite.RespStore = make([]TaskStoreResponse, 10)
}

func (Suite *SuiteStruct) TearDownSuite() {
	Suite.Model.Config.SqlDBConn.Close()
	Suite.RespStore = nil
}

func (Suite *SuiteStruct) SetupTest() {

	//	Suite.RespStore = make([]TaskStoreResponse, 10)
	Suite.RespStore = nil
}

func (Suite *SuiteStruct) TearDownTest() {
	Suite.RespStore = nil
}

// func (Model *ModelStruct) AddTask(Task TaskStoreRequest, WaitGroup *sync.WaitGroup) (TaskStoreResponse, error)
func (Suite *SuiteStruct) TestAddTask() {

	wg := sync.WaitGroup{}
	wg.Add(1)
	resp, err := Suite.Model.AddTask(TaskStoreRequest{
		Title:            "New",
		Task_Description: "This",
		Task_Status:      true,
	}, &wg)
	wg.Wait()

	Suite.Suite.NoError(err, "Error Has occured in a Positive Test Case")
	Suite.Suite.True(resp.ID > 0, "Error Has occured in a Positive Test Case")

	for i := 0; i < 10; i++ {
		wg.Add(1)
		resp, err := Suite.Model.AddTask(TaskStoreRequest{
			Title:            strconv.Itoa(i),
			Task_Description: strconv.Itoa(i * 2),
			Task_Status:      true,
		}, &wg)

		Suite.Suite.NoError(err, "Error Has occured in a Concurent Test Case")

		Suite.RespStore = append(Suite.RespStore, resp)
	}
	wg.Wait()

	for i := 0; i < 10; i++ {
		currResp := Suite.RespStore[i]
		Suite.Suite.True(currResp.ID >= 1, fmt.Sprintf("Error Has occured in a Concurent Test Case %v", currResp.ID))
	}

	wg.Add(1)
	resp, err = Suite.Model.AddTask(TaskStoreRequest{
		Title:            "New Task Negative",
		Task_Description: "This is a new Task Negative.",
		Task_Status:      false,
	}, &wg)
	wg.Wait()

	Suite.Suite.Error(err, "Error Has not occured in a Negative Test Case")

	wg.Add(1)
	resp, err = Suite.Model.AddTask(TaskStoreRequest{
		Title:            "",
		Task_Description: "This is a new Task Negative.",
		Task_Status:      false,
	}, &wg)
	wg.Wait()

	Suite.Suite.Error(err, "Error Has not occured in a Negative Test Case")

	wg.Add(1)
	resp, err = Suite.Model.AddTask(TaskStoreRequest{
		Title:            "New Task Negative",
		Task_Description: "",
		Task_Status:      false,
	}, &wg)
	wg.Wait()

	Suite.Suite.Error(err, "Error Has not occured in a Negative Test Case")

}

// func (Model *ModelStruct) EditTask(Task UpdateTaskStoreRequest, WaitGroup *sync.WaitGroup) (TaskStoreResponse, error) {
func (Suite *SuiteStruct) TestEditTask() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	resp, err := Suite.Model.AddTask(TaskStoreRequest{
		Title:            "New",
		Task_Description: "This",
		Task_Status:      true,
	}, &wg)
	wg.Wait()
	Suite.Suite.NoError(err, "Error occured even before the edit Test Case Start! .")
	newTask := resp.Task
	newTask.Task_Description = strconv.Itoa(100)

	wg.Add(1)
	resp, err = Suite.Model.EditTask(UpdateTaskStoreRequest{
		ID:   resp.ID,
		Task: newTask,
	}, &wg)
	wg.Wait()
	Suite.Suite.NoError(err, "Error occured in First Edit.")

	respList := make([]TaskStoreResponse, 10)

	wg.Wait()

	for i := 0; i < 10; i++ {
		newTask.Task_Description = strconv.Itoa(i)
		curr := UpdateTaskStoreRequest{
			ID:   resp.ID,
			Task: newTask,
		}
		wg.Add(1)
		resp, err = Suite.Model.EditTask(curr, &wg)
		Suite.Suite.NoError(err, "Error occured in Concurent file.")
		if err == nil {
			respList = append(respList, resp)
		}
	}
	wg.Wait()

	for i := 0; i < 9; i++ {
		old := resp
		newres := respList[i+10]
		//	fmt.Println(newres.ID, " f ", old.ID)
		Suite.Suite.Equal(newres.ID, old.ID, "Error Has occurred ID don't match")
		Suite.Suite.NotEqual(newres.Task.Task_Description, old.Task.Task_Description, "Error Has occurred ID don't match")
	}

}

func (Suite *SuiteStruct) TestDeleteTask() {

	wg := sync.WaitGroup{}
	wg.Add(1)
	resp, err := Suite.Model.AddTask(TaskStoreRequest{
		Title:            "New",
		Task_Description: "This",
		Task_Status:      true,
	}, &wg)
	wg.Wait()
	Suite.Suite.NoError(err, "Error occured even before the edit Test Case Start! .")

	wg.Add(1)
	del_resp, err := Suite.Model.DeleteTask(DeleteTaskStoreRequest{
		ID:   resp.ID,
		Task: resp.Task,
	}, &wg)

	wg.Wait()
	Suite.Suite.NoError(err, "Error occue while Positive delete.")
	Suite.Suite.NotNil(del_resp, "Responce was nil  while Positive delete.")

	Suite.Suite.Equal(resp.ID, del_resp.ID, " Response was wrong while Positive delete.")

	wg.Add(1)
	neg_val, err := Suite.Model.DeleteTask(DeleteTaskStoreRequest{
		ID:   del_resp.ID,
		Task: del_resp.Task,
	}, &wg)

	Suite.Suite.Error(err, "Error should not have occured , Negative Test Case Failed!")
	Suite.Suite.NotNil(neg_val, "Responce in wrong Negative Test Case Failed!")

}

func (Suite *SuiteStruct) TestListTask() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	resp, err := Suite.Model.ListTask(ListTaskStore{
		Limit: 10,
		Page:  1,
	}, &wg)
	wg.Wait()

	Suite.Suite.NoError(err, "Error Occured")
	respList := []TaskStoreResponse{}
	respList = append(respList, resp...)

	for _, task := range respList {
		Suite.Suite.True(task.ID > 0, "Defective Data REturned!")
	}

}

func TestSuite(Testor *testing.T) {
	suite.Run(Testor, new(SuiteStruct))
}
