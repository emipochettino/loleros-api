package main

import (
	infrastructure "github.com/emipochettino/loleros-api/internal/infrastructure/adpaters"
)

func main() {
	_ = infrastructure.Route().Run()
}
