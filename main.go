package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
)

type Input struct {
	ColSums  []int      `json:colSums`
	RowSums  []int      `json:rowSums`
	Solution [][]string `json:solution`
	Table    [][]string `json:table`
}

func (input Input) HasAdjecentTree(row, col int) bool {
	var treeAbove, treeRight, treeBelow, treeLeft bool = false, false, false, false

	if row > 0 {
		treeAbove = input.Table[row-1][col] == "t"
	}

	if row < len(input.Table)-1 {
		treeBelow = input.Table[row+1][col] == "t"
	}

	if col > 0 {
		treeLeft = input.Table[row][col-1] == "t"
	}

	if col < len(input.Table[row])-1 {
		treeRight = input.Table[row][col+1] == "t"
	}

	return treeAbove || treeRight || treeBelow || treeLeft
}

func printTable(input Input) {
	colSums := []string{" "}
	rowSums := []string{}

	for _, sum := range input.ColSums {
		colSums = append(colSums, strconv.Itoa(sum))
	}

	for _, sum := range input.RowSums {
		rowSums = append(rowSums, strconv.Itoa(sum))
	}

	log.Println(strings.Join(colSums[:], "|"))
	log.Println(strings.Repeat("-", len(colSums)*2-1))

	for i, row := range input.Table {
		row = append([]string{rowSums[i]}, row...)

		log.Println(strings.Join(row[:], "|"))

		if i < len(input.Table)-1 {
			log.Println(strings.Repeat("-", len(row)*2-1))
		}
	}
}

func main() {
	var input Input

	decoder := json.NewDecoder(os.Stdin)

	err := decoder.Decode(&input)
	if err != nil {
		log.Fatal(err)
	}

	colSums := make([]int, len(input.ColSums))
	rowSums := make([]int, len(input.RowSums))

	copy(colSums, input.ColSums)
	copy(rowSums, input.RowSums)

	for row := range input.Table {
		for col := range input.Table[row] {
			if input.Table[row][col] != " " {
				continue
			}

			if input.HasAdjecentTree(row, col) {
				if rowSums[row] > 0 && colSums[col] > 0 {
					input.Table[row][col] = "T"
					rowSums[row] -= 1
					colSums[col] -= 1
				}
			}
		}
	}

	log.Println("Input:")
	printTable(input)
}
