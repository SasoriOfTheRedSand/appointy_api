package main

import (
        "context"
        "encoding/json"
        "fmt"
        "strconv"
        "errors"
        "path"
        "sync"
//       "log"
        "net/http"
        "time"
      	"go.mongodb.org/mongo-driver/bson"
      	"go.mongodb.org/mongo-driver/bson/primitive"
      	"go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

var lock sync.Mutex

var Defaultskip = int64(0)

var Defaultlimit = int64(10)

var skip = Defaultskip
var limit = Defaultlimit

type Participant struct {
  Name                string              `json:"name,omitempty" bson:"name,omitempty"`
  Email               string              `json:"email,omitempty" bson:"email,omitempty"`
  RSVP                string              `json:"rsvp,omitempty" bson:"rsvp,omitempty"`
}

type Meeting struct {
  ID                  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  Title               string             `json:"title,omitempty" bson:"title,omitempty"`
  Participants         []Participant      `json:"participants" bson:"participants"`
  Start_Time          string             `json:"start_time,omitempty" bson:"start_time,omitempty"`
  End_Time            string             `json:"end_time,omitempty" bson:"end_time,omitempty"`
  Creation_Timestamp  string             `json:"creation_timestamp,omitempty" bson:"creation_timestamp,omitempty"`
}

func (person *Participant) cons() {
	if person.RSVP == "" {
		person.RSVP = "Not Answered"
	}
	if person.Email == "" {
		person.Email = "defaultmail@email.com"
	}
	if person.Name == "" {
		person.Name = person.Email
	}
}

func (obj *Meeting) def() {
	if obj.Title == "" {
		obj.Title = "Untitled Meeting"
	}
	if obj.Start_Time == "" {
		obj.Start_Time = string(time.Now().Format(time.RFC3339))
	}
	if obj.End_Time == "" {
		obj.End_Time = string(time.Now().Local().Add(time.Hour * time.Duration(1)).Format(time.RFC3339))
	}
	if obj.Creation_Timestamp == "" {
		obj.Creation_Timestamp = string(time.Now().Format(time.RFC3339))
	}
	for i := range obj.Participants {
		obj.Participants[i].cons()
	}
}

func MeetingHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		CreateMeetingEndpoint(response, request)
	}
	if request.Method == "GET" {
		ListMeetingTimeFrame(response, request)
	}
}



func ParticipantsBusy(thismeet Meeting) error {
	collection := client.Database("markiv").Collection("meetings")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var meet Meeting
	for _, thisperson := range thismeet.Participants {
//		if thisperson.RSVP lock sync.Mutex == "Yes" {
			filter := bson.M{
				"participants.email": thisperson.Email,
				"participants.rsvp":  "Yes",
				"end_time":            bson.M{"$gt": string(time.Now().Format(time.RFC3339))},
			}
			cursor, _ := collection.Find(ctx, filter)
			for cursor.Next(ctx) {
				cursor.Decode(&meet)
				if (thismeet.Start_Time >= meet.Start_Time && thismeet.Start_Time <= meet.End_Time) ||
					(thismeet.End_Time >= meet.Start_Time && thismeet.End_Time <= meet.End_Time) {
					returnerror := "Error 400: Participant " + thisperson.Name + " RSVP Clash"
					return errors.New(returnerror)
				}
			}
		}
    return nil
//	}
}


func CreateMeetingEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var meet Meeting
	_ = json.NewDecoder(request.Body).Decode(&meet)
	meet.def()
	if meet.Start_Time < meet.Creation_Timestamp {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{ "message": "Meeting cannot start in the past" }`))
		return
	}
	if meet.Start_Time > meet.End_Time {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{ "message": "Invalid time" }`))
		return
	}
	lock.Lock()
	defer lock.Unlock()
	err := ParticipantsBusy(meet)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	collection := client.Database("markiv").Collection("meeting")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, meet)
	meet.ID = result.InsertedID.(primitive.ObjectID)
	json.NewEncoder(response).Encode(meet)
	fmt.Println(meet)
}


