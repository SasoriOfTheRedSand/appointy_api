package main

import (
        "context"
//        "encoding/json"
        "fmt"
//        "log"
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
  Participant
  Start_Time          int                `json:"start_time,omitempty" bson:"start_time,omitempty"`
  End_Time            int                `json:"end_time,omitempty" bson:"end_time,omitempty"`
  Creation_Timestamp  int                `json:"creation_timestamp,omitempty" bson:"creation_timestamp,omitempty"`
}

func CreateMeetingEndpoint(response http.ResponseWriter, request *http.Request) {}
func GetMeetingID(response http.ResponseWriter, request *http.Request) {}
func GetMeetingTimeFrame(response http.ResponseWriter, request *http.Request) {}
func GetMeetingParticipant(response http.ResponseWriter, request *http.Request) {}


func main () {
  fmt.Printf("Hello World.\n")
  fmt.Printf("Starting the api.....")
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
//  mux := http.NewServerMux()
  http.HandleFunc("/meeting", CreateMeetingEndpoint)
  http.HandleFunc("/meeting/{ID}", GetMeetingID)
  http.HandleFunc("/meeting/{Start_Time}{End_Time}", GetMeetingTimeFrame)
  http.HandleFunc("/meeting/{Participant}", GetMeetingParticipant)

  http.ListenAndServe(":8080",nil)
}
