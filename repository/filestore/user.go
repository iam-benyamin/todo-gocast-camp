package filestore

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todo-gocast-camp/constant"
	"todo-gocast-camp/entity"
)

type FileStore struct {
	filePath          string
	serializationMode string
}

// New constructor
func New(path, serializationMode string) FileStore {
	return FileStore{filePath: path, serializationMode: serializationMode}
}

func (f FileStore) Save(u entity.User) {
	f.writeUserToFile(u)
}

func (f FileStore) Load() []entity.User {
	var uStore []entity.User

	file, err := os.Open(f.filePath)
	if err != nil {
		fmt.Println("can't open the file. ", err)
	}

	var data = make([]byte, 1024)

	_, oErr := file.Read(data)
	if oErr != nil {
		fmt.Println("can't read from the file", oErr)
	}

	dataString := string(data)
	userSlice := strings.Split(dataString, "\n")

	for _, u := range userSlice {
		var userStruct entity.User

		switch f.serializationMode {
		case constant.ManDarAvardiSerializationMode:
			userStruct, err = deserializeFromManDaravardiy(u)
			if err != nil {
				fmt.Println("can't deserialize user record to user struct")

				return nil
			}
		case constant.JsonSerializationMode:
			if u[0] != '{' && u[len(u)-1] != '}' {
				continue
			}
			uErr := json.Unmarshal([]byte(u), &userStruct)
			if uErr != nil {
				fmt.Println("can't deserialize user record to user struct with json mode", uErr)

				return nil
			}
		default:
			fmt.Println("invalid serialization mode")

			return nil
		}

		uStore = append(uStore, userStruct)
	}

	return uStore
}

func (f FileStore) writeUserToFile(u entity.User) {
	var file *os.File

	file, err := os.OpenFile(f.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
	if err != nil {
		fmt.Println("can't create or open file.", err)

		return
	}

	var data []byte
	if f.serializationMode == constant.ManDarAvardiSerializationMode {
		fmt.Println("can't marshal ManDarAvardiSerializationMode", constant.ManDarAvardiSerializationMode)

		data = []byte(fmt.Sprintf(
			"%d, %s, %s, %s\n",
			u.ID, u.Name, u.Email, u.Password,
		))
	} else if f.serializationMode == constant.JsonSerializationMode {
		data, err = json.Marshal(u)
		if err != nil {
			fmt.Println("can't marshal user struct to json", err)

			return
		}
	} else {
		fmt.Println("f.serializationMode : ", f.serializationMode)
		fmt.Println("invalid serialization mode")

		return
	}

	data = append(data, []byte("\n")...)
	file.Write(data)
}

func deserializeFromManDaravardiy(userString string) (entity.User, error) {
	userFields := strings.Split(userString, ", ")
	var user entity.User
	user.ID, _ = strconv.Atoi(userFields[0])
	user.Name = userFields[1]
	user.Email = userFields[2]
	user.Password = userFields[3]
	fmt.Println(user)
	return user, nil
}

//func (f FileStore) loadUserFromStorage() entity.User {
//	return f.Load()
//}
