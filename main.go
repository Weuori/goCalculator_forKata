package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type CalcError struct {
	Err    error
	IsExit bool
}

func main() {

	fmt.Println("Строковый калькулятор GoLang")
	Calculator()
}

func Calculator() (result string, err CalcError) {
	task := InputString()
	rawX, rawY, operator := ParseTask(task)
	result, err = Compute(rawX, rawY, operator)

	if err.Err != nil {
		switch err.IsExit {
		case true:
			fmt.Println(err.Err)
			os.Exit(1)
		case false:
			fmt.Println(err.Err)
		}
	} else {
		fmt.Printf("Ответ: %s \n", result)
	}
	return Calculator()
}

func InputString() (task string) {
	fmt.Println("------------------------------")
	fmt.Println("Введите математическую задачу:")
	task, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.ReplaceAll(strings.TrimSpace(task), " ", "")
}

func ParseTask(task string) (rawX []string, rawY []string, operator string) {
	for _, value := range task {
		char := string(value)

		if char == "+" || char == "-" || char == "*" || char == "/" {
			operator += char
			continue
		}

		if len(operator) == 0 {
			rawX = append(rawX, string(value))
		} else {
			rawY = append(rawY, string(value))
		}
	}
	return rawX, rawY, operator
}

func ValidateOperand(rawOperand []string) (operand string, isCypher bool, isRoman bool) {
	operand = strings.Join(rawOperand, "")
	isCypher, _ = regexp.MatchString(`^([1-9]?[0-9]+)$`, operand)
	isRoman, _ = regexp.MatchString(`^M{0,3}(CM|CD|D?C{0,3})?(XC|XL|L?X{0,3})?(IX|IV|V?I{0,3})?$`, operand)

	return operand, isCypher, isRoman
}

func CypherCalc(strX string, strY string, operator string) (result string, err CalcError) {

	var tmp int
	x, _ := strconv.Atoi(strX)
	y, _ := strconv.Atoi(strY)

	if x == 0 || x > 10 || y == 0 || y > 10 {
		return strconv.Itoa(tmp), CalcError{
			Err:    fmt.Errorf("ошибка ввода. Операнд должен быть в диапазоне от 1 до 10"),
			IsExit: true,
		}
	}

	switch operator {
	case "+":
		tmp = x + y
	case "-":
		tmp = x - y
	case "*":
		tmp = x * y
	case "/":
		tmp = x / y
	default:
	}
	return strconv.Itoa(tmp), CalcError{}
}

func RomanCalc(romanX string, romanY string, operator string) (result string, err CalcError) {

	var tmp int
	intX := RomanToInt(romanX)
	intY := RomanToInt(romanY)

	switch {
	case len(romanX) == 0 || len(romanY) == 0:
		return "", CalcError{
			Err:    fmt.Errorf("ошибка ввода. Неправильный формат математической операции"),
			IsExit: true,
		}
	case intX > 10 || intY > 10:
		return "", CalcError{
			Err:    fmt.Errorf("ошибка ввода. Операнд должен быть в диапазоне от I до X"),
			IsExit: true,
		}
	}

	switch operator {
	case "+":
		tmp = intX + intY
	case "-":
		tmp = intX - intY
	case "*":
		tmp = intX * intY
	case "/":
		tmp = intX / intY
	default:
	}

	if tmp < 1 {
		return "", CalcError{
			Err:    fmt.Errorf("ошибка ввода. Результат вычислений не может быть меньше 1"),
			IsExit: false,
		}
	}

	return IntToRoman(tmp), CalcError{}
}

func RomanToInt(romanString string) (romanInt int) {
	romanDict := map[uint]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}
	var currentValue, lastValue int
	for i := len(romanString) - 1; i >= 0; i-- {
		currentValue = romanDict[uint(romanString[i])]
		if currentValue < lastValue {
			romanInt -= currentValue
		} else {
			romanInt += currentValue
		}
		lastValue = currentValue
	}
	return romanInt
}

func IntToRoman(romanInt int) (romanString string) {
	romanDict := map[string]int{
		"M":  1000,
		"CM": 900,
		"D":  500,
		"CD": 400,
		"C":  100,
		"XC": 90,
		"L":  50,
		"XL": 40,
		"X":  10,
		"IX": 9,
		"V":  5,
		"IV": 4,
		"I":  1,
	}

	dictKeys := [...]string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}

	for _, k := range dictKeys {
		for romanInt >= romanDict[k] {
			romanInt -= romanDict[k]
			romanString += k
		}
	}
	return romanString
}

func Compute(rawX []string, rawY []string, operator string) (result string, err CalcError) {

	x, isCypherX, isRomanX := ValidateOperand(rawX)
	y, isCypherY, isRomanY := ValidateOperand(rawY)

	switch {
	case (!isCypherX && !isRomanX) || (!isCypherY && !isRomanY) || (len(operator) != 1) || len(rawX) == 0 || len(rawY) == 0:
		err = CalcError{
			Err:    fmt.Errorf("ошибка ввода. Неправильный формат математической операции"),
			IsExit: false,
		}
	case (isCypherX && isRomanY) || (isRomanX && isCypherY):
		err = CalcError{
			Err:    fmt.Errorf("ошибка ввода. Используются разные системы счисления"),
			IsExit: true,
		}
	case isCypherX && isCypherY:
		result, err = CypherCalc(x, y, operator)
	case isRomanX && isRomanY:
		result, err = RomanCalc(x, y, operator)
	default:
	}
	return result, err
}
