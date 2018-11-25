package main

import "os"
import "fmt"
import "github.com/evert/csv123/view"
import "github.com/evert/csv123/model"
import "github.com/nsf/termbox-go"

func main() {

	if len(os.Args) < 2 {
		usage()
	}

	var sheet = model.ReadFromFile(os.Args[1])

	view.Init(sheet)
	defer view.Close()

	view.Render()

mainloop:
	for {

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowLeft:
				view.Move(-1, 0)
			case termbox.KeyArrowRight:
				view.Move(1, 0)
			case termbox.KeyArrowUp:
				view.Move(0, -1)
			case termbox.KeyArrowDown:
				view.Move(0, +1)
			case termbox.KeyHome, termbox.KeyCtrlA:
				view.MoveHome()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				view.MoveEnd()
			case termbox.KeyPgup:
				view.PageUp()
			case termbox.KeyPgdn:
				view.PageDown()
			}
		case termbox.EventResize:
			view.Render()
		}

	}

}

func usage() {

	fmt.Fprintf(os.Stderr, "usage: %s [inputfile]\n", os.Args[0])
	os.Exit(2)

}
