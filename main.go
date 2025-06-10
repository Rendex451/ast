package main

import (
	"fmt"
)

func main() {
	// Тестовые выражения
	testExpressions := []string{
		"sin(x^2)",
		"x^8 + 3*x^4 - 5*x + 2",
		"x^2 + 2*x + 1",
		"ln(x) * cos(x)",
		"exp(x) / x",
		"x^8",
	}

	// Вывод результатов
	fmt.Println("Expression Analysis Results:")
	fmt.Println("===========================")
	for _, expr := range testExpressions {
		result := AnalyzeExpression(expr)
		fmt.Printf("Expression: %s\n", result["Expression"])
		fmt.Printf("Tree: %s\n", result["Tree"])
		fmt.Printf("Derivative: %s\n", result["Derivative"])
		fmt.Printf("Simplified Derivative: %s\n", result["Simplified Derivative"])
		fmt.Println("---------------------------")
	}
}
