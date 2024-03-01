package main

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	"github.com/go-vgo/robotgo/clipboard"
)

var cores []string
var rscCopiar fyne.Resource = theme.ContentCopyIcon()
var rscExcluir fyne.Resource = theme.NewErrorThemedResource(theme.DeleteIcon())
var caminhoPerfil string

func lerPerfil() {
	dirCache, _ := os.UserCacheDir()
	dirPerfil := filepath.Join(dirCache, "go-cor")
	if _, err := os.Stat(dirPerfil); os.IsNotExist(err) {
		os.MkdirAll(dirPerfil, 0700)
	}
	caminhoPerfil = filepath.Join(dirPerfil, "perfil")
	if _, err := os.Stat(caminhoPerfil); os.IsNotExist(err) {
		file, _ := os.Create(caminhoPerfil)
		defer file.Close()
		cores = []string{}
	} else {
		file, _ := os.Open(caminhoPerfil)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			cores = append(cores, scanner.Text())
		}
	}
}

func salvarPerfil() {
	file, _ := os.OpenFile(caminhoPerfil, os.O_WRONLY, 0644)
	writer := bufio.NewWriter(file)
	defer file.Close()
	defer writer.Flush()
	for _, cor := range cores {
		writer.WriteString(cor)
		writer.WriteString("\n")
	}
}

func main() {
	lerPerfil()

	app := app.New()
	app.SetIcon(theme.ColorPaletteIcon())

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
	containerCoresSalvas.SetMinSize(fyne.Size{Width: 200})

	content := container.NewHBox(containerCoresSalvas, caixaCorAtual)
	janela.SetContent(content)

	janela.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		switch ke.Name {
		case "H":
			caixaCor := colorBox(corAtual.Text[1:], listaCores)
			listaCores.Add(caixaCor)
			clipboard.WriteAll(corAtual.Text)
			cores = append(cores, corAtual.Text[1:])
		}
	})

	go func() {
		for range time.Tick(time.Duration(1) * time.Millisecond) {
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
	rect.SetMinSize(fyne.Size{Height: 20, Width: 60})

	botaoCopiar := widget.NewButtonWithIcon("", rscCopiar, func() { clipboard.WriteAll(fmt.Sprintf("#%v", corHex)) })
	botaoExcluir := widget.NewButtonWithIcon("", rscExcluir, func() {
		listaCores.Remove(caixaCor)
		for i, cor := range cores {
			if cor == corHex {

				cores = append(cores[:i], cores[i+1:]...)
			}
		}
	})

	botaoCopiar.Resize(fyne.Size{Height: 20})
	botaoExcluir.Resize(fyne.Size{Height: 20})

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
	salvarPerfil()
	println("Encerrando...")
}
