package main

import (
	"os"

	"github.com/Faizan2005/Movie-Database/config"
	"github.com/Faizan2005/Movie-Database/controllers"

	"github.com/Faizan2005/Movie-Database/utils"
)

func main() {
	var uri = config.MongoURI

	client, ctx, cancel, err := utils.Connect(uri)

	if err != nil {
		panic(err)
	}

	defer utils.Close(client, ctx, cancel)

	utils.Ping(client, ctx)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := controllers.NewAPIServer(":"+port, client)
	server.Run()
}
