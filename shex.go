package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type function struct {
	params []string
	runs   []string
}

var variables = make(map[string]string)
var functions = make(map[string]function)
var keywords = [8]string{
	"var",
	"set",
	"if",
	"endif",
	"func",
	"endfunc",
	"return",
	"print",
}
var numbers = [10]int{
	0,
	1,
	2,
	3,
	4,
	5,
	6,
	7,
	8,
	9,
}

func parse_text(text string) string {
	parsed_text := ""

	for i := 0; i < len(text); i++ {
		if text[i] == '\\' && text[i+1] == '"' {
			parsed_text += "\""
		} else if text[i] != '"' {
			parsed_text += string(text[i])
		}
	}

	return parsed_text
}

func parse_math(math string) int {
	parsed_math := strings.Split(math, " ")
	result := 0

	for i := 0; i < len(parsed_math); i++ {
		if strings.Contains(parsed_math[i], "(") || strings.Contains(parsed_math[i], ")") {
			parsed_math[i] = strings.Replace(parsed_math[i], "(", "", -1)
			parsed_math[i] = strings.Replace(parsed_math[i], ")", "", -1)
			parsed_math[i] = strings.Replace(parsed_math[i], "\n", "", -1)
		}
	}

	for k, v := range variables {
		for j := 0; j < len(parsed_math); j++ {
			if parsed_math[j] == k {
				parsed_math[j] = v
			}
		}
	}

	for i := 0; i < len(parsed_math)-1; i++ {
		var num1, num2 int

		if i != 0 {
			num1, _ = strconv.Atoi(parsed_math[i-1])
			num2, _ = strconv.Atoi(parsed_math[i+1])

			if parsed_math[i] == "+" {
				result = num1 + num2
			} else if parsed_math[i] == "-" {
				result = num1 - num2
			} else if parsed_math[i] == "*" {
				result = num1 * num2
			} else if parsed_math[i] == "/" {
				result = num1 / num2
			}
		}
	}

	return result
}

