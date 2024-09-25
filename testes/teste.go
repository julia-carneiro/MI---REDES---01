package main

import (
	"encoding/json"
	"fmt"
	"sessao3/cliente/funcoesCliente"
)

var ADRESS string = "servidor-container:22356"

type Request int

const ( //Tipos de mensagens que podem ser enviadas ao servidor
	ROTAS Request = iota
	COMPRA
	CADASTRO
	LERCOMPRAS
)

type Compra struct { //Estrutura de dados de compra
	Cpf     string
	Caminho []string
}

type User struct { //Estrtura de dados de usuario
	Cpf string `json:"Cpf"`
}

type Dados struct { //Estrutura de dados de mensagem para o servidor
	Request      Request `json:"Request"`
	DadosCompra  *Compra `json:"DadosCompra"`
	DadosUsuario *User   `json:"DadosUsuario"`
}

func main() {
	var caminho = []string{"Feira", "Brasília"}
	user := funcoesCliente.User{Cpf: "12345678910"}

	conn := funcoesCliente.ConectarServidor(ADRESS)
	defer conn.Close()
	funcoesCliente.Cadastrar(conn, user.Cpf)

	compra := funcoesCliente.Compra{
		Cpf:     user.Cpf,
		Caminho: caminho,
	}

	dados := funcoesCliente.Dados{
		Request:      funcoesCliente.COMPRA,
		DadosCompra:  &compra,
		DadosUsuario: nil,
	}

	//Converter dados para JSON
	jsonData, err := json.Marshal(dados)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}

	// // Enviar o JSON ao servidor
	// fmt.Println("Enviando dados:", string(jsonData)) // Exibe o JSON como string
	conn.Write(jsonData)
	conn.Write([]byte("\n")) // Enviar uma nova linha para indicar o fim da mensagem

	// Ler a resposta do servidor
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Erro ao ler a resposta do servidor:", err)
		return
	}
	fmt.Println("Resposta do servidor:", string(buffer[:n]))

}
