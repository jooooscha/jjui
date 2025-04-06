package test

import (
	"bufio"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/ui/graph"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	file, _ := os.Open("testdata/output.log")
	rows := graph.ParseRows(file)
	assert.Len(t, rows, 11)
}

func TestParser_Parse_NoCommitId(t *testing.T) {
	file, _ := os.Open("testdata/no-commit-id.log")
	rows := graph.ParseRows(file)
	assert.Len(t, rows, 1)
}

func TestParser_Parse_ShortId(t *testing.T) {
	file, _ := os.Open("testdata/short-id.log")
	rows := graph.ParseRows(file)
	assert.Len(t, rows, 2)
	assert.Equal(t, "X", rows[0].Commit.ChangeId)
	assert.Equal(t, "E", rows[0].Commit.CommitId)
	assert.Equal(t, "T", rows[1].Commit.ChangeId)
	assert.Equal(t, "79", rows[1].Commit.CommitId)
}

func TestParser_Parse_Disconnected(t *testing.T) {
	var lb logBuilder
	lb.write("*   id=abcde author=some@author id=xyrq")
	lb.write("│   some documentation")
	lb.write("~\n")
	lb.write("*   id=abcde author=some@author id=xyrq")
	lb.write("│   another commit")
	lb.write("~\n")
	rows := graph.ParseRows(strings.NewReader(lb.String()))
	assert.Len(t, rows, 2)
}

func TestParser_Parse_Extend(t *testing.T) {
	var lb logBuilder
	lb.write("*   id=abcde author=some@author id=xyrq")
	lb.write("│   some documentation")

	rows := graph.ParseRows(strings.NewReader(lb.String()))
	assert.Len(t, rows, 1)
	row := rows[0]

	extended := row.Lines[1].Extend(row.Indent)
	assert.Len(t, extended.Segments, 1)
}

func TestParser_Parse_WorkingCopy(t *testing.T) {
	var lb logBuilder
	lb.write("*   id=abcde author=some@author id=xyrq")
	lb.write("│   some documentation")
	lb.write("@   id=kdys author=some@author id=12cd")
	lb.write("│   some documentation")

	rows := graph.ParseRows(strings.NewReader(lb.String()))
	assert.Len(t, rows, 2)
	row := rows[1]

	assert.True(t, row.Commit.IsWorkingCopy)
}

type part int

const (
	normal = iota
	id
	author
	bookmark
)

var styles = map[part]lipgloss.Style{
	normal:   lipgloss.NewStyle(),
	id:       lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
	author:   lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
	bookmark: lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
}

type logBuilder struct {
	w strings.Builder
}

func (l *logBuilder) String() string {
	return l.w.String()
}

func (l *logBuilder) write(line string) {
	lipgloss.SetColorProfile(termenv.ANSI)
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "short_id=") {
			text = strings.TrimPrefix(text, "short_id=")
			l.shortId(text)
			continue
		}
		if strings.HasPrefix(text, "id=") {
			text = strings.TrimPrefix(text, "id=")
			l.id(text[:1], text[1:])
			continue
		}
		if strings.HasPrefix(text, "author=") {
			l.author(strings.TrimPrefix(text, "author="))
			continue
		}
		if strings.HasPrefix(text, "bookmarks=") {
			text = strings.TrimPrefix(text, "bookmarks=")
			values := strings.Split(text, ",")
			l.bookmarks(strings.Join(values, " "))
			continue
		}
		l.append(text)
	}
	l.w.WriteString("\n")
}

func (l *logBuilder) append(value string) {
	fmt.Fprintf(&l.w, "%s ", styles[normal].Render(value))
}

func (l *logBuilder) shortId(sid string) {
	fmt.Fprintf(&l.w, " %s ", styles[id].Render(sid))
}

func (l *logBuilder) id(short string, rest string) {
	fmt.Fprintf(&l.w, " %s%s ", styles[id].Render(short), styles[id].Render(rest))
}

func (l *logBuilder) author(value string) {
	fmt.Fprintf(&l.w, " %s ", styles[author].Render(value))
}

func (l *logBuilder) bookmarks(value string) {
	fmt.Fprintf(&l.w, " %s ", styles[bookmark].Render(value))
}
