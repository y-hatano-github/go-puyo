package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

type gameStatus int

const controlTickCnt = 120

const (
	title gameStatus = iota
	newObject
	control
	falling
	chainCheck
	gameOver
	pause
)

type board struct {
	m        [104]int
	cm       [104]int
	score    int
	chainCnt int
	level    int
}

func (b *board) init() {
	for i := 0; i < 104; i++ {
		if i%8 == 0 {
			b.m[i] = 1
		} else {
			if i < 8 || i > 96 {
				b.m[i] = 1
			} else {
				b.m[i] = 0
			}
		}
	}
	b.m[4] = 0
	b.score = 0
	b.chainCnt = 0
	b.level = 1

}
func (b *board) tickCount() int {
	return controlTickCnt - (b.level-1)%10*5
}
func (b *board) set(p *object) {
	b.m[p.x1], b.m[p.x2] = p.c1, p.c2
}

func chain(x, p int, b *board, cnt *int) {
	if x > 0 && x < 104 && b.cm[x] == p && p > 1 {
		*cnt++
		b.cm[x] = 0
		chain(x-8, p, b, cnt)
		chain(x+1, p, b, cnt)
		chain(x+8, p, b, cnt)
		chain(x-1, p, b, cnt)
	}
}

type object struct {
	x1, x2, ex1, ex2, c1, c2, p int
}

func (o *object) init() {
	o.x1, o.x2 = 4, 12
	o.c1, o.c2 = rand.Intn(5)+2, rand.Intn(5)+2
	o.p = 1
}
func (o *object) set(x1, y1, po int) {
	o.x1, o.x2 = x1, y1
	o.p = po
}

func drawCell(x, y int, str string) {
	for _, v := range str {
		z, _ := strconv.Atoi(string(v))
		bg := ([]termbox.Attribute{termbox.ColorBlack,
			termbox.ColorWhite,
			termbox.ColorMagenta,
			termbox.ColorGreen,
			termbox.ColorRed,
			termbox.ColorBlue,
			termbox.ColorYellow,
			termbox.ColorCyan,
		})[z]
		cl, cr := '╺', '╸'
		if z < 2 {
			cl, cr = ' ', ' '
		}

		termbox.SetCell(x, y, cl, termbox.ColorBlack, bg)
		termbox.SetCell(x+1, y, cr, termbox.ColorBlack, bg)
		x += 2
	}
}

func keyEvent(key chan string) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				key <- "esc"
			case termbox.KeySpace:
				key <- " "
			case termbox.KeyEnter:
				key <- "enter"
			default:
				key <- string(ev.Ch)
			}
		}
	}
}

