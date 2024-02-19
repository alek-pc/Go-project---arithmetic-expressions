package server

import (
	orchestrator "arifm_operations/server/orchestrator"
	"encoding/csv"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
)

var storage Storage  // Storage (ну, а что еще)

// структура выражения
type Expression struct {
	Id                   int	   // айдишка
	Start                string    // начало вычислений
	End                  string    // конец вычислений
	Expression           string    // выражение в виде строки
	Separated_expression []string  // выражение разделенное по частям (2+2*2 -> {"2", "+", "2", "*", "2"})
	Status               bool  	   // статус выражения: false - считается, true - посчитано
	Result               int	   // результат вычислений
}
// структура для шаблона страницы
type Response struct {
	Error       string  // сообщение
	Expressions []Expression  // выражения
}
// хранилище выражений
type Storage struct {
	Expressions []Expression
}

// загрузка данных из csv в storage
func (s *Storage) Download() {
	f, err := os.Open("./data/expressions_logs.csv")  // файлик
	if err != nil {
		return
	}
	reader := csv.NewReader(f)
	reader.Comma = ';'  // разделитель
	for {  // проходим по строчкам
		line, err := reader.Read()  // строчка
		if err == io.EOF {  // файл кончился - выходим
			break
		} else if err != nil {
			return
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return
		}
		res, err := strconv.Atoi(line[2])
		if err != nil {
			return
		}
		// записываем новое выражение в storage
		// структура csv: id, выражение, результат, статус, начало вычислений, конец вычислений
		expr := Expression{Id: id, Expression: line[1], Result: res, Status: stringToBool(line[3]), Start: line[4], End: line[5]}
		s.Append(expr)
		if !expr.Status{  // выражение не посчитано
			expr.Separate_expression()  // делим выражение на части
			expre := orchestrator.Expression{Id: expr.Id, Expression: expr.Expression, Separated_expression: expr.Separated_expression}
			go expre.Orchestrator()  // отправляем на вычисления
		}
	}
}
// загрузка данных в csv
func (s *Storage) Upload() {
	f, err := os.Create("./data/expressions_logs.csv")
	if err != nil {
		return
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	defer writer.Flush()
	writer.Comma = ';'
	// построчно записываем данные выражений
	for _, expr := range s.Expressions {
		writer.Write([]string{strconv.Itoa(expr.Id), expr.Expression, strconv.Itoa(expr.Result), boolToString(expr.Status), expr.Start, expr.End})
	}
}
// строка в булевое (для Upload)
func stringToBool(str string) bool{
	return str == "1"
}
// булевое в строку (для Download)
func boolToString(b bool) string{
	if b{
		return "1"
	}
	return "0"
}
// инициализация storage
func Init() *Storage {
	s := Storage{Expressions: make([]Expression, 0)}
	s.Download()
	return &s
}
// добаление выражения в storage
func (s *Storage) Append(exp Expression) {
	s.Expressions = append(s.Expressions, exp)
}
// получение storage из main (там инициализируется)
func GetStorage(stor *Storage) {
	storage = *stor
}

// обработчик главной страницы + check request
func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	resp := RemoveSpaces(r.FormValue("expression"))  // получение выражения
	for _, expr := range orchestrator.Results {  // проверка резултатов вычислений из orchestrator
		if expr.Status && !storage.Expressions[expr.Id].Status {  // статус - тру
			// обновление значение в storage
			storage.Expressions[expr.Id].Result = expr.Result
			storage.Expressions[expr.Id].Status = expr.Status
			storage.Expressions[expr.Id].End = expr.End
			storage.Expressions[expr.Id].Start = expr.Start
			storage.Upload()  // загрузка обновленных данных в csv
		}
	}
	tmpl, err := template.ParseFiles("./templates/main_paper.html")  // шаблонизатор страницы
	if err != nil {
		http.Error(w, "", 500)
	}
	for _, expres := range storage.Expressions {  // проверка на наличие такого выражения
		if expres.Expression == resp { // выражение уже считается
			response := Response{Error: "expression is already computing", Expressions: reverse(storage.Expressions)}
			tmpl.Execute(w, response)  // вывод сообщения и выражений
			return
		}
	}
	expression := Expression{Id: len(storage.Expressions), Expression: resp, Status: false}  // окей, уникальное выражение, создаем объект
	if err := expression.Separate_expression(); err != "" { // ошибка при обработке выражения
		response := Response{Error: err, Expressions: reverse(storage.Expressions)}
		tmpl.Execute(w, response)  // вывод ошибки и списка выражений
		return
	}

	storage.Append(expression)  // добавляем выражение в storage
	storage.Upload()  // загружаем обновление в csv
	response := Response{Error: "all is good. Expression is computing", Expressions: reverse(storage.Expressions)}
	tmpl.Execute(w, response)  // сообщение + список выражений

	// отправка выражения оркестратору
	// создание объекта типа orchestrator.Expression - копии Expression здесь
	expre := orchestrator.Expression{Id: expression.Id, Expression: expression.Expression, Separated_expression: expression.Separated_expression}
	go expre.Orchestrator()  // запуск оркестратора

}

