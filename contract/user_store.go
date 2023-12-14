package contract

import "todo-gocast-camp/entity"

type UserWriteStore interface {
	Save(u entity.User)
}
type UserReadStore interface {
	Load(serializationMode string) []entity.User
}
