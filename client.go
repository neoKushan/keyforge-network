package kfnetwork

import (
	"net"
)

type Client struct {
	Connection net.Conn
	Sequence   uint16
}

func NewClient() *Client {
	client := new(Client)
	return client
}

func (c *Client) Connect(address string) error {
	var e error
	c.Connection, e = net.Dial("tcp", address)

	if e != nil {
		return e
	}

	return nil
}

func (c *Client) SendVersionRequest() error {
	packet := VersionPacket{}
	packet.Sequence = c.Sequence
	packet.Type = PacketTypeVersionRequest
	packet.Version = ProtocolVersion

	e := WritePacket(c.Connection, packet)
	c.Sequence++
	return e
}

func (c *Client) SendExitRequest() error {
	packet := ExitPacket{}
	packet.Sequence = c.Sequence
	packet.Type = PacketTypeExit

	e := WritePacket(c.Connection, packet)
	c.Sequence++
	return e
}

func (c *Client) SendLoginRequest() error {
	packet := LoginRequestPacket{}
	packet.Sequence = c.Sequence
	packet.Type = PacketTypeLoginRequest
	packet.Name = "test"
	packet.ID = GenerateUUID()

	e := WritePacket(c.Connection, packet)
	c.Sequence++
	return e
}

func (c *Client) SendCreateLobbyRequest() error {
	packet := CreateLobbyRequestPacket{}
	packet.Sequence = c.Sequence
	packet.Type = PacketTypeCreateLobbyRequest

	e := WritePacket(c.Connection, packet)
	c.Sequence++
	return e
}

func (c *Client) SendGetCardPile(pile uint8) error {
	packet := CardPileRequestPacket{}
	packet.Sequence = c.Sequence
	packet.Type = PacketTypeCardPileRequest
	packet.Pile = pile

	e := WritePacket(c.Connection, packet)
	c.Sequence++
	return e
}

func (c *Client) SendGetArchivePile() {
	c.SendGetCardPile(CardPileArchive)
}
