package memorystore

import (
	"todo-gocast-camp/entity"
)

type Category struct {
	categories []entity.Category
}

func (c Category) DoseThisUserHaveThisCategoryID(userID, categoryID int) bool {
	isFound := false
	for _, c := range c.categories {
		if c.ID == categoryID && c.UserID == userID {
			isFound = true

			break
		}
	}

	return isFound
}
