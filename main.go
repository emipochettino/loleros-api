package main

import (
	"github.com/emipochettino/loleros-api/internal/application"
	infraAdapters "github.com/emipochettino/loleros-api/internal/infrastructure/adpaters"
	"github.com/emipochettino/loleros-api/internal/infrastructure/providers"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
	"time"
)

func main() {
	ritoToken := os.Getenv("RITO_TOKEN")
	c := cache.New(30*time.Minute, 40*time.Minute)
	ritoProvider, err := providers.NewRitoProvider(application.GetRitoHosts(), ritoToken, c)
	if err != nil {
		log.Fatalf("Something went wrong trying to create rito provider. %s", err)
	}
	matchService := application.NewMatchService(ritoProvider)
	ritoHandler := infraAdapters.RitoHandler{
		MatchService: matchService,
	}

	_ = infraAdapters.NewRouter(ritoHandler).Run()
}
