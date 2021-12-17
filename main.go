package main

import (
	// "errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type ReponseCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Bairro      string `json:"bairro"`
	Complemento string `json:"complemento"`
	Cidade      string `json:"localidade"`
	Estado      string `json:"uf"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	DDD         string `json:"ddd"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
}

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
	// var cep string = ""

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

			const buscaCep = getCep(msgRec)
			txt = buscaCep

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

func getCep(cep string) http.ResponseWriter {

	baseURL := "https://viacep.com.br/ws/" + cep + "/json/"

	response, err := http.Get(baseURL)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	fmt.Println("Status Code: ", response.StatusCode)
	fmt.Println("Content lengh is: ", response.ContentLength)

	content, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(content))
}
