package orkestrator

import (
	set "arifm_operations/server/settings"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Expression struct {
	Id                   int
	Start                string
	End                  string
	Expression           string
	Separated_expression []string
	Status               bool
	Result               int
}

type Operation struct {
	num1   int
	oper   string
	num2   int
	res    int
	status bool // false - computing true - computed
}
var Results []Expression
var Agents int
func (expr *Expression)Orchestrator() {
	expr.Start = getTime()

	expression := expr.Separated_expression

	operation := make([]*Operation, 0)
	for i := 0; i < len(expression)-1; i += 2 {
		num, _ := strconv.Atoi(expression[i])
		operation = append(operation, &Operation{res: num, status: true})
		operation = append(operation, &Operation{oper: expression[i+1], status: true})
	}
	num, _ := strconv.Atoi(expression[len(expression)-1])
	operation = append(operation, &Operation{res: num, status: true})

	wg := sync.WaitGroup{}
	for {
		prev_oper := Operation{status: true}
		next_oper := Operation{status: true}
		expre := make([]*Operation, 0)
		prev_num := &Operation{status: true}

		prev_num_used := false

		for i := 0; i < len(operation)-2; i += 2 {
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
			fmt.Println(prev_num, num1Val, num2Val, prev_num_used, prev_oper, oper, next_oper)
			
			if (!num1Val.status || ((oper.oper == "+" || oper.oper == "-") &&
			(next_oper.status && (next_oper.oper == "*" || next_oper.oper == "/") || !prev_num.status)) || prev_oper.oper == "*" || prev_oper.oper == "/") && 
			(!prev_num_used || !num2Val.status){
	 
				
				expre = append(expre, num1)
				// fmt.Println("added num1", num1)
			 }
			if (prev_oper.oper == "*" || prev_oper.oper == "/") || 
			(next_oper.oper == "*" || next_oper.oper == "/") && (oper.oper == "+" || oper.oper == "-") ||
			!num1Val.status || prev_num_used{
				expre = append(expre, &oper)
				// fmt.Println("added oper", oper)
			}
			
			
			if num1Val.status && num2Val.status && !prev_num_used && prev_oper.status && prev_oper.oper != "*" && prev_oper.oper != "/" && oper.status &&
			(oper.oper == "*" || oper.oper == "/" ||
			  (oper.oper == "+" || oper.oper == "-") && next_oper.status && next_oper.oper != "*" && next_oper.oper != "/"){
			
				prev_num_used = true
				operAg := Operation{num1: num1.res, num2: num2.res, oper: oper.oper, status: false}
				// fmt.Println(operAg)
				expre = append(expre, &operAg)
				wg.Add(1)
				go func(operAgent *Operation) {
					defer wg.Done()
					operAgent.Agent()
					// fmt.Println("prom res", operAgent)
				}(expre[len(expre)-1])

			}else{
				prev_num_used = false
			}
			if i == len(operation) - 3 && !prev_num_used{ 
				expre = append(expre, num2)
				// fmt.Println("added last num", num2)
			 }
			prev_oper = oper
			prev_num = num1
		}
		wg.Wait()

		res := *operation[0]
		operation = []*Operation{}
		for _, v := range expre {
			operation = append(operation, v)
		}

		// for _, v := range operation {
			// fmt.Println(v)
		// }
		// fmt.Println()
		if len(operation) == 0{
			// fmt.Println("end")
			expr.Result = res.res
			expr.End = getTime()
			expr.Status = true
			Results = append(Results, *expr)
			break
		}
	}
}

func getTime() string{
	expr := time.Now()
	y, m, d := strconv.Itoa(expr.Year()), expr.Month().String(), strconv.Itoa(expr.Day())
	h, min, sec := strconv.Itoa(expr.Hour()), strconv.Itoa(expr.Minute()), strconv.Itoa(expr.Second())
	return d + " " + m + " " + y + "	" + h + ":" + min + ":" + sec
}

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
