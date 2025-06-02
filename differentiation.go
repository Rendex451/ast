package main

// Differentiate выполняет символьное дифференцирование по переменной x
func Differentiate(node *Node, variable string) *Node {
	if node == nil {
		return nil
	}

	if isNumber(node.Value) {
		return &Node{Value: "0"}
	}

	if node.Value == variable {
		return &Node{Value: "1"}
	}
	if isVariable(node.Value) {
		return &Node{Value: "0"}
	}

	switch node.Value {
	case "+":
		return &Node{
			Value: "+",
			Left:  Differentiate(node.Left, variable),
			Right: Differentiate(node.Right, variable),
		}
	case "-":
		return &Node{
			Value: "-",
			Left:  Differentiate(node.Left, variable),
			Right: Differentiate(node.Right, variable),
		}
	case "*":
		// (u*v)' = u'*v + u*v'
		return &Node{
			Value: "+",
			Left: &Node{
				Value: "*",
				Left:  Differentiate(node.Left, variable),
				Right: node.Right,
			},
			Right: &Node{
				Value: "*",
				Left:  node.Left,
				Right: Differentiate(node.Right, variable),
			},
		}
	case "/":
		// (u/v)' = (u'*v - u*v')/v^2
		return &Node{
			Value: "/",
			Left: &Node{
				Value: "-",
				Left: &Node{
					Value: "*",
					Left:  Differentiate(node.Left, variable),
					Right: node.Right,
				},
				Right: &Node{
					Value: "*",
					Left:  node.Left,
					Right: Differentiate(node.Right, variable),
				},
			},
			Right: &Node{
				Value: "^",
				Left:  node.Right,
				Right: &Node{Value: "2"},
			},
		}

	case "sin":
		// (sin(u))' = cos(u) * u'
		return &Node{
			Value: "*",
			Left: &Node{
				Value: "cos",
				Left:  node.Left,
			},
			Right: Differentiate(node.Left, variable),
		}
	case "cos":
		// (cos(u))' = -sin(u) * u'
		return &Node{
			Value: "*",
			Left: &Node{
				Value: "-",
				Left:  &Node{Value: "0"},
				Right: &Node{
					Value: "sin",
					Left:  node.Left,
				},
			},
			Right: Differentiate(node.Left, variable),
		}
	case "ln":
		// (ln(u))' = (1/u) * u'
		return &Node{
			Value: "*",
			Left: &Node{
				Value: "/",
				Left:  &Node{Value: "1"},
				Right: node.Left,
			},
			Right: Differentiate(node.Left, variable),
		}
	case "exp":
		// (exp(u))' = exp(u) * u'
		return &Node{
			Value: "*",
			Left: &Node{
				Value: "exp",
				Left:  node.Left,
			},
			Right: Differentiate(node.Left, variable),
		}
	}

	return &Node{Value: "0"} // По умолчанию возвращаем 0
}
