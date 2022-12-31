package main

import "embed"

//go:embed swaggerui
var swaggeruifs embed.FS

//go:embed swagger.yml
var swaggerContent string
