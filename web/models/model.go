package models

import "zl2501-final-project/web/models/storage"

var drivers = make(map[string]Models)

func RegisterDriver(name string, models Models) {
	drivers[name] = models
}

func NewDriver(name string) *Models {
	models := drivers[name]
	return &models
}

type Models struct {
	UserStore *storage.UserStoreInterface
	PostStore *storage.PostStoreInterface
}
