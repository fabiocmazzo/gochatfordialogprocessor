package main

import (
	"github.com/elgs/gostrgen"
	"fmt"
	"log"
	"encoding/json"
	"os"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/yosssi/gmq/mqtt"
)

/**
 * o nome dos atributos do struct tem que começar em Uppercase
 * para o marshal funcionar, já que atributos que iniciam em minusculo
 * são como "privates" e não são visiveis externamente.
 */
type MqttMessage struct {
	ClientId              string `json:"clientId"`
	Msg                   string `json:"msg"`
	AdditionalInformationList []AdditionalInformation `json:"additionalInformationList"`
	Metadata              string `json:"metadata"`
	Stage                 string `json:"stage"`
	ResponseTopic         string `json:"responseTopic"`
	Username              string `json:"username"`
	ClientAction          string `json:"clientAction"`
	ClientForm            string `json:"clientForm"`
	QuestionTipList       []string `json:"questionTipList"`
}

func MqttListener() {

	var clientId, _ = gostrgen.RandGen(20, gostrgen.Lower, "", "")
	opts := MQTT.NewClientOptions().AddBroker("tcp://dev.sales4you.com.br:1883").SetClientID("webchat" + clientId)
	opts.SetCleanSession(true)
	MqttCli := MQTT.NewClient(opts)

	if token := MqttCli.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	var callback MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		fmt.Println(string(msg.Topic()), string(msg.Payload()))
		var userMessage MqttMessage
		log.Printf("Payload %v", msg.Payload())
		json.Unmarshal(msg.Payload(), &userMessage)
		log.Printf("Objeto json %v", userMessage)
		privateMessage <- userMessage
	}
	if token := MqttCli.Subscribe("private/messages/sales4you", 0, callback); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func HandleMessages() {

	var clientId, _ = gostrgen.RandGen(20, gostrgen.Lower, "", "")
	opts := MQTT.NewClientOptions().AddBroker("tcp://dev.sales4you.com.br:1883").SetClientID("webchat" + clientId)
	opts.SetCleanSession(true)
	MqttCli := MQTT.NewClient(opts)

	if token := MqttCli.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {

		if (!MqttCli.IsConnected()) {
			MqttCli := MQTT.NewClient(opts)
			if token := MqttCli.Connect(); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
		}

		msg := <-mqttChannel
		log.Printf("Mensagem sendo preparada para ser enviada %v", msg)
		var jsonMqttMessage, _ = json.Marshal(msg)
		log.Printf("Mensagem mqtt pronta para ser enviada %v", string(jsonMqttMessage))
		token := MqttCli.Publish("private/messages/server", mqtt.QoS0, false, []byte(jsonMqttMessage))
		token.Wait()

		// Teste apenas
		/* var responseMessage = Message{ Message :"Olá recebi sua mensagem, mas o Fabio ainda não me connectou ao meu cérebro :( Sinto muito.", Username: "Atendimento", Email: "robo@sales4you.com.br"}

		var client = clients[msg.wsMapKey]
		err := client.WriteJSON(responseMessage)

		if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, msg.wsMapKey)
		} */
	}
}
