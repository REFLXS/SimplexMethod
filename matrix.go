package main

import (
	"fmt"
	"math/rand"
)

func init() {
	resetState(3, 2)
}

func resetState(rows, cols int) {
	mat := make(Matrix, rows+2)

	for i := range mat {
		mat[i] = make([]float64, cols+1)

		for j := range mat[i] {
			if i < rows {
				if j == 0 {
					mat[i][j] = float64(rand.Intn(15) + 5)
				} else {
					mat[i][j] = float64(rand.Intn(8) + 1)
				}
			} else if i == rows {
				if j == 0 {
					mat[i][j] = 0
				} else {
					mat[i][j] = -float64(rand.Intn(10) + 1)
				}
			} else {
				mat[i][j] = 0
			}
		}
	}

	current = mat

	// ТОЛЬКО x метки!
	rowLabels = make([]string, rows+2)
	for i := 0; i < rows; i++ {
		// Базисные переменные начинаются с x_{cols+1}
		rowLabels[i] = fmt.Sprintf("x%d", cols+i+1)
	}
	rowLabels[rows] = "f"
	rowLabels[rows+1] = "g"

	colLabels = make([]string, cols+1)
	colLabels[0] = "1"
	for i := 1; i <= cols; i++ {
		colLabels[i] = fmt.Sprintf("x%d", i)
	}
}

func cloneMatrix(a Matrix) Matrix {
	copyMat := make(Matrix, len(a))
	for i := range a {
		copyMat[i] = append([]float64{}, a[i]...)
	}
	return copyMat
}

func cloneSlice(s []string) []string {
	return append([]string{}, s...)
}
