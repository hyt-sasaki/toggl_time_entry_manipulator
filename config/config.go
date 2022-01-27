package config

type Config struct {
    TogglConfig TogglConfig `desc:"Toggl config"`
}
type ConfigFile string

type TogglConfig struct {
    APIKey string `desc:"Toggl API key"`
}
