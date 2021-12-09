package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment settings!")
	}

	botToken := os.Getenv("BOT_TELEGRAM_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorize an account %s", bot.Self.UserName)

	// Nome do Bot
	botName := bot.Self.UserName
	fmt.Println("Nome do Bot", botName)

	// Nome de Usuário do Bot
	botUsername := bot.Self.FirstName + "" + bot.Self.LastName
	fmt.Println("Nome de Usuário do Bot", botUsername)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		msgRec := update.Message.Text

		if update.Message == nil {
			log.Printf("[%s] %s", update.Message.From.UserName, msgRec)
			continue

		}

		username := update.Message.Chat.FirstName

		txt := fmt.Sprintf("Seja bem-vindo %s!", username)
		if len(msgRec) == 8 {
			txt = fmt.Sprintf("%s, estamos pesquisando seu CEP", username)
		} else {
			txt = fmt.Sprintf("%s, informe um CEP válido", username)
		}
		// fmt.Println("Será que é um CEP?", cep)
		// Nome do Usuário que está converando com o Bot
		// username := update.Message.Chat.FirstName

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt)

		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)

		// if _, err := bot.Send(msg); err != nil {
		// 	// Observe que os panics são uma maneira ruim de lidar com os erros. Telegram pode
		// 	// ter interrupções no serviço ou erros de rede, você deve tentar enviar novamente
		// 	// mensagens ou lidar de forma mais adequada com as falhas.
		// 	panic(err)
		// }
	}

}
