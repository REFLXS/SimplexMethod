package main

import (
	"fmt"
	"math"
)

func solveSimplex(a Matrix, rLabels, cLabels []string) ([]SimplexStep, *SimplexResult, error) {
	steps := []SimplexStep{}

	workingMatrix := cloneMatrix(a)
	workingRowLabels := cloneSlice(rLabels)
	workingColLabels := cloneSlice(cLabels)

	fIndex, gIndex := -1, -1
	for i, label := range workingRowLabels {
		if label == "f" {
			fIndex = i
		} else if label == "g" {
			gIndex = i
		}
	}

	if fIndex == -1 {
		return nil, nil, fmt.Errorf("не найдена целевая функция f")
	}

	steps = append(steps, SimplexStep{
		Desc:      "Начальная симплекс-таблица",
		Matrix:    cloneMatrix(workingMatrix),
		RowLabels: cloneSlice(workingRowLabels),
		ColLabels: cloneSlice(workingColLabels),
	})

	iteration := 0
	maxIterations := 20

	if gIndex != -1 && !isRowZero(workingMatrix[gIndex]) {
		steps = append(steps, SimplexStep{
			Desc:      "Фаза I: Решение вспомогательной задачи",
			Matrix:    cloneMatrix(workingMatrix),
			RowLabels: cloneSlice(workingRowLabels),
			ColLabels: cloneSlice(workingColLabels),
		})

		for iteration < maxIterations {
			iteration++

			pivotCol := findPivotColumnPhase1(workingMatrix, gIndex)
			if pivotCol == -1 {
				break // Все коэффициенты неотрицательны
			}

			pivotRow := findPivotRow(workingMatrix, pivotCol, gIndex)
			if pivotRow == -1 {
				return steps, nil, fmt.Errorf("задача неограничена в фазе I")
			}

			desc := fmt.Sprintf("Фаза I - Итерация %d. Разрешающий элемент: %.3f[%s][%s]",
				iteration, workingMatrix[pivotRow][pivotCol], workingRowLabels[pivotRow], workingColLabels[pivotCol])

			workingMatrix = pivot(workingMatrix, pivotRow, pivotCol)
			workingRowLabels[pivotRow], workingColLabels[pivotCol] = workingColLabels[pivotCol], workingRowLabels[pivotRow]

			steps = append(steps, SimplexStep{
				Desc:      desc,
				Matrix:    cloneMatrix(workingMatrix),
				RowLabels: cloneSlice(workingRowLabels),
				ColLabels: cloneSlice(workingColLabels),
				PivotRow:  pivotRow,
				PivotCol:  pivotCol,
			})
		}
	}

	iteration = 0
	for iteration < maxIterations {
		iteration++

		pivotCol := findPivotColumnPhase2(workingMatrix, fIndex)
		if pivotCol == -1 {
			solution, value := extractSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex)
			steps = append(steps, SimplexStep{
				Desc:      fmt.Sprintf("Фаза II - Оптимум найден на итерации %d", iteration),
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

		pivotRow := findPivotRow(workingMatrix, pivotCol, fIndex)
		if pivotRow == -1 {
			solution, value := extractSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex)
			steps = append(steps, SimplexStep{
				Desc:      "Задача неограничена",
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

		desc := fmt.Sprintf("Фаза II - Итерация %d. Разрешающий элемент: %.3f[%s][%s]",
			iteration, workingMatrix[pivotRow][pivotCol], workingRowLabels[pivotRow], workingColLabels[pivotCol])

		workingMatrix = pivot(workingMatrix, pivotRow, pivotCol)

		workingRowLabels[pivotRow], workingColLabels[pivotCol] = workingColLabels[pivotCol], workingRowLabels[pivotRow]

		solution, value := extractSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex)

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
	}

	solution, value := extractSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex)
	return steps, &SimplexResult{
		Solution: solution,
		Value:    value,
		Status:   "Достигнут лимит итераций",
	}, nil
}

func isRowZero(row []float64) bool {
	for _, val := range row {
		if math.Abs(val) > 1e-9 {
			return false
		}
	}
	return true
}

func findPivotColumnPhase1(a Matrix, gIndex int) int {
	maxNeg := 0.0
	pivotCol := -1

	for j := 1; j < len(a[0]); j++ {
		if a[gIndex][j] < -1e-9 && math.Abs(a[gIndex][j]) > math.Abs(maxNeg) {
			maxNeg = a[gIndex][j]
			pivotCol = j
		}
	}

	return pivotCol
}

func findPivotColumnPhase2(a Matrix, fIndex int) int {
	minVal := 0.0
	pivotCol := -1

	for j := 1; j < len(a[0]); j++ {
		if a[fIndex][j] < minVal-1e-9 {
			minVal = a[fIndex][j]
			pivotCol = j
		}
	}

	return pivotCol
}

func findPivotRow(a Matrix, pivotCol int, excludeRow int) int {
	minRatio := math.MaxFloat64
	pivotRow := -1

	for i := 0; i < len(a); i++ {
		if i == excludeRow || a[i][pivotCol] <= 1e-9 {
			continue
		}

		ratio := a[i][0] / a[i][pivotCol]
		if ratio >= -1e-9 && ratio < minRatio-1e-9 {
			minRatio = ratio
			pivotRow = i
		}
	}

	return pivotRow
}

func pivot(a Matrix, pivotRow, pivotCol int) Matrix {
	result := cloneMatrix(a)
	pivotElement := a[pivotRow][pivotCol]

	for j := 0; j < len(a[0]); j++ {
		result[pivotRow][j] = a[pivotRow][j] / pivotElement
	}

	for i := 0; i < len(a); i++ {
		if i == pivotRow {
			continue
		}
		factor := a[i][pivotCol]
		for j := 0; j < len(a[0]); j++ {
			result[i][j] = a[i][j] - factor*result[pivotRow][j]

			if math.Abs(result[i][j]) < 1e-10 {
				result[i][j] = 0
			}
		}
	}

	return result
}

func extractSolution(a Matrix, rLabels, cLabels []string, fIndex int) ([]float64, float64) {
	numVars := len(cLabels) - 1
	solution := make([]float64, numVars)

	value := -a[fIndex][0]

	for i := 0; i < len(a); i++ {
		if i == fIndex {
			continue
		}

		if len(rLabels[i]) > 0 && rLabels[i][0] == 'x' {
			varIndex := 0
			fmt.Sscanf(rLabels[i], "x%d", &varIndex)
			if varIndex > 0 && varIndex <= numVars {
				solution[varIndex-1] = a[i][0]
			}
		}
	}

	for j := 1; j < len(cLabels); j++ {
		if len(cLabels[j]) > 0 && cLabels[j][0] == 'x' {
			varIndex := 0
			fmt.Sscanf(cLabels[j], "x%d", &varIndex)
			if varIndex > 0 && varIndex <= numVars {

			}
		}
	}

	return solution, value
}
