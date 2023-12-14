package entity

type Task struct {
	ID       int
	Title    string
	Category int
	Duedate  string
	IsDone   bool
	UserID   int
}
