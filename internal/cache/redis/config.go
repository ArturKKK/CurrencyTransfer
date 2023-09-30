package rediscache

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Expires  int    `yaml:"expires"`
}
