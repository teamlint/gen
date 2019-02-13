package config

// Config for the running of the commands
type Config struct {
	Debug   bool             `toml:"debug,omitempty" json:"debug,omitempty"`
	Prefix  string           `toml:"prefix" json:"prefix"`
	Suffix  string           `toml:"suffix" json:"suffix"`
	DB      DB               `toml:"db" json:"db"`
	Tables  map[string]Table `toml:"tables,omitempty" json:"tables,omitempty"`
	Model   Model            `toml:"model" json:"model"`
	Query   Query            `toml:"query" json:"query"`
	Service Service          `toml:"service" json:"service"`
}
