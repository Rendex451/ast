package main

import (
	"fmt"
	"math"
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

// Parser преобразует инфиксное выражение в синтаксическое дерево
type Parser struct {
	tokens []string
	pos    int
}

// NewParser создаёт новый парсер
func NewParser(expression string) *Parser {
	return &Parser{
		tokens: tokenize(expression),
		pos:    0,
	}
}

func (n *Node) ToInfix() string {
	if n == nil {
		return ""
	}

	if n.Left == nil && n.Right == nil {
		return n.Value
	}

	if contains([]string{"sin", "cos", "tan", "cot", "exp", "ln"}, n.Value) {
		return fmt.Sprintf("%s(%s)", n.Value, n.Left.ToInfix())
	}

	if contains([]string{"-sin", "-cos", "-tan", "-cot"}, n.Value) {
		return fmt.Sprintf("-%s(%s)", n.Value[1:], n.Left.ToInfix())
	}

	if n.Value == "^" {
		leftStr := n.Left.ToInfix()
		if contains([]string{"+", "-", "*", "/"}, n.Left.Value) {
			leftStr = fmt.Sprintf("(%s)", leftStr)
		}
		return fmt.Sprintf("%s^%s", leftStr, n.Right.ToInfix())
	}

	leftStr := n.Left.ToInfix()
	rightStr := n.Right.ToInfix()

	if n.Value == "*" || n.Value == "+" {
		if isAlphanumeric(n.Left.Value) || isNumber(n.Left.Value) || contains([]string{"sin", "cos", "tan", "cot", "exp", "ln", "-sin", "-cos", "-tan", "-cot"}, n.Left.Value) {
			// Нет скобок
		} else {
			leftStr = fmt.Sprintf("(%s)", leftStr)
		}
		if isAlphanumeric(n.Right.Value) || isNumber(n.Right.Value) || contains([]string{"sin", "cos", "tan", "cot", "exp", "ln", "-sin", "-cos", "-tan", "-cot"}, n.Right.Value) {
			// Нет скобок
		} else {
			rightStr = fmt.Sprintf("(%s)", rightStr)
		}
	} else {
		if contains([]string{"+", "-", "*", "/"}, n.Left.Value) {
			leftStr = fmt.Sprintf("(%s)", leftStr)
		}
		if contains([]string{"+", "-", "*", "/"}, n.Right.Value) {
			rightStr = fmt.Sprintf("(%s)", rightStr)
		}
	}

	return fmt.Sprintf("%s %s %s", leftStr, n.Value, rightStr)
}

// tokenize разбивает строку на токены
func tokenize(expression string) []string {
	var tokens []string
	expression = strings.ReplaceAll(expression, " ", "")
	current := ""
	i := 0
	for i < len(expression) {
		char := expression[i]
		if strings.ContainsAny(string(char), "+-*/()^") {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			if char == 'x' && current != "" && isDigit(current[len(current)-1]) {
				tokens = append(tokens, current)
				current = ""
				tokens = append(tokens, "*")
			}
			tokens = append(tokens, string(char))
		} else if isDigit(char) || char == '.' {
			current += string(char)
		} else if isLetter(char) || char == 'x' {
			if char == 'x' && current != "" && isDigit(current[len(current)-1]) {
				tokens = append(tokens, current)
				current = ""
				tokens = append(tokens, "*")
			}
			current += string(char)
		}
		i++
	}
	if current != "" {
		tokens = append(tokens, current)
	}
	return tokens
}

// isDigit проверяет, является ли символ цифрой
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isLetter проверяет, является ли символ буквой
func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// peek возвращает текущий токен
func (p *Parser) peek() string {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return ""
}

// consume извлекает текущий токен и сдвигает позицию
func (p *Parser) consume() string {
	token := p.peek()
	if token != "" {
		p.pos++
	}
	return token
}

// Parse парсит выражение
func (p *Parser) Parse() (*Node, error) {
	node := p.parseExpression()
	if p.pos < len(p.tokens) {
		return nil, fmt.Errorf("лишние токены в выражении")
	}
	return node, nil
}

// parseExpression парсит выражение (высший приоритет)
func (p *Parser) parseExpression() *Node {
	return p.parseAddSub()
}

// parseAddSub парсит сложение и вычитание
func (p *Parser) parseAddSub() *Node {
	node := p.parseMulDiv()
	for p.peek() == "+" || p.peek() == "-" {
		op := p.consume()
		right := p.parseMulDiv()
		node = &Node{Value: op, Left: node, Right: right}
	}
	return node
}

// parseMulDiv парсит умножение и деление
func (p *Parser) parseMulDiv() *Node {
	node := p.parsePower()
	for p.peek() == "*" || p.peek() == "/" {
		op := p.consume()
		right := p.parsePower()
		node = &Node{Value: op, Left: node, Right: right}
	}
	return node
}

// parsePower парсит возведение в степень
func (p *Parser) parsePower() *Node {
	node := p.parseUnary()
	for p.peek() == "^" {
		op := p.consume()
		right := p.parseUnary()
		node = &Node{Value: op, Left: node, Right: right}
	}
	return node
}

// parseUnary парсит унарные операции
func (p *Parser) parseUnary() *Node {
	if p.peek() == "-" {
		p.consume()
		return &Node{
			Value: "*",
			Left:  &Node{Value: "-1"},
			Right: p.parseFactor(),
		}
	}
	return p.parseFactor()
}

// parseFactor парсит факторы (числа, переменные, функции, скобки)
func (p *Parser) parseFactor() *Node {
	token := p.consume()
	if token == "" {
		panic("неожиданный конец выражения")
	}

	if contains([]string{"sin", "cos", "tan", "cot", "exp", "ln"}, token) {
		if p.peek() != "(" {
			panic(fmt.Sprintf("ожидалась '(' после %s", token))
		}
		p.consume()
		arg := p.parseExpression()
		if p.peek() != ")" {
			panic(fmt.Sprintf("ожидалась ')' после аргумента %s", token))
		}
		p.consume()
		return &Node{Value: token, Left: arg}
	}

	if token == "(" {
		expr := p.parseExpression()
		if p.peek() != ")" {
			panic("ожидалась ')'")
		}
		p.consume()
		return expr
	}

	if isNumber(token) || token == "x" || contains([]string{"pi", "e"}, token) {
		return &Node{Value: token}
	}

	panic(fmt.Sprintf("неожиданный токен: %s", token))
}

// Вспомогательные функции
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isConstant(node *Node) bool {
	if node == nil {
		return false
	}
	if isNumber(node.Value) || contains([]string{"pi", "e"}, node.Value) {
		return true
	}
	if node.Value == "^" && isConstant(node.Left) && isConstant(node.Right) {
		return true
	}
	return false
}

func parseNumber(s string) float64 {
	if s == "pi" {
		return math.Pi
	}
	if s == "e" {
		return math.E
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func computeOperation(a, b float64, op string) (float64, bool) {
	switch op {
	case "+":
		return a + b, true
	case "-":
		return a - b, true
	case "*":
		return a * b, true
	case "/":
		if b == 0 {
			return 0, false
		}
		return a / b, true
	case "^":
		return math.Pow(a, b), true
	default:
		return 0, false
	}
}

func formatNumber(num float64) string {
	if math.Abs(num-math.Floor(num)) < 1e-10 {
		return fmt.Sprintf("%d", int(num))
	}
	if math.Abs(num-math.Pi) < 1e-10 {
		return "pi"
	}
	if math.Abs(num-math.E) < 1e-10 {
		return "e"
	}
	return fmt.Sprintf("%v", num)
}
