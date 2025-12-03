package main

import (
	"fmt"
	"math"
)

const eps = 1e-9

func solveSimplex(a Matrix, rLabels, cLabels []string) ([]SimplexStep, *SimplexResult, error) {
	steps := []SimplexStep{}

	workingMatrix := cloneMatrix(a)
	workingRowLabels := cloneSlice(rLabels)
	workingColLabels := cloneSlice(cLabels)

	fIndex := -1
	for i, lab := range workingRowLabels {
		if lab == "f" {
			fIndex = i
			break
		}
	}
	if fIndex == -1 {
		return nil, nil, fmt.Errorf("не найдена целевая функция f")
	}

	steps = append(steps, SimplexStep{
		Desc:      "Шаг 1. Начальная симплекс-таблица",
		Matrix:    cloneMatrix(workingMatrix),
		RowLabels: cloneSlice(workingRowLabels),
		ColLabels: cloneSlice(workingColLabels),
		Solution:  extractCurrentSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex),
		Value:     -workingMatrix[fIndex][0],
	})

	stepNumber := 2
	maxIter := 50

	for iter := 0; iter < maxIter; iter++ {
		allNonNegative := true
		for j := 1; j < len(workingMatrix[0]); j++ {
			if workingMatrix[fIndex][j] < -eps {
				allNonNegative = false
				break
			}
		}

		if allNonNegative {
			solution := extractCurrentSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex)
			value := -workingMatrix[fIndex][0]

			steps = append(steps, SimplexStep{
				Desc:      fmt.Sprintf("Шаг %d. Оптимум найден", stepNumber),
				Matrix:    cloneMatrix(workingMatrix),
				RowLabels: cloneSlice(workingRowLabels),
				ColLabels: cloneSlice(workingColLabels),
				Solution:  solution,
				Value:     value,
			})
			return steps, &SimplexResult{
				Solution: solution,
				Value:    value,
				Status:   "Оптимум найден",
			}, nil
		}

		pivotCol := -1
		minValue := 0.0

		for j := 1; j < len(workingMatrix[0]); j++ {
			coeff := workingMatrix[fIndex][j]
			if coeff < minValue {
				minValue = coeff
				pivotCol = j
			}
		}

		if pivotCol == -1 {
			break
		}

		pivotRow := -1
		minRatio := math.MaxFloat64

		for i := 0; i < fIndex; i++ {
			if workingMatrix[i][pivotCol] <= eps {
				continue
			}
			ratio := workingMatrix[i][0] / workingMatrix[i][pivotCol]
			if ratio >= -eps && ratio < minRatio {
				minRatio = ratio
				pivotRow = i
			}
		}

		if pivotRow == -1 {
			solution := extractCurrentSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex)
			value := -workingMatrix[fIndex][0]

			steps = append(steps, SimplexStep{
				Desc:      fmt.Sprintf("Шаг %d. Задача неограничена", stepNumber),
				Matrix:    cloneMatrix(workingMatrix),
				RowLabels: cloneSlice(workingRowLabels),
				ColLabels: cloneSlice(workingColLabels),
				Solution:  solution,
				Value:     value,
			})
			return steps, &SimplexResult{
				Solution: solution,
				Value:    value,
				Status:   "Задача неограничена",
			}, nil
		}

		desc := fmt.Sprintf("Шаг %d. Разрешающий элемент: %.3f [%s][%s]",
			stepNumber, workingMatrix[pivotRow][pivotCol],
			workingRowLabels[pivotRow], workingColLabels[pivotCol])

		workingMatrix = pivot(workingMatrix, pivotRow, pivotCol)
		workingRowLabels[pivotRow] = workingColLabels[pivotCol]
		solution := extractCurrentSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex)
		value := -workingMatrix[fIndex][0]

		steps = append(steps, SimplexStep{
			Desc:      desc,
			Matrix:    cloneMatrix(workingMatrix),
			RowLabels: cloneSlice(workingRowLabels),
			ColLabels: cloneSlice(workingColLabels),
			PivotRow:  pivotRow,
			PivotCol:  pivotCol,
			Solution:  solution,
			Value:     value,
		})

		stepNumber++
	}

	solution := extractCurrentSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex)
	value := -workingMatrix[fIndex][0]

	return steps, &SimplexResult{
		Solution: solution,
		Value:    value,
		Status:   "Достигнут лимит итераций",
	}, nil
}

func pivot(a Matrix, pivotRow, pivotCol int) Matrix {
	result := cloneMatrix(a)
	pivotElement := a[pivotRow][pivotCol]

	if math.Abs(pivotElement) < eps {
		return result
	}

	for j := 0; j < len(a[0]); j++ {
		result[pivotRow][j] = a[pivotRow][j] / pivotElement
	}

	for i := 0; i < len(a); i++ {
		if i == pivotRow {
			continue
		}
		factor := a[i][pivotCol]
		for j := 0; j < len(a[0]); j++ {
			val := a[i][j] - factor*result[pivotRow][j]
			if math.Abs(val) < 1e-12 {
				val = 0
			}
			result[i][j] = val
		}
	}

	return result
}

func extractCurrentSolution(a Matrix, rLabels, cLabels []string, fIndex int) []float64 {
	numMainVars := 0
	for _, label := range cLabels {
		if len(label) > 0 && label[0] == 'x' {
			var num int
			if _, err := fmt.Sscanf(label, "x%d", &num); err == nil {
				if num > numMainVars {
					numMainVars = num
				}
			}
		}
	}

	solution := make([]float64, numMainVars)

	// Заполняем значения основных переменных
	for i := 0; i < fIndex; i++ {
		label := rLabels[i]
		if len(label) > 0 && label[0] == 'x' {
			var varIndex int
			if _, err := fmt.Sscanf(label, "x%d", &varIndex); err == nil {
				if varIndex > 0 && varIndex <= numMainVars {
					solution[varIndex-1] = a[i][0]
				}
			}
		}
	}

	return solution
}
