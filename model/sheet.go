package model

import "encoding/csv"
import "os"
import "io"
import "bufio"

type Sheet struct {
	data [][]string
}

func NewSheet() Sheet {

	s := Sheet{data: make([][]string, 0)}
	return s

}

func (s Sheet) GetValue(x, y int) string {

	if y < len(s.data) && x < len(s.data[y]) {
		return s.data[y][x]
	} else {
		return ""
	}

}

func (s Sheet) GetMaxX(y int) int {

	return len(s.data[y])

}

func ReadFromFile(fileName string) Sheet {

	s := NewSheet()

	var reader io.Reader
	if os.Args[1] == "-" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		r, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		reader = r

	}

	csv_reader := csv.NewReader(reader)

	for {
		record, err := csv_reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		s.data = append(s.data, record)
	}
	return s

}
