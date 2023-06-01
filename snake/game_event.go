package snake

import (
	"sync"
	"time"
)

type PlayGround struct {
	BusySpace     []*Pos
	BusySpaceLock *sync.Mutex
}

func NewPlayGround() *PlayGround {
	fs := make([]*Pos, 0)
	p := PlayGround{
		BusySpace:     fs,
		BusySpaceLock: &sync.Mutex{},
	}

	return &p
}

type GameEvent struct {
	GameStart chan struct{}
	FoodTouch chan struct{}
	Die       chan struct{}
	Move      chan struct{}
	snakeObj  *Snake
	pg        *PlayGround
	food      *Food
}

func NewGameEvent(sk *Snake, pg *PlayGround, food *Food) *GameEvent {
	g := make(chan struct{})
	f := make(chan struct{})
	d := make(chan struct{})
	m := make(chan struct{})

	event := GameEvent{
		Move:      m,
		GameStart: g,
		FoodTouch: f,
		Die:       d,
		snakeObj:  sk,
		pg:        pg,
		food:      food,
	}

	return &event
}

func (ge *GameEvent) ListenEvent() {

	for {
		select {
		case <-ge.Move:

			ge.pg.BusySpaceLock.Lock()
			fx1 := ge.food.Pos.x1
			fy1 := ge.food.Pos.y1
			fx2 := ge.food.Pos.x2
			fy2 := ge.food.Pos.y2

			skx1 := ge.snakeObj.Body.Pos.x1
			sky1 := ge.snakeObj.Body.Pos.y1
			skx2 := ge.snakeObj.Body.Pos.x2
			sky2 := ge.snakeObj.Body.Pos.y2

			if fx1 == skx1 &&
				fy1 == sky1 &&
				fx2 == skx2 &&
				fy2 == sky2 {

				go func() {
					ge.FoodTouch <- struct{}{}
				}()
			}
			ge.pg.BusySpaceLock.Unlock()

		case <-ge.FoodTouch:
			ge.snakeObj.Grow()
			ge.food.Disappear()
			ge.food.Appear()

		default:
		}
		time.Sleep(50 * time.Millisecond)
	}
}
