package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

//Console displays a UI in the terminal
func Console() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	filesViewX0, filesViewY0 := 0, 0
	filesViewX1, filesViewY1 := int(float32(maxX)*0.25), int(float32(maxY)*0.75)
	filesView, err := g.SetView("files", filesViewX0, filesViewY0, filesViewX1, filesViewY1)

	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	fmt.Fprintln(filesView, "Files!")

	logViewX0, logViewY0 := (filesViewX0 + filesViewX1 + 2), filesViewY0
	logViewX1, logViewY1 := logViewX0+int(float32(maxX)*0.4), filesViewY1
	logView, err := g.SetView("log", logViewX0, logViewY0, logViewX1, logViewY1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	fmt.Fprintln(logView, "Log view!")

	settingsViewX0, settingsViewY0 := logViewX1+2, logViewY0
	settingsViewX1, settingsViewY1 := settingsViewX0+int(float32(maxX)*0.25), filesViewY1
	settingsView, err := g.SetView("settings", settingsViewX0, settingsViewY0, settingsViewX1, settingsViewY1)

	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	fmt.Fprintln(settingsView, "Settings!")
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
