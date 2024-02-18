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

var storage Storage

type Expression struct {
	Id                   int
	Start                string
	End                  string
	Expression           string
	Separated_expression []string
	Status               bool
	Result               int
}
type Response struct {
	Error       string
	Expressions []Expression
}
type Storage struct {
	Expressions []Expression
}

func (s *Storage) Download() {
	f, err := os.Open("./data/expressions_logs.csv")
	if err != nil {
		return
	}
	reader := csv.NewReader(f)
	reader.Comma = ';'
	for {
		line, err := reader.Read()
		if err == io.EOF {
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
		s.Expressions = append(s.Expressions, Expression{Id: id, Expression: line[1], Result: res, Status: stringToBool(line[3]), Start: line[4], End: line[5]})
	}

}
func (s *Storage) Upload() {
	f, err := os.Create("./data/expressions_logs.csv")
	if err != nil {
		return
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	defer writer.Flush()
	writer.Comma = ';'
	for _, expr := range s.Expressions {
		writer.Write([]string{strconv.Itoa(expr.Id), expr.Expression, strconv.Itoa(expr.Result), boolToString(expr.Status), expr.Start, expr.End})
	}
}
func stringToBool(str string) bool{
	return str == "1"
}
func boolToString(b bool) string{
	if b{
		return "1"
	}
	return "0"
}
func Init() *Storage {
	s := Storage{Expressions: make([]Expression, 0)}
	s.Download()
	return &s
}
func (s *Storage) Append(exp Expression) {
	s.Expressions = append(s.Expressions, exp)
}

func GetStorage(stor Storage) {
	storage = stor
}

func GettingResponse(w http.ResponseWriter, r *http.Request) {
	resp := RemoveSpaces(r.FormValue("expression"))
	for _, expr := range orchestrator.Results {
		if expr.Status && !storage.Expressions[expr.Id].Status {
			storage.Expressions[expr.Id].Result = expr.Result
			storage.Expressions[expr.Id].Status = expr.Status
			storage.Expressions[expr.Id].End = expr.End
			storage.Expressions[expr.Id].Start = expr.Start
			storage.Upload()

		}
	}
	tmpl, err := template.ParseFiles("./templates/main_paper.html")
	if err != nil {
		http.Error(w, "", 500)
	}
	for _, expres := range storage.Expressions {
		if expres.Expression == resp { // выражение уже считается
			// w.WriteHeader(http.StatusOK)
			// w.Header().Set("Content-Type", "application/json")
			// res := make(map[string]string)
			// res["info"] = "expression is already computing"
			// jsonResp, _ := json.Marshal(res)
			// w.Write(jsonResp)
			response := Response{Error: "expression is already computing", Expressions: reverse(storage.Expressions)}
			tmpl.Execute(w, response)
			return
		}
	}
	expression := Expression{Id: len(storage.Expressions) - 1, Expression: resp}
	if err := expression.Separate_expression(); err != "" { // ошибка при обработке выражения
		// http.Error(w, "", 400)
		// w.WriteHeader(http.StatusBadRequest)
		// w.Header().Set("Content-Type", "application/json")
		// res := make(map[string]string)
		// res["info"] = err
		// jsonResp, _ := json.Marshal(res)
		// w.Write(jsonResp)
		response := Response{Error: err, Expressions: reverse(storage.Expressions)}
		tmpl.Execute(w, response)
		return
	}
	expression.Id = len(storage.Expressions)
	expression.Status = false

	storage.Append(expression)
	storage.Upload()
	response := Response{Error: "all is good. Expression is computing", Expressions: reverse(storage.Expressions)}
	tmpl.Execute(w, response)

	expre := orchestrator.Expression{Id: expression.Id, Expression: expression.Expression, Separated_expression: expression.Separated_expression}
	go expre.Orchestrator()

}

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
func RemoveSpaces(ex string) string {
	res := ""
	for _, i := range ex {
		if string(i) != " " {
			res += string(i)
		}
	}
	return res
}
