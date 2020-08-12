package storage

import "github.com/danrusei/items-rest-api/pkg/model"

//Storage the behavior for storing and retrieving items
type Storage interface {
	ListGoods() ([]model.Item, error)
	AddGood(...model.Item) (string, error)
	OpenState(string, bool) (string, error)
	DelGood(string) (string, error)
}
