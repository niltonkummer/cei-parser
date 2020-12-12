package main

import (
	"cei_parser/parser"
	"encoding/json"
	"fmt"
	"os"
)

func init() {

}

func main() {

	if len(os.Args) == 1 {
		fmt.Println("invalid name file param")
		os.Exit(1)
	}

	decoder, err := parser.DetectExtension(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	res, err := decoder.Decode()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	json.NewEncoder(os.Stdout).Encode(res)
}
