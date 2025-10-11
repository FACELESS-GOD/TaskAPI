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

	errorChannel := make(chan error)
	requestChannel := make(chan TaskStoreResponse)
	defer close(errorChannel)
	defer close(requestChannel)

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go Suite.Model.AddTask(TaskStoreRequest{
			Title:            strconv.Itoa(100 + i),
			Task_Description: strconv.Itoa(100),
			Task_Status:      true,
		}, &wg, requestChannel, errorChannel)

	}

	go func() {

		for err := range errorChannel {
			Suite.Suite.NoError(err, "Error occured even before the edit Test Case Start! .", err.Error())
		}
	}()

	go func() {

		for res := range requestChannel {
			Suite.Suite.True(res.ID >= 1, fmt.Sprintf("Error Has occured in a Concurent Test Case %v", res.ID))
		}
	}()

	wg.Wait()

}

// func (Model *ModelStruct) EditTask(Task UpdateTaskStoreRequest, WaitGroup *sync.WaitGroup) (TaskStoreResponse, error) {
func (Suite *SuiteStruct) TestEditTask() {

	errorChannel := make(chan error)
	resultChannel := make(chan TaskStoreResponse)
	wg := sync.WaitGroup{}
	resultStore := []TaskStoreResponse{}

	defer close(errorChannel)
	defer close(resultChannel)
	defer wg.Wait()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go Suite.Model.AddTask(TaskStoreRequest{
			Title:            strconv.Itoa(100 + i),
			Task_Description: strconv.Itoa(100),
			Task_Status:      true,
		}, &wg, resultChannel, errorChannel)
	}

	go func() {

		for err := range errorChannel {
			Suite.Suite.NoError(err, "Error occured even before the Add Test Case Start! ."+err.Error())
		}
	}()

	go func() {

		for res := range resultChannel {
			Suite.Suite.True(res.ID >= 1, fmt.Sprintf("Error Has occured in a Add Test Case %v", res.ID))
			resultStore = append(resultStore, res)
		}
	}()

	wg.Wait()

	if len(resultStore) >= 1 {
		for _, prevRes := range resultStore {

			wg.Add(1)
			newTask := TaskStoreRequest{
				Title:            "Hello ",
				Task_Description: "World",
				Task_Status:      true,
			}
			go Suite.Model.EditTask(UpdateTaskStoreRequest{
				ID:   prevRes.ID,
				Task: newTask,
			}, &wg, resultChannel, errorChannel)
		}
	}

	go func() {

		for err := range errorChannel {
			Suite.Suite.NoError(err, "Error occured even before the edit Test Case Start! .")
		}
	}()

	go func() {

		for res := range resultChannel {
			Suite.Suite.True(res.ID >= 1, fmt.Sprintf("Error Has occured in a Concurent Test Case %v", res.ID))
			resultStore = append(resultStore, res)
		}
	}()

	wg.Wait()
}

func (Suite *SuiteStruct) TestDeleteTask() {

	errorChannel := make(chan error)
	resultChannel := make(chan TaskStoreResponse)
	deleteResultChannel := make(chan DeleteTaskStoreResponse)

	wg := sync.WaitGroup{}
	savedTaskID := 0

	defer close(errorChannel)
	defer close(resultChannel)
	defer close(deleteResultChannel)
	defer wg.Wait()

	task := TaskStoreRequest{
		Title:            strconv.Itoa(100),
		Task_Description: strconv.Itoa(100),
		Task_Status:      true,
	}

	wg.Add(1)
	go Suite.Model.AddTask(task, &wg, resultChannel, errorChannel)

	go func() {

		for err := range errorChannel {
			errme := "dError occured even before the Add Test Case Start! ." + err.Error()
			Suite.Suite.NoError(err, errme)
		}
	}()
	go func() {

		for res := range resultChannel {
			Suite.Suite.True(res.ID >= 1, fmt.Sprintf("Error Has occured in a Add Test Case %v", res.ID))
			savedTaskID = int(res.ID)
		}
	}()

	wg.Wait()

	for {

		if savedTaskID >= 1 {
			wg.Add(1)
			go Suite.Model.DeleteTask(DeleteTaskStoreRequest{
				ID:   int64(savedTaskID),
				Task: task,
			}, &wg, deleteResultChannel, errorChannel)

			go func() {

				for err := range errorChannel {
					Suite.Suite.NoError(err, "d2Error occured even before the Add Test Case Start! ."+err.Error())
				}
			}()

			go func() {

				for res := range deleteResultChannel {
					Suite.Suite.True(res.ID >= 1, fmt.Sprintf("Error Has occured in a Add Test Case %v", res.ID))
				}
			}()

			wg.Wait()
			break
		}
	}

}

func (Suite *SuiteStruct) TestListTask() {

	resp, err := Suite.Model.ListTask(ListTaskStore{
		Limit: 10,
		Page:  1,
	})

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
