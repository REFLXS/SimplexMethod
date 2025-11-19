package main

import (
	"fmt"
	"math"
	"strings"
)

const eps = 1e-9

func solveSimplex(a Matrix, rLabels, cLabels []string) ([]SimplexStep, *SimplexResult, error) {
	steps := []SimplexStep{}

	workingMatrix := cloneMatrix(a)
	workingRowLabels := cloneSlice(rLabels)
	workingColLabels := cloneSlice(cLabels)

	fIndex, gIndex := -1, -1
	for i, lab := range workingRowLabels {
		if lab == "f" {
			fIndex = i
		} else if lab == "g" {
			gIndex = i
		}
	}
	if fIndex == -1 {
		return nil, nil, fmt.Errorf("не найдена целевая функция f")
	}
	if gIndex == -1 {
		workingMatrix = append(workingMatrix, make([]float64, len(workingMatrix[0])))
		workingRowLabels = append(workingRowLabels, "g")
		gIndex = len(workingRowLabels) - 1
	}

	for i := 0; i < fIndex; i++ {
		if workingMatrix[i][0] < -eps {
			for j := 0; j < len(workingMatrix[i]); j++ {
				workingMatrix[i][j] = -workingMatrix[i][j]
			}
		}
	}

	steps = append(steps, SimplexStep{
		Desc:      "Начальная симплекс-таблица (после нормализации RHS)",
		Matrix:    cloneMatrix(workingMatrix),
		RowLabels: cloneSlice(workingRowLabels),
		ColLabels: cloneSlice(workingColLabels),
	})

	numCols := len(workingMatrix[0])
	numRows := len(workingMatrix)
	artificialCols := []int{}
	artificialIsBasic := map[int]int{}

	for i := 0; i < fIndex; i++ {
		found := -1
		for j := 1; j < numCols; j++ {
			if math.Abs(workingMatrix[i][j]-1) < eps {
				unit := true
				for ii := 0; ii < fIndex; ii++ {
					if ii == i {
						continue
					}
					if math.Abs(workingMatrix[ii][j]) > eps {
						unit = false
						break
					}
				}
				if unit {
					found = j
					break
				}
			}
		}
		if found != -1 {
			workingRowLabels[i] = workingColLabels[found]
			continue
		}
		for rr := 0; rr < numRows; rr++ {
			if rr == i {
				workingMatrix[rr] = append(workingMatrix[rr], 1.0)
			} else {
				workingMatrix[rr] = append(workingMatrix[rr], 0.0)
			}
		}
		colName := fmt.Sprintf("a%d", i+1)
		workingColLabels = append(workingColLabels, colName)
		newColIndex := len(workingColLabels) - 1
		artificialCols = append(artificialCols, newColIndex)
		artificialIsBasic[newColIndex] = i
		workingRowLabels[i] = colName
		numCols++
	}

	for j := range workingMatrix[gIndex] {
		workingMatrix[gIndex][j] = 0
	}
	for _, col := range artificialCols {
		if row, ok := artificialIsBasic[col]; ok {
			for j := 0; j < len(workingMatrix[row]); j++ {
				workingMatrix[gIndex][j] -= workingMatrix[row][j]
			}
		}
	}

	steps = append(steps, SimplexStep{
		Desc:      "Фаза I: добавлены искусственные переменные; построена g = сумма(искусственных)",
		Matrix:    cloneMatrix(workingMatrix),
		RowLabels: cloneSlice(workingRowLabels),
		ColLabels: cloneSlice(workingColLabels),
	})

	iter := 0
	maxIter := 500
	for iter < maxIter {
		iter++
		pivotCol := findPivotColumnPhaseI(workingMatrix, gIndex)
		if pivotCol == -1 {
			break
		}
		pivotRow := findPivotRowForPhaseI(workingMatrix, pivotCol, fIndex)
		if pivotRow == -1 {
			return steps, nil, fmt.Errorf("вспомогательная задача неограничена (фаза I)")
		}

		desc := fmt.Sprintf("Фаза I - Итерация %d. Разрешающий элемент: %.6f [%s][%s]",
			iter, workingMatrix[pivotRow][pivotCol], workingRowLabels[pivotRow], workingColLabels[pivotCol])

		workingMatrix = pivot(workingMatrix, pivotRow, pivotCol)
		workingRowLabels[pivotRow] = workingColLabels[pivotCol]
		if strings.HasPrefix(workingColLabels[pivotCol], "a") {
			artificialIsBasic[pivotCol] = pivotRow
		} else {
			for col, r := range artificialIsBasic {
				if r == pivotRow && col != pivotCol {
					delete(artificialIsBasic, col)
				}
			}
		}

		steps = append(steps, SimplexStep{
			Desc:      desc,
			Matrix:    cloneMatrix(workingMatrix),
			RowLabels: cloneSlice(workingRowLabels),
			ColLabels: cloneSlice(workingColLabels),
			PivotRow:  pivotRow,
			PivotCol:  pivotCol,
		})
	}

	if math.Abs(workingMatrix[gIndex][0]) > 1e-7 {
		steps = append(steps, SimplexStep{
			Desc:      "Фаза I завершена: несовместна (g != 0)",
			Matrix:    cloneMatrix(workingMatrix),
			RowLabels: cloneSlice(workingRowLabels),
			ColLabels: cloneSlice(workingColLabels),
		})
		return steps, nil, fmt.Errorf("задача несовместна (вспомогательная функция g = %.6f)", workingMatrix[gIndex][0])
	}

	steps = append(steps, SimplexStep{
		Desc:      "Фаза I завершена: допустимое базисное решение найдено (g = 0). Удаление искусственных переменных.",
		Matrix:    cloneMatrix(workingMatrix),
		RowLabels: cloneSlice(workingRowLabels),
		ColLabels: cloneSlice(workingColLabels),
	})

	artSet := make(map[int]bool)
	for _, col := range artificialCols {
		artSet[col] = true
	}

	for col := range artSet {
		row := -1
		for i := 0; i < fIndex; i++ {
			if math.Abs(workingMatrix[i][col]-1) < 1e-8 {
				unit := true
				for ii := 0; ii < fIndex; ii++ {
					if ii == i {
						continue
					}
					if math.Abs(workingMatrix[ii][col]) > 1e-8 {
						unit = false
						break
					}
				}
				if unit {
					row = i
					break
				}
			}
		}
		if row == -1 {
			continue
		}

		replacement := -1
		for j := 1; j < len(workingMatrix[0]); j++ {
			if artSet[j] {
				continue
			}
			if math.Abs(workingMatrix[row][j]) > eps {
				replacement = j
				break
			}
		}
		if replacement != -1 {
			workingMatrix = pivot(workingMatrix, row, replacement)
			workingRowLabels[row] = workingColLabels[replacement]

		} else {

			for ii := 0; ii < len(workingMatrix); ii++ {
				workingMatrix[ii][col] = 0
			}
			workingRowLabels[row] = ""
		}
	}

	newCols := []string{}
	colMap := make([]int, len(workingColLabels))
	newIndex := 0
	for j := 0; j < len(workingColLabels); j++ {
		if artSet[j] {
			colMap[j] = -1
			continue
		}
		colMap[j] = newIndex
		newCols = append(newCols, workingColLabels[j])
		newIndex++
	}

	newMatrix := make(Matrix, len(workingMatrix))
	for i := 0; i < len(workingMatrix); i++ {
		newRow := make([]float64, len(newCols))
		for j := 0; j < len(workingMatrix[i]); j++ {
			if colMap[j] == -1 {
				continue
			}
			newRow[colMap[j]] = workingMatrix[i][j]
		}
		newMatrix[i] = newRow
	}
	workingMatrix = newMatrix
	workingColLabels = newCols

	for i := 0; i < fIndex; i++ {
		if workingRowLabels[i] == "" {
			workingRowLabels[i] = fmt.Sprintf("y%d", i+1)
		}
	}

	steps = append(steps, SimplexStep{
		Desc:      "После удаления искусственных столбцов",
		Matrix:    cloneMatrix(workingMatrix),
		RowLabels: cloneSlice(workingRowLabels),
		ColLabels: cloneSlice(workingColLabels),
	})

	iter = 0
	for iter < maxIter {
		iter++
		pivotCol := findPivotColumnPhaseII(workingMatrix, fIndex)
		if pivotCol == -1 {
			solution, value := extractSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex)
			steps = append(steps, SimplexStep{
				Desc:      fmt.Sprintf("Оптимум найден (Фаза II) на итерации %d", iter),
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
				Desc:      "Задача неограничена (Фаза II)",
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

		desc := fmt.Sprintf("Фаза II - Итерация %d. Разрешающий элемент: %.6f [%s][%s]",
			iter, workingMatrix[pivotRow][pivotCol], workingRowLabels[pivotRow], workingColLabels[pivotCol])

		workingMatrix = pivot(workingMatrix, pivotRow, pivotCol)
		workingRowLabels[pivotRow] = workingColLabels[pivotCol]
		workingColLabels[pivotCol] = ""

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

func findPivotRowForPhaseI(a Matrix, pivotCol int, fIndex int) int {
	minRatio := math.MaxFloat64
	pivotRow := -1
	for i := 0; i < fIndex; i++ {
		coeff := a[i][pivotCol]
		if coeff <= eps {
			continue
		}
		ratio := a[i][0] / coeff
		if ratio >= -eps && ratio < minRatio {
			minRatio = ratio
			pivotRow = i
		}
	}
	return pivotRow
}

func findPivotColumnPhaseI(a Matrix, gIndex int) int {
	pivotCol := -1
	var mostNeg float64 = -eps
	for j := 1; j < len(a[0]); j++ {
		if a[gIndex][j] < mostNeg {
			mostNeg = a[gIndex][j]
			pivotCol = j
		}
	}
	return pivotCol
}

func findPivotColumnPhaseII(a Matrix, fIndex int) int {
	pivotCol := -1
	var mostNeg float64 = -eps
	for j := 1; j < len(a[0]); j++ {
		if a[fIndex][j] < mostNeg {
			mostNeg = a[fIndex][j]
			pivotCol = j
		}
	}
	return pivotCol
}

func findPivotRow(a Matrix, pivotCol int, excludeRow int) int {
	minRatio := math.MaxFloat64
	pivotRow := -1
	limit := excludeRow
	if limit <= 0 || limit > len(a) {
		limit = len(a)
	}
	for i := 0; i < limit; i++ {
		coeff := a[i][pivotCol]
		if coeff <= eps {
			continue
		}
		ratio := a[i][0] / coeff
		if ratio >= -eps && ratio < minRatio {
			minRatio = ratio
			pivotRow = i
		}
	}
	return pivotRow
}

func pivot(a Matrix, pivotRow, pivotCol int) Matrix {
	result := cloneMatrix(a)
	pivotElement := a[pivotRow][pivotCol]
	if math.Abs(pivotElement) < 1e-15 {
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

func extractSolution(a Matrix, rLabels, cLabels []string, fIndex int) ([]float64, float64) {
	numVars := len(cLabels) - 1
	if numVars < 0 {
		numVars = 0
	}
	solution := make([]float64, numVars)
	value := 0.0
	if fIndex >= 0 && fIndex < len(a) {
		value = -a[fIndex][0]
	}
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
	return solution, value
}
