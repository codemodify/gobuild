package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "gen":
			gen()

		case "gen-rm":
			genRM()

		case "clean":
			buildClean()

		case "version":
			fmt.Println("version:", version)
			os.Exit(0)

		default:
			fmt.Println("usage: ", os.Args[0], "[COMMAND], no COMMAND builds it")
			fmt.Println("COMMAND is one of")
			fmt.Println("\tgen - generate config files")
			fmt.Println("\tgen-rm - remove generated config files")
			fmt.Println("\tclean - clean the build")
			fmt.Println("\tversion - print tool version")
			os.Exit(1)
		}

		return
	}

	build()
}

func exitIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
