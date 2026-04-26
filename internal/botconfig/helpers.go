package botconfig

func maskToken(token string) string {
	if len(token) <= 8 {
		return "********"
	}
	return token[:4] + "..." + token[len(token)-4:]
}

func toBotView(bot BotConfig, dbName string) BotView {
	return BotView{
		ID:                   bot.ID,
		Name:                 bot.Name,
		TokenMasked:          maskToken(bot.Token),
		DatabaseID:           bot.DatabaseID,
		DatabaseName:         dbName,
		StartScenario:        bot.StartScenario,
		TelegramAdminUserIDs: bot.TelegramAdminUserIDs,
		AdminOrdersChatID:    bot.AdminOrdersChatID,
		TelegramBotID:        bot.TelegramBotID,
		TelegramUsername:     bot.TelegramUsername,
		TelegramBotName:      bot.TelegramBotName,
		IsEnabled:            bot.IsEnabled,
		UpdatedAt:            bot.UpdatedAt,
	}
}

func toDatabaseProfileView(profile DatabaseProfile) DatabaseProfileView {
	return DatabaseProfileView{
		ID:        profile.ID,
		Name:      profile.Name,
		Driver:    profile.Driver,
		IsEnabled: profile.IsEnabled,
		UpdatedAt: profile.UpdatedAt,
	}
}
