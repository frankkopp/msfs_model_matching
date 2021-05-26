package main

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	out = message.NewPrinter(language.German)
)

func main() {
	out.Println("Hello World")

}
