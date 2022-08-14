package config

type AppConfig struct {
	Env            string
	Port           string
	AuthKey        string
	MongoURI       string
	DbName         string
	UserCollection string
}
