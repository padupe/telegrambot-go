package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const telegramAPI string = "https://api.telegram.org/bot"
const sendMsg string = "/sendMessage"
const botToken string = "BOT_TELEGRAM_TOKEN"

var botAPI string = telegramAPI + os.Getenv(botToken) + sendMsg

const punchCommand string = "/punch"

var lenPunchCommand int = len(punchCommand)

const startCommand string = "/start"

var lenStartCommand int = len(startCommand)

const botTag string = "@QualEAPrevisaoBot"

var lenBotTag int = len(botTag)

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat string `json:"chat"`
}

type Chat struct {
	Id int `json:"id"`
}

// sanitize remove comandos /start /punch ou o nome do próprio Bot
func sanitize(s string) string {
	if len(s) >= lenStartCommand {
		if s[:lenStartCommand] == startCommand {
			s = s[lenStartCommand:]
		}
	}

	if len(s) >= lenPunchCommand {
		if s[:lenPunchCommand] == punchCommand {
			s = s[lenPunchCommand:]
		}
	}
	if len(s) >= lenBotTag {
		if s[:lenBotTag] == botTag {
			s = s[lenBotTag:]
		}
	}
	return s
}

func parseTelegramRequest(req *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(req.Body).Decode(&update); err != nil {
		log.Printf("Error incomening update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

func HandleTelegramWebHook(res http.ResponseWriter, req *http.Request) {

	// Analisar a Entrada da Requisição
	var update, err = parseTelegramRequest(req)
	if err != nil {
		log.Printf("Error parsing update %s", err.Error())
		return
	}

	// Limpando a Entrada
	var cleanSeed = sanitize(update.Message.Text)
	fmt.Println(cleanSeed)
}

func sendTextToTelegramChat(chatId int, text string) (string, error) {
	log.Printf("Sending %s to chat_id: %d", text, chatId)

	response, err := http.PostForm(
		botAPI,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if err != nil {
		log.Printf("Error with posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}
