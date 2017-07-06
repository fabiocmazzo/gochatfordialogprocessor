package main

import (
	"log"
	"net/http"
)

type AdditionalInformation struct {
	Name     string `json:"name"`
	Datatype string `json:"datatype"`
	Value    string `json:"value"`
}

// Define our message object
type Message struct {
	Email                     string `json:"email"`
	Username                  string `json:"username"`
	Message                   string `json:"message"`
	Hidden                    bool   `json:"hidden"`
	Action                    string `json:"action"`
	Form                      string `json:"form"`
	AdditionalInformationList []AdditionalInformation `json:"additionalInformationList"`
	Tips                      string `json:"tips"`
	QuestionTipList           []string `json:"questionTipList"`
	ClientAppId		  string `json:"clientAppId"`
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// O clients é um array global, aqui eu tenho a identificação do cliente, para quando tiver que retornar
	// do MQTT, eu tenho que enviar especificamente para esse client.

	// Vou utilizar o UUID que vou passar do frontend e não mais gerado pelo Go.
	//wsMapKey, err := gostrgen.RandGen(20, gostrgen.Lower, "", "")
	//clients[wsMapKey] = ws

	var wsMapKey string

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object

		err = ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, wsMapKey)
			break
		}

		wsMapKey = msg.ClientAppId
		clients[wsMapKey] = ws
		println("ClientID " + wsMapKey)



		log.Printf("Mensagem recebida %v", msg)

		// enviando a mensagem para o próprio usuário para aparecer no chat
		// caso não seja oculta
		if (!msg.Hidden) {
			ws.WriteJSON(msg)
		}

		// Enviando a mensagem convertida para o Formato MqttMessage
		// para o channel do mqtt
		var mqttMsg MqttMessage

		println(msg.AdditionalInformationList)
		mqttMsg.ClientId = wsMapKey
		mqttMsg.Msg = msg.Message
		mqttMsg.AdditionalInformationList = msg.AdditionalInformationList
		mqttMsg.Metadata = ""
		mqttMsg.Stage = "msg"
		mqttMsg.ResponseTopic = ""
		mqttMsg.Username = "Você"
		mqttMsg.ResponseTopic = "private/messages/sales4you"
		// Por enquanto só temos mensagem nesse Webchat, não temos StartTyping e etc.
		mqttMsg.Stage = "msg"
		mqttChannel <- mqttMsg

	}
}

func ReceiveMessages() {
	for {
		var msg = <-privateMessage
		log.Printf("recuperando a conexao do cleint %v", msg.ClientId)
		var client = clients[msg.ClientId]

		var userMessage Message
		userMessage.Message = msg.Msg
		userMessage.Username = msg.Username
		userMessage.Form = msg.ClientForm
		userMessage.QuestionTipList = msg.QuestionTipList

		err := client.WriteJSON(userMessage)

		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(clients, msg.ClientId)
		}
	}
}
