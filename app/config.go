package app

type Config struct {
	DatabaseConnectionString string
	LogDatabaseQueries       bool
	PhotosPath               string
	DuplicatesPath           string
}
