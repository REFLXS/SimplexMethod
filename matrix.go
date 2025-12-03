package main

import (
	"fmt"
	"math/rand"
)

func init() {
	resetState(3, 2)
}

func resetState(rows, cols int) {
	if rows == 3 && cols == 2 {
		current = Matrix{
			{15, 3, 1, 1, 0, 0},    // x3: 3x1 + x2 + x3 = 15
			{91, 13, 7, 0, 1, 0},   // x4: 13x1 + 7x2 + x4 = 91
			{-15, -5, -3, 0, 0, 1}, // x5: -5x1 - 3x2 + x5 = -15
			{0, -2, -3, 0, 0, 0},   // f = -2x1 - 3x2 → min
		}
		rowLabels = []string{"x3", "x4", "x5", "f"}
		colLabels = []string{"1", "x1", "x2", "x3", "x4", "x5"}
		return
	}

	// Стандартная случайная генерация
	mat := make(Matrix, rows+1)

	for i := range mat {
		mat[i] = make([]float64, cols+1+rows)

		for j := range mat[i] {
			if i < rows {
				if j == 0 {
					mat[i][j] = float64(rand.Intn(15) + 5)
				} else if j <= cols {
					mat[i][j] = float64(rand.Intn(8) + 1)
				} else if j == cols+1+i {
					mat[i][j] = 1.0
				} else {
					mat[i][j] = 0
				}
			} else {
				if j == 0 {
					mat[i][j] = 0
				} else if j <= cols {
					mat[i][j] = -float64(rand.Intn(10) + 1)
				} else {
					mat[i][j] = 0
				}
			}
		}
	}

	current = mat

	rowLabels = make([]string, rows+1)
	for i := 0; i < rows; i++ {
		rowLabels[i] = fmt.Sprintf("x%d", cols+1+i)
	}
	rowLabels[rows] = "f"

	colLabels = make([]string, cols+1+rows)
	colLabels[0] = "1"
	for i := 1; i <= cols; i++ {
		colLabels[i] = fmt.Sprintf("x%d", i)
	}
	for i := 1; i <= rows; i++ {
		colLabels[cols+i] = fmt.Sprintf("x%d", cols+i)
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
