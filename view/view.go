package view

import "fmt"
import "github.com/nsf/termbox-go"
import "strconv"
import "strings"
import "github.com/evert/csv123/model"

const columnWidth = 20
const rowBarWidth = 6

var activeCellX = 0
var activeCellY = 0

var sheetXOffset = 0
var sheetYOffset = 0

var sheet model.Sheet

func Init(sheetData model.Sheet) {

	termbox.Init()
	sheet = sheetData

}

func Close() {

	termbox.Close()

}

func Render() {

	renderLogo()
	renderInput()
	renderRows()
	renderColumns()
	renderSheet()
	renderInput()
	termbox.Flush()

}

func renderLogo() {

	mx, _ := termbox.Size()
	offset := mx - 9
	setCells(offset, 0, "csv x-x-x", termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(offset+4, 0, '1', termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
	termbox.SetCell(offset+6, 0, '2', termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault)
	termbox.SetCell(offset+8, 0, '3', termbox.ColorBlue|termbox.AttrBold, termbox.ColorDefault)

}

func renderInput() {

	mx, _ := termbox.Size()
	str := fmt.Sprintf("%v:%-"+strconv.Itoa(mx-14)+"v",
		cellCoords(),
		sheet.GetValue(activeCellX, activeCellY),
	)
	setCells(0, 0, str, termbox.ColorDefault, termbox.ColorDefault)

}

func renderRows() {

	rows := maxRows()
	for y := 0; y < rows; y++ {

		str := fmt.Sprintf("%4v ", y+1+sheetYOffset) + " "
		setCells(1, y+2, str, termbox.ColorBlack, termbox.ColorCyan)

	}

}
func renderColumns() {

	columnWidth := 20

	columns := maxColumns()

	for x := 0; x < columns; x++ {

		str := center(charForColumn(x+sheetXOffset), columnWidth)
		setCells((x*columnWidth)+rowBarWidth, 1, str, termbox.ColorBlack, termbox.ColorCyan)

	}

}

func renderSheet() {

	columns := maxColumns()
	rows := maxRows()

	for y := 0; y < rows; y++ {

		for x := 0; x < columns; x++ {

			renderCell(x, y)

		}

	}

}

func renderCell(x, y int) {

	real_x := x + sheetXOffset
	real_y := y + sheetYOffset

	active := false
	if real_x == activeCellX && real_y == activeCellY {
		active = true
	}
	val := sheet.GetValue(real_x, real_y)

	renderSheetCell(x, y, val, active)

}

func renderSheetCell(x, y int, value string, active bool) {

	if len(value) > columnWidth-2 {
		value = value[0 : columnWidth-2]
	}
	formatted_string := ""
	if _, err := strconv.Atoi(value); err == nil {
		// Right-justify
		formatted_string = " " + strings.Repeat(" ", columnWidth-2-len(value)) + value + " "
	} else {
		// Left-justify
		formatted_string = " " + value + strings.Repeat(" ", columnWidth-2-len(value)) + " "
	}

	fg := termbox.ColorDefault
	bg := termbox.ColorDefault
	if active {
		fg = termbox.ColorBlack
		bg = termbox.ColorBlue
	}
	setCells(
		(x*columnWidth)+rowBarWidth,
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

// Maximum number of columns based on terminal width
func maxColumns() int {

	mx, _ := termbox.Size()

	// How many columns can we fit?
	return (mx-rowBarWidth)/columnWidth + 1

}

// Maximum number of rows based on terminal height
func maxRows() int {

	_, my := termbox.Size()
	return my - 1

}

// Returns the active cell's coordinates as a string like A1
func cellCoords() string {

	return fmt.Sprintf("%v%v", charForColumn(activeCellX), activeCellY+1)

}

func charForColumn(i int) string {

	return string('A' + i)

}

func center(s string, width int) string {

	pad := (width - len(s)) / 2
	extra := (width - len(s)) % 2
	return strings.Repeat(" ", pad+extra) + s + strings.Repeat(" ", pad)

}

// Move relatively in the sheet
func Move(x, y int) {

	SetActiveCell(activeCellX+x, activeCellY+y)

}

// Move to start of line
func MoveHome() {
	SetActiveCell(0, activeCellY)
}

// Move to end of line
func MoveEnd() {
	SetActiveCell(sheet.GetMaxX(activeCellY), activeCellY)
}

func PageDown() {
	SetActiveCell(activeCellX, activeCellY+maxRows()-2)
}

func PageUp() {
	SetActiveCell(activeCellX, activeCellY-maxRows()+2)
}

func SetActiveCell(x, y int) {

	columns := maxColumns()
	rows := maxRows()
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	prev_x := activeCellX
	prev_y := activeCellY
	activeCellX = x
	activeCellY = y

	full_render := false

	if x < sheetXOffset {
		sheetXOffset = x
		full_render = true
	} else if x > sheetXOffset+(columns-2) {
		sheetXOffset = x - columns + 2
		full_render = true
	} else if y < sheetYOffset {
		sheetYOffset = y
		full_render = true
	} else if y > sheetYOffset+(rows-2) {
		sheetYOffset = y - rows + 2
		full_render = true
	}

	renderInput()
	if full_render {
		Render()
	} else {
		renderCell(prev_x-sheetXOffset, prev_y-sheetYOffset)
		renderCell(x-sheetXOffset, y-sheetYOffset)
		termbox.Flush()
	}

}
