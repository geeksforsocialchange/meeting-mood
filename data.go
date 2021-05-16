package main

import (
	"github.com/hashicorp/go-memdb"
	"log"
)

type UserMoodStruct struct {
	RoomUser string `json:"roomUser"`
	Username string `json:"username"`
	Mood     string `json:"mood"`
	Room	 string `json:"room"`
}

func Database() *memdb.MemDB {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"usermood": &memdb.TableSchema{
				Name: "usermood",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "RoomUser"},
					},
					"username": {
						Name: "username",
						Unique: false,
						Indexer: &memdb.StringFieldIndex{Field: "Username"},
					},
					"mood": {
						Name:    "mood",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Mood"},
					},
					"room": {
						Name: "room",
						Unique: false,
						Indexer: &memdb.StringFieldIndex{Field: "Room"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
	return db
}

func Save(moodStruct UserMoodStruct, db *memdb.MemDB) {
	txn := db.Txn(true)
	txn.Insert("usermood", moodStruct)
	txn.Commit()
}

func Delete(roomuser string, db *memdb.MemDB) {
	log.Println(roomuser)
	txn := db.Txn(true)
	_, err := txn.DeleteAll("usermood", "id", roomuser)
	if err != nil {
		log.Fatal(err)
	}
	txn.Commit()
}

func GetAll(room string, db *memdb.MemDB) []UserMoodStruct {
	txn := db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("usermood", "id")
	if err != nil {
		log.Fatal(err)
	}
	var usermoods []UserMoodStruct
	for obj := it.Next(); obj != nil; obj = it.Next() {
		userMood := obj.(UserMoodStruct)
		if userMood.Room == room {
			usermoods = append(usermoods, userMood)
		}
	}
	return usermoods

}
