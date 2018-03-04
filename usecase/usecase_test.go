package usecase

type inMemoryRepository struct {
	user *inMemoryUserRepository

	transactor
}

type inMemoryUserRepository struct {
}
