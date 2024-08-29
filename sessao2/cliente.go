package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var cidades = []string{"São Paulo", "Salvador", "Recife"}
var rotas = [][]int{
	{0, 1, 1},
	{1, 0, 1},
	{1, 1, 0},
}

func verifica_rota(origem string, destino string) bool {
	var indice_origem, indice_destino int
	existe_origem := false
	existe_destino := false

	// Encontrar os índices das cidades
	for i, cidade := range cidades {
		if origem == cidade {
			existe_origem = true
			indice_origem = i
		}

		if destino == cidade {
			existe_destino = true
			indice_destino = i
		}
	}

	// Verificar se ambas as cidades existem e se há uma rota
	if existe_origem && existe_destino {
		return rotas[indice_origem][indice_destino] == 1
	}

	return false
}

func main() {
	// Conectando ao servidor na porta 8080
	conn, err := net.Dial("tcp", "172.16.103.4:8080")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	// Lendo entrada do usuário
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Digite a cidade de origem: ")
	origem, _ := reader.ReadString('\n')
	origem = strings.TrimSpace(origem)

	fmt.Print("Digite o destino: ")
	destino, _ := reader.ReadString('\n')
	destino = strings.TrimSpace(destino)

	valido := verifica_rota(origem, destino)
	if valido {
		// Enviando mensagem ao servidor
		fmt.Fprintf(conn, "%s,%s\n", origem, destino)
	
		// Recebendo resposta do servidor
		response, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Resposta do servidor:", response)
	}
	else{
		fmt.Println("Rota inválida")
	}
}
