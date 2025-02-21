# Text Editor

In progress.

This is a terminal-based application made with Go using the
[tcell](github.com/gdamore/tcell/v2) library. It is inspired by the vim text
editor. It was created as a fun learning project.

## Run

```
go run . [<filename1> [<filename2>...] ]
```

## Insert mode
- `<esc>`: exit mode

## Command mode
- `<esc>`: exit mode without running command
- `<enter>`: exit mode and run command
- `:<line>`: go to line
- `:next`: open next file
- `:prev`: open prev file
- `:open <filename>`: open and load file if exists
- `:close`: close current file
- `:closeall`: close all files
- `:ls`: print open files to log file - TODO show files in a popup

## Normal mode
- `hjkl`: move left/down/up/right
- `0`: jump to start of line
- `$`: jump to end of line
- `f<char>`: jump to next occurrence of character
- `i`: change to insert mode
- `:`: change to command mode

## Tasks
- DONE feat: open and edit multiple files
- DONE feat: insert utf8 text
- DONE feat: auto-scroll when jumping cursor goes off-screen
- DONE feat: dynamic size for boxes
- DONE feat: fix ui window sizes 
- DONE feat: show line numbers
- DONE core: make separate status line and command line windows
- DONE core: open temporary file for empty buffer
- TODO feat: insert mode - append, CTRL+keys
- TODO core: insert new lines with enter
- TODO object motions - words/lines
- TODO command mode - run external commands on current line
- TODO command mode - run external commands on entire file
- TODO feat: add visual mode
- TODO impl popup to show longer info messages to user
- TODO keep user command history
- TODO improv: add colors and enhance UI
