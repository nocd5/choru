package main

import (
	"fmt"
	"github.com/nocd5/choru"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"
)

const (
	animal_words_url = "https://raw.githubusercontent.com/hzlzh/Domain-Name-List/master/Animal-words.txt"
)

func main() {
	words := generateWordList()

	c := choru.New()
	c.MaxHeight = 10
	c.HeaderFg = choru.FgGreen
	c.CursorFg = choru.FgBlack
	c.CursorBg = choru.BgYellowBright
	c.Footer = fmt.Sprintf("--------------------")
	c.FooterFg = choru.FgGreen

	rand.Seed(time.Now().UnixNano())
	letter := 'a' + rune(rand.Intn(26))

	score := 0
	for {
		var l []string
		for _, v := range words {
			if v == "" {
				continue
			}
			if unicode.ToLower(rune(v[0])) == letter {
				l = append(l, v)
			}
		}
		if len(l) > 0 {
			c.Header = fmt.Sprintf("--- Score : %d ---", score)
			if i, v := c.Choose(l); i >= 0 {
				letter = rune(v[len(v)-1])
				words = remove(words, v)
				score += len(v)
			} else {
				break
			}
		} else {
			fmt.Println("< GAME OVER >")
			fmt.Printf("    There is no words beginning letter the '%c'.\n", unicode.ToUpper(letter))
			break
		}
	}
	fmt.Printf("    Score : %d\n", score)
}

func generateWordList() []string {
	var words []string
	if resp, err := http.Get(animal_words_url); err == nil {
		defer resp.Body.Close()
		content, _ := ioutil.ReadAll(resp.Body)
		words = strings.Split((string(content)), "\n")
	} else {
		fmt.Fprint(os.Stderr, err)
	}
	return words
}

func remove(l []string, e string) []string {
	var r []string
	for _, v := range l {
		if v != e {
			r = append(r, v)
		}
	}
	return r
}
