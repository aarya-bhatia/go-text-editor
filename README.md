# Text Editor

In progress.

This is a terminal-based application made with Go using the
[tcell](github.com/gdamore/tcell/v2) library. It is inspired by the vim text
editor. It was created as a fun learning project.

## Build

```
go build .
```

## Run

```
./go-editor [<filename1> [<filename2>...] ]
```

## Insert mode
- `<esc>`: exit mode

## Command mode
- `<esc>`: exit mode without running command
- `<enter>`: exit mode and run command
- `:<line>`: go to line
- `:next`: open next file
- `:prev`: open prev file

## Normal mode
- `hjkl`: move left/down/up/right
- `0`: jump to start of line
- `$`: jump to end of line
- `f<char>`: jump to next occurrence of character
- `i`: change to insert mode
- `:`: change to command mode

## Feature List
- DONE open and edit multiple files
- DONE insert utf8 text
- DONE auto-scroll when jumping cursor goes off-screen
- TODO fix ui window sizes 
- TODO make separate status line and command line windows
- TODO handle screen refresh
- TODO unit tests
- TODO insert mode - append, CTRL+keys
- TODO insert new lines with enter
- TODO object motions - words/lines
- TODO command mode - run external commands on current line
- TODO command mode - run external commands on entire file
- TODO add visual mode

