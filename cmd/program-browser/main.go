package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/go-redis/redis"
	"github.com/m1dugh/program-browser/internal/bugcrowd"
	"github.com/m1dugh/program-browser/internal/config"
	pbTypes "github.com/m1dugh/program-browser/pkg/types"
	"gopkg.in/yaml.v3"
)

type programFoundCB func(pbTypes.Program)

func runBugcrowd(cfg config.InputConfig, cb programFoundCB) error {

	log.Println("Starting bugcrowd fetching")

	opts := bugcrowd.Options{
		Filters: cfg.Filters,
	}
	api := bugcrowd.NewBugcrowdApi(&opts)

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

		go cb(*prog)
	}

	return nil
}

func prepareCallback(cfg config.OutputConfig) (programFoundCB, error) {

	var allCallbacks []func(pbTypes.Program) error

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
		allCallbacks = append(allCallbacks, cb)
	}

	if cfg.File != nil {
		var err error
		var target *os.File = os.Stdout

		if cfg.File.Target != "" {
			if target, err = os.Create(cfg.File.Target); err != nil {
				return nil, err
			}
		}

		var marshaller func(v any) ([]byte, error) = nil

		switch cfg.File.Format {
		case "json":
			marshaller = json.Marshal
		case "yaml":
			marshaller = yaml.Marshal
		default:
			return nil, fmt.Errorf("Unknown format %s", cfg.File.Format)
		}

		cb := func(prog pbTypes.Program) error {
			res, err := marshaller(prog)
			if err != nil {
				return err
			}
			if _, err = target.Write(res); err != nil {
				return err
			}
			_, err = target.WriteString("\n")
			return err
		}
		allCallbacks = append(allCallbacks, cb)
	}

	var merged programFoundCB = func(prog pbTypes.Program) {

		for _, cb := range allCallbacks {
			if err := cb(prog); err != nil {
				log.Print(err)
			}
		}
	}

	return merged, nil
}

func runProviders(cfg config.InputConfig, cb programFoundCB) error {

	var err error = nil
	var mut sync.Mutex
	var wg sync.WaitGroup

	if cfg.Bugcrowd != nil && cfg.Bugcrowd.Enable {
		wg.Go(func() {
			localErr := runBugcrowd(cfg, cb)
			if err == nil {
				return
			}

			mut.Lock()
			err = errors.Join(err, localErr)
			mut.Unlock()
		})
	}

	if len(cfg.ExtraEntries) > 0 {
		wg.Go(func() {
			for _, prog := range cfg.ExtraEntries {
				cb(prog)
			}
		})
	}

	wg.Wait()

	return err
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

	cb, err := prepareCallback(cfg.Output)
	if err != nil {
		log.Panic(err)
	}

	err = runProviders(cfg.Input, cb)
	if err != nil {
		log.Println(err)
	}
}
