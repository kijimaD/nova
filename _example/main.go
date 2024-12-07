// Copyright 2023 The Ebitengine Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"embed"
	_ "embed"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"golang.org/x/text/language"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kijimaD/nova/event"
	"github.com/kijimaD/nova/lexer"
	"github.com/kijimaD/nova/parser"
)

const (
	screenWidth  = 720
	screenHeight = 720
	fontSize     = 26
	padding      = 40
)

var japaneseFaceSource *text.GoTextFaceSource
var eventQ event.Queue

//go:embed input.sce
var input []byte

//go:embed file
var FS embed.FS

type Game struct {
	bgImage     *ebiten.Image
	promptImage *ebiten.Image
	startTime   time.Time
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		eventQ.Run()
	}

	select {
	case v := <-eventQ.NotifyChan:
		switch event := v.(type) {
		case *event.ChangeBg:
			eimg, err := loadImage(event.Source)
			if err != nil {
				log.Fatal(err)
			}
			g.bgImage = eimg
		}
	default:
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	{
		// 背景画像
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, screenHeight/4)
		screen.DrawImage(g.bgImage, op)
	}

	{
		// 背景色
		black := color.RGBA{0x00, 0x00, 0x00, 0x80}
		vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, black, false)
	}

	f := &text.GoTextFace{
		Source:   japaneseFaceSource,
		Size:     fontSize,
		Language: language.Japanese,
	}

	// 待ち状態表示
	if eventQ.OnAnim {
		elapsed := time.Since(g.startTime).Seconds()
		offsetY := 4 * math.Sin(elapsed*4) // sin関数で上下に動かす
		bounds := g.promptImage.Bounds()
		bounds.Min.Y = int(20 + offsetY) // 初期位置 + オフセット
		bounds.Max.Y = bounds.Min.Y + bounds.Dy()

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(720-float64(bounds.Max.X*2), 720-float64(bounds.Min.Y*2))
		screen.DrawImage(g.promptImage, op)
	}

	{
		japaneseText := eventQ.Display()
		const lineSpacing = fontSize + 4
		x, y := padding, padding
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		op.LineSpacing = lineSpacing
		text.Draw(screen, japaneseText, f, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{startTime: time.Now()}
	{
		font, err := FS.ReadFile("file/JF-Dot-Kappa20B.ttf")
		if err != nil {
			log.Fatal(err)
		}
		s, err := text.NewGoTextFaceSource(bytes.NewReader(font))
		if err != nil {
			log.Fatal(err)
		}
		japaneseFaceSource = s
	}

	l := lexer.NewLexer(string(input))
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	if err != nil {
		log.Fatal(err)
	}
	e := event.NewEvaluator()
	e.Eval(program)
	eventQ = event.NewQueue(e)
	eventQ.Start()

	{
		eimg, err := loadImage("file/black.png")
		if err != nil {
			log.Fatal(err)
		}
		game.bgImage = eimg
	}
	{
		eimg, err := loadImage("file/prompt.png")
		if err != nil {
			log.Fatal(err)
		}
		game.promptImage = eimg
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("demo")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func loadImage(filename string) (*ebiten.Image, error) {
	bs, err := FS.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	dec, _, err := image.Decode(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	img := ebiten.NewImageFromImage(dec)
	if err != nil {
		return nil, err
	}
	return img, nil
}
