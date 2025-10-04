package types

type DBConfig struct {
	Database struct {
		User string `yaml:"user"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Name string `yaml:"name"`
	} `yaml:"database"`
}