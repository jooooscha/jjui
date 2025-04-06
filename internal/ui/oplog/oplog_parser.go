package oplog

import (
	"github.com/idursun/jjui/internal/screen"
	"io"
)

func ParseRows(reader io.Reader) []Row {
	var rows []Row
	var row Row
	rawSegments := screen.ParseFromReader(reader)

	for segmentedLine := range screen.BreakNewLinesIter(rawSegments) {
		rowLine := NewRowLine(segmentedLine)
		if opIdIdx := rowLine.FindIdIndex(); opIdIdx != -1 {
			if row.OperationId != "" {
				rows = append(rows, row)
			}
			row = Row{OperationId: rowLine.Segments[opIdIdx].Text}
		}
		row.Lines = append(row.Lines, &rowLine)
	}
	rows = append(rows, row)
	return rows
}
