package config

// using driver yaml
type Configuration struct {
	Server ServerConfig `yaml:"Server"`
	DB     DBConfig     `yaml:"DB"`
	Cert   CertConfig   `yaml:"Cert"`
}

type (
	ServerConfig struct {
		Port                    string `yaml:"Port"`
		BasePath                string `yaml:"BasePath"`
		GracefulTimeoutInSecond int    `yaml:"GracefulTimeout"`
		ReadTimeoutInSecond     int    `yaml:"ReadTimeout"`
		WriteTimeoutInSecond    int    `yaml:"WriteTimeout"`
		APITimeout              int    `yaml:"APITimeout"`
	}
	DBConfig struct {
		RetryInterval int    `yaml:"RetryInterval"`
		MaxIdleConn   int    `yaml:"MaxIdleConn"`
		MaxConn       int    `yaml:"MaxConn"`
		DSN           string `yaml:"DSN"`
	}
	CertConfig struct {
		Public  string `yaml:"Public"`
		Private string `yaml:"Private"`
	}
)
