package main

import "testing"

func TestCreateMeetingEndpoint (t *testing.T) {

	var jsonStr = []byte(`{"_id":1,"title":"FIRST MEETINGS","participant":"as,as@gmail.com,yes","start_time":"11:00PM","end_time":"12:00AM","creation_timestamp":"09:00PM"}`)

	req, err := http.NewRequest("POST", "/meeting", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateEntry)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"_id":1,"title":"FIRST MEETINGS","participant":"as,as@gmail.com,yes","start_time":"11:00PM","end_time":"12:00AM","creation_timestamp":"09:00PM"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}


func TestGetMeetingID(t *testing.T) {
	req, err := http.NewRequest("GET", "/meeting", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("_id", "1")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetMeetingID)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"_id":1,"title":"FIRST MEETINGS","participant":"as,as@gmail.com,yes","start_time":"11:00PM","end_time":"12:00AM","creation_timestamp":"09:00PM"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}


func TestGetMeetingIDNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/meeting", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("_id", "123")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetMeetingID)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}



func TestListMeetingParticipant(t *testing.T) {
	req, err := http.NewRequest("GET", "/meeting/{email}/id", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("email", "hello@gmail.com")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListMeetingParticipant)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"_id":1,"_id":2,"_id":15}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestListMeetingParticipantNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/meeting", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("email", "hekkkksksk@gmail.com")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListMeetingParticipant)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
