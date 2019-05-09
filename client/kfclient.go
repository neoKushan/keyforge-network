package main

import (
	"bufio"
	"errors"
	"fmt"
	kfnetwork "keyforge-network"
	"keyforge-network/vault"
	"os"
	"strings"
)

var client *kfnetwork.Client
var connected bool

func main() {
	client = kfnetwork.NewClient()

	for {
		input, e := prompt()

		if e != nil {
			fmt.Println(e)
		}

		sanitized := strings.Trim(input, "\n")
		sanitized = strings.Trim(sanitized, "\r")

		parts := strings.Split(sanitized, " ")
		command := parts[0]
		args := parts[1:]

		routeCommand(command, args)
	}
}

func prompt() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("kfclient> ")

	return reader.ReadString('\n')
}

func routeCommand(command string, args []string) {
	switch command {
	case "quit":
		quit()
	case "connect":
		connect(args)
	case "login":
		login(args)
	case "global":
		global(args)
	case "lobby":
		createLobby(args)
	case "who":
		who()
	default:
		fmt.Println("Command not found.")
	}
}

func connect(args []string) error {
	if connected {
		return errors.New("already connected")
	}

	if len(args) < 1 {
		return errors.New("not enough arguments provided")
	}

	fmt.Println("Attempting to connect to", args[0])
	e := client.Connect(args[0])

	if e != nil {
		fmt.Println("Unable to connect to server:", e.Error())
		return e
	}

	fmt.Println("Connected to server at address", args[0])
	connected = true

	go readLoop()
	return nil
}

func login(args []string) error {
	if !connected {
		fmt.Println("Please connect before attempting to login.")
		return errors.New("attempted to login prior to connecting")
	}

	if len(args) < 2 {
		return errors.New("not enough arguments provided")
	}

	username := args[0]
	password := args[1]

	user, e := vault.Login(username, password)

	if e != nil {
		return e
	}

	e = client.SendVersionRequest()

	if e != nil {
		return e
	}

	e = client.SendLoginRequest(user.UserName, user.ID, user.Token)

	if e != nil {
		return e
	}

	fmt.Println("Logged in as user", username)
	return nil
}

func createLobby(args []string) {
	if len(args) < 1 {
		return
	}

	name := strings.Join(args, " ")
	client.SendCreateLobbyRequest(name)
}

func who() {
	client.SendPlayerListRequest()
}

func quit() {
	fmt.Println("Quitting.")
	os.Exit(0)
}

func global(args []string) {
	if len(args) < 1 {
		return
	}
	message := strings.Join(args, " ")
	client.SendGlobalChatRequest(message)
}

func readLoop() {
	for {
		packet, e := kfnetwork.ReadPacket(client.Connection)

		if e != nil {
			logEntry := fmt.Sprintf("ReadPacket: %s", e.Error())
			fmt.Println(logEntry)
			return
		}

		handlePacket(packet)
	}
}

func handlePacket(packet kfnetwork.Packet) {
	switch packet.GetHeader().Type {
	case kfnetwork.PacketTypePlayerListResponse:
		playerListResponse(packet.(kfnetwork.PlayerListResponsePacket))
	case kfnetwork.PacketTypeGlobalChatResponse:
		globalChatResponse(packet.(kfnetwork.GlobalChatResponsePacket))
	}
}

func playerListResponse(packet kfnetwork.PlayerListResponsePacket) {
	fmt.Println("")
	for _, entry := range packet.Players {
		fmt.Println("ID:", entry.ID, "Name:", entry.Name)
	}

	fmt.Println(packet.Count, "players online.")
}

func globalChatResponse(packet kfnetwork.GlobalChatResponsePacket) {
	message := fmt.Sprintf("[Global Chat] %s: %s", packet.Name, packet.Message)
	fmt.Println(message)
}
