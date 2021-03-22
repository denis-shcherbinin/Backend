package main

import (
	"github.com/PolyProjectOPD/Backend/internal/api"
	_ "github.com/lib/pq"
)

const configPath = "configs/config"

func main() {
	api.Run(configPath)
}
