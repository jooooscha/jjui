package jj

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

const RootChangeId = "zzzzzzzz"

type Commit struct {
	ChangeIdShort string
	ChangeId      string
	Parents       []string
	IsWorkingCopy bool
	Author        string
	Timestamp     string
	Bookmarks     []string
	Description   string
	Immutable     bool
	Conflict      bool
	Empty         bool
	Index         int
}

func (c Commit) IsRoot() bool {
	return c.ChangeId == RootChangeId
}

func GetCommits(location string) *Dag {
	cmd := exec.Command("jj", "log", "--reversed", "--template", TEMPLATE)
	cmd.Dir = location
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	d := Parse(bytes.NewReader(output))
	return &d
}

func Parse(reader io.Reader) Dag {
	d := NewDag()
	all, err := io.ReadAll(reader)
	if err != nil {
		return d
	}
	lines := strings.Split(string(all), "\n")
	stack := make([]*Node, 0)
	stack = append(stack, nil)
	levels := make([]int, 0)
	levels = append(levels, -1)
	seen := make(map[string]bool)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if line == "" || line == "~" {
			continue
		}
		index := strings.IndexAny(line, "○◆@×")
		if index == -1 {
			continue
		}
		_, after, _ := strings.Cut(line[index:], " ")
		parts := strings.Split(after, ";")
		commit := Commit{
			ChangeIdShort: strings.TrimSpace(parts[0]),
		}
		seen[commit.ChangeIdShort] = true
		if len(parts) > 1 {
			commit.ChangeId = parts[1]
		}
		edgeType := DirectEdge
		if len(parts) > 2 {
			commit.Parents = strings.Split(parts[2], ",")
			for _, parent := range commit.Parents {
				if _, ok := seen[parent]; !ok {
					edgeType = IndirectEdge
				}
			}
		}
		if len(parts) > 3 && parts[3] != "." {
			commit.Bookmarks = strings.Split(parts[3], ",")
		}
		if len(parts) > 4 {
			commit.IsWorkingCopy = parts[4] == "true"
		}
		if len(parts) > 5 {
			commit.Immutable = parts[5] == "true"
		}
		if len(parts) > 6 {
			commit.Conflict = parts[6] == "true"
		}
		if len(parts) > 7 {
			commit.Empty = parts[7] == "true"
		}
		if len(parts) > 8 {
			commit.Author = parts[8]
		}
		if len(parts) > 9 {
			commit.Timestamp = parts[9]
		}
		if len(parts) > 10 {
			commit.Description = parts[10]
		}
		node := d.AddNode(&commit)
		if index < levels[len(levels)-1] {
			levels = levels[:len(levels)-1]
			stack = stack[:len(stack)-1]
		}
		if stack[len(stack)-1] != nil {
			stack[len(stack)-1].AddEdge(node, edgeType)
		}
		if index == levels[len(levels)-1] {
			stack[len(stack)-1] = node
		}
		if index > levels[len(levels)-1] {
			levels = append(levels, index)
			stack = append(stack, node)
		}
		if commit.ChangeId == RootChangeId {
			commit.Conflict = false
			commit.Parents = nil
			commit.Immutable = false
			commit.Author = ""
			commit.Bookmarks = nil
			commit.Description = ""
		}
	}
	return d
}
