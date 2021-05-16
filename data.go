package main

import (
	"github.com/hashicorp/go-memdb"
	"log"
)

type UserMoodStruct struct {
	Username string `json:"username"`
	Mood string `json:"mood"`
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
						Indexer: &memdb.StringFieldIndex{Field: "Username"},
										},
					"mood": {
						Name:    "mood",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Mood"},
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

func Delete(username string, db *memdb.MemDB) {
	log.Println(username)
	txn := db.Txn(true)
	_, err := txn.DeleteAll("usermood", "id", username)
	if err != nil {
		log.Fatal(err)
	}
	txn.Commit()
}

func GetAll(db *memdb.MemDB) []UserMoodStruct {
	txn := db.Txn(false)
	defer txn.Abort()

	it, _ := txn.Get("usermood", "id")
	var usermoods []UserMoodStruct
	for obj := it.Next(); obj != nil; obj = it.Next() {
		userMood := obj.(UserMoodStruct)
		usermoods = append(usermoods, userMood)
	}
	return usermoods

}