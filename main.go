package main

import (
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

	server := controllers.NewAPIServer(":8000", client)
	server.Run()
}
