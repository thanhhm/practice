package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

func main() {
	// Read excel file
	excelFileName := os.Args[1]
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println("Read excel file error: ", err.Error())
		return
	}

	tmpHeader := make(map[string]bool)
	var headers []string
	countryMap := make(map[string]map[string][]string)
	for _, sheet := range xlFile.Sheets {
		for rIdx, row := range sheet.Rows {
			for cIdx, cell := range row.Cells {
				if rIdx == 0 { // Read Header
					headers = append(headers, cell.String())
					tmpHeader[cell.String()] = true
				}
				if cIdx == 2 {
					continue
				}

				country := row.Cells[2].String()
				if len(countryMap[country]) == 0 {
					countryMap[country] = make(map[string][]string)
				}
				countryMap[country][headers[cIdx]] = append(countryMap[country][headers[cIdx]], cell.String())
			}
		}
	}

	// country -> header -> cell multiple value
	resultMap := make(map[string]map[string]map[string]int)
	metaHeaderMap := make(map[string]map[string]bool)

	for country, header := range countryMap {
		if len(resultMap[country]) == 0 {
			resultMap[country] = make(map[string]map[string]int)
		}

		for h, cells := range header {
			if len(metaHeaderMap[h]) == 0 {
				metaHeaderMap[h] = make(map[string]bool)
			}

			if len(resultMap[country][h]) == 0 {
				resultMap[country][h] = make(map[string]int)
			}

			// Count cell multi value
			for _, cell := range cells {
				vals := strings.Split(cell, ";")
				for _, v := range vals {
					if !tmpHeader[v] {
						resultMap[country][h][v]++

						metaHeaderMap[h][v] = false
						// fmt.Print(" ", v)
					}

				}
			}
		}
	}

	metaHeader := make(map[string][]string) // header and multi sub header
	var mainHeader []string
	for h, subHeaders := range metaHeaderMap {
		for sh := range subHeaders {
			metaHeader[h] = append(metaHeader[h], sh)
		}
		mainHeader = append(mainHeader, h)
	}

	// main header
	fmt.Print(" |")
	for _, mh := range mainHeader {
		s := mh
		subHeaders := metaHeader[mh]
		for range subHeaders {
			s += " |"
		}
		fmt.Print(s)
	}
	fmt.Println()

	// sub header
	fmt.Print(" |")
	for _, mh := range mainHeader {
		subHeaders := metaHeader[mh]
		for _, v := range subHeaders {
			fmt.Printf("%v |", v)
		}
	}
	fmt.Println()

	for country, ctHeaderMap := range resultMap { // country
		s := country

		for _, mh := range mainHeader { // meta header
			subHeaders := metaHeader[mh]
			for _, v := range subHeaders {

				ctSubHeaderMap, ok := ctHeaderMap[mh]
				if !ok {
					s += " |" + "0"
				}

				s += " |" + strconv.Itoa(ctSubHeaderMap[v])
			}
		}

		fmt.Println(s)

	}
}
