package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"unicode"
	"unicode/utf8"
)

var usage = "%s [OPTION]... [FILE]...\n"

var (
	lines     = flag.Bool("l", false, "counts lines")
	help      = flag.Bool("h", false, "display this message")
	words     = flag.Bool("w", false, "counts words")
	chars     = flag.Bool("c", false, "counts characters")
	maxLength = flag.Bool("L", false, "print the length of the longest line")

	inField = false
)

func main() {
	flag.Parse()

	if *help {
		fmt.Printf(usage, os.Args[0])
		flag.PrintDefaults()
		return
	}

	var err error
	var f string

	buf := make([]byte, 0)

	if len(flag.Args()) > 0 {
		f = flag.Arg(0)

		buf, err = ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
	} else {

		buf, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
	}

	var r func(rune) bool

	if *words {
		r = isWord
	} else if *lines {
		r = isNewLine
	} else if *chars {
		r = isAll
	} else if *maxLength {
		result := getLongestLineLength(buf)
		fmt.Printf("%d %s\n", result, f)
		return
	} else {
		chars := decode(buf, isAll)
		words := decode(buf, isWord)
		lines := decode(buf, isNewLine)

		fmt.Printf("%d %d %d %s\n", lines, words, chars, f)
		return
	}

	results := decode(buf, r)

	fmt.Printf("%d %s\n", results, f)
}

func decode(d []byte, f func(rune) bool) int {
	var count int

	for len(d) > 0 {
		r, size := utf8.DecodeRune(d)
		if f(r) {
			count++
		}

		d = d[size:]
	}

	return count
}

func isWord(c rune) bool {
	wasInField := inField
	inField = !unicode.IsSpace(c)
	if inField && !wasInField {
		return true
	}

	return false
}

func isNewLine(c rune) bool { return c == '\n' }
func isAll(c rune) bool     { return true }

func getLongestLineLength(f []byte) int {
	fields := bytes.FieldsFunc(f, isNewLine)

	lens := make([]int, len(fields))
	for i, l := range fields {
		lens[i] = len(l)
	}

	sort.Ints(lens)

	return lens[len(lens)-1]
}
