package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestRoomHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/12345", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(roomHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestNewRoom(t *testing.T) {
	req, err := http.NewRequest("GET", "/12345/all", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(allHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	wantedBody := "null\n"
	if content := rr.Body.String(); content != wantedBody {
		t.Errorf("all api returned incorrect data for new room: got %v want %v",
			content, wantedBody)
	}
}

func TestSetDeleteMood(t *testing.T) {
	var jsonStr = []byte(`{"username": "alice", "mood": "Agree", "room": "12345", "roomUser": "12345alice"}`)
	req, err := http.NewRequest("POST", "/12345/mood", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addMoodHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	person, err := getFirstRecord()
	if err != nil {
		t.Fatal(err)
	}
	if person.Mood != "Agree" || person.Room != "12345" || person.Username != "alice" || person.RoomUser != "12345alice" {
		t.Errorf("saved item didn't match: got %v", person)
	}


	Delete("12345alice", db)
	_, err = getFirstRecord()
	if err == nil {
		t.Errorf("item wasn't deleted")
	}
}

func getFirstRecord() (UserMoodStruct, error) {
	var txn = db.Txn(false)
	defer txn.Abort()
	it, err := txn.Get("usermood", "id")
	if err != nil {
		return UserMoodStruct{}, err
	}
	obj := it.Next()
	if obj != nil {
		return obj.(UserMoodStruct), nil
	} else {
		return UserMoodStruct{}, errors.New("no record returned")
	}
}