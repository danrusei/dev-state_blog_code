package model

import (
	"strings"
	"time"
)

const LayoutRO = "02-01-2006"

/*
Good struct holds the data of an item type:
ID -- autogenarated
Name -- Goods name
Manufactured -- Goods manufactured date
ExpDate -- Goods validity
ExpOpen -- Goods validity if it is opened
*/
type Good struct {
	Name    string    `json:"name" firestore:"Name"`
	ExpDate Timestamp `json:"expdate" firestore:"ExpDate"`
	ExpOpen int       `json:"expopen" firestore:"ExpOpen"`
	Comment string    `json:"comment" firestore:"Comment"`
}

/*
Item struct holds the data of the instance of the goods:
ID -- autogenerated
Type -- the type of product
IsOpen -- True if the product is opened
Opened -- The date when it was opeened
IsValid -- Is the item still in validity or has expired
*/
type Item struct {
	ID      string `json:"id" firestore:"Id"`
	Created Timestamp
	Good
	TargetAge string    `json:"targetage" firestore:"TargetAge"`
	IsOpen    bool      `json:"isopen" firestore:"IsOpen"`
	Opened    Timestamp `json:"opened,omitempty" firestore:"Opened"`
	IsValid   bool      `json:"isvalid" firestore:"IsValid"`
	DaysValid int       `json:"daysvalid" firestore:"DaysValid"`
}

//Timestamp create custom type to be able to handle different date format
type Timestamp struct {
	Time time.Time
}

//In order to satisfy the Unmarshaler interface, we define a single method on timestamp called UnmarshalJSON.
func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	// Convert to string and remove quotes
	s := strings.Trim(string(b), "\"")

	// Parse the time using the layout
	t, err := time.Parse(LayoutRO, s)
	if err != nil {
		return err
	}
	// Assign the parsed time to our variable
	ts.Time = t
	return nil
}

func (ts Timestamp) MarshalJSON() ([]byte, error) {
	// The +2 is to take account of the quotation marks
	b := make([]byte, 0, len(LayoutRO)+2)

	// Write the JSON output
	b = append(b, '"')
	b = ts.Time.AppendFormat(b, LayoutRO)
	b = append(b, '"')

	return b, nil
}
