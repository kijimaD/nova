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
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kijimaD/nov/event"
	"github.com/kijimaD/nov/lexer"
	"github.com/kijimaD/nov/parser"
)

const input = `*start
親譲りの無鉄砲で小供の時から損ばかりして居る。
小学校に居る時分学校の二階から飛び降りて一週間程腰を抜かした事がある。
なぜそんな無闇をしたと聞く人があるかも知れぬ。
別段深い理由でもない。新築の二階から首を出して居たら、
同級生の一人が冗談に、いくら威張つても、そこから飛び降りる事は出来まい。
弱虫やーい。と囃したからである。
小使に負ぶさつて帰つて来た時、おやぢが大きな眼をして二階位から
飛び降りて腰を抜かす奴があるかと云つたから、
此次は抜かさずに飛んで見せますと答へた。
[jump target="start"]`

var japaneseFaceSource *text.GoTextFaceSource
var eventQ event.Queue

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	japaneseFaceSource = s

	l := lexer.NewLexer(input)
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

const (
	screenWidth  = 640
	screenHeight = 640
	fontSize     = 12
)

type Game struct {
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		eventQ.Run()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	gray := color.RGBA{0x80, 0x80, 0x80, 0xff}

	{
		japaneseText := eventQ.Display()
		f := &text.GoTextFace{
			Source:   japaneseFaceSource,
			Size:     fontSize,
			Language: language.Japanese,
		}
		const lineSpacing = fontSize + 4
		x, y := 0, 280
		w, h := text.Measure(japaneseText, f, lineSpacing)
		// The left upper point is not x but x-w, since the text runs in the rigth-to-left direction as the secondary direction.
		vector.DrawFilledRect(screen, float32(x)-float32(w), float32(y), float32(w), float32(h), gray, false)
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
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Text I18N (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
