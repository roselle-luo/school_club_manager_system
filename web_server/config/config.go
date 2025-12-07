package config

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret  string
	Expires int64
}

type Config struct {
	DB     DBConfig
	JWT    JWTConfig
	Server ServerConfig
}

func Default() Config {
	return Config{
		DB:     DBConfig{Host: "8.138.158.24", Port: 3306, User: "user", Password: "zlsmh123456.", Name: "dachuang"},
		JWT:    JWTConfig{Secret: "replace", Expires: 86400},
		Server: ServerConfig{Addr: ":9000", BaseURL: "http://localhost:8080", PublicDir: "public", UploadDir: "uploads"},
	}
}

type ServerConfig struct {
	Addr      string
	BaseURL   string
	PublicDir string
	UploadDir string
}
