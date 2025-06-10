package main

import (
	"fmt"
	"strconv"
)

// Differentiate выполняет символьное дифференцирование по переменной x
func Differentiate(node *Node) *Node {
	if node == nil {
		return nil
	}

	if node.Value == "x" {
		return &Node{Value: "1"}
	}
	if isNumber(node.Value) || contains([]string{"pi", "e"}, node.Value) {
		return &Node{Value: "0"}
	}

	switch node.Value {
	case "sin":
		innerDiff := Differentiate(node.Left)
		return &Node{
			Value: "*",
			Left:  &Node{Value: "cos", Left: node.Left},
			Right: innerDiff,
		}
	case "cos":
		innerDiff := Differentiate(node.Left)
		return &Node{
			Value: "*",
			Left:  &Node{Value: "-sin", Left: node.Left},
			Right: innerDiff,
		}
	case "tan":
		innerDiff := Differentiate(node.Left)
		return &Node{
			Value: "*",
			Left: &Node{
				Value: "^",
				Left:  &Node{Value: "/", Left: &Node{Value: "1"}, Right: &Node{Value: "cos", Left: node.Left}},
				Right: &Node{Value: "2"},
			},
			Right: innerDiff,
		}
	case "cot":
		innerDiff := Differentiate(node.Left)
		return &Node{
			Value: "*",
			Left: &Node{
				Value: "-",
				Left:  &Node{Value: "0"},
				Right: &Node{
					Value: "^",
					Left:  &Node{Value: "/", Left: &Node{Value: "1"}, Right: &Node{Value: "sin", Left: node.Left}},
					Right: &Node{Value: "2"},
				},
			},
			Right: innerDiff,
		}
	case "exp":
		innerDiff := Differentiate(node.Left)
		return &Node{
			Value: "*",
			Left:  &Node{Value: "exp", Left: node.Left},
			Right: innerDiff,
		}
	case "ln":
		innerDiff := Differentiate(node.Left)
		return &Node{
			Value: "*",
			Left:  &Node{Value: "/", Left: &Node{Value: "1"}, Right: node.Left},
			Right: innerDiff,
		}
	case "+":
		return &Node{
			Value: "+",
			Left:  Differentiate(node.Left),
			Right: Differentiate(node.Right),
		}
	case "-":
		return &Node{
			Value: "-",
			Left:  Differentiate(node.Left),
			Right: Differentiate(node.Right),
		}
	case "*":
		if isConstant(node.Left) {
			return &Node{
				Value: "*",
				Left:  node.Left,
				Right: Differentiate(node.Right),
			}
		}
		if isConstant(node.Right) {
			return &Node{
				Value: "*",
				Left:  node.Right,
				Right: Differentiate(node.Left),
			}
		}
		return &Node{
			Value: "+",
			Left: &Node{
				Value: "*",
				Left:  Differentiate(node.Left),
				Right: node.Right,
			},
			Right: &Node{
				Value: "*",
				Left:  node.Left,
				Right: Differentiate(node.Right),
			},
		}
	case "/":
		return &Node{
			Value: "/",
			Left: &Node{
				Value: "-",
				Left: &Node{
					Value: "*",
					Left:  Differentiate(node.Left),
					Right: node.Right,
				},
				Right: &Node{
					Value: "*",
					Left:  node.Left,
					Right: Differentiate(node.Right),
				},
			},
			Right: &Node{
				Value: "^",
				Left:  node.Right,
				Right: &Node{Value: "2"},
			},
		}
	case "^":
		if node.Right != nil && isNumber(node.Right.Value) {
			exponent, _ := strconv.ParseFloat(node.Right.Value, 64)
			innerDiff := Differentiate(node.Left)
			return &Node{
				Value: "*",
				Left:  &Node{Value: fmt.Sprintf("%d", int(exponent))},
				Right: &Node{
					Value: "*",
					Left: &Node{
						Value: "^",
						Left:  node.Left,
						Right: &Node{Value: fmt.Sprintf("%d", int(exponent-1))},
					},
					Right: innerDiff,
				},
			}
		} else if isConstant(node.Left) {
			innerDiff := Differentiate(node.Right)
			return &Node{
				Value: "*",
				Left: &Node{
					Value: "^",
					Left:  node.Left,
					Right: node.Right,
				},
				Right: &Node{
					Value: "*",
					Left:  &Node{Value: "ln", Left: node.Left},
					Right: innerDiff,
				},
			}
		} else {
			return &Node{
				Value: "*",
				Left: &Node{
					Value: "^",
					Left:  node.Left,
					Right: node.Right,
				},
				Right: &Node{
					Value: "+",
					Left: &Node{
						Value: "*",
						Left:  Differentiate(node.Right),
						Right: &Node{Value: "ln", Left: node.Left},
					},
					Right: &Node{
						Value: "*",
						Left:  &Node{Value: "/", Left: node.Right, Right: node.Left},
						Right: Differentiate(node.Left),
					},
				},
			}
		}
	}

	return node
}
