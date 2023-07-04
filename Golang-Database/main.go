package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

const Version = "1.0.0"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type options struct {
	Logger
}

func New() {

}

func Write() error {

}

func Read() error {

}

func ReadAll() {

}

func Delete() error {

}

func getOrCreateMutex() *sync.Mutex {

}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func main() {
	dir := "./" //where collection gets stored
	db, err := New(dir, nil)
	if err != nil {
		panic(err)
	}

	employees := []User{
		{"Aditya", "24", "6586755763", "Google", Address{"Hyderabad", "Telangana", "India", "234566"}},
		{"Bhavya", "22", "1537153678", "Microsoft", Address{"Bangalore", "Karnataka", "India", "743659"}},
		{"Chirag", "24", "2347658916", "Apple", Address{"Hyderabad", "Telangana", "India", "542301"}},
		{"Daman", "23", "2548201662", "Samsung", Address{"Noida", "Delhi NCR", "India", "102459"}},
		{"Garv", "25", "2364492016", "Amazon", Address{"Pune", "Maharashtra", "India", "102342"}},
	}

	for _, val := range employees {
		db.Write("users", val.Name, User{
			Name:    val.Name,
			Age:     val.Age,
			Contact: val.Contact,
			Company: val.Company,
			Address: val.Address,
		})
	}

	records, err := db, ReadAll("users")
	if err != nil {
		panic(err)
	}
	fmt.Println(records)

	allUsers := []User{}

	for _, f := range records {
		employeeFound := User{}
		if err := json.Unmarshal([]byte(f), &employeeFound); err != nil {
			panic(err)
		}
		allUsers = append(allUsers, employeeFound)
	}
	fmt.Println(allUsers)

}
