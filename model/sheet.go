package model

import "encoding/csv"
import "os"
import "io"
import "bufio"

type Sheet = [][]string

func ReadFromFile(fileName string) Sheet {

	var sheet Sheet
	sheet = make(Sheet, 0)

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
		sheet = append(sheet, record)
	}
	return sheet

}
