package main

import (
	"flag"
	"log"
	"os"

	"github.com/xconstruct/stark/core"
	slog "github.com/xconstruct/stark/log"
	"github.com/xconstruct/stark/proto"
	"gopkg.in/qml.v1"
)

var verbose = flag.Bool("v", false, "verbose debug output")

type App struct {
	ctx     *core.Context
	client  *proto.Client
	history string
	window  *qml.Window
}

func main() {
	flag.Parse()

	app := &App{}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func (app *App) Run() error {
	go app.initProto()
	if err := qml.Run(app.runGui); err != nil {
		return err
	}
	return nil
}

func (app *App) runGui() error {
	engine := qml.NewEngine()
	engine.On("quit", func() {
		os.Exit(0)
	})

	root, err := engine.LoadString("qml.go", qmlMain)
	if err != nil {
		return err
	}

	app.window = root.CreateWindow(nil)
	app.window.On("publish", app.PublishText)

	app.window.Show()
	app.window.Wait()
	return nil
}

func (app *App) initProto() {
	stark, err := core.NewApp("stark")
	stark.Must(err)
	defer stark.Close()
	if *verbose {
		stark.Log.SetLevel(slog.LevelDebug)
	}

	app.ctx = stark.NewContext()
	app.client = proto.NewClient("desktop-"+proto.GenerateId(), app.ctx.Proto)

	app.client.Subscribe("", "self", app.HandleIncoming)
}

func (app *App) HandleIncoming(msg proto.Message) {
	text := msg.PayloadGetString("text")
	if text == "" {
		text = msg.Action + " from " + msg.Source
	}
	app.history += text + "<br>"
	app.window.Call("setHistory", app.history)
}

func (app *App) PublishText(text string) {
	app.ctx.Log.Debugln("publishing:", text)
	app.client.Publish(proto.Message{
		Action: "natural/handle",
		Payload: map[string]interface{}{
			"text": text,
		},
	})
}