func CheckMeetingwithID(id primitive.ObjectID) (Meeting, error) {
	var meet Meeting
	collection := client.Database("markiv").Collection("meeting")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, Meeting{ID: id}).Decode(&meet)
	if meet.ID != id {
		err = errors.New("Error 400: ID not present")
	}
	return meet, err
}


func GetMeetingID(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response.Header().Set("content-type", "application/json")
		fmt.Println(path.Base(request.URL.Path))
		id, _ := primitive.ObjectIDFromHex(path.Base(request.URL.Path))
		meetingwithID, err := CheckMeetingwithID(id)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(response).Encode(meetingwithID)
	}
}

func CheckMeetingWithTime(CheckStartTime string, CheckEndTime string) []Meeting {
	collection := client.Database("markiv").Collection("meeting")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "start_time", Value: 1}})
	opts.Skip = &skip
	opts.Limit = &limit
	filter := bson.D{
		{Key: "start_time", Value: bson.M{"$gt": CheckStartTime}},
		{Key: "end_time", Value: bson.M{"$lt": CheckEndTime}},
	}
	cursor, _ := collection.Find(ctx, filter, opts)
	var meetingsreturn []Meeting
	var meet Meeting
	for cursor.Next(ctx) {
		cursor.Decode(&meet)
		meetingsreturn = append(meetingsreturn, meet)
	}
	return meetingsreturn
}


func ListMeetingTimeFrame(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	fmt.Println((request.URL.Query()["start"][0]))
	fmt.Println((request.URL.Query()["end"][0]))
	CheckStartTime := request.URL.Query()["start"][0]
	CheckEndTime := request.URL.Query()["end"][0]
	if len(request.URL.Query()["limit"]) != 0 {
		limit, _ = strconv.ParseInt(request.URL.Query()["limit"][0], 0, 64)
	}
	if len(request.URL.Query()["ofset"]) != 0 {
		skip, _ = strconv.ParseInt(request.URL.Query()["offset"][0], 0, 64)
	}
	meetingswithtime := CheckMeetingWithTime(CheckStartTime, CheckEndTime)
	json.NewEncoder(response).Encode(meetingswithtime)
	skip = Defaultskip
	limit = Defaultlimit
}

func CheckParticipant(email string) []Meeting {
	collection := client.Database("markiv").Collection("meeting")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "start_time", Value: 1}})
	opts.Skip = &skip
	opts.Limit = &limit
	cursor, _ := collection.Find(ctx, bson.D{
		{Key: "participant.email", Value: email},
	}, opts)
	var meetingsreturn []Meeting
	var meet Meeting
	for cursor.Next(ctx) {
		cursor.Decode(&meet)
		meetingsreturn = append(meetingsreturn, meet)
	}
	return meetingsreturn
}


func ListMeetingParticipant(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response.Header().Set("content-type", "application/json")
		fmt.Println((request.URL.Query()["participant"][0]))
		if len(request.URL.Query()["limit"]) != 0 {
			limit, _ = strconv.ParseInt(request.URL.Query()["limit"][0], 0, 64)
		}
		if len(request.URL.Query()["ofset"]) != 0 {
			skip, _ = strconv.ParseInt(request.URL.Query()["offset"][0], 0, 64)
		}
		email := request.URL.Query()["participant"][0]
		participantmeetings := CheckParticipant(email)
		if len(participantmeetings) == 0 {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(`{ "message": "Participant not present" }`))
			return
		}
		json.NewEncoder(response).Encode(participantmeetings)
		skip = Defaultskip
		limit = Defaultlimit
	}
}

func main () {
  fmt.Printf("Hello World.\n")
  fmt.Printf("Starting the api.....")
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
//  mux := http.NewServerMux()
  http.HandleFunc("/meetings", MeetingHandler)
  http.HandleFunc("/meeting/", GetMeetingID)
  http.HandleFunc("/participants/", ListMeetingParticipant)

  http.ListenAndServe(":8080",nil)
}
