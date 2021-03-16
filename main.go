package main

import (
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

const MaxInt = int(^uint(0) >> 1)

type DisjointSet struct {
	Collection *list.List
}

func NewDisjointSet() DisjointSet {
	ds := DisjointSet{}

	ds.Collection = list.New()

	return ds
}

func (ds DisjointSet) String() string {
	var sb strings.Builder

	sb.WriteString("\nds: [\n")
	for e := ds.Collection.Front(); e != nil; e = e.Next() {
		sb.WriteString("  ")
		sb.WriteString(listString(e.Value.(*list.List)))
		if e.Next() != nil {
			sb.WriteString(",\n")
		}
	}
	sb.WriteString("\n]")

	return sb.String()
}

func (ds *DisjointSet) MakeSet(member interface{}) *list.List {
	set := list.New()

	set.PushBack(member)

	ds.Collection.PushBack(set)

	return set
}

func listString(list *list.List) string {
	if list == nil {
		return "listString (nil)"
	}

	var sb strings.Builder

	sb.WriteString("[")

	for e := list.Front(); e != nil; e = e.Next() {
		if e != list.Front() {
			sb.WriteString(", ")
		}
		sb.WriteString(strconv.Itoa(e.Value.(int)))
	}

	sb.WriteString("]")

	return sb.String()
}

func (ds *DisjointSet) Union(x, y int) error {
	xSet := ds.FindSet(x)
	if xSet == nil {
		return fmt.Errorf("DisjointSet.Union: No x set for value: %d", x)
	}

	ySet := ds.FindSet(y)
	if ySet == nil {
		return fmt.Errorf("DisjointSet.Union: No y set for value: %d", y)
	}

	if xSet == ySet {
		// Nothing to do
		return nil
	}

	ySet.PushBackList(xSet)

	for e := ds.Collection.Front(); e != nil; e = e.Next() {
		if e.Value == xSet {
			ds.Collection.Remove(e)
		}
	}

	return nil
}

func (ds DisjointSet) FindSet(member interface{}) *list.List {
	for e0 := ds.Collection.Front(); e0 != nil; e0 = e0.Next() {
		set := e0.Value.(*list.List)

		for e1 := set.Front(); e1 != nil; e1 = e1.Next() {
			if e1.Value == member {
				return set
			}
		}
	}

	return nil
}

type Input struct {
	ColSums  []int      `json:"colSums"`
	RowSums  []int      `json:"rowSums"`
	Solution [][]string `json:"solution"`
	Table    [][]string `json:"inputTable"`
}

func (input Input) Verify() bool {
	for row := range input.Table {
		for col := range input.Table[row] {
			if input.Table[row][col] != input.Solution[row][col] {
				return false
			}
		}
	}

	return true
}

func (input Input) String() string {
	table, err := json.Marshal(input.Table)
	if err != nil {
		log.Fatalf(err.Error())
	}

	solution, err := json.Marshal(input.Solution)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return fmt.Sprintf("Input table %s\nInput solution %s\n", string(table), string(solution))
}

const (
	Grass = iota
	Tent
	Tree
	Undecided
)

type Square struct {
	Camped bool
	Col    int
	Row    int
	Kind   int
	Label  int
}

type Grid struct {
	ColSums []int
	RowSums []int

	Squares [][]Square
}

func (grid Grid) Size() int {
	return len(grid.Squares)
}

func (grid Grid) Row(r int) []Square {
	return grid.Squares[r]
}

func (grid Grid) Col(c int) []Square {
	col := make([]Square, len(grid.Squares))

	for row := range grid.Squares {
		col = append(col, grid.Squares[row][c])
	}

	return col
}

func (grid Grid) AdjecentTreeSquares(row, col int) []Square {
	var squares = make([]Square, 0)

	offsets := []struct {
		rowOffset int
		colOffset int
	}{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	}

	for i := range offsets {
		_row := row + offsets[i].rowOffset
		_col := col + offsets[i].colOffset

		if _row >= 0 && _row < len(grid.Squares) &&
			_col >= 0 && _col < len(grid.Squares[_row]) {

			if grid.Squares[_row][_col].Kind == Tree {
				squares = append(squares, grid.Squares[_row][_col])
			}
		}
	}

	return squares
}

func (grid Grid) AdjecentUncampedTreeSquares(row, col int) []Square {
	adjecentTreeSquares := grid.AdjecentTreeSquares(row, col)
	squares := make([]Square, 0)

	for i := range adjecentTreeSquares {
		if !adjecentTreeSquares[i].Camped {
			squares = append(squares, adjecentTreeSquares[i])
		}
	}

	return squares
}

func (grid Grid) AdjecentTentSquares(row, col int) []Square {
	var squares = make([]Square, 0)

	offsets := []struct {
		rowOffset int
		colOffset int
	}{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for i := range offsets {
		_row := row + offsets[i].rowOffset
		_col := col + offsets[i].colOffset

		if _row >= 0 && _row < len(grid.Squares) &&
			_col >= 0 && _col < len(grid.Squares[_row]) &&
			grid.Squares[_row][_col].Kind == Tent {
			squares = append(squares, grid.Squares[_row][_col])
		}
	}

	return squares
}

func (grid Grid) Count(kind int) ([]int, []int) {
	byRow := make([]int, len(grid.Squares))
	byCol := make([]int, len(grid.Squares))

	for row := range grid.Squares {
		for col := range grid.Squares[row] {
			if grid.Squares[row][col].Kind == kind {
				byRow[row] += 1
				byCol[col] += 1
			}
		}
	}

	return byRow, byCol
}

func (grid Grid) UndecidedLabels() ([]map[int]int, []map[int]int) {
	byRow := make([]map[int]int, len(grid.Squares))
	byCol := make([]map[int]int, len(grid.Squares))

	for row := range grid.Squares {
		for col := range grid.Squares[row] {
			if grid.Squares[row][col].Label == 0 {
				continue
			}

			// TODO It's kindo of wasteful to allocate for map that might not be used
			if byRow[row] == nil {
				byRow[row] = make(map[int]int)
			}
			if byCol[col] == nil {
				byCol[col] = make(map[int]int)
			}

			byRow[row][grid.Squares[row][col].Label] += 1
			byCol[col][grid.Squares[row][col].Label] += 1
		}
	}

	return byRow, byCol
}

func (grid Grid) Camp(row, col int) error {
	isUndecided := grid.Squares[row][col].Kind == Undecided

	if !isUndecided {
		return fmt.Errorf("Camp failed. Square at %d,%d is not undecided.", col, row)
	}

	adjecentTreeSquares := grid.AdjecentTreeSquares(row, col)
	hasAdjecentTree := len(adjecentTreeSquares) > 0

	if !hasAdjecentTree {
		return fmt.Errorf("Camp failed. Square at %d,%d does not have adjecent tree.", col, row)
	}

	hasAdjecentTent := len(grid.AdjecentTentSquares(row, col)) > 0

	if hasAdjecentTent {
		return fmt.Errorf("Camp failed. Square at %d,%d has an adjecent tent.", col, row)
	}

	grid.Squares[row][col].Kind = Tent

	return nil
}

func (grid Grid) IsValid() bool {
	for row := range grid.Squares {
		for col := range grid.Squares[row] {
			switch grid.Squares[row][col].Kind {
			case Grass:
			case Tree:
				break
			case Tent:
				if len(grid.AdjecentTreeSquares(row, col)) == 0 ||
					len(grid.AdjecentTentSquares(row, col)) != 0 {
					return false
				}
			case Undecided:
				return false
			default:
				log.Fatalf("Unknown grid value at %d,%d: %s", row, col, grid.Squares[row][col].Kind)
			}
		}
	}

	return true
}

func (grid Grid) MarkGrass() {
	tByRow, tByCol := grid.Count(Tent)

	for row := range grid.Squares {
		for col := range grid.Squares[row] {
			square := grid.Squares[row][col]

			if square.Kind == Undecided {
				// Squares where row sum is zero must be grass
				if grid.RowSums[row] == 0 {
					grid.Squares[row][col].Kind = Grass
					continue
				}

				// Squares where col sum is zero must be grass
				if grid.ColSums[col] == 0 {
					grid.Squares[row][col].Kind = Grass
					continue
				}

				// If number of tents in the row or column is already equal to the sum,
				// the square must be grass
				if grid.RowSums[row] == tByRow[row] || grid.ColSums[col] == tByCol[col] {
					grid.Squares[row][col].Kind = Grass
					continue
				}

				// If there's no adjecent uncamped trees it must be grass, since a tent
				// can only be placed horizontally or vertically adjacent to the tree
				// that doesn't already have a tent.
				if len(grid.AdjecentUncampedTreeSquares(row, col)) == 0 {
					grid.Squares[row][col].Kind = Grass
					continue
				}

				// Tents are never placed next to each other, neither vertically,
				// horizontally, or diagonally. If there's an adjecent tent to the
				// square, the square must grass.
				if len(grid.AdjecentTentSquares(row, col)) > 0 {
					grid.Squares[row][col].Kind = Grass
				}
			}
		}
	}
}

func (grid Grid) MarkTents() {
	tByRow, tByCol := grid.Count(Tent)
	uLabelByRow, uLabelByCol := grid.UndecidedLabels()

	for row := range grid.Squares {
		for col := range grid.Squares[row] {
			square := grid.Squares[row][col]
			rGroupSize := uLabelByRow[row][square.Label]
			cGroupSize := uLabelByCol[col][square.Label]

			if square.Kind != Undecided {
				continue
			}

			if grid.RowSums[row]-tByRow[row] == len(uLabelByRow[row]) &&
				rGroupSize == 1 {
				grid.Camp(row, col)
				continue
			}

			if grid.ColSums[col]-tByCol[col] == len(uLabelByCol[col]) &&
				cGroupSize == 1 {
				grid.Camp(row, col)
				continue
			}
		}
	}
}

func (grid Grid) CampTrees() {
	offsets := []struct {
		rowOffset int
		colOffset int
	}{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for row := range grid.Squares {
		for col := range grid.Squares[row] {
			if grid.Squares[row][col].Kind == Tree {
				uAroundCount := 0

				for i := range offsets {
					_row := row + offsets[i].rowOffset
					_col := col + offsets[i].colOffset

					if _row >= 0 && _row < len(grid.Squares) &&
						_col >= 0 && _col < len(grid.Squares[_row]) &&
						grid.Squares[_row][_col].Kind == Undecided {
						uAroundCount += 1
					}
				}

				if uAroundCount == 0 {
					grid.Squares[row][col].Camped = true
				}
			}

			if grid.Squares[row][col].Kind == Tent {
				adjecentUncampedTreeSquares := grid.AdjecentUncampedTreeSquares(row, col)

				// If there's only a single adjecent uncamped tree, we can mark it camped
				if len(adjecentUncampedTreeSquares) == 1 {
					treeSquare := adjecentUncampedTreeSquares[0]
					grid.Squares[treeSquare.Row][treeSquare.Col].Camped = true
				}
			}

		}
	}
}

func (grid Grid) String(showLabels bool) string {
	iColSums := make([]interface{}, 0)
	iRowSums := make([]interface{}, 0)

	iColSums = append(iColSums, "")

	for _, sum := range grid.ColSums {
		iColSums = append(iColSums, strconv.Itoa(sum))
	}

	for _, sum := range grid.RowSums {
		iRowSums = append(iRowSums, strconv.Itoa(sum))
	}

	t := table.NewWriter()

	if showLabels {
		t.AppendRow(table.Row{""})
		t.AppendSeparator()

		for row := range grid.Squares {
			tableRow := table.Row{""}

			for col := range grid.Squares[row] {
				tableRow = append(tableRow, grid.Squares[row][col].Label)
			}

			t.AppendRow(tableRow)
			t.AppendSeparator()
		}
	} else {
		t.AppendRow(iColSums)
		t.AppendSeparator()

		for row := range grid.Squares {
			tableRow := table.Row{iRowSums[row]}

			for col := range grid.Squares[row] {
				switch grid.Squares[row][col].Kind {
				case Grass:
					tableRow = append(tableRow, "🟩")
					break
				case Tent:
					tableRow = append(tableRow, "⛺")
					break
				case Tree:
					if grid.Squares[row][col].Camped {
						tableRow = append(tableRow, "🌳✓")
					} else {
						tableRow = append(tableRow, "🌳")
					}
					break
				case Undecided:
					tableRow = append(tableRow, "❔")
					break
				default:
					log.Fatalf("Unknown grid value at %d,%d: %s", row, col, grid.Squares[row][col].Kind)
				}
			}

			t.AppendRow(tableRow)
			t.AppendSeparator()
		}
	}

	return t.Render()
}

func (grid Grid) FlatString() string {
	var sb strings.Builder

	sb.WriteString("Grid solution ")

	for row := range grid.Squares {
		sb.WriteString("[")
		for col := range grid.Squares[row] {
			if col > 0 {
				sb.WriteString(",")
			}

			switch grid.Squares[row][col].Kind {
			case Grass:
				sb.WriteString("\" \"")
				break
			case Tent:
				sb.WriteString("\"T\"")
				break
			case Tree:
				sb.WriteString("\"t\"")
			default:
				log.Fatalf("Unknown grid value at %d,%d: %s", row, col, grid.Squares[row][col].Kind)
			}
		}
		sb.WriteString("],")
	}

	return sb.String()
}

func NewGrid(input Input) Grid {
	grid := Grid{}
	grid.Squares = make([][]Square, len(input.Table))

	grid.ColSums = make([]int, grid.Size())
	grid.RowSums = make([]int, grid.Size())

	copy(grid.ColSums, input.ColSums)
	copy(grid.RowSums, input.RowSums)

	for row := range input.Table {
		grid.Squares[row] = make([]Square, len(input.Table[row]))

		for col := range input.Table[row] {
			grid.Squares[row][col].Row = row
			grid.Squares[row][col].Col = col

			switch input.Table[row][col] {
			case "t":
				grid.Squares[row][col].Kind = Tree
				break
			default:
				grid.Squares[row][col].Kind = Undecided
				break
			}
		}
	}

	return grid
}

func (grid Grid) Filled() bool {
	for row := range grid.Squares {
		for col := range grid.Squares[row] {
			if grid.Squares[row][col].Kind == Undecided {
				return false
			}
		}
	}

	return true
}

func (grid Grid) CopySolution(input Input) {
	for row := range grid.Squares {
		for col := range grid.Squares[row] {
			switch grid.Squares[row][col].Kind {
			case Grass:
				input.Table[row][col] = " "
				break
			case Tent:
				input.Table[row][col] = "T"
			case Tree:
				input.Table[row][col] = "t"
				break
			}
		}
	}
}

// Connected component labeling for Undecided squares
func (grid Grid) UndecidedCCL() {
	label := 1
	ds := NewDisjointSet()

	for row := range grid.Squares {
		for col := range grid.Squares[row] {
			grid.Squares[row][col].Label = 0
		}
	}

	for row := range grid.Squares {
		for col := range grid.Squares[row] {
			priorNeighbors := make([]*Square, 0)
			currentSquare := &grid.Squares[row][col]

			if currentSquare.Kind != Undecided {
				continue
			}

			if row > 0 && grid.Squares[row-1][col].Kind == Undecided {
				priorNeighbors = append(priorNeighbors, &grid.Squares[row-1][col])
			}

			if col > 0 && grid.Squares[row][col-1].Kind == Undecided {
				priorNeighbors = append(priorNeighbors, &grid.Squares[row][col-1])
			}

			if len(priorNeighbors) == 0 {
				currentSquare.Label = label

				ds.MakeSet(label)

				label += 1
			} else {
				minLabel := MaxInt
				labels := make([]int, 0)

				for i := range priorNeighbors {
					labels = append(labels, priorNeighbors[i].Label)

					if priorNeighbors[i].Label < minLabel {
						minLabel = priorNeighbors[i].Label
					}
				}

				currentSquare.Label = minLabel

				for i := range labels {
					if labels[i] != minLabel {
						err := ds.Union(minLabel, labels[i])
						if err != nil {
							log.Println(err)
							return
						}
					}
				}
			}
		}
	}

	for row := range grid.Squares {
		for col := range grid.Squares[row] {
			if grid.Squares[row][col].Kind == Undecided {
				dSet := ds.FindSet(grid.Squares[row][col].Label)

				if dSet != nil {
					grid.Squares[row][col].Label = dSet.Front().Value.(int)
				}
			}
		}
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s [OPTION] FILE\n", os.Args[0])
		flag.PrintDefaults()
	}

	verbose := flag.Bool("v", false, "Enable verbose output")

	flag.Parse()

	if !*verbose {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	// ---------------------------------------------------------------------------

	var input Input
	var decoder *json.Decoder
	var grid Grid

	args := flag.Args()

	if len(args) > 0 {
		inputFile, err := os.Open(args[0])
		if err != nil {
			log.Println("Failed to open input file", args)
			log.Fatal(err)
		}

		decoder = json.NewDecoder(inputFile)
	} else {
		decoder = json.NewDecoder(os.Stdin)
	}

	decoder.DisallowUnknownFields()

	err := decoder.Decode(&input)
	if err != nil {
		log.Println("Failed to decode input")
		log.Fatal(err)
	}

	grid = NewGrid(input)

	// Make a copy of the solution
	input.Solution = make([][]string, len(input.Table))
	for row := range input.Table {
		input.Solution[row] = make([]string, len(input.Table[row]))
		copy(input.Solution[row], input.Table[row])
	}

	// Clear out solution from input
	for row := range input.Table {
		for col := range input.Table[row] {
			if input.Table[row][col] == "T" {
				input.Table[row][col] = " "
			}
		}
	}

	log.Printf("Input grid:\n%s\n\n", grid.String(false))

	step := 0

	for {
		grid.MarkGrass()
		grid.UndecidedCCL()
		grid.MarkTents()
		grid.CampTrees()

		log.Printf("Step %d:\n%s\n\n", step, grid.String(false))

		step += 1

		if grid.Filled() {
			break
		}
	}

	grid.CopySolution(input)

	if input.Verify() {
		log.Println("Solution is correct! 👍")
	} else {
		log.Println("Solution is wrong! 👎")
	}
}
