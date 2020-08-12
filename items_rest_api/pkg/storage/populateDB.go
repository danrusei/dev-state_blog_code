package storage

import (
	"encoding/json"
	"fmt"

	"github.com/danrusei/items-rest-api/pkg/model"
)

// PopulateItems insert the items in empty database
func PopulateItems(db Storage) {
	defaultItems := []byte(`[{
			"name":"Milk-(False)",
			"expdate":"23-08-2019",
			"expopen":10,
			"comment":"utilize once per day",
			"targetage":"Adult",
			"isopen":true,
			"opened":"30-07-2019"},{
			"name":"Milk2-(False)",
			"manufactured":"23-07-2019",
			"expdate":"23-12-2019",
			"expopen":10,
			"comment":"2 pils per day after meals",
			"targetage":"Child",
			"isopen":true,
			"opened":"30-07-2019"},{
			"name":"CannedFish-(True)",
			"manufactured":"15-11-2018",
			"expdate":"10-10-2020",
			"expopen":30,
			"comment":"at need, no more than 20 ml per day",
			"targetage":"Child",
			"isopen":true,
			"opened":"20-08-2019"},{
			"name":"Butter-(False)",
			"manufactured":"15-07-2019",
			"expdate":"23-08-2019",
			"expopen":20,
			"comment":"no prescription",
			"targetage":"Adult",
			"isopen":false},{
			"name":"CannedBeans-(True)",
			"manufactured":"24-02-2019",
			"expdate":"10-08-2020",
			"expopen":5,
			"comment":"30 minutes after first pill",
			"targetage":"Child",
			"isopen":false}]`)

	data := make([]model.Item, 4)
	if err := json.Unmarshal(defaultItems, &data); err != nil {
		fmt.Println("Could not unmarshal data:", err)
		return
	}
	db.AddGood(data...)
}
