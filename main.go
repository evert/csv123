package main

import "github.com/nsf/termbox-go"
import "time"
import "fmt"
import "strings"
import "encoding/csv"
import "os"
import "io"
import "strconv"

const column_width = 20
const row_bar_width = 5

type SheetData = [][]string

var sheet SheetData

func main() {

	if len(os.Args) < 2 {
		usage()
	}

	sheet = make(SheetData, 0)
	read_file()

	termbox.Init()

	defer termbox.Close()
	render()

	termbox.Sync()

	time.Sleep(time.Second * 2)

}

func render() {

	render_rows()
	render_columns()
	render_sheet()

}

func max_columns() int {

	mx, _ := termbox.Size()

	// How many columns can we fit?
	return (mx - row_bar_width) / column_width

}
func max_rows() int {

	_, my := termbox.Size()
	return my - 1

}

func render_rows() {

	termbox.SetCell(1, 1, '@', termbox.ColorBlue, termbox.ColorGreen)
	rows := max_rows()
	for y := 0; y < rows; y++ {

		str := fmt.Sprintf("%4v ", y+1)
		setCells(1, y+1, str, termbox.ColorBlack, termbox.ColorCyan)

	}

}
func render_columns() {

	column_width := 20

	columns := max_columns()

	for x := 0; x < columns; x++ {

		str := center(char_for_column(x), column_width)
		setCells((x*column_width)+row_bar_width, 1, str, termbox.ColorBlack, termbox.ColorCyan)

	}

}

func render_sheet() {

	columns := max_columns()
	rows := max_rows()

	for y := 0; y < rows; y++ {

		if y >= len(sheet) {
			break
		}

		for x := 0; x < columns; x++ {
			if x >= len(sheet[y]) {
				break
			}

			render_sheet_cell(x, y, sheet[y][x])

		}

	}

}

func render_sheet_cell(x, y int, value string) {

	if len(value) > column_width-2 {
		value = value[0 : column_width-2]
	}
	formatted_string := ""
	if _, err := strconv.Atoi(value); err == nil {
		// Right-justify
		formatted_string = " " + strings.Repeat(" ", column_width-2-len(value)) + value + " "
	} else {
		// Left-justify
		formatted_string = " " + value + strings.Repeat(" ", column_width-2-len(value)) + " "
	}
	setCells(
		(x*column_width)+row_bar_width,
		y+2,
		formatted_string,
		termbox.ColorDefault,
		termbox.ColorDefault,
	)

}

func setCells(x, y int, str string, fg, bg termbox.Attribute) {

	for off := 0; off < len(str); off++ {

		termbox.SetCell(x+off, y, rune(str[off]), fg, bg)

	}

}

func read_file() {

	reader, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
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

}

func char_for_column(i int) string {

	return string('A' + i)

}

func center(s string, width int) string {

	pad := (width - len(s)) / 2
	extra := (width - len(s)) % 2
	return strings.Repeat(" ", pad+extra) + s + strings.Repeat(" ", pad)

}

func usage() {

	fmt.Fprintf(os.Stderr, "usage: %s [inputfile]\n", os.Args[0])
	os.Exit(2)

}
