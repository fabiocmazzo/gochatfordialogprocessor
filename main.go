package main

import (
	"log"
	"net/http"
	"github.com/elgs/gostrgen"
	"github.com/gorilla/websocket"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var clients = make(map[string]*websocket.Conn) // connected clients
var broadcast = make(chan Message)           // broadcast channel
var privateMessage = make(chan MqttMessage)      // privateMessage channel
var mqttChannel    = make(chan MqttMessage)

var  MqttCli MQTT.Client

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}



func main() {
	// Criando um WebServer para o Webchat
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// A rota do Websocket
	http.HandleFunc("/ws", HandleConnections)


	var clientId, _ = gostrgen.RandGen(20, gostrgen.Lower, "", "")
	opts := MQTT.NewClientOptions().AddBroker("tcp://dev.sales4you.com.br:1883").SetClientID("webchat" + clientId)
	opts.SetCleanSession(true)
	MqttCli := MQTT.NewClient(opts)

	if token := MqttCli.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// Inicia o handleMessage que vai funcionar como um listener das mensagens
	go HandleMessages()
	go MqttListener()
	go ReceiveMessages()

	// Iniciando na porta 8000
	log.Println("Webchat Sales4You iniciado na porta:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


