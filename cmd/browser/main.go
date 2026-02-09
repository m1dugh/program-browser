package main

import (
	"fmt"
	"log"

	pbTypes "github.com/m1dugh/program-browser/pkg/types"
	"github.com/m1dugh/program-browser/internal/bugcrowd"
)

func runBugcrowd() error {

	log.Println("Starting bugcrowd fetching")

	api := bugcrowd.NewBugcrowdApi()

	progs, err := api.FetchPrograms()
	if err != nil {
		return err
	}

	var prog *pbTypes.Program
	for {
		prog = <- progs
		if prog == nil {
			break
		}
		fmt.Println(prog.Name)
		fmt.Println("allowed")
		for _, entry := range prog.Scope.AllowedEndpoints {
			fmt.Println(entry)
		}
		fmt.Println("denied")
		for _, entry := range prog.Scope.DeniedEndpoints {
			fmt.Println(entry)
		}
		fmt.Println()
	}

	return nil
}

func main() {
	log.SetFlags(log.LstdFlags)

	err := runBugcrowd()

	if err != nil {
		log.Println(err)
	}
}
