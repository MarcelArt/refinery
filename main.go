/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/MarcelArt/refinery/cmd"
	_ "github.com/MarcelArt/refinery/docs"
)

// @title Refinery API
// @version 0.0.1
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @securityDefinitions.apikey WebhookKey
// @in header
// @name X-Webhook-Key
func main() {
	cmd.Execute()
}
