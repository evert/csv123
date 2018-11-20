package main

import "github.com/nsf/termbox-go"
import "fmt"
import "strings"
import "encoding/csv"
import "os"
import "io"
import "strconv"

const column_width = 20
const row_bar_width = 6

var active_cell_x = 0
var active_cell_y = 0

var sheet_x_offset = 0
var sheet_y_offset = 0

type SheetData = [][]string

var sheet SheetData

func main() {

	if len(os.Args) < 2 {
		usage()
	}

	sheet = make(SheetData, 0)
	read_file()

	termbox.Init()

	setCells(0, 0, "csv x-x-x", termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(4, 0, '1', termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
	termbox.SetCell(6, 0, '2', termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault)
	termbox.SetCell(8, 0, '3', termbox.ColorBlue|termbox.AttrBold, termbox.ColorDefault)

	defer termbox.Close()

	render()

mainloop:
	for {

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowLeft:
				set_active_cell(active_cell_x-1, active_cell_y)
			case termbox.KeyArrowRight:
				set_active_cell(active_cell_x+1, active_cell_y)
			case termbox.KeyArrowUp:
				set_active_cell(active_cell_x, active_cell_y-1)
			case termbox.KeyArrowDown:
				set_active_cell(active_cell_x, active_cell_y+1)
			case termbox.KeyHome, termbox.KeyCtrlA:
				set_active_cell(0, active_cell_y)
			case termbox.KeyEnd, termbox.KeyCtrlE:
				set_active_cell(len(sheet[0]), active_cell_y)
			case termbox.KeyPgup:
				set_active_cell(active_cell_x, active_cell_y-max_rows()+2)
			case termbox.KeyPgdn:
				set_active_cell(active_cell_x, active_cell_y+max_rows()-2)
			}
		case termbox.EventResize:
			render()
		}

	}

}

func render() {

	render_rows()
	render_columns()
	render_sheet()
	termbox.Flush()

}

func max_columns() int {

	mx, _ := termbox.Size()

	// How many columns can we fit?
	return (mx-row_bar_width)/column_width + 1

}
func max_rows() int {

	_, my := termbox.Size()
	return my - 1

}

func render_rows() {

	rows := max_rows()
	for y := 0; y < rows; y++ {

		str := fmt.Sprintf("%4v ", y+1+sheet_y_offset) + " "
		setCells(1, y+2, str, termbox.ColorBlack, termbox.ColorCyan)

	}

}
func render_columns() {

	column_width := 20

	columns := max_columns()

	for x := 0; x < columns; x++ {

		str := center(char_for_column(x+sheet_x_offset), column_width)
		setCells((x*column_width)+row_bar_width, 1, str, termbox.ColorBlack, termbox.ColorCyan)

	}

}

func render_sheet() {

	columns := max_columns()
	rows := max_rows()

	for y := 0; y < rows; y++ {

		for x := 0; x < columns; x++ {

			render_cell(x, y)

		}

	}

}

func render_cell(x, y int) {

	real_x := x + sheet_x_offset
	real_y := y + sheet_y_offset

	active := false
	if real_x == active_cell_x && real_y == active_cell_y {
		active = true
	}
	val := ""
	if real_y < len(sheet) && real_x < len(sheet[real_y]) {
		val = sheet[real_y][real_x]
	}

	render_sheet_cell(x, y, val, active)

}

func render_sheet_cell(x, y int, value string, active bool) {

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

	fg := termbox.ColorDefault
	bg := termbox.ColorDefault
	if active {
		fg = termbox.ColorBlack
		bg = termbox.ColorBlue
	}
	setCells(
		(x*column_width)+row_bar_width,
		y+2,
		formatted_string,
		fg,
		bg,
	)

}

func setCells(x, y int, str string, fg, bg termbox.Attribute) {

	for off := 0; off < len(str); off++ {

		termbox.SetCell(x+off, y, rune(str[off]), fg, bg)

	}

}

func set_active_cell(x, y int) {

	columns := max_columns()
	rows := max_rows()
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	prev_x := active_cell_x
	prev_y := active_cell_y
	active_cell_x = x
	active_cell_y = y

	full_render := false

	if x < sheet_x_offset {
		sheet_x_offset = x
		full_render = true
	} else if x > sheet_x_offset+(columns-2) {
		sheet_x_offset = x - columns + 2
		full_render = true
	} else if y < sheet_y_offset {
		sheet_y_offset = y
		full_render = true
	} else if y > sheet_y_offset+(rows-2) {
		sheet_y_offset = y - rows + 2
		full_render = true
	}

	if full_render {
		render()
	} else {
		render_cell(prev_x-sheet_x_offset, prev_y-sheet_y_offset)
		render_cell(x-sheet_x_offset, y-sheet_y_offset)
		termbox.Flush()
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
