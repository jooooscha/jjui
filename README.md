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
You can squash revisions into one revision, by pressing `S`. The following revision will be automatically selected. However, you can change the selection using `j` and `k`.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_squash.gif)

### Show revision details

Pressing `l` (as in going right into the details of a revision) will open the details view of the revision you selected.

In this mode, you can:
- Split selected files using `s`
- Restore selected files using `r`
- View diffs of the highlighted by pressing `d`

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_details.gif)

For detailed information, see [Details](https://github.com/idursun/jjui/wiki/Details) wiki page.

### Bookmarks
You can move bookmarks to the revision you selected.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_bookmarks.gif)


### Op Log
You can switch to op log view by pressing `o`. Pressing `r` restores the selected operation. For more information, see [Op log](https://github.com/idursun/jjui/wiki/Oplog) wiki page.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_oplog.gif)

### Preview
You can open the preview window by pressing `p`. If the selected item is a revision, then the output of `jj show` command is displayed. Similarly, `jj diff` output is displayed for selected files,  and `jj op show` output is displayed for selected operations. 

While the preview window is showing, you can press; `ctrl+n` to scroll one line down, `ctrl+p` to scroll one line up, `ctrl+d` to scroll half a page down, `ctrl+u` to scroll half a page up. 

Additionally, you can press `d` to show the contents of preview in diff view. 

For detailed information, see [Preview](https://github.com/idursun/jjui/wiki/Preview) wiki page.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_preview.gif)

Additionally,
* View the diff of a revision by pressing `d`.
* Edit the description of a revision by pressing `D`
* Create a _new_ revision by pressing `n`
* Split a revision by pressing `s`.
* Abandon a revision by pressing `a`.
* Absorb a revision by pressing `A`.
* _Edit_ a revision by pressing `e`
* Git _push_/_fetch_ by pressing `g` 
* Undo the last change by pressing `u`
* Show evolog of a revision by pressing `v`

## Configuration

See [configuration](https://github.com/idursun/jjui/wiki/Configuration) section in the wiki.

## Installation

### Nix

You can install `jjui` using nix from the unstable channel.

```shell
nix-env -iA nixpkgs.jjui
```

### From go install

To install the latest released (or pre-released) version:

```shell
go install github.com/idursun/jjui/cmd/jjui@latest
```

To install the latest commit from `main`:

```shell
go install github.com/idursun/jjui/cmd/jjui@HEAD
```
To install the latest commit from `main` bypassing the local cache:

```shell
GOPROXY=direct go install github.com/idursun/jjui/cmd/jjui@HEAD
```

### From source

You can build `jjui` from source.

```shell
git clone https://github.com/idursun/jjui.git
cd jjui
go install ./...
```


### From pre-built binaries
You can download pre-built binaries from the [releases](https://github.com/idursun/jjui/releases) page.

## Compatibility

Minimum supported `jj` version is **v0.21**+.

## Contributing

Feel free to submit a pull request.
