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
