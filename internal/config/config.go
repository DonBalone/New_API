package config

type Config struct {
	Env         string `yaml:"env"`
	StorageInfo `yaml:"storage_info"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address  string `yaml:"address"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
type StorageInfo struct {
	host     string `yaml:"localhost"`
	port     string `yaml:"5432"`
	sslmode  string `yaml:"disable"`
	dbname   string `yaml:"new_api"`
	user     string `yaml:"postgres"`
	password string `yaml:"12345"`
}
