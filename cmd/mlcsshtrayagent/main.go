package main

import (
	"github.com/mlc911/mlcsshtrayagent/internal/app"
	"github.com/mlc911/mlcsshtrayagent/internal/monitor"
)

func main() {
	app.New(monitor.NewProvider()).Run()
}
