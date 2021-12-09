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

	botName := bot.Self.UserName
	fmt.Println("Nome do Bot", botName)

	botUsername := bot.Self.FirstName + "" + bot.Self.LastName
	fmt.Println("Nome de Usuário do Bot", botUsername)

	// Cria uma nova estrutura UpdateConfig com um deslocamento de 0. Offsets são usados
	// para garantir que o Telegram saiba que tratamos dos valores anteriores e não
	// precisa deles repetidos.
	updateConfig := tgbotapi.NewUpdate(0)
	fmt.Println("1", updateConfig)

	// Diga ao Telegram que devemos esperar até 30 segundos em cada solicitação de um
	// update. Dessa forma, podemos obter informações tão rapidamente quanto fazer muitas
	// solicitações frequentes sem ter que enviar quase a mesma quantidade.
	updateConfig.Timeout = 30
	fmt.Println("2", updateConfig.Timeout)

	// Comece a sondar o Telegram para verificar se houveram atualizações.
	updates := bot.GetUpdatesChan(updateConfig)

	// Vamos examinar cada atualização que recebemos do Telegram.
	for update := range updates {
		// O Telegram pode enviar muitos tipos de atualizações, dependendo do que o seu Bot
		// está preparado para fazer. Queremos apenas olhar as mensagens por enquanto, para que possamos
		// descartar quaisquer outras atualizações.
		if update.Message == nil {
			continue
		}

		username := update.Message.Chat.FirstName
		fmt.Println("Username:", username)

		// Agora que sabemos que recebemos uma nova mensagem, podemos construir uma
		// resposta! Pegaremos o ID do bate-papo e o texto da mensagem recebida
		// e usaremos para criar uma nova mensagem.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		fmt.Println("Id do Usuário e Texto da Mensagem", msg)
		// Também diremos que esta mensagem é uma resposta à mensagem anterior.
		// Para quaisquer outras especificações além de ID de bate-papo ou Texto, você precisará
		// define campos em `MessageConfig`.

		msg.ReplyToMessageID = update.Message.MessageID
		fmt.Println("ID da Mensagem", msg.ReplyToMessageID)
		fmt.Println("Responder para Mensagem", update.Message.MessageID)

		// Ok, estamos enviando nossa mensagem! Não nos importamos com a mensagem
		// acabamos de enviar, então vamos descartá-la.
		bot.Send(msg)
		// tgbotapi.NewDeleteMessage(update.Message.Chat.ID)
		/*
					TODO - Configurar o Bot para responder a comandos específicos
					if !update.Message.IsCommand() { // ignore any non-command Messages
			            continue
			        }

			        // Crie um novo MessageConfig. Ainda não temos texto,
			        // então o deixamos vazio.
			        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			        // Extraia o comando da Mensagem.
			        switch update.Message.Command() {
			        case "help":
			            msg.Text = "I understand /sayhi and /status."
			        case "sayhi":
			            msg.Text = "Hi :)"
			        case "status":
			            msg.Text = "I'm ok."
			        default:
			            msg.Text = "I don't know that command"
			        }
		*/

		// if _, err := bot.Send(msg); err != nil {
		// 	// Observe que os panics são uma maneira ruim de lidar com os erros. Telegram pode
		// 	// ter interrupções no serviço ou erros de rede, você deve tentar enviar novamente
		// 	// mensagens ou lidar de forma mais adequada com as falhas.
		// 	panic(err)
		// }
	}

}
