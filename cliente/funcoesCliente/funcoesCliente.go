package funcoesCliente

import (
	"encoding/json"
	"fmt"
	"net"
)

var indiceCidade = map[int]string{
	0: "São Paulo",
	1: "Salvador",
	2: "Recife",
}
var vagas = map[string][]int{
	"São Paulo": {0, 2, 3},
	"Salvador":  {1, 0, 1},
	"Recife":    {1, 1, 0},
}

var rotas = map[string][]int{
	"São Paulo": {0, 1, 1},
	"Salvador":  {1, 1, 1},
	"Recife":    {1, 1, 1},
}

type Request int

const (
	GET Request = iota
	POST
)

// type Compra struct {
// 	nome    string
// 	cpf     string
// 	trechos [][]string
// }

type Compra struct {
	nome    string
	cpf     string
	origem  string
	destino string
}

type User struct {
	nome string
	cpf  string
}

// var metodo Request = GET
// fmt.Println(metodo.String())

type Dados struct {
	request      Request //Tipo de requisição
	dadosCompra  *Compra //Caso seja um post, as informações da viagem e do passageiro
	dadosUsuario *User   // Dados para cadastro do usuario

}

func Menu() {
	var operacao int
	// Lendo entrada do usuário
	fmt.Println("O que deseja fazer?\n1- Cadastrar usuário\n2- Ver rotas\n3- Comprar passagens")
	fmt.Scanf(&operacao)
}

func SolicitarDados(conn net.Conn) {
	dados := Dados{
		request:      GET,
		dadosCompra:  nil,
		dadosUsuario: nil,
	}
	// Converter os dados para JSON
	jsonData, err := json.Marshal(dados)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}

	// Enviar o JSON ao servidor
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

func VerificarRota(origem string, destino string) bool {
	var indice_destino int
	existe_origem := false
	existe_destino := false

	// Encontrar os índices das cidades
	for cidade, _ := range rotas {
		if origem == cidade {
			existe_origem = true

		}

		if destino == cidade {
			existe_destino = true
			indice_destino = Buscarindice((cidade))
		}
	}

	// Verificar se ambas as cidades existem e se há uma rota
	if existe_origem && existe_destino {
		return rotas[origem][indice_destino] == 1
	}

	return false
}

func Cadastrar(conn net.Conn, nome string, cpf string) {
	user := User{
		nome: nome,
		cpf:  cpf,
	}
	dados := Dados{
		request:      POST,
		dadosCompra:  nil,
		dadosUsuario: &user,
	}

	jsonData, err := json.Marshal(dados)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}

	// Enviar o JSON ao servidor
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

func Comprar(caminhos [][]string, user User) {

}

// Encontrar todos os caminhos possíveis
func EncontrarCaminho(origem, destino string) [][]string {
	caminhos := [][]string{}
	visitados := make(map[string]bool)
	fila := [][]string{{origem}}

	for len(fila) > 0 {
		caminhoAtual := fila[0]
		fila = fila[1:]

		ultimo := caminhoAtual[len(caminhoAtual)-1]

		if ultimo == destino {
			caminhos = append(caminhos, caminhoAtual)
			continue
		}

		if visitados[ultimo] {
			continue
		}
		visitados[ultimo] = true

		for i, rota := range rotas[ultimo] {
			if rota == 1 {
				proximaCidade := indiceCidade[i]
				if !visitados[proximaCidade] {
					fila = append(fila, append(caminhoAtual, proximaCidade))
				}
			}
		}
	}

	return caminhos
}

func Buscarindice(cidade string) int {
	for i, cidadeincice := range indiceCidade { //busca o indice da cidade
		if cidadeincice == cidade {
			return i
		}
	}
	return -1
}
