package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/mw491/godo/internal/model"
	"github.com/objectbox/objectbox-go/objectbox"
	gc "github.com/rthornton128/goncurses"
)

var Todos []*model.Todo

var regular_pair int16 = 1
var sel_pair int16 = 2
var sel_todo = 0

func initObjectBox() *objectbox.ObjectBox {
	objectBox, err := objectbox.NewBuilder().Model(model.ObjectBoxModel()).Build()
	if err != nil {
		panic(err)
	}
	return objectBox
}

func refreshTodos(box *model.TodoBox) {
	var err error
	Todos, err = box.GetAll()
	if err != nil {
		panic(err)
	}
}

func newTodo(stdscr *gc.Window, box *model.TodoBox) {
	var text []string
	stdscr.Move(len(Todos)+3, 5)
	gc.Cursor(1)

	refresh := func(key gc.Key, addtext bool) {
		stdscr.Clear()
		printScreen(stdscr)
		stdscr.Move(len(Todos)+3, 5)
		if addtext {
			text = append(text, string(key))
		}
		for _, v := range text {
			stdscr.Print(v)
		}
		stdscr.Refresh()
	}
	exit := func() {
		gc.Cursor(0)
		text = nil
		refresh(0, false)
		sel_pair = 2
	}

	for {
		key := stdscr.GetChar()
		switch key {
		case 27:
			exit()
			return
		case gc.KEY_BACKSPACE, '\b', '\x7f':
			if len(text) > 0 {
				text = text[:len(text)-1]
				refresh(key, false)
			}
		case gc.KEY_RETURN:
			todo := &model.Todo{Task: strings.Join(text, ""), Done: false}
			_, err := box.Put(todo)
			if err != nil {
				panic(err)
			}
			refreshTodos(box)
			exit()
			sel_todo = len(Todos) - 1
			return
		default:
			refresh(key, true)
		}
	}

}

func printScreen(stdscr *gc.Window) {
	calcPos := func(msg string) (row, col int) {
		row, col = stdscr.MaxYX()
		row, col = (row/2)-1, (col-len(msg))/2
		return row, col
	}
	msg := "TodoList"
	_, col := calcPos(msg)
	stdscr.MovePrint(2, col, msg)
	if len(Todos) == 0 {
		msg := "You are FREE!! Hit O to add a new todo."
		row, col := calcPos(msg)
		stdscr.MovePrint(row, col, msg)
	}
	for i, todo := range Todos {
		var color_pair int16
		if i == sel_todo {
			color_pair = sel_pair
		} else {
			color_pair = regular_pair
		}

		checkbox := func() string {
			if todo.Done {
				return "[X]"
			} else {
				return "[ ]"
			}
		}()

		stdscr.ColorOn(color_pair)
		stdscr.MovePrint(i+3, 5, fmt.Sprintf("%s %s", checkbox, todo.Task))
		stdscr.ColorOff(color_pair)
	}
	stdscr.Refresh()
}

func main() {
	stdscr, err := gc.Init()
	if err != nil {
		log.Fatal("init", err)
	}
	defer gc.End()

	gc.Echo(false)
	gc.CBreak(true)
	gc.Cursor(0)
	stdscr.Keypad(true)

	if !gc.HasColors() {
		log.Fatal("Example requires a colour capable terminal")
	}

	if err := gc.StartColor(); err != nil {
		log.Fatal(err)
	}

	if err := gc.InitPair(regular_pair, gc.C_WHITE, gc.C_BLACK); err != nil {
		log.Fatal("InitPair failed: ", err)
	}
	gc.InitPair(sel_pair, gc.C_BLACK, gc.C_WHITE)

	// INITIALISE DATABASE
	ob := initObjectBox()
	defer ob.Close()
	box := model.BoxForTodo(ob)

	refreshTodos(box)

	quit := false

	for !quit {
		printScreen(stdscr)

		key := stdscr.GetChar()

		switch key {
		case 'q':
			quit = true
		case 'k':
			if sel_todo > 0 {
				sel_todo -= 1
			}
		case 'j':
			if sel_todo+1 < len(Todos) {
				sel_todo += 1
			}
		case 'd':
			todo := Todos[sel_todo]
			box.Remove(todo)
			refreshTodos(box)
			stdscr.Clear()
			if sel_todo == len(Todos) {
				sel_todo = len(Todos) - 1
			}
			printScreen(stdscr)
		case 'o':
			sel_pair = regular_pair
			newTodo(stdscr, box)
		case ' ':
			todo := Todos[sel_todo]
			todo.Done = !todo.Done
			box.Put(todo)
		}

	}
}
