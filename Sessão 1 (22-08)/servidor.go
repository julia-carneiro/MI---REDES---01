package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	// Escutar na porta 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Servidor está escutando na porta 8080...")

	// Aceitar conexões de clientes
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Erro ao aceitar conexão:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Cliente conectado.")

	// Ler e escrever mensagens
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Mensagem recebida:", string(message))

		// Responder ao cliente
		newMessage := strings.ToUpper(message)
		conn.Write([]byte(newMessage + "\n"))
	}
}
