package main

import (
	"fmt"
	"math"
)

const eps = 1e-9

type ConstraintType int

const (
	LessEqual ConstraintType = iota
	GreaterEqual
	Equal
)

func solveSimplex(a Matrix, rLabels, cLabels []string) ([]SimplexStep, *SimplexResult, error) {
	steps := []SimplexStep{}

	// Создаем рабочую копию
	workingMatrix := cloneMatrix(a)
	workingRowLabels := cloneSlice(rLabels)
	workingColLabels := cloneSlice(cLabels)

<<<<<<< Updated upstream
	// Находим индексы f и g
	fIndex, gIndex := -1, -1
=======
	fIndex := -1
>>>>>>> Stashed changes
	for i, lab := range workingRowLabels {
		if lab == "f" {
			fIndex = i
			break
		}
	}
	if fIndex == -1 {
		return nil, nil, fmt.Errorf("не найдена целевая функция f")
	}
<<<<<<< Updated upstream

	numConstraints := fIndex

	// Анализируем ограничения и определяем нужны ли искусственные переменные
	constraintTypes := make([]ConstraintType, numConstraints)
	needsArtificial := false
	artificialVars := []int{}

	// Определяем типы ограничений
	for i := 0; i < numConstraints; i++ {
		// Проверяем знак свободного члена
		if workingMatrix[i][0] < -eps {
			// Отрицательный свободный член - нужно искусственная переменная
			constraintTypes[i] = GreaterEqual
			needsArtificial = true
		} else {
			// Положительный свободный член - предполагаем ≤
			constraintTypes[i] = LessEqual
		}
	}

	// Добавляем недостающие переменные и искусственные переменные если нужно
	if needsArtificial {
		// Добавляем искусственные переменные
		for i := 0; i < numConstraints; i++ {
			if constraintTypes[i] == GreaterEqual || workingMatrix[i][0] < -eps {
				// Добавляем искусственную переменную
				for r := 0; r < len(workingMatrix); r++ {
					if r == i {
						workingMatrix[r] = append(workingMatrix[r], 1.0) // Искусственная переменная
					} else if r == gIndex {
						workingMatrix[r] = append(workingMatrix[r], -1.0) // Включаем в вспомогательную функцию
					} else {
						workingMatrix[r] = append(workingMatrix[r], 0.0)
					}
				}
				artificialCol := len(workingColLabels)
				artificialVars = append(artificialVars, artificialCol)
				workingColLabels = append(workingColLabels, fmt.Sprintf("a%d", len(artificialVars)))
				workingRowLabels[i] = workingColLabels[artificialCol]
			}
		}

		// Пересчитываем вспомогательную функцию g
		for j := range workingMatrix[gIndex] {
			workingMatrix[gIndex][j] = 0
		}
		for _, col := range artificialVars {
			workingMatrix[gIndex][col] = -1
		}
	}

	steps = append(steps, SimplexStep{
		Desc:      "Начальная таблица",
=======

	steps = append(steps, SimplexStep{
		Desc:      "Шаг 1. Начальная симплекс-таблица",
>>>>>>> Stashed changes
		Matrix:    cloneMatrix(workingMatrix),
		RowLabels: cloneSlice(workingRowLabels),
		ColLabels: cloneSlice(workingColLabels),
		Solution:  extractCurrentSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex),
		Value:     -workingMatrix[fIndex][0],
	})

<<<<<<< Updated upstream
	// ФАЗА I: Решение вспомогательной задачи если есть искусственные переменные
	if len(artificialVars) > 0 {
		steps = append(steps, SimplexStep{
			Desc:      "Фаза I: Начало решения вспомогательной задачи",
			Matrix:    cloneMatrix(workingMatrix),
			RowLabels: cloneSlice(workingRowLabels),
			ColLabels: cloneSlice(workingColLabels),
		})

		// Решаем вспомогательную задачу минимизации g
		iter := 0
		maxIter := 100
		for iter < maxIter {
			iter++

			// Ищем разрешающий столбец (максимальный положительный в строке g)
			pivotCol := -1
			maxCoeff := -math.MaxFloat64
			for j := 1; j < len(workingMatrix[gIndex]); j++ {
				// Исключаем искусственные переменные из выбора
				isArtificial := false
				for _, ac := range artificialVars {
					if j == ac {
						isArtificial = true
						break
					}
				}
				if !isArtificial && workingMatrix[gIndex][j] > maxCoeff {
					maxCoeff = workingMatrix[gIndex][j]
					pivotCol = j
				}
			}

			if pivotCol == -1 || maxCoeff < eps {
				break // Оптимум достигнут
			}

			// Ищем разрешающую строку
			pivotRow := -1
			minRatio := math.MaxFloat64
			for i := 0; i < fIndex; i++ {
				if workingMatrix[i][pivotCol] > eps {
					ratio := workingMatrix[i][0] / workingMatrix[i][pivotCol]
					if ratio >= 0 && ratio < minRatio {
						minRatio = ratio
						pivotRow = i
					}
				}
			}

			if pivotRow == -1 {
				return steps, nil, fmt.Errorf("вспомогательная задача неограничена")
			}

			// Выполняем поворот
			workingMatrix = pivot(workingMatrix, pivotRow, pivotCol)
			workingRowLabels[pivotRow] = workingColLabels[pivotCol]

			steps = append(steps, SimplexStep{
				Desc:      fmt.Sprintf("Фаза I - Итерация %d", iter),
				Matrix:    cloneMatrix(workingMatrix),
				RowLabels: cloneSlice(workingRowLabels),
				ColLabels: cloneSlice(workingColLabels),
				PivotRow:  pivotRow,
				PivotCol:  pivotCol,
			})
		}

		// Проверяем результат фазы I
		if math.Abs(workingMatrix[gIndex][0]) > eps {
			return steps, nil, fmt.Errorf("задача несовместна, невозможно найти начальное допустимое решение")
		}

		// Удаляем искусственные переменные из базиса
		for i := 0; i < fIndex; i++ {
			if strings.HasPrefix(workingRowLabels[i], "a") {
				// Ищем неискусственную переменную для ввода в базис
				found := false
				for j := 1; j < len(workingColLabels); j++ {
					isArtificial := false
					for _, ac := range artificialVars {
						if j == ac {
							isArtificial = true
							break
						}
					}
					if !isArtificial && math.Abs(workingMatrix[i][j]) > eps {
						workingMatrix = pivot(workingMatrix, i, j)
						workingRowLabels[i] = workingColLabels[j]
						found = true
						break
					}
				}
				if !found {
					// Если не нашли замену, оставляем искусственную переменную (вырожденный случай)
					continue
				}
			}
		}
	}

	// Удаляем столбцы искусственных переменных и строку g
	if len(artificialVars) > 0 {
		// Создаем новую матрицу без искусственных переменных
		newCols := []string{}
		colMap := make(map[int]int)
		newIdx := 0

		for j := 0; j < len(workingColLabels); j++ {
			isArtificial := false
			for _, ac := range artificialVars {
				if j == ac {
					isArtificial = true
					break
				}
			}
			if !isArtificial {
				colMap[j] = newIdx
				newCols = append(newCols, workingColLabels[j])
				newIdx++
			}
		}

		newMatrix := make(Matrix, len(workingMatrix)-1) // Удаляем строку g
		for i := 0; i < len(workingMatrix)-1; i++ {
			newRow := make([]float64, len(newCols))
			for j := range workingMatrix[i] {
				if newJ, exists := colMap[j]; exists {
					newRow[newJ] = workingMatrix[i][j]
				}
			}
			newMatrix[i] = newRow
		}

		workingMatrix = newMatrix
		workingColLabels = newCols
		workingRowLabels = workingRowLabels[:len(workingRowLabels)-1] // Удаляем "g"
		fIndex = len(workingMatrix) - 1                               // f теперь последняя строка
	}

	steps = append(steps, SimplexStep{
		Desc:      "Начало Фазы II",
		Matrix:    cloneMatrix(workingMatrix),
		RowLabels: cloneSlice(workingRowLabels),
		ColLabels: cloneSlice(workingColLabels),
	})

	// ФАЗА II: Решение основной задачи
	iter := 0
	maxIter := 100
	for iter < maxIter {
		iter++

		// Ищем разрешающий столбец (минимальный отрицательный коэффициент в f)
		pivotCol := -1
		minVal := 0.0
		for j := 1; j < len(workingMatrix[fIndex]); j++ {
			if workingMatrix[fIndex][j] < minVal-eps {
				minVal = workingMatrix[fIndex][j]
				pivotCol = j
			}
		}

		if pivotCol == -1 {
			// Оптимум найден
			solution, value := extractSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex)
			steps = append(steps, SimplexStep{
				Desc:      "Оптимум найден",
=======
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
>>>>>>> Stashed changes
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

<<<<<<< Updated upstream
		// Ищем разрешающую строку
		pivotRow := -1
		minRatio := math.MaxFloat64
		for i := 0; i < fIndex; i++ {
			if workingMatrix[i][pivotCol] > eps {
				ratio := workingMatrix[i][0] / workingMatrix[i][pivotCol]
				if ratio >= -eps && ratio < minRatio {
					minRatio = ratio
					pivotRow = i
				}
=======
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
>>>>>>> Stashed changes
			}
		}

		if pivotRow == -1 {
<<<<<<< Updated upstream
			solution, value := extractSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex)
=======
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
>>>>>>> Stashed changes
			return steps, &SimplexResult{
				Solution: solution,
				Value:    value,
				Status:   "Функция неограничена сверху",
			}, nil
		}

<<<<<<< Updated upstream
		// Выполняем поворот
		workingMatrix = pivot(workingMatrix, pivotRow, pivotCol)
		workingRowLabels[pivotRow] = workingColLabels[pivotCol]
=======
		desc := fmt.Sprintf("Шаг %d. Разрешающий элемент: %.3f [%s][%s]",
			stepNumber, workingMatrix[pivotRow][pivotCol],
			workingRowLabels[pivotRow], workingColLabels[pivotCol])

		workingMatrix = pivot(workingMatrix, pivotRow, pivotCol)
		workingRowLabels[pivotRow] = workingColLabels[pivotCol]
		solution := extractCurrentSolution(workingMatrix, workingRowLabels, workingColLabels, fIndex)
		value := -workingMatrix[fIndex][0]
>>>>>>> Stashed changes

		steps = append(steps, SimplexStep{
			Desc:      fmt.Sprintf("Фаза II - Итерация %d", iter),
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

<<<<<<< Updated upstream
func extractSolution(a Matrix, rLabels, cLabels []string, fIndex int) ([]float64, float64) {
	numVars := len(cLabels) - 1
	solution := make([]float64, numVars)

	// Все переменные по умолчанию 0
	for i := range solution {
		solution[i] = 0
	}

	// Заполняем базисные переменные
	for i := 0; i < fIndex; i++ {
		label := rLabels[i]
		if strings.HasPrefix(label, "x") {
			varIndex := 0
			fmt.Sscanf(label, "x%d", &varIndex)
			if varIndex > 0 && varIndex <= numVars {
				solution[varIndex-1] = a[i][0]
=======
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
>>>>>>> Stashed changes
			}
		}
	}

<<<<<<< Updated upstream
	// Значение целевой функции
	value := a[fIndex][0]

	return solution, value
=======
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
>>>>>>> Stashed changes
}

func pivot(a Matrix, pivotRow, pivotCol int) Matrix {
	result := cloneMatrix(a)
	pivotElement := a[pivotRow][pivotCol]

	// Нормализуем разрешающую строку
	for j := range result[pivotRow] {
		result[pivotRow][j] /= pivotElement
	}

	// Обновляем остальные строки
	for i := range result {
		if i == pivotRow {
			continue
		}
		factor := a[i][pivotCol]
		for j := range result[i] {
			result[i][j] -= factor * result[pivotRow][j]
			if math.Abs(result[i][j]) < 1e-12 {
				result[i][j] = 0
			}
		}
	}

	return result
}
