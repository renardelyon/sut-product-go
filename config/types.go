package config

type Config struct {
	Port      string `mapstructure:"PORT"`
	NotifHost string `mapstructure:"NOTIF_HOST"`
	DBUrl     string `mapstructure:"DB_URL"`
}
