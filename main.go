package main

import (
        "context"
        "encoding/json"
        "fmt"
       "log"
        "net/http"
        "time"
//      	"go.mongodb.org/mongo-driver/bson"
      	"go.mongodb.org/mongo-driver/bson/primitive"
      	"go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Participant struct {
  Name                string              `json:"name,omitempty" bson:"name,omitempty"`
  Email               string              `json:"email,omitempty" bson:"email,omitempty"`
  RSVP                string              `json:"rsvp,omitempty" bson:"rsvp,omitempty"`
}

type Meeting struct {
  ID                  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  Titile              string             `json:"title,omitempty" bson:"title,omitempty"`
  Participant         []Participant      `json:"participant" bson:"participant"`
  Start_Time          string             `json:"start_time,omitempty" bson:"start_time,omitempty"`
  End_Time            string             `json:"end_time,omitempty" bson:"end_time,omitempty"`
  Creation_Timestamp  string             `json:"creation_timestamp,omitempty" bson:"creation_timestamp,omitempty"`
}

func CreateMeetingEndpoint(response http.ResponseWriter, request *http.Request) {
  response.Header().Set("content-type", "application/json")
  var meeting Meeting
  _ = json.NewDecoder(request.Body).Decode(&meeting)
  collection := client.Database("markiv").Collection("meeting")
  ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, meeting)
	json.NewEncoder(response).Encode(result)
}

func GetMeetingID(response http.ResponseWriter, request *http.Request) {
  response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var meeting Meeting
	collection := client.Database("markiv").Collection("meeting")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Meeting{ID: id}).Decode(&meeting)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meeting)
}

func ListMeetingTimeFrame(response http.ResponseWriter, request *http.Request) {

}

func ListMeetingParticipant(response http.ResponseWriter, request *http.Request) {
  response.Header().Set("content-type", "application/json")
	var meeting []Meetings
	collection := client.Database("thepolyglotdeveloper").Collection("meeting")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var meeting Meeting
		cursor.Decode(&meeting)
		people = append(meeting, meeting)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meeting)
}


func main () {
  fmt.Printf("Hello World.\n")
  fmt.Printf("Starting the api.....")
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
//  mux := http.NewServerMux()
  http.HandleFunc("/meeting", CreateMeetingEndpoint)
  http.HandleFunc("/meeting/{ID}", GetMeetingID)
  http.HandleFunc("/meeting/{Start_Time}{End_Time}/{ID}", ListMeetingTimeFrame)
  http.HandleFunc("/meeting/{email}/{ID}", ListMeetingParticipant)

  http.ListenAndServe(":8080",nil)
}
