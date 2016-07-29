package main

import (
	"fmt"
	"github.com/nocd5/choru"
)

func main() {
	fmt.Println("--------------------------------")
	fmt.Println("j:down / k:up / g:top / G:bottom")
	fmt.Println("--------------------------------")

	c1 := choru.New()
	if i, v := c1.Choose([]string{"abcde", "fghij", "klmno", "pqrst", "uvwxyz"}); i >= 0 {
		fmt.Println(fmt.Sprintf("%d : %s", i, v))
	}

	items := []string{
		"0-0 FOO", "0-1 bar", "0-2 hoge", "0-3 fuga", "0-4 moge",
		"1-0 foo", "1-1 BAR", "1-2 hoge", "1-3 fuga", "1-4 moge",
		"2-0 foo", "2-1 bar", "2-2 HOGE", "2-3 fuga", "2-4 moge",
		"3-0 foo", "3-1 bar", "3-2 hoge", "3-3 FUGA", "3-4 moge",
		"4-0 foo", "4-1 bar", "4-2 hoge", "4-3 fuga", "4-4 MOGE",
		`"5-0 FOO", "5-1 BAR", "5-2 HOGE", "5-3 FUGA", "5-4 MOGE"`,
	}
	c2 := choru.New()
	c2.LineFg = choru.FgWhiteBold
	c2.CursorFg = choru.FgBlack
	c2.CursorBg = choru.BgYellowBright
	c2.MaxHeight = 10
	c2.Header = "==================== Header ===================="
	c2.HeaderFg = choru.FgBlack
	c2.HeaderBg = choru.BgGreenBright
	c2.Footer = "==================== Footer ===================="
	c2.FooterFg = choru.FgBlack
	c2.FooterBg = choru.BgRedBright
	if i, v := c2.Choose(items); i >= 0 {
		fmt.Println(fmt.Sprintf("%d : %s", i, v))
	}
}
