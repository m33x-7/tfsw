module github.com/m33x-7/tfsw

go 1.17

require tfsw/cmd v0.0.1

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jedib0t/go-pretty/v6 v6.4.3 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/schollz/progressbar/v3 v3.8.2 // indirect
	github.com/spf13/cobra v1.2.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/term v0.2.0 // indirect
	tfsw/internal/utils v0.0.1 // indirect
)

replace (
	tfsw/cmd => ./cmd
	tfsw/internal/utils => ./internal/utils
)
