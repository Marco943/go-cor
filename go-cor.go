package main

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
)

func main() {
	app := app.New()
	janela := app.NewWindow("Go-Cor")
	janela.Resize(fyne.Size{Width: 250, Height: 250})

	corAtual := widget.NewLabel("")
	caixaCorAtual := canvas.NewRectangle(color.NRGBA{255, 255, 255, 255})

	content := container.NewAdaptiveGrid(1, widget.NewLabel("Cores Salvas"), corAtual, caixaCorAtual)
	janela.SetContent(content)

	go func() {
		for range time.Tick(time.Millisecond) {
			atualizarCor(corAtual, caixaCorAtual)
		}
	}()

	janela.ShowAndRun()
	encerrar()
}

func HexColor(hex string) color.NRGBA {
	values, _ := strconv.ParseUint(string(hex), 16, 32)
	return color.NRGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}

func atualizarCor(cor *widget.Label, caixaCor *canvas.Rectangle) {
	x, y := robotgo.Location()
	corHex := robotgo.GetPixelColor(x, y)
	cor.SetText(corHex)
	caixaCor.FillColor = HexColor(corHex)
}

func encerrar() {
	fmt.Println("Encerrando...")
}
