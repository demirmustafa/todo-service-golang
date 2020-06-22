package model

type Todo struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type UpdateTodoRequest struct {
	Title string `json:"title"`
	Completed bool `json:"completed"`
}
