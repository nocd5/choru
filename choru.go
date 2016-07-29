package choru

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-runewidth"
	"github.com/mattn/go-tty"
)

const (
	VtDefault = "0"

	FgDefault     = VtDefault
	FgBlack       = "30"
	FgRed         = "31"
	FgGreen       = "32"
	FgYellow      = "33"
	FgBlue        = "34"
	FgMagenta     = "35"
	FgCyan        = "36"
	FgWhite       = "37"
	FgBlackBold   = "30;1"
	FgRedBold     = "31;1"
	FgGreenBold   = "32;1"
	FgYellowBold  = "33;1"
	FgBlueBold    = "34;1"
	FgMagentaBold = "35;1"
	FgCyanBold    = "36;1"
	FgWhiteBold   = "37;1"

	BgBlack         = "40"
	BgRed           = "41"
	BgGreen         = "42"
	BgYellow        = "43"
	BgBlue          = "44"
	BgMagenta       = "45"
	BgCyan          = "46"
	BgWhite         = "47"
	BgDefault       = "49"
	BgBlackBright   = "100"
	BgRedBright     = "101"
	BgGreenBright   = "102"
	BgYellowBright  = "103"
	BgBlueBright    = "104"
	BgMagentaBright = "105"
	BgCyanBright    = "106"
	BgWhiteBright   = "107"

	vt_end      = "\033[0m"
	cursor_hide = "\033[?25l"
	cursor_show = "\033[?25h"

	vk_return = 0x0d
	vk_esc    = 0x1b

	erase_from_cursor = 0
	erase_to_cursor   = 1
	erase_all         = 2
)

type Choru struct {
	LineFg    string
	LineBg    string
	CursorFg  string
	CursorBg  string
	MaxHeight int
	Header    string
	HeaderFg  string
	HeaderBg  string
	Footer    string
	FooterFg  string
	FooterBg  string
}

func New() *Choru {
	return &Choru{
		LineFg:    VtDefault,
		LineBg:    BgDefault,
		CursorFg:  VtDefault,
		CursorBg:  BgMagenta,
		MaxHeight: -1,
		Header:    "",
		HeaderFg:  VtDefault,
		HeaderBg:  BgDefault,
		Footer:    "",
		FooterFg:  VtDefault,
		FooterBg:  BgDefault,
	}
}

func (c *Choru) Choose(items []string) (int, string) {
	tty, err := tty.Open()
	if err != nil {
		panic(err)
	}
	out := colorable.NewColorable(tty.Output())

	cleanup := func() {
		out.Write([]byte(vt_end))
		if c.Header != "" {
			out.Write([]byte(cursorPreviousLine(1)))
		}
		out.Write([]byte(eraseDisplay(erase_from_cursor)))
		out.Write([]byte(cursor_show))
	}

	index := -1
	label := ""
	out.Write([]byte(cursor_hide))
	defer cleanup()

	// Ctrl-C handler
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	go func() {
		<-sig
		cleanup()
		os.Exit(2)
	}()

	offset := 0
	row := 0
	dirty := make([]bool, len(items))
	for i := 0; i < len(dirty); i++ {
		dirty[i] = true
	}

	width, height, err := tty.Size()
	if err != nil {
		width = 80
		height = 25
	}
	if c.MaxHeight > 0 {
		height = min(height, c.MaxHeight)
	}

	count := len(items)
	if c.Header != "" {
		count++
	}
	if c.Footer != "" {
		count++
	}
	height = min(height, count)

	// reserve display area
	for i := 0; i < height-1; i++ {
		fmt.Fprintln(out, "")
	}
	out.Write([]byte(cursorPreviousLine(height - 1)))

	if c.Header != "" {
		out.Write([]byte(getAnsiColor(c.HeaderFg, c.HeaderBg)))
		out.Write([]byte(truncate(c.Header, width)))
		out.Write([]byte(vt_end))
		out.Write([]byte(cursorNextLine(1)))
		height--
	}
	if c.Footer != "" {
		out.Write([]byte(cursorNextLine(height - 1)))
		out.Write([]byte(getAnsiColor(c.FooterFg, c.FooterBg)))
		out.Write([]byte(truncate(c.Footer, width)))
		out.Write([]byte(vt_end))
		out.Write([]byte(cursorPreviousLine(height - 1)))
		height--
	}

	var pos int
EVENT_LOOP:
	for {
		// draw items
		pos = 0
		for i, line := range items[offset:] {
			if dirty[offset+i] {
				out.Write([]byte(eraseLine(erase_from_cursor)))
				if offset+i == row {
					out.Write([]byte(getAnsiColor(c.CursorFg, c.CursorBg)))
				} else {
					out.Write([]byte(getAnsiColor(c.LineFg, c.LineBg)))
				}
				out.Write([]byte(truncate(line, width)))
				out.Write([]byte(vt_end))
				dirty[offset+i] = false
			}

			// pos is 0 oriented
			if (pos + 1) >= height {
				break
			}

			out.Write([]byte(cursorNextLine(1)))
			pos++
		}
		// reset position
		out.Write([]byte(cursorPreviousLine(pos)))

		// key event
		r, err := tty.ReadRune()
		if err != nil {
			panic(err)
		}
		switch r {
		case 'j':
			if row < len(items)-1 {
				dirty[row], dirty[row+1] = true, true
				row++
				if row-offset >= height {
					offset++
					for i := 0; i < len(dirty); i++ {
						dirty[i] = true
					}
				}
			}
		case 'k':
			if row > 0 {
				dirty[row], dirty[row-1] = true, true
				row--
				if row < offset {
					offset--
					for i := 0; i < len(dirty); i++ {
						dirty[i] = true
					}
				}
			}
		case 'g':
			if row > 0 {
				dirty[row], dirty[0] = true, true
				row = 0
				if row < offset {
					offset = 0
					for i := 0; i < len(dirty); i++ {
						dirty[i] = true
					}
				}
			}
		case 'G':
			if row < len(items)-1 {
				dirty[row], dirty[len(items)-1] = true, true
				row = len(items) - 1
				if row-offset >= height {
					offset = len(items) - height
					for i := 0; i < len(dirty); i++ {
						dirty[i] = true
					}
				}
			}
		case vk_return:
			index = row
			label = items[row]
			break EVENT_LOOP
		case vk_esc, 'q':
			break EVENT_LOOP
		case '':
			sig <- syscall.SIGINT
		}
	}
	return index, label
}

func getAnsiColor(fg, bg string) string {
	return "\033[" + fg + ";" + bg + "m"
}

func eraseLine(p int) string {
	return fmt.Sprintf("\033[%dK", p)
}

func eraseDisplay(p int) string {
	return fmt.Sprintf("\033[%dJ", p)
}

func cursorPreviousLine(n int) string {
	return fmt.Sprintf("\033[%dF", n)
}

func cursorNextLine(n int) string {
	return fmt.Sprintf("\033[%dE", n)
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func truncate(s string, l int) string {
	return runewidth.Truncate(s, l, "...")
}