func drawString(x, y int, str string) {
	runes := []rune(str)
	for i, v := range runes {
		termbox.SetCell(x+i, y, v, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func updateConsole(s gameStatus, b *board) {
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)

	switch s {
	case title:
		drawCell(1, 2, "23456234562345623456")
		drawString(1, 5, "               Go-Puyo")
		drawString(1, 8, "         Hit Enter to play game.")
		drawString(1, 9, "         Hit ESC to exit.")
		drawCell(1, 12, "23456234562345623456")
	default:
		drawString(20, 3, "[ESC]   EXIT")
		drawString(20, 4, "[SPACE] DROP")
		drawString(20, 5, "[a]     LEFT")
		drawString(20, 6, "[d]     RIGHT")
		drawString(20, 7, "[s]     ROTATE RIGHT")
		drawString(20, 8, "[w]     ROTATE LEFT")
		drawString(20, 9, "[p]     PAUSE/RESUME")
		str := ""
		l := 0
		for i, v := range b.m {
			if i%8 == 0 && i > 0 {
				drawCell(1, 2+l, str+"1")
				str = ""
				l++
			}
			str += strconv.Itoa(v)
		}
		drawCell(1, 14, "111111111")
		drawString(1, 16, "SCORE:"+strconv.Itoa(b.score))
		drawString(1, 17, "LEVEL:"+strconv.Itoa(b.level))

		if s == pause {
			drawString(8, 8, "PAUSE")
		}

		if s == gameOver {
			drawString(1, 18, "*** GAME OVER ***")
			drawString(1, 19, "Hit [ENTER] to restart.")
			drawString(1, 20, "Hit [ESC] to exit.")
		}
	}
	termbox.Flush()
}

func execGame(key chan string) {
	s := title

	b := &board{}
	b.init()

	o := &object{}
	o.init()

	t := 0

MAINLOOP:
	for {
		startTime := time.Now().UnixNano() / int64(time.Millisecond)
		updateConsole(s, b)

		x1 := o.x1
		x2 := o.x2
		p := o.p

		select {
		case k := <-key:
			if k == "esc" {
				break MAINLOOP
			}
			switch s {
			case title:
				if k == "enter" {
					s = control
				}

			case control:
				if k == "p" {
					s = pause
					break
				}
				if k == " " {
					s = falling
					t = 0
					break
				}
				b.m[x1], b.m[x2] = 0, 0
				if t < b.tickCount() {
					x2 += b2i(k == "d")*b2i(b.m[o.x1+1] == 0)*b2i(b.m[o.x2+1] == 0) -
						b2i(k == "a")*b2i(b.m[o.x1-1] == 0)*b2i(b.m[o.x2-1] == 0)

					p += b2i(k == "s")*b2i(p%4 == 1)*b2i(o.x2+1 > 0)*b2i(b.m[o.x2+1] == 0) +
						b2i(k == "s")*b2i(p%4 == 2)*b2i(o.x2+8 > 0)*b2i(b.m[o.x2+8] == 0) +
						b2i(k == "s")*b2i(p%4 == 3)*b2i(o.x2-1 > 0)*b2i(b.m[o.x2-1] == 0) +
						b2i(k == "s")*b2i(p%4 == 0)*b2i(o.x2-8 > 0)*b2i(b.m[o.x2-8] == 0)

					p -= b2i(k == "w")*b2i(p%4 == 1)*b2i(o.x2-1 > 0)*b2i(b.m[o.x2-1] == 0) +
						b2i(k == "w")*b2i(p%4 == 2)*b2i(o.x2-8 > 0)*b2i(b.m[o.x2-8] == 0) +
						b2i(k == "w")*b2i(p%4 == 3)*b2i(o.x2+1 > 0)*b2i(b.m[o.x2+1] == 0) +
						b2i(k == "w")*b2i(p%4 == 0)*b2i(o.x2+8 > 0)*b2i(b.m[o.x2+8] == 0)
					if p < 0 {
						p = 4
					}

					x1 = x2 + b2i(p%4 == 1)*(-8) +
						b2i(p%4 == 2) +
						b2i(p%4 == 3)*8 +
						b2i(p%4 == 0)*(-1)
					o.set(x1, x2, p)
					b.set(o)
				}

			case gameOver:
				if k == "enter" {
					b.init()
					o.init()
					t = 0
					s = control
				}

			case pause:
				if k == "p" {
					s = control
				}
			}
		default:
			switch s {
			case newObject:
				o.init()
				b.chainCnt = 0
				s = control

			case control:
				if t == b.tickCount() {
					b.m[x1], b.m[x2] = 0, 0
					if b.m[x1+8] == 0 && b.m[x2+8] == 0 {
						x2 += 8
						x1 += 8
						o.set(x1, x2, p)
					} else {
						s = falling
						t = 0
					}
					b.set(o)
					if x1 == 4 {
						s = gameOver
					}
				}

			case falling:
				if t%30 == 0 {
					f := true
					for i := 96; i > 16; i-- {
						if b.m[i] == 0 && b.m[i-8] != 0 {
							b.m[i] = b.m[i-8]
							b.m[i-8] = 0
							f = false
						}
					}
					if f == true {
						s = chainCheck
					}
				}
			case chainCheck:
				s = newObject
				for i := 0; i < 96; i++ {
					b.cm = b.m
					cnt := 0
					chain(i, b.m[i], b, &cnt)
					if cnt > 3 {
						b.m = b.cm
						b.score += cnt
						b.score += 10 * b.chainCnt * b.level
						b.chainCnt++
						s = falling

						if (b.level * 10) <= b.score {
							b.level++
						}
					}
				}
			}
		}
		if s != title && s != gameOver {
			t++
			wait := true
			for wait {
				wait = (time.Now().UnixNano()/int64(time.Millisecond))-startTime <= 1
			}
			if t > b.tickCount() {
				t = 0
			}
		}
	}
}

func b2i(b bool) int {
	if b == true {
		return 1
	}
	return 0
}

func main() {

	rand.Seed(time.Now().UnixNano())

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	key := make(chan string)
	go keyEvent(key)

	execGame(key)
}
