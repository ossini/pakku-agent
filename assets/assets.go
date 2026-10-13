package assets

import (
	_ "embed"

	"github.com/fatih/color"
)

//go:embed pakku.txt
var pakku string

func PrintBanner() {
	color.Yellow(pakku)
}