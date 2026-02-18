package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/m1dugh/program-browser/internal/bugcrowd"
	"github.com/m1dugh/program-browser/internal/config"
	pbTypes "github.com/m1dugh/program-browser/pkg/types"
)

type programFoundCB func(pbTypes.Program) error

func runBugcrowd(cb programFoundCB) error {

	log.Println("Starting bugcrowd fetching")

	api := bugcrowd.NewBugcrowdApi()

	progs, err := api.FetchPrograms()
	if err != nil {
		return err
	}

	var prog *pbTypes.Program
	for {
		prog = <-progs
		if prog == nil {
			break
		}

		go func() {
			if err := cb(*prog); err != nil {
				log.Print(err)
			}
		}()
	}

	return nil
}

func prepareCallback(cfg config.Config) (programFoundCB, error) {

	if cfg.Redis != nil {
		client := redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password(),
			DB:       cfg.Redis.DB,
		})

		cb := func(prog pbTypes.Program) error {
			val, err := json.Marshal(prog)
			if err != nil {
				return err
			}

			err = client.RPush(cfg.Redis.Name, string(val)).Err()

			if err != nil {
				return err
			}

			return nil
		}
		return cb, nil
	}

	return func(prog pbTypes.Program) error {
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

		return nil
	}, nil
}

func main() {
	log.SetFlags(log.LstdFlags)

	var configFile string
	flag.StringVar(&configFile, "config", "", "The path to the config")

	flag.Parse()

	cfg, err := config.NewConfig(configFile)

	if err != nil {
		log.Panic(err)
	}

	cb, err := prepareCallback(cfg)
	if err != nil {
		log.Panic(err)
	}

	err = runBugcrowd(cb)

	if err != nil {
		log.Println(err)
	}
}
