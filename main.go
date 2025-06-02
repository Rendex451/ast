package main

import (
	"fmt"
)

func main() {
	// Примеры использования
	exprs := []string{
		"sin(x^2)",
		"x^8 + 3*x^4 - 5*x + 2",
		"x^2 + 2*x + 1",
		"ln(x) * cos(x)",
		"exp(x) / x",
		"x^8",
	}

	for _, expr := range exprs {
		fmt.Printf("Исходное выражение: %s\n", expr)
		tree := InfixToTree(expr)
		derivative := Differentiate(tree, "x")
		simplified := Simplify(derivative)
		fmt.Printf("Производная: %s\n", TreeToInfix(derivative))
		fmt.Printf("Упрощённая производная: %s\n\n", TreeToInfix(simplified))
		fmt.Println("========================================")
	}
}
