package initializers

import "github.com/neel1996/saul/repository"

var (
	userRepository repository.UserRepository
)

func InitializeRepositories() {
	userRepository = repository.NewUserRepository(dynamoDBClient)
}
