package template

import "embed"

//go:embed *.html

var Template embed.FS

//go:embed assets/logo.png
var Logo []byte