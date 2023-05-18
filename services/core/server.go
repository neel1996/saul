package main

import (
	"core/initializers"
)

func main() {
	config := initializers.InitializeConfiguration()
	initializers.Bootstrap(config)

	r := initializers.InitializeRoutes(config)

	err := r.Run(":" + config.Port)
	if err != nil {
		panic(err)
	}
}
