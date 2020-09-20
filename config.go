package sento

// Config for a sento-powered discord bot
type Config struct {
	Token      string
	InviteLink string
	Prefix     string
}

var (
	// DefaultConfig for a discord bot
	DefaultConfig = Config{
		Prefix: "#!",
	}
)
