package main

import (
	"fmt"
	"log"
	"os"
	"snake/snake"

	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorRed

	g.SetManagerFunc(layout)

	pg := snake.NewPlayGround()
	snakeObj, e := snake.InitSpirit(g, pg)
	if e != nil {
		panic(e)
	}

	food := snake.NewFood(g, pg)
	food.Appear()

	gEvent := snake.NewGameEvent(snakeObj, pg, food)

	go gEvent.ListenEvent()
	go snakeObj.Move(gEvent.Move)

	if err := initKeybindings(g, snakeObj); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, _ := g.Size()
	v, err := g.SetView("help", maxX-25, 0, maxX-1, 9)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v, "← ↑ → ↓: Move")
		fmt.Fprintln(v, "^C: Exit")
	}
	return nil
}

func initKeybindings(g *gocui.Gui, s *snake.Snake) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			s.MoveLeft()
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			s.MoveRight()
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			s.MoveDown()
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			s.MoveUp()
			return nil
		}); err != nil {
		return err
	}

	return nil
}
