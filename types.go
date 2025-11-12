package main

type Matrix [][]float64

type Step struct {
	Desc      string
	Matrix    Matrix
	RowLabels []string
	ColLabels []string
	PivotRow  int
	PivotCol  int
}

type PageData struct {
	Matrix         Matrix
	RowLabels      []string
	ColLabels      []string
	Steps          []Step
	Error          string
	SolutionVector []string
	FinalMatrix    Matrix
	FinalRowLabels []string
	FinalColLabels []string
	IsSimplex      bool
	SimplexSteps   []SimplexStep
	SimplexResult  *SimplexResult
}

type SimplexStep struct {
	Desc      string
	Matrix    Matrix
	RowLabels []string
	ColLabels []string
	PivotRow  int
	PivotCol  int
	Solution  []float64
	Value     float64
}

type SimplexResult struct {
	Solution []float64
	Value    float64
	Status   string
}

var (
	current   Matrix
	rowLabels []string
	colLabels []string
)
