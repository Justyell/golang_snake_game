package snake

import (
	"strconv"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
)

const (
	S_UP = iota
	S_DOWN
	S_LEFT
	S_RIGHT

	BODY_SIZE        = 2
	BODY_INIT_LENGTH = 3
	BODY_HEAD_NAME   = "snake_head"
)

var snakeOpLock sync.Mutex

type Snake struct {
	Direction  int
	Alive      bool
	Body       *Body
	playground *PlayGround
	GrowTimes  int
	gui        *gocui.Gui
}

type Body struct {
	Pos      *Pos
	IsHead   bool
	NextBody *Body
	BodyName string
}

type Pos struct {
	x1, y1, x2, y2 int
}

// func initLog() {
// 	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// 	if err != nil {
// 		log.Fatalln("打开文件失败：", err)
// 	}
// 	log.SetOutput(file)
// }

func InitSpirit(g *gocui.Gui, p *PlayGround) (*Snake, error) {

	// initLog()

	maxX, maxY := g.Size()
	init_x1 := maxX / 2
	init_y1 := maxY / 2
	init_x2 := maxX/2 + BODY_SIZE
	init_y2 := maxY/2 + BODY_SIZE - 1

	//防止坐标位基数，导致与food无法匹配上
	if init_x1%2 != 0 {
		init_x1--
		init_x2--
	}

	o := Snake{
		Direction:  S_LEFT,
		Alive:      true,
		playground: p,
		GrowTimes:  0,
		gui:        g,
	}

	b, e := o.generateTailBody(init_x1, init_y1, init_x2, init_y2, nil, BODY_HEAD_NAME, g, true)
	if e != nil {
		return nil, e
	}
	o.Body = b
	tail_x1 := init_x2
	tail_y1 := init_y1
	tail_x2 := tail_x1 + BODY_SIZE
	tail_y2 := init_y2
	prevBody := b
	for i := 0; i < BODY_INIT_LENGTH-1; i++ {
		iStr := strconv.Itoa(i)
		b, e := o.generateTailBody(tail_x1, tail_y1, tail_x2, tail_y2, prevBody, "snake_tail_"+iStr, g, false)
		if e != nil {
			return nil, e
		}
		tail_x1 = tail_x2
		tail_x2 = tail_x1 + BODY_SIZE
		prevBody = b
	}

	return &o, nil

}

func (s *Snake) generateTailBody(x1, y1, x2, y2 int, pBody *Body, bodyName string, g *gocui.Gui, isHead bool) (*Body, error) {

	s.gui.SetView(bodyName, x1, y1, x2, y2)

	pos := Pos{x1, y1, x2, y2}
	b := Body{
		Pos:      &pos,
		IsHead:   isHead,
		BodyName: bodyName,
	}
	if pBody != nil {
		pBody.NextBody = &b
	}
	return &b, nil
}

func (s *Snake) generateHeadBody(x1, y1, x2, y2 int) (*Body, error) {

	s.GrowTimes++

	growTimesStr := strconv.Itoa(s.GrowTimes)

	bodyName := "grow_" + growTimesStr

	s.gui.SetView(bodyName, x1, y1, x2, y2)

	pos := Pos{x1, y1, x2, y2}
	h := Body{
		Pos:      &pos,
		IsHead:   true,
		BodyName: bodyName,
	}

	oriHead := s.Body
	oriHead.IsHead = false
	s.Body = &h
	h.NextBody = oriHead
	return &h, nil
}

