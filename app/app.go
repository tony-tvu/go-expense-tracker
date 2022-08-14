package app

type AppConfigs struct {
	Env            string
	Database       string
	DBTimeout      int
	UserCollection string
	AuthKey        []byte
}
