package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
)

var tmpl = template.Must(template.New("page.html").Funcs(template.FuncMap{
	"formatFloat": func(f float64) string {
		if math.Abs(f) < 1e-9 {
			return "0"
		}
		return strconv.FormatFloat(f, 'f', 3, 64)
	},
	"add": func(a, b int) int {
		return a + b
	},
	"abs": func(f float64) float64 {
		return math.Abs(f)
	},
}).ParseFiles("page.html"))

func handler(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		Matrix:    current,
		RowLabels: rowLabels,
		ColLabels: colLabels,
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		action := r.FormValue("action")
		if action != "random" {
			for i := range current {
				for j := range current[i] {
					name := "cell_" + strconv.Itoa(i) + "_" + strconv.Itoa(j)
					if v := r.FormValue(name); v != "" {
						v = strings.Replace(v, ",", ".", -1)
						if num, err := strconv.ParseFloat(v, 64); err == nil {
							current[i][j] = num
						}
					}
				}
			}
		}

		switch action {
		case "random":
			mainRows := len(current) - 2
			cols := len(current[0]) - 1
			resetState(mainRows, cols)
		case "addrow":
			mainRows := len(current) - 2
			cols := len(current[0]) - 1
			newRow := make([]float64, len(current[0]))
			temp := make(Matrix, len(current)+1)
			copy(temp, current[:mainRows])
			temp[mainRows] = newRow
			copy(temp[mainRows+1:], current[mainRows:])
			current = temp

			newRowLabels := make([]string, len(rowLabels)+1)
			copy(newRowLabels, rowLabels[:mainRows])
			// ТОЛЬКО x метки!
			newRowLabels[mainRows] = fmt.Sprintf("x%d", cols+mainRows+1)
			copy(newRowLabels[mainRows+1:], rowLabels[mainRows:])
			rowLabels = newRowLabels

		case "delrow":
			if len(current) > 3 {
				mainRows := len(current) - 2
				if mainRows > 1 {
					current = append(current[:mainRows-1], current[mainRows:]...)
					rowLabels = append(rowLabels[:mainRows-1], rowLabels[mainRows:]...)
				}
			}
		case "addcol":
			for i := range current {
				current[i] = append(current[i], 0)
			}
			cols := len(colLabels) - 1
			colLabels = append(colLabels, fmt.Sprintf("x%d", cols+1))

		case "delcol":
			if len(current[0]) > 2 {
				for i := range current {
					current[i] = current[i][:len(current[i])-1]
				}
				colLabels = colLabels[:len(colLabels)-1]

				// Обновляем метки строк после удаления столбца
				cols := len(colLabels) - 1
				for i := 0; i < len(rowLabels)-2; i++ {
					rowLabels[i] = fmt.Sprintf("x%d", cols+i+1)
				}
			}
		case "simplex":
			steps, result, err := solveSimplex(current, rowLabels, colLabels)
			if err != nil {
				pageData.Error = err.Error()
			} else {
				pageData.SimplexSteps = steps
				pageData.SimplexResult = result
				pageData.IsSimplex = true
			}
		case "test":
			testMatrix := Matrix{
				{15, 3, 1, 1, 0, 0},
				{91, 13, 7, 0, 1, 0},
				{15, 5, 3, 0, 0, 1},
				{0, -2, -3, 0, 0, 0},
				{0, 0, 0, 0, 0, 0},
			}
			current = testMatrix
			rowLabels = []string{"x6", "x7", "x8", "f", "g"}        // ТОЛЬКО x!
			colLabels = []string{"1", "x1", "x2", "x3", "x4", "x5"} // ТОЛЬКО x!
		}

		pageData.Matrix = current
		pageData.RowLabels = rowLabels
		pageData.ColLabels = colLabels
	}

	tmpl.Execute(w, pageData)
}
