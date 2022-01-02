package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("FILE is a required argument")
		os.Exit(0)
	} else if len(os.Args) > 2 {
		fmt.Println("Too many arguments; 1 expected")
		os.Exit(1)
	}
	filepath := os.Args[1]
	contents := readFile(filepath)
	file := ReadPCapFile(contents)
	fmt.Println(len(file.Packets))
}

func readFile(filepath string) []byte {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return contents
}
