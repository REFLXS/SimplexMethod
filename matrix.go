package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	resetState(3, 3)
}

func resetState(rows, cols int) {
	mat := make(Matrix, rows+2) // +2 для строк f и g

	for i := range mat {
		mat[i] = make([]float64, cols+1) // +1 для правой части

		for j := range mat[i] {
			if i < rows {
				if j == 0 {
					mat[i][j] = float64(rand.Intn(15) + 5)
				} else {
					mat[i][j] = float64(rand.Intn(8) + 1)
				}
			} else if i == rows {
				if j == 0 {
					mat[i][j] = 0 // Свободный член целевой функции
				} else {
					mat[i][j] = -float64(rand.Intn(10) + 1)
				}
			} else {
				mat[i][j] = 0
			}
		}
	}

	for i := 0; i < rows; i++ {
		if i < cols {
			mat[i][i+1] = 1
		}
	}

	current = mat
	rowLabels = make([]string, rows+2)
	for i := 0; i < rows; i++ {
		rowLabels[i] = fmt.Sprintf("y%d", i+1)
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
