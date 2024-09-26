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
	_ "embed"
	"image/color"
	"log"

	"golang.org/x/text/language"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kijimaD/nov/event"
	"github.com/kijimaD/nov/lexer"
	"github.com/kijimaD/nov/parser"
)

const (
	screenWidth  = 720
	screenHeight = 720
	fontSize     = 26
)

var japaneseFaceSource *text.GoTextFaceSource
var eventQ event.Queue

//go:embed JF-Dot-Kappa20B.ttf
var font []byte

//go:embed input.sce
var input []byte

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(font))
	if err != nil {
		log.Fatal(err)
	}
	japaneseFaceSource = s

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
}

type Game struct {
	bgImage *ebiten.Image
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		eventQ.Run()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0)
		screen.DrawImage(g.bgImage, op)
	}

	{
		black := color.RGBA{0x00, 0x00, 0x00, 0x50}
		vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, black, false)
	}

	{
		japaneseText := eventQ.Display()
		f := &text.GoTextFace{
			Source:   japaneseFaceSource,
			Size:     fontSize,
			Language: language.Japanese,
		}
		const lineSpacing = fontSize + 4
		const padding = 20
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
	img, _, err := ebitenutil.NewImageFromFile("./_example/forest.jpg")
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Text I18N (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{bgImage: img}); err != nil {
		log.Fatal(err)
	}
}
