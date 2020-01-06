package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

// Query database and compare with dependency api
func naiveFunc(id int) bool {
	db, _ := gorm.Open("mysql", "user:password@127.0.0.1:3306/test")
	defer db.Close()

	var productID int
	db.First(&productID, id)

	client := &http.Client{}
	resp, _ := client.Get("dependency.hell")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	depID, _ := strconv.Atoi(string(body))

	return productID == depID
}

type DApi struct{}

func (d DApi) getID() int { return rand.Intn(1) }

// Inject dependency as paramters
func goodFunc1(id int, db *gorm.DB, deps *DApi) bool {
	var productID int
	db.First(&productID, id)

	depID := deps.getID()

	return productID == depID
}

type IDB interface {
	getDBID(id int) int
}

type dbImpl struct{ db *gorm.DB }

func (di dbImpl) getDBID(id int) int {
	var productID int
	di.db.First(&productID, id)
	return productID
}

type IAPIs interface {
	getDepsID() int
}

type apiImpl struct{}

func (ai apiImpl) getDepsID() int {
	return rand.Intn(1)
}

// Wrap dependency by interfaces and inject as parameters
func goodFunc2(id int, db IDB, deps IAPIs) bool {
	return db.getDBID(id) == deps.getDepsID()
}

func main() {
	conn, _ := gorm.Open("mysql", "user:password@127.0.0.1:3306/test")
	idb := dbImpl{db: conn}
	iapi := apiImpl{}

	fmt.Println(goodFunc2(rand.Intn(1), idb, iapi))
}
