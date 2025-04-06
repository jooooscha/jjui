package screen

import (
	"fmt"
	"strconv"
	"strings"
)

type Segment struct {
	Text     string
	Params   string
	Reversed bool
}

func splitString(str, searchString string) []string {
	index := strings.Index(str, searchString)
	if index == -1 {
		return []string{str}
	}

	before := str[:index]
	after := str[index+len(searchString):]

	return []string{before, searchString, after}
}

func (s Segment) Reverse(text string) []*Segment {
	ret := make([]*Segment, 0)
	for _, part := range splitString(s.Text, text) {
		if part == "" {
			continue
		}
		ret = append(ret, &Segment{
			Text:     part,
			Params:   s.Params,
			Reversed: part == text,
		})
	}
	return ret
}

func (s Segment) String() string {
	if s.Text == "\n" {
		return s.Text
	}
	if s.Params == "" {
		return s.Text
	}
	if s.Reversed {
		return fmt.Sprintf("\x1b[%sm\x1b[7m%s\x1b[0m", s.Params, s.Text)
	}
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", s.Params, s.Text)
}

func (s Segment) WithBackground(bg string) *Segment {
	var newParts []string
	parts := strings.Split(s.Params, ";")

	i := 0
	for i < len(parts) {
		part := parts[i]
		num, err := strconv.Atoi(part)
		if err != nil {
			i++
			continue
		}
		p := num

		isBg := false
		if (p >= 40 && p <= 49) || (p >= 100 && p <= 109) {
			isBg = true
		}

		if !isBg {
			newParts = append(newParts, part)
		}
		i++
	}

	for _, part := range strings.Split(bg, ";") {
		if _, err := strconv.Atoi(part); err != nil {
			panic(fmt.Sprintf("invalid background parameter %q", part))
		}
		newParts = append(newParts, part)
	}

	return &Segment{
		Text:   s.Text,
		Params: strings.Join(newParts, ";"),
	}
}

func (s Segment) StyleEqual(other Segment) bool {
	return s.Params == other.Params
}

// BreakNewLinesIter group segments into lines by breaking segments at new lines
func BreakNewLinesIter(rawSegments <-chan *Segment) <-chan []*Segment {
	output := make(chan []*Segment)
	go func() {
		defer close(output)
		currentLine := make([]*Segment, 0)
		for rawSegment := range rawSegments {
			idx := strings.IndexByte(rawSegment.Text, '\n')
			for idx != -1 {
				text := rawSegment.Text[:idx]
				if len(text) > 0 {
					currentLine = append(currentLine, &Segment{
						Text:   text,
						Params: rawSegment.Params,
					})
				}
				output <- currentLine
				currentLine = make([]*Segment, 0)
				rawSegment.Text = rawSegment.Text[idx+1:]
				idx = strings.IndexByte(rawSegment.Text, '\n')
			}
			if len(rawSegment.Text) > 0 {
				currentLine = append(currentLine, rawSegment)
			}
		}
		if len(currentLine) > 0 {
			output <- currentLine
		}
	}()
	return output
}
