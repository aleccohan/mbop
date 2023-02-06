package config

import "os"

type mbopConfig struct {
	MailerModule string
	JwtModule    string
	JwkURL       string
	UsersModule  string
}

var conf *mbopConfig

func Get() *mbopConfig {
	if conf != nil {
		return conf
	}

	c := &mbopConfig{}
	c.UsersModule = fetchWithDefault("USERS_MODULE", "")
	c.JwtModule = fetchWithDefault("JWT_MODULE", "")
	c.JwkURL = fetchWithDefault("JWK_URL", "")
	c.MailerModule = fetchWithDefault("MAILER_MODULE", "print")

	conf = c
	return conf
}

func fetchWithDefault(name, defaultValue string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}

	return defaultValue
}
