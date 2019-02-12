package config

// Config for the running of the commands
type Config struct {
	DB      DB               `toml:"db" json:"db"`
	Debug   bool             `toml:"debug,omitempty" json:"debug,omitempty"`
	Tables  map[string]Table `toml:"tables,omitempty" json:"tables,omitempty"`
	Model   Model            `toml:"model" json:"model"`
	Query   Query            `toml:"query" json:"query"`
	Service Service          `toml:"service" json:"service"`
}
