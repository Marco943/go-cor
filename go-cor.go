package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	"github.com/go-vgo/robotgo/clipboard"
)

var cores = []string{"ffffff", "ff0000", "ff00ff", "00ff00", "bb7700", "77dd00", "928274", "ffffff", "ff0000", "ff00ff", "00ff00", "bb7700", "77dd00", "928274"}

func main() {
	app := app.New()
	janela := app.NewWindow("Go-Cor")
	janela.Resize(fyne.Size{Width: 250, Height: 250})

	corAtual := canvas.NewText("", color.White)
	corAtual.TextStyle = fyne.TextStyle{Monospace: true}
	corAtual.TextSize = 14

	rectCorAtual := canvas.NewRectangle(color.NRGBA{255, 255, 255, 255})
	rectCorAtual.CornerRadius = 10
	caixaCorAtual := container.NewStack(rectCorAtual, corAtual)

	listaCores := container.NewVBox(widget.NewLabel("Cores Salvas"))

	for _, corHex := range cores {
		caixaCor := colorBox(corHex, listaCores)

		listaCores.Add(caixaCor)
	}

	containerCoresSalvas := container.NewScroll(listaCores)

	content := container.NewGridWithColumns(2, containerCoresSalvas, caixaCorAtual)
	janela.SetContent(content)

	go func() {
		for range time.Tick(time.Millisecond) {
			atualizarCor(corAtual, rectCorAtual)
		}
	}()

	janela.ShowAndRun()
	encerrar()
}

func colorBox(corHex string, listaCores *fyne.Container) *fyne.Container {
	var caixaCor *fyne.Container

	corHex = strings.ToUpper(corHex)

	corRGB := HexRGB(&corHex)
	corFonte := corContraste(&corRGB)

	texto := canvas.NewText(fmt.Sprintf("#%v", corHex), corFonte)
	texto.TextStyle = fyne.TextStyle{Monospace: true}
	texto.TextSize = 12

	rect := canvas.NewRectangle(corRGB)
	rect.CornerRadius = 5
	rect.SetMinSize(fyne.NewSize(20, 30))

	botaoCopiar := widget.NewButton("Copiar", func() { clipboard.WriteAll(fmt.Sprintf("#%v", corHex)) })
	botaoExcluir := widget.NewButton("Excluir", func() { listaCores.Remove(caixaCor) })

	caixaCor = container.NewHBox(container.NewStack(rect, container.NewCenter(texto)), botaoCopiar, botaoExcluir)

	return caixaCor
}

func HexRGB(hex *string) color.NRGBA {
	values, _ := strconv.ParseUint(string(*hex), 16, 32)
	return color.NRGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}

func corContraste(cor *color.NRGBA) color.Color {
	if avaliacao := (0.299*float32(cor.R) + 0.587*float32(cor.G) + 0.114*float32(cor.B)) / 255; avaliacao > 0.5 {
		return color.Black
	} else {
		return color.White
	}
}

func atualizarCor(cor *canvas.Text, rectCor *canvas.Rectangle) {
	x, y := robotgo.Location()
	corHex := strings.ToUpper(robotgo.GetPixelColor(x, y))
	corRGB := HexRGB(&corHex)

	cor.Text = fmt.Sprintf("#%v", corHex)
	cor.Color = corContraste(&corRGB)
	cor.Refresh()
	rectCor.FillColor = corRGB
}

func encerrar() {
	fmt.Println("Encerrando...")
}
