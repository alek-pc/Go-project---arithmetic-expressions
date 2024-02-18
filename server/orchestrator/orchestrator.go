package orkestrator

import (
	set "arifm_operations/server/settings"
	// "fmt"
	"strconv"
	"sync"
	"time"
)

// структура выражения - копия server.Expression
type Expression struct {
	Id                   int
	Start                string
	End                  string
	Expression           string
	Separated_expression []string
	Status               bool
	Result               int
}

// структура операций
type Operation struct {
	num1   int	// первое число
	oper   string // операция
	num2   int	// второе число
	res    int	// результат вычислений
	status bool // false - computing true - computed
}
var Results []Expression  // результаты всех вычислений (для получения из server)
var Agents int

// великий и могучий ОРКЕСТРАТОР
func (expr *Expression)Orchestrator() {
	expr.Start = getTime()  // получаем время начала вычислений

	expression := expr.Separated_expression  // берем разделенное выражение

	operation := make([]*Operation, 0)  // массив операций (будем работать с ним)
	// переводим выражение в операции
	for i := 0; i < len(expression)-1; i += 2 {
		num, _ := strconv.Atoi(expression[i])
		operation = append(operation, &Operation{res: num, status: true})  // число
		operation = append(operation, &Operation{oper: expression[i+1], status: true})  // операция
	}
	num, _ := strconv.Atoi(expression[len(expression)-1])
	operation = append(operation, &Operation{res: num, status: true})  // последнее число не попало в for

	wg := sync.WaitGroup{}  // будем ожидать окончания всех операций
	for {
		prev_oper := Operation{status: true}  // предыдущая операция
		next_oper := Operation{status: true}	// следующая операция
		expre := make([]*Operation, 0)			// здесь будем собирать следующие действия
		prev_num := &Operation{status: true}	// предыдущее число

		prev_num_used := false		// предыдущее число использовалось в выражении

		for i := 0; i < len(operation)-2; i += 2 {  // проходимся по operations
			// получаем число сейчас, следующее число, следующую операцию
			num1 := operation[i]
			num1Val := *num1
			if i < len(operation)-3 {
				next_oper = *operation[i+3]
			}else{
				next_oper = Operation{status: true}
			}
			oper := *operation[i+1]
			num2 := operation[i+2]
			num2Val := *num2

			// fmt.Println(num1, oper, num2, prev_oper, next_oper)
			// fmt.Println(prev_num, num1Val, num2Val, prev_num_used, prev_oper, oper, next_oper)
			
			// проверка на добавление числа на следующей итерации (в предыдущей операции не участвовало, сейчас тоже не будет)
			if (!num1Val.status || ((oper.oper == "+" || oper.oper == "-") &&
			(next_oper.status && (next_oper.oper == "*" || next_oper.oper == "/") || !prev_num.status)) || prev_oper.oper == "*" || prev_oper.oper == "/") && 
			(!prev_num_used || !num2Val.status){
	 
				
				expre = append(expre, num1)
				// fmt.Println("added num1", num1)
			 }
			 // добавление знака на следующую итерацию
			if (prev_oper.oper == "*" || prev_oper.oper == "/") || 
			(next_oper.oper == "*" || next_oper.oper == "/") && (oper.oper == "+" || oper.oper == "-") ||
			!num1Val.status || prev_num_used{
				expre = append(expre, &oper)
				// fmt.Println("added oper", oper)
			}
			
			// проведение операции
			if num1Val.status && num2Val.status && !prev_num_used && prev_oper.status && prev_oper.oper != "*" && prev_oper.oper != "/" && oper.status &&
			(oper.oper == "*" || oper.oper == "/" ||
			  (oper.oper == "+" || oper.oper == "-") && next_oper.status && next_oper.oper != "*" && next_oper.oper != "/"){
			
				prev_num_used = true  // num2 (в на следующей итерации num1) использовалось
				operAg := Operation{num1: num1.res, num2: num2.res, oper: oper.oper, status: false}  // создание новой операции
				// fmt.Println(operAg)
				expre = append(expre, &operAg)  // добавляем операцию
				wg.Add(1)
				// запускаем агента
				go func(operAgent *Operation) {
					defer wg.Done()
					operAgent.Agent()
					// fmt.Println("prom res", operAgent)
				}(expre[len(expre)-1])

			}else{
				prev_num_used = false  // число не использовалось
			}
			// последнее число, оторое не использовалось
			if i == len(operation) - 3 && !prev_num_used{ 
				expre = append(expre, num2)
				// fmt.Println("added last num", num2)
			 }
			prev_oper = oper  // обновляем предыдущих
			prev_num = num1
		}
		wg.Wait()  // ждем завершения всех операций

		res := *operation[0]  // берем первую операцию в operations
		operation = []*Operation{}  // обновляем operations
		for _, v := range expre {
			operation = append(operation, v)
		}

		// for _, v := range operation {
			// fmt.Println(v)
		// }
		// fmt.Println()
		if len(operation) == 0{  // в operatins ничего не осталось - посчиталось
			// fmt.Println("end")
			expr.Result = res.res // записываем результат
			expr.End = getTime()  // время конца вычислений
			expr.Status = true		// статус - посчитано
			Results = append(Results, *expr)  // добавляем в Results - чтобы server увидел
			break
		}
	}
}

// получение времени сейчас в нужном мне формате
func getTime() string{
	expr := time.Now()
	y, m, d := strconv.Itoa(expr.Year()), expr.Month().String(), strconv.Itoa(expr.Day())
	h, min, sec := strconv.Itoa(expr.Hour()), strconv.Itoa(expr.Minute()), strconv.Itoa(expr.Second())
	return d + " " + m + " " + y + "	" + h + ":" + min + ":" + sec
}

// агент
func (oper *Operation) Agent() {
	switch {
	case oper.oper == "-" || oper.num2 < 0:
		time.Sleep(time.Second * time.Duration(set.GetSettings().MinusTime))
		oper.res = oper.num1 + oper.num2

	case oper.oper == "+":
		time.Sleep(time.Second * time.Duration(set.GetSettings().PlusTime))
		oper.res = oper.num1 + oper.num2

	case oper.oper == "*":
		time.Sleep(time.Second * time.Duration(set.GetSettings().MultiplicationTime))
		oper.res = oper.num1 * oper.num2

	case oper.oper == "/":
		time.Sleep(time.Second * time.Duration(set.GetSettings().DivisionTime))
		oper.res = oper.num1 / oper.num2

	}
	oper.status = true

	// fmt.Println("gor is over")
}
