package main

//go:generate go generate github.com/sk4zuzu/topspin/micro/topspin/api
//go:generate go generate github.com/sk4zuzu/topspin/micro/topspin/srv1
//go:generate go generate github.com/sk4zuzu/topspin/micro/topspin/srv2

import (
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/sk4zuzu/topspin/micro/topspin/api"
	"github.com/sk4zuzu/topspin/micro/topspin/srv1"
	"github.com/sk4zuzu/topspin/micro/topspin/srv2"
)

type serviceDef struct {
	Run bool
}

var (
	serviceDefs = map[string]*serviceDef{
		"api": &serviceDef{
			Run: false,
		},
		"srv1": &serviceDef{
			Run: false,
		},
		"srv2": &serviceDef{
			Run: false,
		},
	}
)

func parseServiceDefs() (enabled int) {
	RUN, ok := os.LookupEnv("RUN")
	if !ok {
		log.Fatal("`RUN` envvar not found")
	}
	for _, key := range strings.Split(RUN, ",") {
		if serviceDef, ok := serviceDefs[key]; ok {
			serviceDef.Run = true
			enabled++
		}
	}
	return
}

func main() {
	count := parseServiceDefs()
	if ok := count > 0; !ok {
		log.Fatal("no microservice specified")
	}

	var wg sync.WaitGroup
	wg.Add(count)

	if ok := serviceDefs["api"].Run; ok {
		go func() {
			defer wg.Done()
			api.Run()
		}()
	}
	if ok := serviceDefs["srv1"].Run; ok {
		go func() {
			defer wg.Done()
			srv1.Run()
		}()
	}
	if ok := serviceDefs["srv2"].Run; ok {
		go func() {
			defer wg.Done()
			srv2.Run()
		}()
	}

	wg.Wait()
}
