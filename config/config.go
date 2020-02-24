package config

import "os"

// AppConfig contains the application configuration
type AppConfig struct {
	AppSecret       string
	ReCaptchaSecret string
	IPStackKey      string
	Version         string
	AppHost         string
}

// LoadSecrets loads secret config values from env vars
func (c *AppConfig) LoadSecrets() {
	c.AppSecret = getEnv("APP_SECRET", c.AppSecret)
	c.ReCaptchaSecret = getEnv("RECAPTCHA_SECRET", c.ReCaptchaSecret)
	c.IPStackKey = getEnv("IPSTACK_KEY", c.IPStackKey)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
