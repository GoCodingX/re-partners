package config

type Config struct {
	Port           string   `env:"PORT"                    envDefault:"8080"`
	DatabaseUrl    string   `env:"DATABASE_URL,required"`
	MigrationsDir  string   `env:"MIGRATIONS_DIR,required"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS"`
}