// переворот массива
func reverse(exps []Expression) []Expression {
	res := make([]Expression, len(exps))
	for i := len(exps) - 1; i >= 0; i-- {
		res[len(exps)-1-i] = exps[i]
	}
	return res
}

// разделение выражения на состовляющие
func (ex *Expression) Separate_expression() string {
	expression := ex.Expression
	if expression == "" {
		return "it is empty"
	}
	oper_counter := 0
	separatedExpression := make([]string, 0) // разделенное выражение: 2 + 2 + 4 -> {2, "+", 2, "+", 4}
	for i, sign := range expression {
		if in(string(sign), "1234567890") { // символ - цифра
			if i == 0 { // первое значение - новый элемент массива
				if string(sign) == "0" { // начинать с нуля - плохо
					return "I don't love '0' in the beginning of expression"
				}
				separatedExpression = append(separatedExpression, string(sign))
			} else if in(string(expression[i-1]), "1234567890") { // предыдущий символ - цифра | "-"
				separatedExpression[len(separatedExpression)-1] += string(sign) // прибавляем к последнему элементу массива
			} else if string(expression[i-1]) == "-" {
				separatedExpression[len(separatedExpression)-1] += string(sign)
				separatedExpression = append(separatedExpression, separatedExpression[len(separatedExpression)-1])
				separatedExpression[len(separatedExpression)-2] = "+"
			} else {
				separatedExpression = append(separatedExpression, string(sign))
			}
		} else if in(string(sign), "+-/*") {
			if i == len(expression)-1 {
				return "last sign is " + string(sign)
			}
			if string(sign) == "-" && i == 0 {
				oper_counter++
				separatedExpression = append(separatedExpression, "-")
			} else if in(string(sign), "/*") && i == 0 { // начать с / or *
				return "Bad expression (/ or * in the beginning)"
			} else if i != 0 && in(string(expression[i-1]), "+-/*") {
				return "Bad expression: two operands together"
			} else {
				oper_counter++
				separatedExpression = append(separatedExpression, string(sign))
			}

		} else { // неизвестный знак
			return "Unknown sign " + string(sign) + "in expression"
		}
	}
	if oper_counter == 0 {
		return "there isn't any operation"
	}
	ex.Separated_expression = separatedExpression
	return ""
}

// проверка наличия подстроки s1 в s2
func in(s1 string, s2 string) bool {
	for _, i := range s2 {
		if s1 == string(i) {
			return true
		}
	}
	return false
}
// стирание пробелов в строке
func RemoveSpaces(ex string) string {
	res := ""
	for _, i := range ex {
		if string(i) != " " {
			res += string(i)
		}
	}
	return res
}
