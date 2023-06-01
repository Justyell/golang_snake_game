package snake

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jroimartin/gocui"
)

type Food struct {
	Pos        *Pos
	Name       string
	FreeSpace  []*Pos
	playground *PlayGround
	gui        *gocui.Gui
}

func NewFood(g *gocui.Gui, pg *PlayGround) *Food {
	p := Pos{}
	f := Food{
		Pos:        &p,
		Name:       "food",
		playground: pg,
		gui:        g,
	}

	return &f
}

func (f *Food) Disappear() {
	f.gui.DeleteView(f.Name)
}

func (f *Food) Appear() {
	maxX, maxY := f.gui.Size()
	rand.Seed(time.Now().UnixNano())
	tc := make(map[string]struct{})
	bc := make(map[string]struct{})

	f.playground.BusySpaceLock.Lock()
	f.FreeSpace = nil
	for _, pos := range f.playground.BusySpace {
		tc[fmt.Sprintf("%d,%d", pos.x1, pos.y1)] = struct{}{}
		bc[fmt.Sprintf("%d,%d", pos.x2, pos.y2)] = struct{}{}
	}
	for ix := 0; ix <= maxX-BODY_SIZE; ix++ {
		if ix%2 != 0 {
			continue
		}
		for iy := 0; iy <= maxY-BODY_SIZE-1; iy++ {

			if _, ok := tc[fmt.Sprintf("%d,%d", ix, iy)]; ok {
				continue
			}
			if _, ok := bc[fmt.Sprintf("%d,%d", ix, iy)]; ok {
				continue
			}

			pos := Pos{
				ix, iy, ix + BODY_SIZE, iy + BODY_SIZE - 1,
			}
			f.FreeSpace = append(f.FreeSpace, &pos)
		}
	}
	rand.Seed(time.Now().UnixNano())
	randFreeIndex := rand.Intn(len(f.FreeSpace) - 1)
	appearPos := f.FreeSpace[randFreeIndex]

	f.Pos.x1 = appearPos.x1
	f.Pos.y1 = appearPos.y1
	f.Pos.x2 = appearPos.x2
	f.Pos.y2 = appearPos.y2

	v, _ := f.gui.SetView(f.Name, appearPos.x1, appearPos.y1, appearPos.x2, appearPos.y2)

	v.Highlight = true
	v.BgColor = gocui.ColorRed
	f.playground.BusySpaceLock.Unlock()
}
