[![Build & Test](https://github.com/idursun/jjui/actions/workflows/go.yml/badge.svg)](https://github.com/idursun/jjui/actions/workflows/go.yml)

# Jujutsu UI

`jjui` is a terminal user interface for working with [Jujutsu version control system](https://github.com/jj-vcs/jj). I have built it according to my own needs and will keep adding new features as I need them. I am open to feature requests and contributions.

## Features

Currently, you can:

### Change revset with auto-complete
You can change revset while enjoying auto-complete and signature help while typing.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_revset.gif)

### Rebase
You can rebase a revision or a branch onto another revision in the revision tree.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_rebase.gif)

See [Rebase](https://github.com/idursun/jjui/wiki/Rebase) for detailed information.

### Squash
You can squash revisions into one revision, by pressing `S`. The following revision will be automatically selected. However, you can change the selection by using `j` and `k`.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_squash.gif)

### Show revision details

Pressing `l` (as in going right into details of a revision) will open the details view of the revision you selected.

In this mode, you can:
- Split selected files using `s`
- Restore selected files using `r`
- View diffs of the highlighted by pressing `d`

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_details_diff.gif)

For detailed information, see [Details](https://github.com/idursun/jjui/wiki/Details) wiki page.

### Abandon
You can abandon a revision.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_abandon.gif)

### Bookmarks
You can move bookmarks to the revision you selected.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_bookmarks.gif)

### Split
You can split revisions.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_split.gif)

### Preview
You can open the preview window by pressing `p`. Preview window displays output of the `jj show` command of the selected revision. If the currenlty selected item is a file, then `jj diff` output is displayed. 

While the preview window is showing, you can press; `ctrl+n` to scroll one line down, `ctrl+p` to scroll one line up, `ctrl+n` to scroll half page down, `ctrl+u` to scroll half page up. 

Additionally, you can press `tab` to focus in and out of the preview window. Once in the focus mode, you can use normal (e.g. `j`, `k`, `d`, `u`) navigation keys as they are not bound to the revision tree view. 

For detailed information, see [Preview](https://github.com/idursun/jjui/wiki/Preview) wiki page.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_preview.gif)

### Diffs
You can see diffs of revisions.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_diff.gif)

Additionally,
* Edit description of a revision by pressing `D`
* Create a _new_ revision by pressing `n`
* _Edit_ a revision by pressing `e`
* Git _push_/_fetch_ by pressing `g`, followed by `p` or `f`
* Undo last change by pressing `u`

## Installation

### Nix

You can install `jjui` using nix from the unstable channel.

```shell
nix-env -iA nixpkgs.jjui
```

### From source

You can build `jjui` from source.

```
git clone https://github.com/idursun/jjui.git
cd jjui
go install cmd/jjui/main.go
```

### From go install
```
go install github.com/idursun/jjui/cmd/jjui@latest
```

### From pre-built binaries
You can download pre-built binaries from the [releases](https://github.com/idursun/jjui/releases) page.

## Compatibility

It's compatible with jj **v0.21**+.

## Contributing

Feel free to submit a pull request.
