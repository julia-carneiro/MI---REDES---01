package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Conectar ao servidor na porta 8080
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	// Ler entrada do usuário e enviar ao servidor
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Digite uma mensagem: ")
		text, _ := reader.ReadString('\n')

		// Enviar ao servidor
		fmt.Fprintf(conn, text+"\n")

		// Receber a resposta do servidor
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Resposta do servidor:", message)
	}
}
