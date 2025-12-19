package main

import (
	"embed"
	"encoding/json"
	"os"

	"microblog/backend"
	"microblog/backend/pkg/args"
	"microblog/backend/pkg/logger"
	"microblog/backend/pkg/util"
)

//go:embed package.json
var embeddedVersion []byte

//go:embed dist/*
var embeddedFiles embed.FS

func main() {
	util.LoadEnv()
	logger.InitLogrus()

	if args.Install() != nil ||
		args.Version(embeddedVersion) != nil {
		return
	}
	var info map[string]any
	if json.Unmarshal(embeddedVersion, &info); info["name"] != "" {
		os.Setenv("APP_NAME", info["name"].(string))
	}
	backend.StartServer(embeddedFiles)
}
