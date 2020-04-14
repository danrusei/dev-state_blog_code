package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

//Item holds the info retrieve from the API
type Item struct {
	ID        string `json:"id"`
	Created   string `json:"Created"`
	Name      string `json:"name"`
	Expdate   string `json:"expdate"`
	Expopen   int    `json:"expopen"`
	Comment   string `json:"comment"`
	Targetage string `json:"targetage"`
	Isopen    bool   `json:"isopen"`
	Opened    string `json:"opened"`
	Isvalid   bool   `json:"isvalid"`
	Daysvalid int    `json:"daysvalid"`
}

//F is the entry point for cloud functions
func F(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	req := dialogflow.WebhookRequest{}
	if err := jsonpb.Unmarshal(r.Body, &req); err != nil {
		log.Println("Couldn't Unmarshal request to jsonpb")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	log.Println("A number of longs below to test")
	log.Println("Get parameters: ", req.GetQueryResult().GetParameters().GetFields())
	log.Println("Get Actions: ", req.GetQueryResult().GetAction())
	log.Println("Get Intent: ", req.GetQueryResult().GetIntent())

	mapfields := req.GetQueryResult().GetParameters().GetFields()

	items := []map[string]string{}

	switch req.QueryResult.Intent.DisplayName {
	case "itemstock":
		log.Println("It was requested something about itemstock")
		age := mapfields["AgeType"].GetStringValue()
		med := mapfields["meds"].GetStringValue()
		items = SelectItems(age, med)

	case "listitems":
		log.Println("It was requested something about listitems")
		age := mapfields["AgeType"].GetStringValue()
		medAll := ""
		items = SelectItems(age, medAll)
	}

	fullfilement := []string{}

	for _, item := range items {
		text := fmt.Sprintf("%s that expires in %s days.", item["name"], item["expiredays"])
		fullfilement = append(fullfilement, text)
	}

	response := dialogflow.WebhookResponse{}

	if len(fullfilement) > 0 {
		response = dialogflow.WebhookResponse{
			FulfillmentText: fmt.Sprintf("You have these meds back home: " + strings.Join(fullfilement, "; \n")),
		}
	} else {
		response = dialogflow.WebhookResponse{
			FulfillmentText: fmt.Sprintf("Sorry, you don't have these meds back home"),
		}
	}

	data, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		log.Printf("Can't marshall the data %v: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

//SelectItems select the items from the list
func SelectItems(age string, med string) []map[string]string {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	url, err := url.Parse("https://yours_backend_server/")
	if err != nil {
		log.Fatalln("could not parse the URL")
	}

	items, err := GetAllItems(&client, url)
	if err != nil {
		log.Fatalf("something happened while trying to retrive data: %v", err)
	}

	listItems := []map[string]string{}

	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Name), strings.ToLower(med)) {
			findItem := make(map[string]string)
			if item.Targetage == age {
				findItem["name"] = item.Name
				findItem["expiredays"] = strconv.Itoa(item.Daysvalid)

				listItems = append(listItems, findItem)
			}
			continue
		}
		continue
	}

	return listItems
}

//GetAllItems retrieve the items from the service
func GetAllItems(client *http.Client, url *url.URL) ([]Item, error) {

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var items = []Item{}

	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		return nil, err
	}

	return items, err
}
