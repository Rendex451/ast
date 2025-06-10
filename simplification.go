package main

import (
	"fmt"
	"math"
	"strconv"
)

// Simplify упрощает синтаксическое дерево
func Simplify(node *Node) *Node {
	if node == nil {
		return nil
	}

	node.Left = Simplify(node.Left)
	node.Right = Simplify(node.Right)

	// Упрощение функций
	if node.Value == "ln" && node.Left != nil && node.Left.Value == "e" {
		return &Node{Value: "1"}
	}
	if node.Value == "exp" && node.Left != nil && node.Left.Value == "ln" {
		return node.Left.Left
	}
	if node.Value == "ln" && node.Left != nil && node.Left.Value == "exp" {
		return node.Left.Left
	}

	// sin(u)/cos(u) -> tan(u)
	if node.Value == "/" && node.Left != nil && node.Left.Value == "sin" &&
		node.Right != nil && node.Right.Value == "cos" && node.Left.Left.ToInfix() == node.Right.Left.ToInfix() {
		return &Node{Value: "tan", Left: node.Left.Left}
	}

	// (cos(u)*cos(u))/(cos(u)^2) -> (1/cos(u))^2
	if node.Value == "/" && node.Left != nil && node.Left.Value == "*" &&
		node.Right != nil && node.Right.Value == "^" &&
		node.Right.Left != nil && node.Right.Left.Value == "cos" &&
		node.Right.Right != nil && node.Right.Right.Value == "2" &&
		node.Left.Left != nil && node.Left.Left.Value == "cos" &&
		node.Left.Right != nil && node.Left.Right.Value == "cos" &&
		node.Left.Left.Left.ToInfix() == node.Left.Right.Left.ToInfix() &&
		node.Left.Left.Left.ToInfix() == node.Right.Left.Left.ToInfix() {
		return &Node{
			Value: "^",
			Left:  &Node{Value: "/", Left: &Node{Value: "1"}, Right: &Node{Value: "cos", Left: node.Left.Left.Left}},
			Right: &Node{Value: "2"},
		}
	}

	// (-sin(u)*sin(u))/(cos(u)^2) -> (1/cos(u))^2
	if node.Value == "/" && node.Left != nil && node.Left.Value == "*" &&
		node.Right != nil && node.Right.Value == "^" &&
		node.Right.Left != nil && node.Right.Left.Value == "cos" &&
		node.Right.Right != nil && node.Right.Right.Value == "2" &&
		node.Left.Left != nil && node.Left.Left.Value == "-" &&
		node.Left.Left.Right != nil && node.Left.Left.Right.Value == "sin" &&
		node.Left.Right != nil && node.Left.Right.Value == "sin" &&
		node.Left.Left.Right.Left.ToInfix() == node.Left.Right.Left.ToInfix() &&
		node.Left.Left.Right.Left.ToInfix() == node.Right.Left.Left.ToInfix() {
		return &Node{
			Value: "^",
			Left:  &Node{Value: "/", Left: &Node{Value: "1"}, Right: &Node{Value: "cos", Left: node.Left.Left.Right.Left}},
			Right: &Node{Value: "2"},
		}
	}

	// Константы
	if contains([]string{"+", "-", "*", "/", "^"}, node.Value) && node.Left != nil && node.Right != nil {
		if isConstant(node.Left) && isConstant(node.Right) {
			a := parseNumber(node.Left.Value)
			b := parseNumber(node.Right.Value)
			if result, ok := computeOperation(a, b, node.Value); ok {
				return &Node{Value: formatNumber(result)}
			}
		}
	}

	// Упрощения для *
	if node.Value == "*" && node.Left != nil && node.Left.Value == "1" {
		return node.Right
	}
	if node.Value == "*" && node.Right != nil && node.Right.Value == "1" {
		return node.Left
	}
	if node.Value == "*" && node.Left != nil && node.Left.Value == "0" {
		return &Node{Value: "0"}
	}
	if node.Value == "*" && node.Right != nil && node.Right.Value == "0" {
		return &Node{Value: "0"}
	}
	if node.Value == "*" && node.Right != nil && node.Right.Value == "/" && node.Right.Left != nil && node.Right.Left.Value == "1" {
		return &Node{Value: "/", Left: node.Left, Right: node.Right.Right}
	}
	if node.Value == "*" && node.Left != nil && node.Left.Value == "/" && node.Left.Left != nil && node.Left.Left.Value == "1" {
		return &Node{Value: "/", Left: node.Right, Right: node.Left.Right}
	}

	// Упрощения для /
	if node.Value == "/" && node.Left != nil && node.Right != nil && node.Left.ToInfix() == node.Right.ToInfix() {
		return &Node{Value: "1"}
	}
	if node.Value == "/" && node.Left != nil && node.Left.Value == "0" && node.Right != nil && node.Right.Value != "0" {
		return &Node{Value: "0"}
	}

	// Упрощения для +
	if node.Value == "+" && node.Left != nil && node.Left.Value == "0" {
		return node.Right
	}
	if node.Value == "+" && node.Right != nil && node.Right.Value == "0" {
		return node.Left
	}

	// Упрощения для -
	if node.Value == "-" && node.Right != nil && node.Right.Value == "0" {
		return node.Left
	}

	// Упрощения для ^
	if node.Value == "^" && node.Right != nil && node.Right.Value == "1" {
		return node.Left
	}
	if node.Value == "^" && node.Right != nil && node.Right.Value == "0" {
		return &Node{Value: "1"}
	}
	if node.Value == "^" && node.Left != nil && node.Left.Value == "1" {
		return &Node{Value: "1"}
	}
	if node.Value == "^" && node.Right != nil && isNumber(node.Right.Value) {
		exponent, _ := strconv.ParseFloat(node.Right.Value, 64)
		if math.Abs(exponent-math.Floor(exponent)) < 1e-10 && exponent >= 0 {
			node.Right = &Node{Value: fmt.Sprintf("%d", int(exponent))}
			if exponent == 1 {
				return node.Left
			}
		}
	}

	// Упрощение чисел
	if isNumber(node.Value) {
		num, _ := strconv.ParseFloat(node.Value, 64)
		if math.Abs(num-math.Floor(num)) < 1e-10 {
			return &Node{Value: fmt.Sprintf("%d", int(num))}
		}
	}

	// -1 * sin(u) -> -sin(u)
	if node.Value == "*" && node.Left != nil && node.Left.Value == "-1" && node.Right != nil &&
		contains([]string{"sin", "cos", "tan", "cot"}, node.Right.Value) {
		return &Node{Value: "-" + node.Right.Value, Left: node.Right.Left}
	}

	return node
}

// InfixToTree преобразует инфиксное выражение в дерево
func InfixToTree(expression string) (*Node, error) {
	parser := NewParser(expression)
	return parser.Parse()
}

// AnalyzeExpression анализирует выражение
func AnalyzeExpression(expr string) map[string]string {
	result := make(map[string]string)
	result["Expression"] = expr
	defer func() {
		if r := recover(); r != nil {
			result["Tree"] = fmt.Sprintf("Error: %v", r)
			result["Derivative"] = "None"
			result["Simplified Derivative"] = "None"
		}
	}()

	tree, err := InfixToTree(expr)
	if err != nil {
		result["Tree"] = fmt.Sprintf("Error: %v", err)
		result["Derivative"] = "None"
		result["Simplified Derivative"] = "None"
		return result
	}

	tree = Simplify(tree)
	derivative := Differentiate(tree)
	simplifiedDerivative := Simplify(derivative)

	result["Tree"] = tree.ToInfix()
	result["Derivative"] = derivative.ToInfix()
	result["Simplified Derivative"] = simplifiedDerivative.ToInfix()
	return result
}
