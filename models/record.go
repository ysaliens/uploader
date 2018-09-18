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

// Create - Save a file of records to db keeping the same connection
func (mdb MongoDBConnection) Create(records []Record) {
	mdb.session = mdb.GetSession()
	defer mdb.session.Close()

	c := mdb.session.DB(database).C(collection)
	bulk := c.Bulk()
	// bulk.Unordered() //Further optimizes speed if order is not important in records

	// Check if first record in file is in db, if it is, upsert ALL others since they exist
	if _, err := mdb.FetchOneByNameYearCode(records[0].Name, records[0].Year, records[0].BudgetCode); err == nil {
		for _, r := range records {
			update := bson.M{"$setOnInsert": bson.M{"opex": r.Opex, "category": r.Category, "budgetdesc": r.BudgetDesc,
				"jan": r.Jan, "feb": r.Feb, "mar": r.Mar, "apr": r.Apr, "may": r.May, "jun": r.Jun,
				"jul": r.Jul, "aug": r.Aug, "sep": r.Sep, "oct": r.Oct, "nov": r.Nov,
				"dec": r.Dec}}
			selector := bson.M{"name": r.Name, "year": r.Year, "budgetcode": r.BudgetCode}
			bulk.Upsert(selector, update)
		}
	} else { // If first record in a file is NOT found, file hasn't been processed, insert all
		for _, r := range records {
			bulk.Insert(r)
		}
	}
	// Run bulk operation on db
	_, err := bulk.Run()
	if err != nil {
		panic(err)
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
		Timeout:  60 * time.Second,
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
