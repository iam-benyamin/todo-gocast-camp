package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"todo-gocast-camp/constant"
	"todo-gocast-camp/contract"
	"todo-gocast-camp/entity"
	"todo-gocast-camp/repository/filestore"
	"todo-gocast-camp/repository/memorystore"
	"todo-gocast-camp/service/task"
	//"golang.org/x/crypto/bcrypt"
)

var (
	userStorage     []entity.User
	categoryStorage []entity.Category
	// taskStorage []entity.Task

	authenticatedUser *entity.User
	serializationMode string
)

const (
	userStoragePath string = "user.csv"
)

func main() {
	taskMemoryRepo := memorystore.NewTaskStore()
	taskService := task.NewService(taskMemoryRepo)

	serializeMode := flag.String("serialize-mode", constant.ManDarAvardiSerializationMode, "serialization mode to write data to")
	command := flag.String("command", "no command", "command to run")
	flag.Parse()

	fmt.Println("Hello to TODO app ")

	scanner := bufio.NewScanner(os.Stdin)

	switch *serializeMode {
	case constant.ManDarAvardiSerializationMode:
		serializationMode = constant.ManDarAvardiSerializationMode
	default:
		serializationMode = constant.JsonSerializationMode
	}

	var userFileStore = filestore.New(userStoragePath, serializationMode)
	users := userFileStore.Load()
	userStorage = append(userStorage, users...)

	for {
		runCommand(userFileStore, *command, &taskService)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()
	}
}

func runCommand(store contract.UserWriteStore, command string, taskService *task.Service) {
	if command != "login" && command != "register-user" && command != "exit" && authenticatedUser == nil {
		login()
	}

	switch command {
	case "create-task":
		createTask(taskService)
	case "list-task":
		listTask(taskService)
	case "create-category":
		createCategory()
	case "register-user":
		registerUser(store)
	case "login":
		login()
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("command is not valid", command)
	}
}

func createTask(taskService *task.Service) {
	fmt.Println("authenticated user is ", authenticatedUser.Name)
	scanner := bufio.NewScanner(os.Stdin)
	var title, duedate, category string

	fmt.Println("Please Enter the title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("Please Enter the task category")
	scanner.Scan()
	category = scanner.Text()

	fmt.Println("Please Enter the task due date")
	scanner.Scan()
	duedate = scanner.Text()

	categoryID, err := strconv.Atoi(category)
	if err != nil {
		fmt.Printf("category id is not valid integer %v\n", err)
		return
	}

	respone, cErr := taskService.Create(task.CreateRequest{
		Title:               title,
		Category:            categoryID,
		Duedate:             duedate,
		AuthenticatedUserID: authenticatedUser.ID,
	})
	if cErr != nil {
		fmt.Println("error:", cErr)

		return
	}
	fmt.Println("created task:", respone.Task)
}

func createCategory() {
	scanner := bufio.NewScanner(os.Stdin)
	var color, title string

	fmt.Println("Please Enter the task title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("Please Enter the task color")
	scanner.Scan()
	color = scanner.Text()

	category := entity.Category{
		ID:     len(categoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}

	categoryStorage = append(categoryStorage, category)

	fmt.Println("category:", category, color)
}

func listTask(taskService *task.Service) {
	userTasks, lErr := taskService.List(task.ListRequest{UserID: authenticatedUser.ID})
	if lErr != nil {
		fmt.Println("error: ", lErr)

		return
	}

	fmt.Println("user tasks", userTasks)
}

func registerUser(store contract.UserWriteStore) {
	scanner := bufio.NewScanner(os.Stdin)
	var id int
	var name, email, password string

	fmt.Println("Please Enter the name")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("Please Enter the email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("Please Enter the password")
	scanner.Scan()
	password = scanner.Text()

	id = len(userStorage) + 1

	user := entity.User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: hashThePassword(password),
	}

	userStorage = append(userStorage, user)

	store.Save(user)
}

func hashThePassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func login() {
	scanner := bufio.NewScanner(os.Stdin)
	var id, email, password string

	fmt.Println("Please Enter the email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("Please Enter the password")
	scanner.Scan()
	password = scanner.Text()

	id = email

	for _, user := range userStorage {
		if user.Email == email && user.Password == hashThePassword(password) {
			authenticatedUser = &user
			fmt.Println("okay!", id)

			break
		}
	}
	if authenticatedUser == nil {
		fmt.Println("unauthorized!")

		os.Exit(0)
	}
}
