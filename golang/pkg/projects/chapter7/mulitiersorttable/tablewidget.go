package mulitiersorttable

import (
	"fmt"
	"sort"
)

type TableWidget struct {
	TableRows            []*tableRow // Slice of pointers to []string
	TableHeaders         []*tableHeader
	currentSortParameter *headersClicksAndIndex
}

/*
Init(data [][]string) error initializes a table and returns any encountered errors otherwise returns nil
*/
func (tw *TableWidget) Init(data [][]string) error {
	// Whether data is empty
	if len(data) == 0 {
		return fmt.Errorf("creating table with %v; no data", data)
	}
	tw.TableRows = make([]*tableRow, len(data)-1)

	// Are there columns in data
	var width int
	if width = tableWidth(data); width == 0 {
		return fmt.Errorf("creating table with %v; the number of columns is 0", data)
	}
	tw.TableHeaders = make([]*tableHeader, tableWidth(data))

	// Set headers
	for i := 0; i < len(tw.TableHeaders); i++ {
		var headersTitle string
		var empty bool
		if i < len(data[0]) {
			headersTitle = data[0][i]
			empty = false
		} else {
			headersTitle = ""
			empty = true
		}

		tw.TableHeaders[i] = &tableHeader{Title: headersTitle, clickedCount: 0, isEmpty: empty}
	}

	// Set rows
	for i := 1; i < len(data); i++ {
		newRow := make([]string, len(data[i]))
		copy(newRow, data[i])
		newTableRow := tableRow(newRow)

		tw.TableRows[i-1] = &newTableRow
	}

	return nil
}

/*
(tw *TableWidget) String() returns string representation of the table.
*/
func (tw *TableWidget) String() string {
	var table string
	for _, header := range tw.TableHeaders {
		table += header.Title + " "
	}
	table += "\n"

	for _, row := range tw.TableRows {
		for _, columnValue := range *row {
			table += columnValue + " "
		}
		table += "\n"
	}

	return table
}

/*
tableWidth(data [][]string) int returns columns number
*/
func tableWidth(data [][]string) int {
	var width int
	for _, row := range data {
		if len(row) > width {
			width = len(row)
		}
	}
	return width
}

/*
setSortParameter(column, clicks int) error allows user to set clicks on the chosen header
*/
func (tw *TableWidget) setSortParameter(columnIndex int) error {
	if tw.TableHeaders[columnIndex].Title != "" {
		tw.currentSortParameter = &headersClicksAndIndex{clicksCount: 1, index: columnIndex}
		return nil
	}

	return fmt.Errorf("setting sort parameter; title of column №%d is empty", columnIndex)
}

/*
Sort() sorts the table in the following order: the first parameter is most "clicked" header, the next parameter is
the second most "clicked" header and so on
*/
func (tw *TableWidget) Sort() {
	sort.Sort(tw)
}

func (tw *TableWidget) Len() int {
	return len(tw.TableRows)
}

func (tw *TableWidget) Swap(i, j int) {
	tw.TableRows[i], tw.TableRows[j] = tw.TableRows[j], tw.TableRows[i]
}

func (tw *TableWidget) Less(i, j int) bool {
	if (*tw.TableRows[i])[tw.currentSortParameter.index] != (*tw.TableRows[j])[tw.currentSortParameter.index] {
		return (*tw.TableRows[i])[tw.currentSortParameter.index] < (*tw.TableRows[j])[tw.currentSortParameter.index]
	}

	return false
}

/*
(tw *TableWidget) sortClicks() Pairs
1) Combines info about header's clicks and its index for each header
2) Puts these pairs into a slice
3) Sorts the slice in the DESC order
4) Returns the sorted slice.
*/
func (tw *TableWidget) sortClicks() Pairs {
	pairs := make(Pairs, 0, len(tw.TableHeaders))
	var pairsCount int
	// Put clicksCount/headerIndex pairs into a slice
	for i, header := range tw.TableHeaders {
		if !header.isEmpty {
			pairs = append(pairs, &headersClicksAndIndex{clicksCount: header.clickedCount, index: i})
			pairsCount++
		}
	}

	// Sort the slice of pairs in the DESC order
	sort.Sort(sort.Reverse(&pairs))

	return pairs
}
