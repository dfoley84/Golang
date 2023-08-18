package main

import (
	"fmt"
	"log"
	"strconv"
	"github.com/xuri/excelize/v2"
)

func main() {
	xlFile, err := excelize.OpenFile("./<DIR>/<File>.xlsx")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	
	rows, err := xlFile.GetRows(xlFile.GetSheetName(0))
	if err != nil {
		log.Fatalf("Error getting rows: %v", err)
	}

	for i := 0; i < len(rows); i += 5000 {
		newFile := excelize.NewFile()
		sheetName := newFile.GetSheetName(0)
		for j, row := range rows[i:min(i+5000, len(rows))] {
			for k, colCell := range row {
				newFile.SetCellValue(sheetName, fmt.Sprintf("%s%d", toAlphaString(k), j+1), colCell)
			}
		}

		if err := newFile.SaveAs(fmt.Sprintf("./<DIR>/<File>%s.xlsx", strconv.Itoa(i))); err != nil {
			log.Fatalf("Failed to save file: %v", err)
		}
	}
}

func toAlphaString(i int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if i < 26 {
		return string(letters[i])
	}
	return toAlphaString(i/26-1) + string(letters[i%26])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
