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
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	"github.com/go-vgo/robotgo/clipboard"
)

type TappableBox struct {
	widget.BaseWidget
	cor *Cor
}

func (b *TappableBox) CreateRenderer() fyne.WidgetRenderer {
	texto := canvas.NewText(b.cor.Hex, b.cor.CorFonte)
	texto.TextStyle = fyne.TextStyle{Monospace: true}
	texto.TextSize = 12

	rect := canvas.NewRectangle(b.cor.Rgb)
	rect.CornerRadius = 5
	rect.SetMinSize(fyne.Size{Height: 16, Width: 60})

	iconeFixado := widget.NewCheckWithData("", binding.BindBool(&b.cor.fixo))

	c := container.NewHBox(iconeFixado, container.NewStack(rect, container.NewCenter(texto)))

	return widget.NewSimpleRenderer(c)
}

func (b *TappableBox) Tapped(*fyne.PointEvent) {
	clipboard.WriteAll(b.cor.Hex)
}

func NewTappableBox(cor *Cor) *TappableBox {
	box := &TappableBox{
		cor: cor,
	}
	box.ExtendBaseWidget(box)
	return box
}

type Cor struct {
	Hex      string
	Rgb      color.NRGBA
	CorFonte color.Color
	fixo     bool
}

func NovaCor(hex string, fixo bool) *Cor {
	hex = strings.ToUpper(hex)
	rgb := HexRGB(&hex)
	return &Cor{
		Hex:      hex,
		Rgb:      rgb,
		CorFonte: corContraste(&rgb),
		fixo:     fixo,
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

var cores []*Cor
var x, y int
var corAtual Cor
var executando bool = true
var travaPosicao bool = false
var caminhoPerfil string

func lerPerfil() {
	dirCache, _ := os.UserCacheDir()
	dirPerfil := filepath.Join(dirCache, "go-cor")
	os.MkdirAll(dirPerfil, 0700)

	caminhoPerfil = filepath.Join(dirPerfil, "perfil")
	if _, err := os.Stat(caminhoPerfil); os.IsNotExist(err) {
		file, _ := os.Create(caminhoPerfil)
		defer file.Close()
		cores = []*Cor{}
	} else {
		file, _ := os.Open(caminhoPerfil)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			cores = append(cores, NovaCor(scanner.Text(), true))
		}
	}
}

func salvarPerfil() {
	file, _ := os.OpenFile(caminhoPerfil, os.O_WRONLY|os.O_TRUNC, 0644)
	writer := bufio.NewWriter(file)
	defer file.Close()
	defer writer.Flush()
	for i, cor := range cores {
		if !cor.fixo {
			continue
		}
		writer.WriteString(cor.Hex)
		if i != len(cores)-1 {
			writer.WriteString("\n")
		}
	}
}

func setTopMost() {
	janelas, err := robotgo.FindIds("Go-Cor")
	if err != nil {
		return
	}
	setTopMostWindows(janelas[0])
}

func main() {
	lerPerfil()

	app := app.New()
	app.SetIcon(theme.ColorPaletteIcon())

	janela := app.NewWindow("Go-Cor")
	janela.Resize(fyne.Size{Width: 220, Height: 250})
	janela.SetFixedSize(true)

	corAtual = *NovaCor("#FFFFFF", false)

	textoCorAtual := canvas.NewText(corAtual.Hex, corAtual.CorFonte)
	textoCorAtual.TextStyle = fyne.TextStyle{Monospace: true}
	textoCorAtual.TextSize = 14

	rectCorAtual := canvas.NewRectangle(corAtual.Rgb)
	rectCorAtual.CornerRadius = 10
	rectCorAtual.SetMinSize(fyne.NewSize(60, 60))

	caixaCorAtual := container.NewStack(rectCorAtual, container.NewCenter(textoCorAtual))

	listaCores := container.NewVBox(widget.NewLabel("Cores salvas"))

	for _, corHex := range cores {
		cor := corHex
		caixaCor := NewTappableBox(cor)
		listaCores.Add(caixaCor)
	}

	containerCoresSalvas := container.NewScroll(listaCores)

	iconePause := widget.NewIcon(theme.MediaPauseIcon())
	iconePause.Hidden = executando
	botaoPause := widget.NewButton("Pausar", func() {
		executando = !executando
		iconePause.Hidden = executando
	})

	textPos := canvas.NewText(fmt.Sprintf("(%v, %v)", x, y), color.White)
	textPos.TextSize = 11
	textPos.TextStyle.Monospace = true

	iconeTrava := widget.NewIcon(theme.MediaStopIcon())
	iconeTrava.Hidden = !travaPosicao

	containerDireita := container.NewVBox(caixaCorAtual, botaoPause, iconePause, container.NewHBox(textPos, iconeTrava))

	content := container.NewGridWithColumns(2, containerCoresSalvas, containerDireita)
	janela.SetContent(content)

	janela.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		switch ke.Name {
		case "H":
			cor := corAtual
			cores = append(cores, &cor)
			caixaCor := NewTappableBox(&cor)
			listaCores.Add(caixaCor)
			clipboard.WriteAll(corAtual.Hex)
		case "P":
			executando = !executando
			iconePause.Hidden = executando
		case "T":
			travaPosicao = !travaPosicao
			iconeTrava.Hidden = !travaPosicao
		}

	})

	// if desk, ok := app.(desktop.App); ok {
	// 	m := fyne.NewMenu("Go-Cors",
	// 		fyne.NewMenuItem("Abrir", func() {
	// 			executando = true
	// 			janela.Show()
	// 		}),
	// 	)
	// 	desk.SetSystemTrayMenu(m)
	// }

	// janela.SetCloseIntercept(func() {
	// 	executando = false
	// 	janela.Hide()
	// })

	go func() {
		for range time.Tick(time.Duration(1) * time.Millisecond) {
			if executando {
				atualizarCor(textoCorAtual, rectCorAtual, textPos)
			}
		}
	}()

	var i int
	go func() {
		for range time.Tick(time.Duration(1) * time.Second) {
			if i == 0 {
				setTopMost()
			}
			i++
		}
	}()

	janela.Show()

	app.Run()

	encerrar()
}

func atualizarCor(textoCor *canvas.Text, rectCor *canvas.Rectangle, textPos *canvas.Text) {
	if !travaPosicao {
		x, y = robotgo.Location()
		textPos.Text = fmt.Sprintf("(%v, %v)", x, y)
		textPos.Refresh()
	}
	corAtual = *NovaCor(fmt.Sprintf("#%v", robotgo.GetPixelColor(x, y)), false)

	textoCor.Text = corAtual.Hex
	textoCor.Color = corAtual.CorFonte
	rectCor.FillColor = corAtual.Rgb
	textoCor.Refresh()
}

func encerrar() {
	salvarPerfil()
}
