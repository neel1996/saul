package initializers

import "core/repository"

var (
	userRepository repository.UserRepository
)

func InitializeRepositories() {
	userRepository = repository.NewUserRepository(dynamoDBClient)
}
