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
    Email      EmailConfig
}

func LoadConfig() Config {
	return Config{
		LichessAPIBase: "https://lichess.org/api/games/user/",
		Users:          []string{"itsspriyansh"},
		Email: EmailConfig{
            From:     "pranjaljha00@gmail.com",
            To:       "21bcs161@iiitdmj.ac.in",
            Password: "Pranjal@99",
            SMTPHost: "smtp.gmail.com",
            SMTPPort: 587,
        },
	}
}