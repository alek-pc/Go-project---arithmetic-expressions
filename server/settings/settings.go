package settings

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Settings struct {
	PlusTime           int
	MinusTime          int
	DivisionTime       int
	MultiplicationTime int
	WorkersNum int
}
type Response struct {
	Plus     int
	Minus    int
	Division int
	Multi    int
	WorkersNum int
}

func (s *Settings) Upload() {
	f, err := os.Create("./data/settings.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	defer writer.Flush()
	writer.Comma = ';'
	writer.Write([]string{strconv.Itoa(s.PlusTime), strconv.Itoa(s.MinusTime), strconv.Itoa(s.DivisionTime), strconv.Itoa(s.MultiplicationTime), strconv.Itoa(s.WorkersNum)})
}
func (s *Settings) Download() {
	f, err := os.Open("./data/settings.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.Comma = ';'
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return
		}
		s.PlusTime, err = strconv.Atoi(line[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		s.MinusTime, err = strconv.Atoi(line[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		s.DivisionTime, err = strconv.Atoi(line[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		s.MultiplicationTime, err = strconv.Atoi(line[3])
		if err != nil {
			fmt.Println(err)
			return
		}
		s.WorkersNum, err = strconv.Atoi(line[4])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
func Init() *Settings {
	set := Settings{PlusTime: 10, MinusTime: 10, DivisionTime: 10, MultiplicationTime: 10}
	set.Download()
	return &set
}

var settings Settings

func SendSettings(set Settings) {
	settings = set
}
func GetSettings()*Settings{
	return &settings
}
func SettingsPage(w http.ResponseWriter, r *http.Request) {
	plus_val := r.FormValue("plus_set")
	if plus_val != "" {
		plus, err := strconv.Atoi(plus_val)
		if err != nil {
			http.Error(w, "", 500)
			fmt.Println(err)
			return
		}
		minus, err := strconv.Atoi(r.FormValue("minus_set"))
		if err != nil {
			http.Error(w, "", 500)
			return
		}
		division, err := strconv.Atoi(r.FormValue("division_set"))
		if err != nil {
			http.Error(w, "", 500)
			return
		}
		multiplaction, err := strconv.Atoi(r.FormValue("multiplaction_set"))
		if err != nil {
			http.Error(w, "", 500)
			return
		}
		workersNum, err := strconv.Atoi(r.FormValue("numOfWorkers"))
		if err != nil{
			http.Error(w, "", 500)
		}
		settings.PlusTime = plus
		settings.MinusTime = minus
		settings.DivisionTime = division
		settings.MultiplicationTime = multiplaction
		settings.WorkersNum = workersNum

		settings.Upload()
	}
	tmpl, err := template.ParseFiles("./templates/settings_page.html")
	if err != nil {
		http.Error(w, "", 500)
		fmt.Println(err)
		return
	}
	response := Response{Plus: settings.PlusTime, Minus: settings.MinusTime,
		Division: settings.DivisionTime, Multi: settings.MultiplicationTime, WorkersNum: settings.WorkersNum}
	tmpl.Execute(w, response)
}
