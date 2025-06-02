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

	// Рекурсивно упрощаем поддеревья
	node.Left = Simplify(node.Left)
	node.Right = Simplify(node.Right)

	if node.Type == "number" || node.Type == "variable" {
		return node
	}

	if node.Type == "operator" {
		// Упрощение числовых операций
		if node.Left.Type == "number" && node.Right.Type == "number" {
			leftVal, _ := strconv.ParseFloat(node.Left.Value, 64)
			rightVal, _ := strconv.ParseFloat(node.Right.Value, 64)
			switch node.Value {
			case "+":
				return &Node{Type: "number", Value: fmt.Sprintf("%g", leftVal+rightVal)}
			case "-":
				return &Node{Type: "number", Value: fmt.Sprintf("%g", leftVal-rightVal)}
			case "*":
				return &Node{Type: "number", Value: fmt.Sprintf("%g", leftVal*rightVal)}
			case "/":
				if rightVal != 0 {
					return &Node{Type: "number", Value: fmt.Sprintf("%g", leftVal/rightVal)}
				}
			case "^":
				return &Node{Type: "number", Value: fmt.Sprintf("%g", math.Pow(leftVal, rightVal))}
			}
		}

		// Упрощения вида 0 + x, x + 0, x * 1, 1 * x, x * 0, 0 * x
		if node.Value == "+" {
			if node.Left.Type == "number" && node.Left.Value == "0" {
				return node.Right
			}
			if node.Right.Type == "number" && node.Right.Value == "0" {
				return node.Left
			}
		}
		if node.Value == "-" {
			if node.Right.Type == "number" && node.Right.Value == "0" {
				return node.Left
			}
			if node.Left.Type == "number" && node.Left.Value == "0" {
				return &Node{
					Type:  "operator",
					Value: "*",
					Left:  &Node{Type: "number", Value: "-1"},
					Right: node.Right,
				}
			}
		}
		if node.Value == "*" {
			if node.Left.Type == "number" && node.Left.Value == "1" {
				return node.Right
			}
			if node.Right.Type == "number" && node.Right.Value == "1" {
				return node.Left
			}
			if (node.Left.Type == "number" && node.Left.Value == "0") ||
				(node.Right.Type == "number" && node.Right.Value == "0") {
				return &Node{Type: "number", Value: "0"}
			}
			// Упрощение -1 * x или x * -1
			if node.Left.Type == "number" && node.Left.Value == "-1" {
				return &Node{
					Type:  "operator",
					Value: "*",
					Left:  &Node{Type: "number", Value: "-1"},
					Right: node.Right,
				}
			}
			if node.Right.Type == "number" && node.Right.Value == "-1" {
				return &Node{
					Type:  "operator",
					Value: "*",
					Left:  &Node{Type: "number", Value: "-1"},
					Right: node.Left,
				}
			}
			// Упрощение x^a * x^b -> x^(a+b)
			if node.Left.Type == "operator" && node.Left.Value == "^" && node.Right.Type == "operator" && node.Right.Value == "^" {
				if node.Left.Left.Type == "variable" && node.Right.Left.Type == "variable" && node.Left.Left.Value == node.Right.Left.Value {
					leftExp, _ := strconv.ParseFloat(node.Left.Right.Value, 64)
					rightExp, _ := strconv.ParseFloat(node.Right.Right.Value, 64)
					return &Node{
						Type:  "operator",
						Value: "^",
						Left:  node.Left.Left,
						Right: &Node{Type: "number", Value: fmt.Sprintf("%g", leftExp+rightExp)},
					}
				}
			}
			// Упрощение x^a * (k/x^b) -> k*x^(a-b)
			if node.Left.Type == "operator" && node.Left.Value == "^" && node.Right.Type == "operator" && node.Right.Value == "/" {
				if node.Left.Left.Type == "variable" && node.Right.Right.Type == "variable" && node.Left.Left.Value == node.Right.Right.Value {
					leftExp, _ := strconv.ParseFloat(node.Left.Right.Value, 64)
					if node.Right.Left.Type == "number" {
						coef, _ := strconv.ParseFloat(node.Right.Left.Value, 64)
						return &Node{
							Type:  "operator",
							Value: "*",
							Left:  &Node{Type: "number", Value: fmt.Sprintf("%g", coef)},
							Right: &Node{
								Type:  "operator",
								Value: "^",
								Left:  node.Left.Left,
								Right: &Node{Type: "number", Value: fmt.Sprintf("%g", leftExp-1)},
							},
						}
					}
				}
			}
			// Упрощение x * (k/x) -> k
			if node.Left.Type == "variable" && node.Right.Type == "operator" && node.Right.Value == "/" {
				if node.Right.Right.Type == "variable" && node.Left.Value == node.Right.Right.Value {
					if node.Right.Left.Type == "number" {
						return node.Right.Left
					}
				}
			}
			// Упрощение k * (1/x) -> k/x
			if node.Left.Type == "number" && node.Right.Type == "operator" && node.Right.Value == "/" {
				if node.Right.Left.Type == "number" && node.Right.Left.Value == "1" && node.Right.Right.Type == "variable" {
					return &Node{
						Type:  "operator",
						Value: "/",
						Left:  node.Left,
						Right: node.Right.Right,
					}
				}
			}
			// Упрощение (1/x) * k -> k/x
			if node.Right.Type == "number" && node.Left.Type == "operator" && node.Left.Value == "/" {
				if node.Left.Left.Type == "number" && node.Left.Left.Value == "1" && node.Left.Right.Type == "variable" {
					return &Node{
						Type:  "operator",
						Value: "/",
						Left:  node.Right,
						Right: node.Left.Right,
					}
				}
			}
		}
		if node.Value == "^" {
			if node.Right.Type == "number" && node.Right.Value == "1" {
				return node.Left
			}
			if node.Right.Type == "number" && node.Right.Value == "0" {
				return &Node{Type: "number", Value: "1"}
			}
			if node.Left.Type == "number" && node.Left.Value == "1" {
				return &Node{Type: "number", Value: "1"}
			}
			// Упрощение x^1 -> x
			if node.Right.Type == "number" && node.Right.Value == "1" {
				return node.Left
			}
			// Упрощение x^0 -> 1
			if node.Right.Type == "number" && node.Right.Value == "0" {
				return &Node{Type: "number", Value: "1"}
			}
		}
		if node.Value == "/" {
			if node.Left.Type == "number" && node.Left.Value == "0" {
				return &Node{Type: "number", Value: "0"}
			}
			if node.Right.Type == "number" && node.Right.Value == "1" {
				return node.Left
			}
		}
	}

	if node.Type == "function" {
		if node.Left.Type == "number" {
			val, _ := strconv.ParseFloat(node.Left.Value, 64)
			switch node.Value {
			case "sin":
				return &Node{Type: "number", Value: fmt.Sprintf("%g", math.Sin(val))}
			case "cos":
				return &Node{Type: "number", Value: fmt.Sprintf("%g", math.Cos(val))}
			case "exp":
				return &Node{Type: "number", Value: fmt.Sprintf("%g", math.Exp(val))}
			case "ln":
				if val > 0 {
					return &Node{Type: "number", Value: fmt.Sprintf("%g", math.Log(val))}
				}
			case "tg":
				return &Node{Type: "number", Value: fmt.Sprintf("%g", math.Tan(val))}
			case "ctg":
				sinVal := math.Sin(val)
				if sinVal != 0 {
					return &Node{Type: "number", Value: fmt.Sprintf("%g", math.Cos(val)/sinVal)}
				}
			}
		}
	}

	// Упрощение (-1)*sin(x) -> -sin(x)
	if node.Type == "operator" && node.Value == "*" && node.Left.Type == "number" && node.Left.Value == "-1" && node.Right.Type == "function" {
		return &Node{
			Type:  "operator",
			Value: "*",
			Left:  &Node{Type: "number", Value: "-1"},
			Right: node.Right,
		}
	}

	return node
}
