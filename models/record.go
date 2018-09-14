package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// @TODO Move database variables to proper config file
const (
	hosts      = "localhost"
	database   = "uploader"
	username   = ""
	password   = ""
	collection = "seaspan"
)

// Record - Database record
// @TODO Discuss changing numbers to be stored as ints in db
type Record struct {
	Name       string `json:"name"`
	Year       string `json:"year"`
	Opex       string `json:"opex"`
	Category   string `json:"category"`
	BudgetCode string `json:"budgetCode"`
	BudgetDesc string `json:"budgetDesc"`
	Jan        string `json:"jan"`
	Feb        string `json:"feb"`
	Mar        string `json:"mar"`
	Apr        string `json:"apr"`
	May        string `json:"may"`
	Jun        string `json:"jun"`
	Jul        string `json:"jul"`
	Aug        string `json:"aug"`
	Sep        string `json:"sep"`
	Oct        string `json:"oct"`
	Nov        string `json:"nov"`
	Dec        string `json:"dec"`
	TTL        string `json:"ttl"`
}

// MongoDBConnection Encapsulates a connection to a database.
type MongoDBConnection struct {
	session *mgo.Session
}

// Create - Save a sheet of records to db keeping the same connection
func (mdb MongoDBConnection) Create(records []Record) {
	mdb.session = mdb.GetSession()
	defer mdb.session.Close()

	for _, r := range records {
		if _, err := mdb.FetchOneByNameYearCode(r.Name, r.Year, r.BudgetCode); err == nil {
			continue // If record exists, continue to next one
		}
		c := mdb.session.DB(database).C(collection)
		c.Insert(r)
	}
}

// FetchByName - Get all records for a given vessel name
// @TODO Combine/Rework fetch methods, this is just for example
func (mdb MongoDBConnection) FetchByName(name string) (results []Record, err error) {
	mdb.session = mdb.GetSession()
	defer mdb.session.Close()

	c := mdb.session.DB(database).C(collection)
	err = c.Find(bson.M{"name": name}).All(&results)
	return results, err
}

// FetchByNameYear - Get all records for a given vessel name and year
func (mdb MongoDBConnection) FetchByNameYear(name string, year string) (results []Record, err error) {
	mdb.session = mdb.GetSession()
	defer mdb.session.Close()

	c := mdb.session.DB(database).C(collection)
	err = c.Find(bson.M{"year": year}).All(&results)
	return results, err
}

// FetchByNameYearCode - Get a record for a given vessel, year, and budget code
func (mdb MongoDBConnection) FetchByNameYearCode(name string, year string, budgetCode string) (results []Record, err error) {
	mdb.session = mdb.GetSession()
	defer mdb.session.Close()

	c := mdb.session.DB(database).C(collection)
	err = c.Find(bson.M{"name": name, "year": year, "budgetcode": budgetCode}).All(&results)
	return results, err
}

// FetchOneByNameYearCode - Get a single record for a given vessel, year, and budget code (uniqueness constraint)
func (mdb MongoDBConnection) FetchOneByNameYearCode(name string, year string, budgetCode string) (result Record, err error) {
	mdb.session = mdb.GetSession()
	defer mdb.session.Close()

	c := mdb.session.DB(database).C(collection)
	err = c.Find(bson.M{"name": name, "year": year, "budgetcode": budgetCode}).One(&result)
	return result, err
}

// GetSession return a new session if there is no previous one.
// Remove hardcoded localhost if database is ever not local
func (mdb *MongoDBConnection) GetSession() *mgo.Session {
	if mdb.session != nil {
		return mdb.session.Copy()
	}

	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  10 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}
