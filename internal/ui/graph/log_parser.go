package graph

import (
	"github.com/idursun/jjui/internal/screen"
	"io"
	"log"
	"unicode/utf8"
)

func ParseRows(reader io.Reader) []Row {
	var rows []Row
	var row Row
	rawSegments := screen.ParseFromReader(reader)

	for segmentedLine := range screen.BreakNewLinesIter(rawSegments) {
		rowLine := NewGraphRowLine(segmentedLine)
		if changeIdIdx := rowLine.FindIdIndex(0); changeIdIdx != -1 {
			rowLine.Flags = Revision | Highlightable
			previousRow := row
			row = NewGraphRow()
			if previousRow.Commit != nil {
				rows = append(rows, previousRow)
				row.Previous = &previousRow
			}
			for j := 0; j < changeIdIdx; j++ {
				row.Indent += utf8.RuneCountInString(rowLine.Segments[j].Text)
			}
			rowLine.ChangeIdIdx = changeIdIdx
			row.Commit.ChangeIdShort = rowLine.Segments[changeIdIdx].Text
			row.Commit.ChangeId = row.Commit.ChangeIdShort + rowLine.Segments[changeIdIdx+1].Text
			commitIdIdx := rowLine.FindIdIndex(changeIdIdx + 2)
			if commitIdIdx != -1 {
				rowLine.CommitIdIdx = commitIdIdx
				row.Commit.CommitIdShort = rowLine.Segments[commitIdIdx].Text
				row.Commit.CommitId = row.Commit.CommitIdShort + rowLine.Segments[commitIdIdx+1].Text
			} else {
				log.Fatalln("commit id not found")
			}
		}
		row.AddLine(&rowLine)
	}
	if row.Commit != nil {
		rows = append(rows, row)
	}
	return rows
}
