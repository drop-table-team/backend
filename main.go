package main

import (
	"backend/module"
	"backend/util"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
)

func main() {
	moduleConfigPath := util.MaybeEnv("MODULE_CONFIG_PATH")

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// register interrupt handler to handle keyboard interrupts (ctrl-c)
	initInterruptHandler()

	moduleManager, err := initModules(moduleConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	stopFunctions = append(stopFunctions, func() { _ = moduleManager.StopAll() })

	<-make(chan int)
}

// contains functions that should run before the program stops
var stopFunctions []func()

func initInterruptHandler() {
	c := make(chan os.Signal, 3)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		log.Print("received interrupt signal, stopping gracefully")
		go func() {
			for _, stopFunction := range stopFunctions {
				stopFunction()
			}
			os.Exit(0)
		}()
		<-c
		<-c
		log.Fatal("received 3 interrupt signals, aborting immediately")
	}()
}

func initModules(moduleConfigPath *string) (*module.ModuleManager, error) {
	var config module.ModuleConfig
	if moduleConfigPath == nil {
		log.Print("no config file given")
	} else {
		// check if module config file is valid
		if _, err := os.Stat(*moduleConfigPath); err != nil && errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("module config file doesn't exist")
		} else if err != nil {
			return nil, err
		}

		configFileContent, err := os.ReadFile(*moduleConfigPath)
		if err != nil {
			return nil, err
		}
		config, err := module.ParseServiceConfig(configFileContent)
		if err != nil {
			return nil, err
		}
		log.Printf("parsed module config: %v", string(util.UnwrapError(json.Marshal(config))))
	}

	module, err := module.NewModuleManager(config)
	if err != nil {
		return nil, err
	}

	return &module, module.StartAll()
}
