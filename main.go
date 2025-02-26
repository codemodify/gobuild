package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "gen":
			gen()
		case "ver":
			log.Println("version: ", version)
		default:
			log.Println("unknown command")
		}

		return
	}

	build()
}

func exitIfErr(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
