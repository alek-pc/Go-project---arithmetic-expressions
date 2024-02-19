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

// структура настроек
type Settings struct {
	PlusTime           int // время сложения
	MinusTime          int // время вычитания
	DivisionTime       int // время деления
	MultiplicationTime int // время умножения
	WorkersNum         int // кол-во воркеров
}

type Agent struct {
	MaxNum int
	CurNum int
	Busy   bool
}
type Agents_struct struct {
	Agents []Agent
}

var Agents Agents_struct

func (ags *Agents_struct) FreeAgent() bool {
	for _, ag := range ags.Agents {
		if ag.CurNum > 0 {
			return true
		}
	}
	return false
}
func (ags *Agents_struct) AddAgent(ag Agent) {
	ags.Agents = append(ags.Agents, ag)
}
func (ags *Agents_struct) AddOperation() bool {
	for ag := range ags.Agents {
		if ags.Agents[ag].CurNum > 0 {
			ags.Agents[ag].CurNum--
			ags.Agents[ag].Busy = true
			return true
		}
	}
	return false
}
func (ags *Agents_struct) OperationMade() {
	for i := len(ags.Agents) - 1; i >= 0; i-- {
		if ags.Agents[i].CurNum < ags.Agents[i].MaxNum {
			ags.Agents[i].CurNum++
			if ags.Agents[i].CurNum == ags.Agents[i].MaxNum {
				ags.Agents[i].Busy = false
			}
			break
		}
	}
}

// структура для шаблона страницы (думаю по названиям все понятно)
type Response struct {
	Plus       int
	Minus      int
	Division   int
	Multi      int
	WorkersNum int
}

// загрузка данных в csv
func (s *Settings) Upload() {
	for {
		if len(Agents.Agents) < settings.WorkersNum {
			Agents.AddAgent(Agent{MaxNum: 5, CurNum: 5})
		} else if len(Agents.Agents) > settings.WorkersNum {
			Agents.Agents = Agents.Agents[:len(Agents.Agents)-1]
		} else {
			break
		}
	}
	f, err := os.Create("./data/settings.csv") // берем файлик
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	defer writer.Flush()
	writer.Comma = ';' // разделитель
	// структура csv: время плюса;время минуса;время деления;время умножения;кол-во воркеров
	writer.Write([]string{strconv.Itoa(s.PlusTime), strconv.Itoa(s.MinusTime), strconv.Itoa(s.DivisionTime), strconv.Itoa(s.MultiplicationTime), strconv.Itoa(s.WorkersNum)})
}

// загрузка данных из csv
func (s *Settings) Download() {
	f, err := os.Open("./data/settings.csv") // открываем файлик
	// каждую ошибку выводим в консоль
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.Comma = ';'
	for { // проходимся по всем строкам
		line, err := reader.Read() // чтение строки
		if err == io.EOF {         // строки закончились - выходим
			break
		} else if err != nil {
			return
		}
		// берем каждое значение из строки
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

// инициализация Settings
func Init() *Settings {
	set := Settings{PlusTime: 10, MinusTime: 10, DivisionTime: 10, MultiplicationTime: 10} // данные по умолчанию
	set.Download()                                                                         // загрузка из СУБД
	Agents = Agents_struct{Agents: make([]Agent, 0)}
	for {
		if len(Agents.Agents) < set.WorkersNum {
			Agents.AddAgent(Agent{MaxNum: 5, CurNum: 5})
		} else if len(Agents.Agents) > set.WorkersNum {
			Agents.Agents = Agents.Agents[:len(Agents.Agents)-1]
		} else {
			break
		}
	}
	return &set
}

var settings Settings // настройки

// отпрвака Settings сюда из main
func SendSettings(set Settings) {
	settings = set
}

// получение настроек из других пакетов
func GetSettings() *Settings {
	return &settings
}

// обработчик страницы настроек
func SettingsPage(w http.ResponseWriter, r *http.Request) {
	plus_val := r.FormValue("plus_set") // берем значение сложения
	// если форма пустая - не чекаем
	if plus_val != "" {
		// берем каждое значение из формы
		// ошибки выводим
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
		if err != nil {
			http.Error(w, "", 500)
		}
		// обновляем данные
		settings.PlusTime = plus
		settings.MinusTime = minus
		settings.DivisionTime = division
		settings.MultiplicationTime = multiplaction
		settings.WorkersNum = workersNum

		settings.Upload() // загрузка в csv
	}
	tmpl, err := template.ParseFiles("./templates/settings_page.html") // шаблон страницы
	if err != nil {
		http.Error(w, "", 500)
		fmt.Println(err)
		return
	}
	// создаем Response - для шаблона страницы
	response := Response{Plus: settings.PlusTime, Minus: settings.MinusTime,
		Division: settings.DivisionTime, Multi: settings.MultiplicationTime, WorkersNum: settings.WorkersNum}
	tmpl.Execute(w, response) // вывод страницы
}
