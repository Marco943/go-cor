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

	corAtual := canvas.NewText("", color.White)
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

func HexRGB(hex string) color.NRGBA {
	values, _ := strconv.ParseUint(string(hex), 16, 32)
	return color.NRGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}

func corContraste(cor *color.NRGBA) color.Color {
	if avaliacao := (0.299*float32(cor.R) + 0.587*float32(cor.G) + 0.114*float32(cor.B)) / 255; avaliacao > 0.5 {
		return color.Black
	} else {
		return color.White
	}
}

func atualizarCor(cor *canvas.Text, caixaCor *canvas.Rectangle) {
	x, y := robotgo.Location()
	corHex := robotgo.GetPixelColor(x, y)
	corRGB := HexRGB(corHex)

	cor.Text = corHex
	cor.Color = corContraste(&corRGB)
	cor.Refresh()
	caixaCor.FillColor = corRGB
}

func encerrar() {
	fmt.Println("Encerrando...")
}
