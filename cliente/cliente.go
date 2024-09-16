package main

import (
	//"bufio"
	"fmt"
	"net"

	//"os"
	//"strings"
	"sessao3/cliente/funcoesCliente"
)

var ADRESS string = "localhost:22355"

func main() {
	// Conectando ao servidor na porta 8080
	conn, err := net.Dial("tcp", ADRESS)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	funcoesCliente.Menu(conn)

	// // Lendo entrada do usuário
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Digite a cidade de origem: ")
	// origem, _ := reader.ReadString('\n')
	// origem = strings.TrimSpace(origem)

	// fmt.Print("Digite o destino: ")
	// destino, _ := reader.ReadString('\n')
	// destino = strings.TrimSpace(destino)

	// valido := funcoesCliente.VerificarRota(origem, destino)
	// if valido {
	// 	// Enviando mensagem ao servidor
	// 	fmt.Fprintf(conn, "%s,%s\n", origem, destino)

	// 	// Recebendo resposta do servidor
	// 	response, _ := bufio.NewReader(conn).ReadString('\n')
	// 	fmt.Println("Resposta do servidor:", response)
	// }else{
	// 	fmt.Println("Rota inválida")
	// }
}
