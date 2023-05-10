package main

import (
	"github.com/neel1996/saul/initializers"
)

func main() {
	config := initializers.InitializeConfiguration()
	initializers.InitializeObjects(config)

	r := initializers.InitializeRoutes(config)

	err := r.Run(":" + config.Port)
	if err != nil {
		panic(err)
	}
}
