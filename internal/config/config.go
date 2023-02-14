package config

import "os"

type MbopConfig struct {
	MailerModule string
	JwtModule    string
	JwkURL       string
	UsersModule  string
}

var conf *MbopConfig

func Get() *MbopConfig {
	if conf != nil {
		return conf
	}

	c := &MbopConfig{}
	c.UsersModule = fetchWithDefault("USERS_MODULE", "")
	c.JwtModule = fetchWithDefault("JWT_MODULE", "")
	c.JwkURL = fetchWithDefault("JWK_URL", "")
	c.MailerModule = fetchWithDefault("MAILER_MODULE", "print")
	c.UsersModule = fetchWithDefault("USERS_MODULE", "")

	conf = c
	return conf
}

func fetchWithDefault(name, defaultValue string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}

	return defaultValue
}

// TO BE USED FROM TESTING ONLY.
func Reset() {
	conf = nil
}
