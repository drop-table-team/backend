package storage

type Config struct {
	Endpoint    string
	Bucket      string
	Credentials Credentials
}

type Credentials struct {
	AccessKey string
	SecretKey string
}
