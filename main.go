package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"time"

	"github.com/xconstruct/stark/core"
	slog "github.com/xconstruct/stark/log"
	"github.com/xconstruct/stark/proto"
	"github.com/xconstruct/stark-desktop/assets"
	"gopkg.in/qml.v1"
)

var verbose = flag.Bool("v", false, "verbose debug output")

type History struct {
	Type string
	Time time.Time
	Text string
	Message proto.Message
}

type App struct {
	ctx     *core.Context
	client  *proto.Client
	history []History
	historyText string
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

	root, err := engine.LoadString("qml.go", assets.QmlMainWindow)
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
	app.AddHistory(History{
		Type: "status",
		Time: time.Now(),
		Text: "Connected",
	})
}

func (app *App) HandleIncoming(msg proto.Message) {
	text := msg.PayloadGetString("text")
	if text == "" {
		text = msg.Action + " from " + msg.Source
	}
	htype := "in"
	if msg.Source == "" {
		htype = "out"
	}

	app.AddHistory(History{
		Type: htype,
		Time: time.Now(),
		Text: text,
		Message: msg,
	})
}

func (app *App) AddHistory(h History) {
	app.history = append(app.history, h)

	hist := ""
	for _, h := range app.history {
		style := `<span style="color: #0000ff">%s</span> %s<br>`
		switch h.Type {
		case "out":
			style = `<span style="color: #0000ff">%s</span> <b>%s</b><br>`
		case "status":
			style = `<span style="color: #0000ff">%s</span> <i>%s</i><br>`
		}
		hist += fmt.Sprintf(style, h.Time.Format("15:04:05"), h.Text)
	}
	app.window.Call("setHistory", hist)
}

func (app *App) PublishText(text string) {
	app.ctx.Log.Debugln("publishing:", text)
	msg := proto.Message{
		Action: "natural/handle",
		Payload: map[string]interface{}{
			"text": text,
		},
	}
	app.client.Publish(msg)
	app.HandleIncoming(msg)
}
