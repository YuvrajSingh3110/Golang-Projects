package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
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

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)
	opt := Options{}
	if options != nil {
		opt = *options
	}

	if opt.Logger == nil {
		opt.Logger = lumber.NewConsoleLogger((lumber.Info))
	}
	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opt.Logger,
	}
	if _, err := os.Stat(dir); err == nil {
		opt.Logger.Debug("Using %s (database already exists)\n", dir)
		return &driver, nil
	}
	opt.Logger.Debug("Creating database at %s...\n", dir)
	return &driver, os.MkdirAll(dir, 0755) //0755 is access permission
}

func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing collection - no place to store records")
	}
	if resource == "" {
		return fmt.Errorf("Missing resource - no name")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()
	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"

	if err := os.Mkdir(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	b = append(b, byte('\n'))
	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, fnlPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing collection - unable to read")
	}
	if resource == "" {
		return fmt.Errorf("Missing resource - no name")
	}
	record := filepath.Join(d.dir, collection, resource)
	if _, err := stat(record); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(record + ".json")
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)
}

func (d *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("Missing collection - unable to read")
	}

	dir := filepath.Join(d.dir, collection)
	if _, err := stat(dir); err != nil {
		return nil, err
	}

	files, _ := ioutil.ReadDir(dir)
	var record []string
	for _, file := range files{
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil{
			return nil, err
		}
		record = append(record, string(b))
	}
	return record, nil
}

func (d *Driver) Delete(collection, resource string) error {
	path := filepath.Join(collection, resource)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)
	switch fi, err := stat(dir);{
	case fi==nil, err!=nil:
		return fmt.Errorf("Unable to find file or directory %v", path)
	case fi.Mode().IsDir():
		os.RemoveAll(dir)
	case fi.Mode().IsDir():
		os.RemoveAll(dir + ".json")
	}
	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}
	return m
}


//to check if the collection or the diectory exists
func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
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

	records, err := db.ReadAll("users")
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
