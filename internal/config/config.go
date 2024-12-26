package config

type EmailConfig struct {
	From     string
	To       string
	Password string
	SMTPHost string
	SMTPPort int
}

type Config struct {
	LichessAPIBase string
	Users          []string
	Email          EmailConfig
}

func LoadConfig() Config {
	return Config{
		LichessAPIBase: "https://lichess.org/api/games/user/",
		// Users:          []string{"itsspriyansh"},
		Users: []string{"praneki_li","itsspriyansh"},
		Email: EmailConfig{
			From: "pranjaljha00@gmail.com",
			To:   "idupidu00@gmail.com",
			Password: "zcppikkpqnnoolut",
			SMTPHost: "smtp.gmail.com",
			SMTPPort: 587,
		},
	}
}
