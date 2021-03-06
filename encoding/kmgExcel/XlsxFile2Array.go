package kmgExcel

import (
	"strings"

	"github.com/tealeg/xlsx"
	//"io"
	//"archive/zip"
)

// get all raw data from excel
// output index mean=> sheetIndex ,row ,cell ,value
// not remove any cells
func XlsxFile2Array(path string) ([][][]string, error) {
	file, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, err
	}
	return xlsx2ArrayXlsxFile(file)
}

//get raw data from one sheet of excel
//output index mean=> row ,cell ,value
// not remove any cells
func XlsxFileSheetIndex2Array(path string, index int) ([][]string, error) {
	file, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, err
	}
	output := [][]string{}
	sheet := file.Sheets[index]
	for _, row := range sheet.Rows {
		if row == nil {
			continue
		}
		r := []string{}
		for _, cell := range row.Cells {
			r = append(r, cell.String())
		}
		output = append(output, r)
	}
	return output, nil
}

//get data from first sheet of excel
//output index mean=> row ,cell ,value
// remove all right and bottom blank cells
func XlsxFileFirstSheet2ArrayTrim(path string) (output [][]string, err error) {
	output, err = XlsxFileSheetIndex2Array(path, 0)
	if err != nil {
		return
	}
	output = Trim2DArray(output)
	return
}

//get data from first column of first sheet of excel file
//output index name=>row index,value
//remove any blank cells.
// output will not include blank cell in middle.
func XlsxFileFirstSheet2FirstColumnTrim(path string) (output []string, err error) {
	outputArray, err := XlsxFileFirstSheet2ArrayTrim(path)
	if err != nil {
		return
	}
	output = make([]string, len(outputArray))
	i := 0
	for _, row := range outputArray {
		v := strings.TrimSpace(row[0])
		if v == "" {
			continue
		}
		output[i] = v
		i++
	}
	output = output[:i]
	return
}

func xlsx2ArrayXlsxFile(file *xlsx.File) (output [][][]string, err error) {
	output = [][][]string{}
	for _, sheet := range file.Sheets {
		s := [][]string{}
		for _, row := range sheet.Rows {
			if row == nil {
				continue
			}
			r := []string{}
			for _, cell := range row.Cells {
				r = append(r, cell.String())
			}
			s = append(s, r)
		}
		output = append(output, s)
	}
	return output, nil
}
