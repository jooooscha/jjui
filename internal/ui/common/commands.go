package common

import (
	"github.com/charmbracelet/bubbletea"
	"jjui/internal/jj"
	"strings"
)

type Commands struct {
	jj jj.Commands
}

func (c Commands) GitFetch() tea.Cmd {
	f := func() tea.Msg {
		output, err := c.jj.GitFetch().CombinedOutput()
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
	return tea.Sequence(CommandRunning("jj git fetch"), f)
}

func (c Commands) GitPush() tea.Cmd {
	f := func() tea.Msg {
		output, err := c.jj.GitPush().CombinedOutput()
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
	return tea.Sequence(CommandRunning("jj git push"), f)
}

func (c Commands) Rebase(from, to string, operation Operation) tea.Cmd {
	rebase := c.jj.RebaseCommand
	if operation == RebaseBranchOperation {
		rebase = c.jj.RebaseBranchCommand
	}
	cmd := rebase(from, to)
	return ShowOutput(cmd)
}

func (c Commands) SetDescription(revision string, description string) tea.Cmd {
	return ShowOutput(c.jj.SetDescription(revision, description))
}

func (c Commands) MoveBookmark(revision string, bookmark string) tea.Cmd {
	return ShowOutput(c.jj.MoveBookmark(revision, bookmark))
}

func (c Commands) DeleteBookmark(bookmark string) tea.Cmd {
	return ShowOutput(c.jj.DeleteBookmark(bookmark))
}

func (c Commands) FetchRevisions(revset string) tea.Cmd {
	return func() tea.Msg {
		graphLines, err := c.jj.GetCommits(revset)
		if err != nil {
			return UpdateRevisionsFailedMsg(err)
		}
		return UpdateRevisionsMsg(graphLines)
	}
}

func (c Commands) FetchBookmarks(revision string, op Operation) tea.Cmd {
	return func() tea.Msg {
		cmd := c.jj.ListBookmark(revision)
		//TODO: handle error
		output, _ := cmd.CombinedOutput()
		var bookmarks []string
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			bookmarks = append(bookmarks, line)
		}
		return UpdateBookmarksMsg{
			Bookmarks: bookmarks,
			Revision:  revision,
			Operation: op,
		}
	}
}

func (c Commands) SetBookmark(revision string, name string) tea.Cmd {
	return ShowOutput(c.jj.SetBookmark(revision, name))
}

func (c Commands) GetDiff(revision string) tea.Cmd {
	return func() tea.Msg {
		output, _ := c.jj.Diff(revision).CombinedOutput()
		return ShowDiffMsg(output)
	}
}

func (c Commands) Edit(revision string) tea.Cmd {
	return func() tea.Msg {
		output, err := c.jj.Edit(revision).CombinedOutput()
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
}

func (c Commands) DiffEdit(revision string) tea.Cmd {
	return tea.ExecProcess(c.jj.DiffEdit(revision), func(err error) tea.Msg {
		return RefreshMsg{SelectedRevision: revision}
	})
}

func (c Commands) Split(revision string) tea.Cmd {
	return tea.ExecProcess(c.jj.Split(revision), func(err error) tea.Msg {
		return RefreshMsg{SelectedRevision: revision}
	})
}

func (c Commands) Abandon(revision string) tea.Cmd {
	return func() tea.Msg {
		output, err := c.jj.Abandon(revision).CombinedOutput()
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
}

func (c Commands) NewRevision(from string) tea.Cmd {
	return ShowOutput(c.jj.New(from))
}

func NewCommands(jj jj.Commands) Commands {
	return Commands{jj}
}
