package sento

// Bot is a sento-powered bot application
type Bot interface {
	New()
	Start() error
	Stop() error
}
