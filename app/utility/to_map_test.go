package utility

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"
)

const timeLayout = "2006-01-02 15:04:05"

type Profile struct {
	Experience string    `map:"experience"`
	Date       time.Time `map:"time"`
}

// its own toMap method
func (p Profile) StructToMap() (key string, value interface{}) {
	return "time", p.Date.Format(timeLayout)
}

type MySlice []int

// its own toMap method
func (a MySlice) StructToMap() (string, interface{}) {
	key := "array"
	b := strings.Builder{}
	if len(a) == 0 {
		return key, b.String()
	}
	for i := 0; i < len(a); i++ {
		b.WriteString(strconv.Itoa(a[i]) + ",")
	}
	return key, b.String()
}

// alias type
type Gender int

const (
	male   Gender = 1
	female Gender = 2
)

// Dive struct
type GithubPage struct {
	URL  string `map:"url"`
	Star int    `map:"star"`
}

type StructNoDive struct {
	NoDive int `map:"no_dive"`
}

// User is used for demonstration
type User struct {
	Name        string       `map:"name,omitempty,wildcard"` // string
	Email       *string      `map:"email_ptr,omitempty"`     // pointer
	MyGender    Gender       `map:"gender,omitempty"`        // type alias
	Github      GithubPage   `map:"github,dive,omitempty"`   // struct dive
	NoDive      StructNoDive `map:"no_dive,omitempty"`       // no dive struct
	MyProfile   Profile      `map:"my_profile,omitempty"`    // struct implements its own method
	Arr         []int        `map:"arr,omitempty"`           // normal slice
	MyArr       MySlice      `map:"my_arr,omitempty"`        // slice implements its own method
	IgnoreFiled string       `map:"-"`
}

func newUser() User {
	name := "user"
	email := "yaopei.liang@foxmail.com"
	myGender := male
	MyGithub := GithubPage{
		URL:  "https://github.com/liangyaopei",
		Star: 1,
	}
	NoDive := StructNoDive{NoDive: 1}
	dateStr := "2020-07-21 12:00:00"
	date, _ := time.Parse(timeLayout, dateStr)
	profile := Profile{
		Experience: "my experience",
		Date:       date,
	}
	arr := []int{1, 2, 3}
	myArr := MySlice{11, 12, 13}
	return User{
		Name:        name,
		Email:       &email,
		MyGender:    myGender,
		Github:      MyGithub,
		NoDive:      NoDive,
		MyProfile:   profile,
		Arr:         arr,
		MyArr:       myArr,
		IgnoreFiled: "ignore",
	}
}

func TestStructToMap(t *testing.T) {
	user := newUser()
	tag := "map"
	methodName := "StructToMap"
	res, err := StructToMap(&user, tag, methodName)
	if err != nil {
		t.Errorf("struct to map:%s", err.Error())
		return
	}
	for k, v := range res {
		t.Logf("k:%v,v:%v", k, v)
	}
}

type benchmarkUser struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

func newBenchmarkUser() benchmarkUser {
	return benchmarkUser{
		Name:    "name",
		Age:     18,
		Address: "github address",
		Contact: "github contact",
	}
}

func BenchmarkStructToMapByJson(b *testing.B) {
	user := newBenchmarkUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, _ := json.Marshal(&user)
		m := make(map[string]interface{})
		json.Unmarshal(data, &m)
	}
}

func BenchmarkStructToMapByToMap(b *testing.B) {
	user := newBenchmarkUser()
	tag := "json"
	methodName := ""
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StructToMap(&user, tag, methodName)
	}
}
