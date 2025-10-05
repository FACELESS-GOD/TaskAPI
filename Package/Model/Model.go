package Model

type TaskStoreRequest struct {
	Title            string
	Task_Description string
	Task_Status      bool
}

type TaskStoreResponse struct {
	ID   int64
	Task TaskStoreRequest
}

type UpdateTaskStoreRequest struct {
	ID   int64
	Task TaskStoreRequest
}

type DeleteTaskStoreRequest struct {
	ID   int64
	Task TaskStoreRequest
}

type DeleteTaskStoreResponse struct {
	Status bool
	ID     int64
	Task   TaskStoreRequest
}

type ListTaskStore struct {
	Limit  int64
	Page   int64
	Offset int64
}
