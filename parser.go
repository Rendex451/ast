package main

import (
	"strconv"
	"strings"
)

// Node представляет узел синтаксического дерева
type Node struct {
	Type  string // "number", "variable", "operator", "function"
	Value string
	Left  *Node
	Right *Node
}

// tokenize разбивает инфиксное выражение на токены
func tokenize(expr string) []string {
	expr = strings.ReplaceAll(expr, " ", "")
	tokens := []string{}
	num := ""
	for i := 0; i < len(expr); i++ {
		c := expr[i]
		if isDigit(c) || c == '.' {
			num += string(c)
		} else {
			if num != "" {
				tokens = append(tokens, num)
				num = ""
			}
			if isOperator(c) || c == '(' || c == ')' || c == ',' {
				tokens = append(tokens, string(c))
			} else {
				fun := ""
				for i < len(expr) && isLetter(expr[i]) {
					fun += string(expr[i])
					i++
				}
				tokens = append(tokens, fun)
				i--
			}
		}
	}
	if num != "" {
		tokens = append(tokens, num)
	}
	return tokens
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isOperator(c byte) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '^'
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	case "^":
		return 3
	}
	return 0
}

func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

func isVariable(token string) bool {
	return len(token) == 1 && isLetter(token[0])
}

func isFunction(token string) bool {
	return token == "sin" || token == "cos" || token == "exp" || token == "ln" ||
		token == "tg" || token == "ctg"
}

func processOperator(stack *[]*Node, opStack *[]string) {
	if len(*stack) == 0 || len(*opStack) == 0 {
		return
	}
	op := (*opStack)[len(*opStack)-1]
	*opStack = (*opStack)[:len(*opStack)-1]

	node := &Node{Value: op}
	if isFunction(op) {
		node.Type = "function"
		if len(*stack) > 0 {
			node.Left = (*stack)[len(*stack)-1]
			*stack = (*stack)[:len(*stack)-1]
		}
	} else {
		node.Type = "operator"
		if len(*stack) > 1 {
			node.Right = (*stack)[len(*stack)-1]
			*stack = (*stack)[:len(*stack)-1]
			node.Left = (*stack)[len(*stack)-1]
			*stack = (*stack)[:len(*stack)-1]
		}
	}
	*stack = append(*stack, node)
}

// TreeToInfix преобразует синтаксическое дерево обратно в инфиксную запись
func TreeToInfix(node *Node) string {
	if node == nil {
		return ""
	}
	switch node.Type {
	case "number", "variable":
		return node.Value
	case "function":
		return node.Value + "(" + TreeToInfix(node.Left) + ")"
	case "operator":
		left := TreeToInfix(node.Left)
		right := TreeToInfix(node.Right)

		if node.Left.Type == "operator" && precedence(node.Left.Value) < precedence(node.Value) {
			left = "(" + left + ")"
		}
		if node.Right.Type == "operator" && precedence(node.Right.Value) <= precedence(node.Value) {
			right = "(" + right + ")"
		}
		return left + node.Value + right
	}
	return ""
}

// InfixToTree преобразует инфиксное выражение в синтаксическое дерево
func InfixToTree(expr string) *Node {
	tokens := tokenize(expr)
	var stack []*Node
	var opStack []string

	for _, token := range tokens {
		if isNumber(token) {
			stack = append(stack, &Node{Type: "number", Value: token})
		} else if isVariable(token) {
			stack = append(stack, &Node{Type: "variable", Value: token})
		} else if isFunction(token) {
			opStack = append(opStack, token)
		} else if token == "(" {
			opStack = append(opStack, token)
		} else if token == ")" {
			for len(opStack) > 0 && opStack[len(opStack)-1] != "(" {
				processOperator(&stack, &opStack)
			}
			if len(opStack) > 0 {
				opStack = opStack[:len(opStack)-1] // Удаляем '('
			}
			if len(opStack) > 0 && isFunction(opStack[len(opStack)-1]) {
				processOperator(&stack, &opStack)
			}
		} else if isOperator(token[0]) {
			for len(opStack) > 0 && opStack[len(opStack)-1] != "(" &&
				precedence(opStack[len(opStack)-1]) >= precedence(token) {
				processOperator(&stack, &opStack)
			}
			opStack = append(opStack, token)
		}
	}

	for len(opStack) > 0 {
		processOperator(&stack, &opStack)
	}

	if len(stack) == 0 {
		return nil
	}
	return stack[0]
}
