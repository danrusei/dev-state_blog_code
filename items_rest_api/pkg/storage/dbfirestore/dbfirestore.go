package dbfirestore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/danrusei/items-rest-api/pkg/model"
	"google.golang.org/api/iterator"
)

var (
	//ErrIterate informs if iteration errors
	ErrIterate = errors.New("can't iterate over the colection documents")

	//ErrExtractDataToStruct informs if unable to extract firestore data to struct
	ErrExtractDataToStruct = errors.New("can't extract the data into a struct with DataTo")

	//ErrDelete informs if unable to delete the player
	ErrDelete = errors.New("can't delete the player from database")
)

//FirestoreDB store the client
type FirestoreDB struct {
	dbClient *firestore.Client
}

//NewFirestoreDB instantiate the client
func NewFirestoreDB(client *firestore.Client) *FirestoreDB {
	return &FirestoreDB{
		dbClient: client,
	}
}

//ListGoods show all the items
func (db *FirestoreDB) ListGoods() ([]model.Item, error) {

	var singleItem model.Item
	var allItems []model.Item
	var listItems []model.Item

	ctx := context.Background()
	itemsDoc := db.dbClient.Collection("meds")
	q := itemsDoc.OrderBy("Id", firestore.Desc)
	iter := q.Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, ErrIterate
		}
		if err := doc.DataTo(&singleItem); err != nil {
			return nil, ErrExtractDataToStruct
		}
		allItems = append(allItems, singleItem)
	}

	for _, item := range allItems {
		valid, days := checkValidity(item)
		item.IsValid = valid
		item.DaysValid = days
		listItems = append(listItems, item)
	}

	return listItems, nil
}

//AddGood add items to database
func (db *FirestoreDB) AddGood(items ...model.Item) (string, error) {

	ctx := context.Background()
	for _, item := range items {

		addtime := time.Now().Format(model.LayoutRO)
		addtime1, err := time.Parse(model.LayoutRO, addtime)
		if err != nil {
			log.Printf("Can't parse the date, %v", err)
			return "", fmt.Errorf("can't parse the date: %v", err)
		}

		item.Created = model.Timestamp{Time: addtime1}

		ref := db.dbClient.Collection("meds").NewDoc()
		item.ID = ref.ID
		wr, err := ref.Create(ctx, item)

		if err != nil {
			return "", err
		}

		ops := "Item " + item.Name + "was created at " + wr.UpdateTime.String()
		log.Println(ops)

	}

	return "", nil
}

//OpenState switch the open state status
func (db *FirestoreDB) OpenState(docName string, status bool) (string, error) {

	return "", nil
}

//DelGood delete an item from database
func (db *FirestoreDB) DelGood(docName string) (string, error) {

	ctx := context.Background()

	wr, err := db.dbClient.Collection("meds").Doc(docName).Delete(ctx)
	if err != nil {
		return "", ErrDelete
	}

	ops := "Player " + docName + " has been removed at " + wr.UpdateTime.String()
	log.Println(ops)

	return ops, nil
}

func checkValidity(i model.Item) (bool, int) {
	var exp1date, exp2date, expdate time.Duration
	t := time.Now()
	i.IsValid = true
	exp1date = t.Sub(i.ExpDate.Time)
	if exp1date > 0 {
		i.IsValid = false
		exp1date = 0
	}

	if i.IsOpen {
		exp2date = t.Sub(i.Opened.Time.AddDate(0, 0, i.ExpOpen))
		if exp2date > 0 {
			i.IsValid = false
			exp2date = 0
		}
	}

	if i.IsValid {
		expdate = max(exp1date, exp2date)
	}

	return i.IsValid, int(-expdate.Hours() / 24)
}

func max(x, y time.Duration) time.Duration {
	if y < 0 {
		if x > y {
			return x
		}
		return y
	}
	return x
}
