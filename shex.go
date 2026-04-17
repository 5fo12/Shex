package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
)

type function struct {
	params map[string]parameter
	runs   []string
}

type parameter struct {
	uname string
	vtype string
}

type variable struct {
	value string
	vtype string
}

var vars = make(map[string]variable)
var funcs = make(map[string]function)
var keywords = [9]string{
	"var",
	"set",
	"if",
	"endif",
	"func",
	"endfunc",
	"return",
	"print",
	"read",
}

func parseText(text string) string {
	parsedText := ""

	for i := 0; i < len(text); i++ {
		if text[i] == '\\' && text[i+1] == '"' {
			parsedText += "\""
		} else if text[i] != '"' {
			parsedText += string(text[i])
		}
	}

	return parsedText
}

func parseMath(math string) float64 {
	parsedMath := strings.Split(math, " ")
	var result float64

	for i := 0; i < len(parsedMath); i++ {
		if strings.Contains(parsedMath[i], "(") || strings.Contains(parsedMath[i], ")") {
			parsedMath[i] = strings.Replace(parsedMath[i], "(", "", -1)
			parsedMath[i] = strings.Replace(parsedMath[i], ")", "", -1)
			parsedMath[i] = strings.Replace(parsedMath[i], "\n", "", -1)
		}
	}

	for k, v := range vars {
		for j := 0; j < len(parsedMath); j++ {
			if parsedMath[j] == k && v.vtype == "number" {
				parsedMath[j] = v.value
			}
		}
	}

	for i := 0; i < len(parsedMath)-1; i++ {
		var num1, num2 float64

		if i != 0 {
			num1, _ = strconv.ParseFloat(parsedMath[i-1], 64)
			num2, _ = strconv.ParseFloat(parsedMath[i+1], 64)

			switch parsedMath[i] {
			case "+":
				result = num1 + num2
			case "-":
				result = num1 - num2
			case "*":
				result = num1 * num2
			case "/":
				result = num1 / num2
			}
		}
	}

	return result
}

func parseExp(exp string) bool {
	newExp := strings.Replace(exp, "[", "", -1)
	newExp = strings.Replace(newExp, "]", "", -1)
	newExp = strings.Replace(newExp, "\n", "", -1)
	parts := strings.Split(newExp, " ")

	for i := 0; i < len(parts); i++ {
		for k, v := range vars {
			parts[i] = strings.Replace(parts[i], k, v.value, -1)
		}
	}

	result := false
	num1, _ := strconv.ParseFloat(parts[0], 64)
	num2, _ := strconv.ParseFloat(parts[2], 64)

	switch parts[1] {
	case "<":
		result = num1 < num2
	case ">":
		result = num1 > num2
	case "=":
		result = num1 == num2
	case "<=":
		result = num1 <= num2
	case ">=":
		result = num1 >= num2
	}

	return result
}

