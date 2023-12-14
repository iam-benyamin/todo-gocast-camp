package task

import (
	"fmt"
	"todo-gocast-camp/entity"
)

type ServiceRepository interface {
	//DoseThisUserHaveThisCategoryID(userID, categoryID int) bool
	CreateNewTask(t entity.Task) (entity.Task, error)
	ListUserTasks(userID int) ([]entity.Task, error)
}

type Service struct {
	repository ServiceRepository
}

func NewService(repo ServiceRepository) Service {
	return Service{repository: repo}
}

type CreateRequest struct {
	Title               string
	Category            int
	Duedate             string
	AuthenticatedUserID int
}

type CreateResponse struct {
	Task entity.Task
}

type ListRequest struct {
	UserID int
}

type ListResponse struct {
	Tasks []entity.Task
}

func (t Service) Create(req CreateRequest) (CreateResponse, error) {
	// if !t.repository.DoseThisUserHaveThisCategoryID(req.AuthenticatedUserID, req.Category) {
	//     return CreateResponse{}, fmt.Errorf("user dose not have this category: %d", req.Category)
	// }

	createdTask, cErr := t.repository.CreateNewTask(entity.Task{
		Title:    req.Title,
		Duedate:  req.Duedate,
		Category: req.Category,
		IsDone:   false,
		UserID:   req.AuthenticatedUserID,
	})

	if cErr != nil {
		return CreateResponse{}, fmt.Errorf("can't create new task %v", cErr)
	}
	return CreateResponse{Task: createdTask}, nil
}

func (t Service) List(req ListRequest) (ListResponse, error) {
	tasks, lErr := t.repository.ListUserTasks(req.UserID)
	if lErr != nil {
		return ListResponse{}, fmt.Errorf("can't list tasks: %v", lErr)
	}

	return ListResponse{Tasks: tasks}, nil
}
