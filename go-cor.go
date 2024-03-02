package main

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"slices"
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

type Cor struct {
	Hex      string
	Rgb      color.NRGBA
	CorFonte color.Color
}

func NovaCor(hex string) *Cor {
	hex = strings.ToUpper(hex)
	rgb := HexRGB(&hex)
	return &Cor{
		Hex:      hex,
		Rgb:      rgb,
		CorFonte: corContraste(&rgb),
	}
}

func HexRGB(hex *string) color.NRGBA {
	values, _ := strconv.ParseUint(string((*hex)[1:]), 16, 32)
	return color.NRGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}

func corContraste(cor *color.NRGBA) color.Color {
	if avaliacao := (0.299*float32(cor.R) + 0.587*float32(cor.G) + 0.114*float32(cor.B)) / 255; avaliacao > 0.5 {
		return color.Black
	} else {
		return color.White
	}
}

var cores []Cor
var x, y int
var corAtual Cor
var executando bool = true
var travaPosicao bool = false
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
		cores = []Cor{}
	} else {
		file, _ := os.Open(caminhoPerfil)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			cores = append(cores, *NovaCor(scanner.Text()))
		}
	}
}

func salvarPerfil() {
	file, _ := os.OpenFile(caminhoPerfil, os.O_WRONLY|os.O_TRUNC, 0644)
	writer := bufio.NewWriter(file)
	defer file.Close()
	defer writer.Flush()
	for i, cor := range cores {
		writer.WriteString(cor.Hex)
		if i != len(cores)-1 {
			writer.WriteString("\n")
		}
	}
}

func main() {
	lerPerfil()

	app := app.New()
	app.SetIcon(theme.ColorPaletteIcon())

	janela := app.NewWindow("Go-Cor")
	janela.Resize(fyne.Size{Width: 250, Height: 250})

	corAtual = *NovaCor("#FFFFFF")

	textoCorAtual := canvas.NewText(corAtual.Hex, corAtual.CorFonte)
	textoCorAtual.TextStyle = fyne.TextStyle{Monospace: true}
	textoCorAtual.TextSize = 14

	rectCorAtual := canvas.NewRectangle(corAtual.Rgb)
	rectCorAtual.CornerRadius = 10
	caixaCorAtual := container.NewStack(rectCorAtual, textoCorAtual)

	listaCores := container.NewVBox(widget.NewLabel("Cores salvas"))

	for _, corHex := range cores {
		caixaCor := colorBox(corHex, listaCores)
		listaCores.Add(caixaCor)
	}

	containerCoresSalvas := container.NewScroll(listaCores)
	containerCoresSalvas.SetMinSize(fyne.Size{Width: 200})

	iconePause := widget.NewIcon(theme.MediaPauseIcon())
	iconePause.Hidden = executando

	iconeTrava := widget.NewLabel("Travado")
	iconeTrava.Hidden = !travaPosicao

	content := container.NewHBox(containerCoresSalvas, container.NewVBox(iconePause, iconeTrava, caixaCorAtual))
	janela.SetContent(content)

	janela.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		switch ke.Name {
		case "H":
			caixaCor := colorBox(corAtual, listaCores)
			listaCores.Add(caixaCor)
			clipboard.WriteAll(corAtual.Hex)
			cores = append(cores, corAtual)
		case "P":
			executando = !executando
			iconePause.Hidden = executando
		case "T":
			travaPosicao = !travaPosicao
			iconeTrava.Hidden = !travaPosicao
		}

	})

	go func() {
		for range time.Tick(time.Duration(1) * time.Millisecond) {
			if executando {
				atualizarCor(textoCorAtual, rectCorAtual)
			}
		}
	}()

	janela.ShowAndRun()
	encerrar()
}

func colorBox(cor Cor, listaCores *fyne.Container) *fyne.Container {
	var caixaCor *fyne.Container

	texto := canvas.NewText(cor.Hex, cor.CorFonte)
	texto.TextStyle = fyne.TextStyle{Monospace: true}
	texto.TextSize = 12

	rect := canvas.NewRectangle(cor.Rgb)
	rect.CornerRadius = 5
	rect.SetMinSize(fyne.Size{Height: 20, Width: 60})

	botaoCopiar := widget.NewButtonWithIcon("", rscCopiar, func() { clipboard.WriteAll(cor.Hex) })
	botaoExcluir := widget.NewButtonWithIcon("", rscExcluir, func() {
		listaCores.Remove(caixaCor)
		deletarCor(&cor)
	})

	botaoCopiar.Resize(fyne.Size{Height: 20})
	botaoExcluir.Resize(fyne.Size{Height: 20})

	caixaCor = container.NewHBox(container.NewStack(rect, container.NewCenter(texto)), botaoCopiar, botaoExcluir)

	return caixaCor
}

func atualizarCor(textoCor *canvas.Text, rectCor *canvas.Rectangle) {
	if !travaPosicao {
		x, y = robotgo.Location()
	}
	corAtual = *NovaCor(fmt.Sprintf("#%v", robotgo.GetPixelColor(x, y)))

	textoCor.Text = corAtual.Hex
	textoCor.Color = corAtual.CorFonte
	rectCor.FillColor = corAtual.Rgb
	textoCor.Refresh()
}

func deletarCor(corDeletada *Cor) {
	for i, cor := range cores {
		if cor == *corDeletada {
			cores = slices.Delete(cores, i, i+1)
			return
		}
	}
}

func encerrar() {
	salvarPerfil()
	println("Encerrando...")
}