func parse(text []string) {
	for i := 0; i < len(text); i++ {
		if strings.HasPrefix(text[i], "!") {
			continue
		} else if strings.HasPrefix(text[i], "<") {
			j := i + 1

			for !strings.HasSuffix(text[j], ">") {
				j += 1
			}

			i = j
		} else if strings.HasPrefix(text[i], "var") {
			cleaned := strings.Replace(text[i], "\n", "", -1)
			cleaned = strings.Replace(text[i], ":", " ", -1)
			parts := strings.Split(cleaned, " ")
			varIsKeyword := false

			for j := 0; j < len(keywords); j++ {
				if parts[1] == keywords[j] {
					varIsKeyword = true
				}
			}

			if varIsKeyword == false {
				vars[parts[1]] = variable{"null", parts[2]}
			}
		} else if strings.HasPrefix(text[i], "set") {
			parts := strings.Split(text[i], " ")

			if len(parts) > 3 {
				for j := 3; j < len(parts); j++ {
					parts[2] += " " + parts[j]
				}
			}

			for k, _ := range vars {
				if parts[1] == k {
					if strings.HasPrefix(parts[2], "\"") {
						vars[k] = variable{parseText(parts[2]), "text"}
					} else if strings.HasPrefix(parts[2], "(") {
						vars[k] = variable{strconv.FormatFloat(parseMath(parts[2]), 'f', -1, 64), "number"}
					} else if strings.HasPrefix(parts[2], "{") {
						value := strings.Replace(parts[2], "{", "", -1)
						value = strings.Replace(value, "}", "", -1)
						parse([]string{value})
						vars[k] = vars["return"]
					} else {
						vars[k] = variable{parts[2], "number"}
					}
				}
			}
		} else if strings.HasPrefix(text[i], "if") {
			insideLines := []string{}

			j := i + 1
			for !strings.HasPrefix(text[j], "endif") {
				insideLines = append(insideLines, text[j])
				j++
			}

			noKeyword, _ := strings.CutPrefix(text[i], "if ")

			if parseExp(noKeyword) {
				parse(insideLines)
			}

			i = j
		} else if strings.HasPrefix(text[i], "func") {
			insideLines := []string{}
			funcIsKeyword := false

			j := i + 1

			for !strings.HasPrefix(text[j], "endfunc") {
				insideLines = append(insideLines, text[j])
				j++
			}

			parts := strings.Split(text[i], " ")
			params := make(map[string]parameter)

			i = j

			for k := 2; k < len(parts); k++ {
				param := strings.Split(parts[k], ":")
				paramName := strings.ToUpper(param[0]) + strconv.Itoa(rand.IntN(1000))
				params[paramName] = parameter{param[0], param[1]}

			}

			for k := 0; k < len(keywords); k++ {
				if parts[1] == keywords[k] {
					funcIsKeyword = true
				}
			}

			for k := 0; k < len(insideLines); k++ {
				for name, param := range params {
					insideLines[k] = strings.Replace(insideLines[k], strings.ToLower(param.uname), name, -1)
				}
			}

			if funcIsKeyword == false {
				funcs[parts[1]] = function{params, insideLines}
			}
		} else if strings.HasPrefix(text[i], "return") {
			noKeyword, _ := strings.CutPrefix(text[i], "return ")

			if strings.HasPrefix(noKeyword, "\"") {
				vars["return"] = variable{parseText(noKeyword), "text"}
			} else if strings.HasPrefix(noKeyword, "(") {
				vars["return"] = variable{strconv.FormatFloat(parseMath(noKeyword), 'f', -1, 64), "number"}
			} else {
				vars["return"] = variable{noKeyword, "number"}
			}
		} else if strings.HasPrefix(text[i], "print") {
			noKeyword, _ := strings.CutPrefix(text[i], "print ")

			for k, v := range vars {
				if strings.Contains(noKeyword, "$"+k+"$") {
					noKeyword = strings.Replace(noKeyword, "$"+k+"$", v.value, -1)
				}
			}

			if strings.HasPrefix(noKeyword, "\"") {
				fmt.Println(parseText(noKeyword))
			} else if strings.HasPrefix(noKeyword, "(") {
				fmt.Printf("%.4v\n", parseMath(noKeyword))
			} else if strings.HasPrefix(noKeyword, "{") {
				code := strings.Replace(noKeyword, "{", "", -1)
				code = strings.Replace(code, "}", "", -1)
				parse([]string{code})
				fmt.Println(vars["return"].value)
			} else {
				fmt.Println(noKeyword)
			}
		} else if strings.HasPrefix(text[i], "read") {
			var input string
			fmt.Scan(&input)
			vars["return"] = variable{input, "text"}
		} else {
			for funcName := range funcs {
				if strings.HasPrefix(text[i], funcName) {
					noKeyword, _ := strings.CutPrefix(text[i], "func ")
					noKeyword, _ = strings.CutPrefix(noKeyword, funcName+" ")
					parts := strings.Split(noKeyword, " ")
					values := []string{}

					for k := 0; k < len(parts); k++ {
						if strings.HasPrefix(parts[k], "\"") && !strings.HasSuffix(parts[k], "\"") {
							stringValue := ""
							m := k
							for m < len(parts) {
								stringValue += parts[m]
								if strings.HasSuffix(parts[m], "\"") {
									break
								} else {
									stringValue += " "
								}
								m++
							}
							k = m
							values = append(values, stringValue)
						} else {
							values = append(values, parts[k])
						}
					}

					k := 0
					for name, param := range funcs[funcName].params {
						vars[name] = variable{values[k], param.vtype}
						k++
					}

					parse(funcs[funcName].runs)
				}
			}
		}
	}
}

func main() {
	rawSource, _ := os.ReadFile(os.Args[1])
	stringSource := strings.Replace(string(rawSource), "\r\n", "\n", -1)
	source := strings.Split(stringSource, "\n")
	parse(source)
}
