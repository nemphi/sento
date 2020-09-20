package sento

// Send `message` into the given channel
func (bot *Bot) Send(info HandleInfo, message string) (err error) {
	msg, err := bot.Sess.ChannelMessageSend(info.ChannelID, message)
	handleInfo := HandleInfo{
		AuthorID:  msg.Author.ID,
		ChannelID: msg.ChannelID,
		GuildID:   msg.GuildID,
	}
	bot.EmitEvent(EventMessageSent, handleInfo)
	return
}
