package argparse

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/escribano/termbox-go"
)

// extractOptions will extract all options from the slice of arguments provided,
// returning one slice of invididual options, and a slice for all other arguments
// present.
func extractOptions(allArgs ...string) (options, args []string) {
	count := 0
	max := len(allArgs)

	for count < max {
		a := allArgs[count]

		// If we have option-escape string, assume the next arg is supposed
		// to be normal text instead of potentially being a option.
		if a == "--" && len(allArgs) > count+1 {
			args = append(args, allArgs[count+1])
			count = count + 2
			continue
		}

		// Using a option regex, check if we have a normal param or a option.
		optionRegex := regexp.MustCompile(`^-{1,2}[a-zA-Z]+$`)
		if !optionRegex.MatchString(a) {
			args = append(args, a)
			count++
			continue
		}

		// Okay, we must have a option. Which type?
		isShort := true
		if len(a) > 2 && a[:2] == "--" {
			isShort = false
		}

		// If short-option, grab all letters individual options.
		if isShort == true {
			for _, c := range a[1:] {
				options = append(options, string(c))
			}
		} else {
			options = append(options, a[2:])
		}
		count++
	}

	return options, args
}

// getScreenWidth returns the width of the screen the program is executed within.
func getScreenWidth() int {
	if err := termbox.Init(); err != nil {
		panic(err) // TODO: This should really be made to return an error.
	}
	w, _ := termbox.Size()
	termbox.Close()

	return w
}

// join will join the provided strings by the specified delimiter. The delimiter
// does not have to be limited to a single character; any string can be a delimiter.
func join(delimiter string, args ...string) string {
	var join bytes.Buffer
	num := len(args)

	if num == 0 {
		return ""
	}

	for index, val := range args {
		join.WriteString(val)
		if index < num-1 {
			join.WriteString(delimiter)
		}
	}

	return join.String()
}

// spacer provides a string containing only space-characters of the
// exact number specified.
func spacer(length int) string {
	count := 0
	char := " "
	var buff bytes.Buffer

	for count < length {
		buff.WriteString(char)
		count++
	}

	return buff.String()
}

// wordWrap breaks the provided string down into an array of strings with
// character-counts not exceeding the specified max length.
func wordWrap(text string, max int) []string {
	var lines []string
	var line []string

	if len(text) <= max {
		return []string{text}
	}

	split := strings.Split(text, " ")
	length := 0

	if len(split) <= 1 {
		return split
	}

	for _, word := range split {
		if len(word)+length+len(line) > max {
			lines = append(lines, join(" ", line...))
			line = []string{word}
			length = len(word)
		} else {
			length = length + len(word)
			line = append(line, word)
		}
	}
	lines = append(lines, join(" ", line...))

	return lines
}
