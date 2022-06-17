package frontend

import "embed"

//go:embed dist/*
var Data embed.FS

//go:embed style.css
var CSS []byte
