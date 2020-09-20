package sento

// Send `message` into the given channel
func (bot *Bot) Send(info HandleInfo, message string) (err error) {
	_, err = bot.Sess.ChannelMessageSend(info.ChannelID, message)
	return
}
