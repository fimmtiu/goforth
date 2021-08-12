package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	log.Println("Hello World!")
	parser := NewParser(bufio.NewReader(os.Stdin))
	fmt.Println(parser.NextToken())
}