func (s *Snake) Move(e chan struct{}) {
	var prev_x1, prev_y1, prev_x2, prev_y2 int
	for {
		time.Sleep(500 * time.Millisecond)
		b := s.Body
		prev_x1 = b.Pos.x1
		prev_y1 = b.Pos.y1
		prev_x2 = b.Pos.x2
		prev_y2 = b.Pos.y2

		s.playground.BusySpaceLock.Lock()
		s.playground.BusySpace = nil
		s.playground.BusySpace = make([]*Pos, 0)

		for {
			var tmp_prev_x1, tmp_prev_y1, tmp_prev_x2, tmp_prev_y2 int

			snakeOpLock.Lock()
			direction := s.Direction
			snakeOpLock.Unlock()

			drawHead := func(b *Body) {
				var wg sync.WaitGroup
				wg.Add(1)
				s.gui.Update(func(g *gocui.Gui) error {
					g.SetView(b.BodyName, b.Pos.x1, b.Pos.y1, b.Pos.x2, b.Pos.y2)
					wg.Done()
					return nil
				})
				wg.Wait()
				//mark busy space
				pos := Pos{
					x1: b.Pos.x1,
					x2: b.Pos.x2,
					y1: b.Pos.y1,
					y2: b.Pos.y2,
				}
				s.playground.BusySpace = append(s.playground.BusySpace, &pos)
			}
			if b.IsHead {
				switch direction {
				case S_LEFT:

					b.Pos.x1 = b.Pos.x1 - BODY_SIZE
					b.Pos.x2 = b.Pos.x2 - BODY_SIZE
					drawHead(b)
					if b.NextBody != nil {
						b = b.NextBody
						continue
					} else {
						b = s.Body
						break
					}

				case S_RIGHT:

					b.Pos.x1 = b.Pos.x1 + BODY_SIZE
					b.Pos.x2 = b.Pos.x2 + BODY_SIZE
					drawHead(b)
					if b.NextBody != nil {
						b = b.NextBody
						continue
					} else {
						b = s.Body
						break
					}

				case S_UP:

					b.Pos.y1 = b.Pos.y1 - BODY_SIZE + 1
					b.Pos.y2 = b.Pos.y2 - BODY_SIZE + 1
					drawHead(b)
					if b.NextBody != nil {
						b = b.NextBody
						continue
					} else {
						b = s.Body
						break
					}

				case S_DOWN:
					b.Pos.y1 = b.Pos.y1 + BODY_SIZE - 1
					b.Pos.y2 = b.Pos.y2 + BODY_SIZE - 1
					drawHead(b)
					if b.NextBody != nil {
						b = b.NextBody
						continue
					} else {
						b = s.Body
						break
					}

				}
			}

			tmp_prev_x1 = b.Pos.x1
			tmp_prev_y1 = b.Pos.y1
			tmp_prev_x2 = b.Pos.x2
			tmp_prev_y2 = b.Pos.y2

			var wg sync.WaitGroup
			wg.Add(1)
			s.gui.Update(func(g *gocui.Gui) error {
				s.gui.SetView(b.BodyName, prev_x1, prev_y1, prev_x2, prev_y2)
				wg.Done()
				return nil
			})

			wg.Wait()
			b.Pos.x1 = prev_x1
			b.Pos.y1 = prev_y1
			b.Pos.x2 = prev_x2
			b.Pos.y2 = prev_y2

			prev_x1 = tmp_prev_x1
			prev_y1 = tmp_prev_y1
			prev_x2 = tmp_prev_x2
			prev_y2 = tmp_prev_y2

			//mark busy space
			pos := Pos{
				x1: b.Pos.x1,
				x2: b.Pos.x2,
				y1: b.Pos.y1,
				y2: b.Pos.y2,
			}
			s.playground.BusySpace = append(s.playground.BusySpace, &pos)

			if b.NextBody != nil {
				b = b.NextBody
			} else {
				break
			}

		}
		s.playground.BusySpaceLock.Unlock()

		e <- struct{}{}
		// go func() {
		// }()

	}
}

func (s *Snake) Grow() {
	var nx1, nx2, ny1, ny2 int

	switch s.Direction {
	case S_LEFT:
		nx1 = s.Body.Pos.x1 - BODY_SIZE
		nx2 = s.Body.Pos.x2 - BODY_SIZE
		ny1 = s.Body.Pos.y1
		ny2 = s.Body.Pos.y2
	case S_RIGHT:
		nx1 = s.Body.Pos.x1 + BODY_SIZE
		nx2 = s.Body.Pos.x2 + BODY_SIZE
		ny1 = s.Body.Pos.y1
		ny2 = s.Body.Pos.y2
	case S_UP:
		nx1 = s.Body.Pos.x1
		nx2 = s.Body.Pos.x2
		ny1 = s.Body.Pos.y1 - BODY_SIZE
		ny2 = s.Body.Pos.y2 - BODY_SIZE
	case S_DOWN:
		nx1 = s.Body.Pos.x1
		nx2 = s.Body.Pos.x2
		ny1 = s.Body.Pos.y1 + BODY_SIZE
		ny2 = s.Body.Pos.y2 + BODY_SIZE
	}

	s.generateHeadBody(nx1, ny1, nx2, ny2)
}

func (s *Snake) MoveLeft() {
	snakeOpLock.Lock()
	defer snakeOpLock.Unlock()

	if s.Direction == S_UP || s.Direction == S_DOWN {
		s.Direction = S_LEFT
	}

}

func (s *Snake) MoveRight() {
	snakeOpLock.Lock()
	defer snakeOpLock.Unlock()
	if s.Direction == S_UP || s.Direction == S_DOWN {
		s.Direction = S_RIGHT
	}

}

func (s *Snake) MoveUp() {
	snakeOpLock.Lock()
	defer snakeOpLock.Unlock()
	if s.Direction == S_LEFT || s.Direction == S_RIGHT {
		s.Direction = S_UP
	}
}

func (s *Snake) MoveDown() {
	snakeOpLock.Lock()
	defer snakeOpLock.Unlock()
	if s.Direction == S_LEFT || s.Direction == S_RIGHT {
		s.Direction = S_DOWN
	}
}