func parse_exp(exp string) bool {
	new_exp := strings.Replace(exp, "[", "", -1)
	new_exp = strings.Replace(new_exp, "]", "", -1)
	new_exp = strings.Replace(new_exp, "\n", "", -1)
	parts := strings.Split(new_exp, " ")

	for i := 0; i < len(parts); i++ {
		for k, v := range variables {
			if strings.Contains(parts[i], k) {
				parts[i] = strings.Replace(parts[i], k, v, -1)
			}
		}
	}

	result := false
	var num1, num2 int

	for i := 0; i < len(numbers); i++ {
		if strings.HasPrefix(parts[0], strconv.Itoa(numbers[i])) {
			num1, _ = strconv.Atoi(parts[0])
		}

		if strings.HasPrefix(parts[2], strconv.Itoa(numbers[i])) {
			num2, _ = strconv.Atoi(parts[0])
		}
	}

	if num1 == 0 && num2 == 0 {
		if parts[1] == "<" {
			result = num1 < num2
		} else if parts[1] == ">" {
			result = num1 > num2
		} else if parts[1] == "=" {
			result = num1 == num2
		} else if parts[1] == "<=" {
			result = num1 <= num2
		} else if parts[1] == ">=" {
			result = num1 >= num2
		}
	} else {
		if parts[1] == "<" {
			result = parts[0] < parts[2]
		} else if parts[1] == ">" {
			result = parts[0] > parts[2]
		} else if parts[1] == "=" {
			result = parts[0] == parts[2]
		} else if parts[1] == "<=" {
			result = parts[0] <= parts[2]
		} else if parts[1] == ">=" {
			result = parts[0] >= parts[2]
		}
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
			parts := strings.Split(cleaned, " ")
			var_is_keyword := false

			for j := 0; j < len(keywords); j++ {
				if parts[1] == keywords[j] {
					var_is_keyword = true
				}
			}

			if var_is_keyword == false {
				variables[parts[1]] = "null"
			}
		} else if strings.HasPrefix(text[i], "set") {
			parts := strings.Split(text[i], " ")

			if len(parts) > 3 {
				for j := 3; j < len(parts); j++ {
					parts[2] += " " + parts[j]
				}
			}

			for k, _ := range variables {
				if parts[1] == k {
					if strings.HasPrefix(parts[2], "\"") {
						variables[k] = parse_text(parts[2])
					} else if strings.HasPrefix(parts[2], "(") {
						variables[k] = strconv.Itoa(parse_math(parts[2]))
					} else if strings.HasPrefix(parts[2], "{") {
						value := strings.Replace(parts[2], "{", "", -1)
						value = strings.Replace(value, "}", "", -1)
						parse([]string{value})
						variables[k] = variables["return"]
					} else {
						variables[k] = parts[2]
					}
				}
			}
		} else if strings.HasPrefix(text[i], "if") {
			inside_lines := []string{}

			j := i + 1
			for !strings.HasPrefix(text[j], "endif") {
				inside_lines = append(inside_lines, text[j])
				j++
			}

			no_keyword, _ := strings.CutPrefix(text[i], "if ")

			if parse_exp(no_keyword) {
				parse(inside_lines)
			}

			i = j
		} else if strings.HasPrefix(text[i], "func") {
			inside_lines := []string{}
			func_is_keyword := false

			j := i + 1

			for !strings.HasPrefix(text[j], "endfunc") {
				inside_lines = append(inside_lines, text[j])
				j++
			}

			parts := strings.Split(text[i], " ")
			params := []string{}

			i = j

			for k := 2; k < len(parts); k++ {
				params = append(params, parts[k])
			}

			for k := 0; k < len(keywords); k++ {
				if parts[1] == keywords[k] {
					func_is_keyword = true
				}
			}

			if func_is_keyword == false {
				functions[parts[1]] = function{params, inside_lines}
			}
		} else if strings.HasPrefix(text[i], "return") {
			no_keyword, _ := strings.CutPrefix(text[i], "return ")

			if strings.HasPrefix(no_keyword, "\"") {
				variables["return"] = parse_text(no_keyword)
			} else if strings.HasPrefix(no_keyword, "(") {
				variables["return"] = strconv.Itoa(parse_math(no_keyword))
			}
		} else if strings.HasPrefix(text[i], "print") {
			no_keyword, _ := strings.CutPrefix(text[i], "print ")

			for k, v := range variables {
				if strings.Contains(no_keyword, "$"+k+"$") {
					no_keyword = strings.Replace(no_keyword, "$"+k+"$", v, -1)
				}
			}

			if strings.HasPrefix(no_keyword, "\"") {
				fmt.Println(parse_text(no_keyword))
			} else if strings.HasPrefix(no_keyword, "(") {
				fmt.Println(parse_math(no_keyword))
			} else if strings.HasPrefix(no_keyword, "{") {
				code := strings.Replace(no_keyword, "{", "", -1)
				code = strings.Replace(code, "}", "", -1)
				parse([]string{code})
				fmt.Println(variables["return"])
			} else {
				fmt.Println(no_keyword)
			}
		} else {
			for key, _ := range functions {
				if strings.HasPrefix(text[i], key) {
					no_keyword, _ := strings.CutPrefix(text[i], "func ")
					no_keyword, _ = strings.CutPrefix(no_keyword, key+" ")
					parts := strings.Split(no_keyword, " ")
					values := []string{}

					for k := 0; k < len(parts); k++ {
						if strings.HasPrefix(parts[k], "\"") && !strings.HasSuffix(parts[k], "\"") {
							string_value := ""
							m := k
							for m < len(parts) {
								string_value += parts[m]
								if strings.HasSuffix(parts[m], "\"") {
									break
								} else {
									string_value += " "
								}
								m++
							}
							k = m
							values = append(values, string_value)
						} else {
							values = append(values, parts[k])
						}
					}

					for k := 0; k < len(functions[key].params); k++ {
						variables[functions[key].params[k]] = values[k]
					}

					parse(functions[key].runs)
				}
			}
		}
	}
}

func main() {
	raw_source, _ := os.ReadFile(os.Args[1])
	string_source := strings.Replace(string(raw_source), "\r\n", "\n", -1)
	source := strings.Split(string_source, "\n")
	parse(source)
}
