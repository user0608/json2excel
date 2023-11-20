package json2excel

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/xuri/excelize/v2"
)

var ErrNoDataForExcel = errors.New("no data available to generate the Excel")
var ErrNilRequestProvided = errors.New("no valid request has been provided")

type JSON2ExcelConverter interface {
	Excel(req *request) (io.Reader, error)
	SaveExcel(req *request, filePath string) error
}

func NewJSON2ExcelConverter() JSON2ExcelConverter {
	return &j2x{}
}

type j2x struct{}

func (*j2x) removeNumbers(input string) string {
	re := regexp.MustCompile("[0-9]")
	result := re.ReplaceAllString(input, "")
	return result
}

func (x *j2x) WriteHeader(file *excelize.File, sheetName string, headers []string) (
	startcell string, endcell string, err error) {

	for i, header := range headers {
		col := fmt.Sprintf("%c", 'A'+i)
		cell := fmt.Sprintf("%s1", col)
		if err = file.SetCellValue(sheetName, cell, header); err != nil {
			return
		}
		if i == 0 {
			startcell = cell
		}
		endcell = cell
	}

	var startCol = x.removeNumbers(startcell)
	var endCol = x.removeNumbers(endcell)
	if err := file.SetColWidth(sheetName, startCol, endCol, 30); err != nil {
		return startcell, endcell, err
	}
	style, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	if err != nil {
		return startcell, endcell, err
	}
	if err = file.SetCellStyle(sheetName, startcell, endcell, style); err != nil {
		return startcell, endcell, err
	}
	return
}

func (*j2x) WriteData(file *excelize.File, sheetName string, data [][]any) (startcell string, endcell string, err error) {
	for i, rowData := range data {
		rowIndex := i + 2
		for j, cellData := range rowData {
			cell := fmt.Sprintf("%c%d", 'A'+j, rowIndex)
			if i == 0 && j == 0 {
				startcell = cell
			}
			if err = file.SetCellStr(sheetName, cell, fmt.Sprint(cellData)); err != nil {
				return
			}
			endcell = cell
		}
	}
	return
}

func (jx *j2x) Excel(req *request) (io.Reader, error) {
	if req == nil {
		return nil, ErrNilRequestProvided
	}
	if len(req.Data.Columns) == 0 {
		return nil, ErrNoDataForExcel
	}
	file := excelize.NewFile()
	defer file.Close()
	var sheetName = "Sheet1"
	_, err := file.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	startCell, _, err := jx.WriteHeader(file, sheetName, req.Data.Columns)
	if err != nil {
		return nil, err
	}
	_, endCell, err := jx.WriteData(file, sheetName, req.Data.RowsValues)
	if err != nil {
		return nil, err
	}
	if err := file.AddTable(sheetName, &excelize.Table{
		Range:     fmt.Sprintf("%s:%s", startCell, endCell),
		Name:      "table",
		StyleName: "TableStyleMedium2",
	}); err != nil {
		return nil, err
	}
	var buff bytes.Buffer
	if err := file.Write(&buff); err != nil {
		return nil, err
	}
	return &buff, nil
}

func (jx *j2x) SaveExcel(req *request, nameFile string) error {
	if req == nil {
		return ErrNilRequestProvided
	}
	excel, err := jx.Excel(req)
	if err != nil {
		return err
	}
	file, err := os.Create(nameFile)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := io.Copy(file, excel); err != nil {
		return err
	}
	return nil
}
