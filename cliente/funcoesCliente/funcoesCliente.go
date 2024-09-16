package funcoesCliente

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
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
	ROTAS Request = iota
	COMPRA
	CADASTRO
)

type Compra struct {
	nome    string
	cpf     string
	origem  string
	destino string
}

type User struct {
	Nome string `json:"Nome"`
	Cpf  string `json:"Cpf"`
}

type Dados struct {
	Request      Request `json:"Request"`
	DadosCompra  *Compra `json:"DadosCompra"`
	DadosUsuario *User   `json:"DadosUsuario"`
}

func Menu(conn net.Conn) {
	var operacao int
	

	// Lendo entrada do usuário
	fmt.Println("O que deseja fazer?\n1- Cadastrar usuário\n2- Ver rotas\n3- Comprar passagens")
	fmt.Scanf("%d\n", &operacao)

	switch operacao {
	case 1:
		var nome, cpf string
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Digite seu nome:")
		nome, _ = reader.ReadString('\n')
		nome = strings.TrimSpace(nome)
	
		fmt.Println("Digite seu CPF:")
		cpf, _ = reader.ReadString('\n')
		cpf = strings.TrimSpace(cpf)

		// fmt.Println("Digite seu nome:")
		// fmt.Scanf("%s\n"&nome)
		// fmt.Println("Digite seu CPF:")
		// fmt.Scanln(&cpf)

		Cadastrar(conn, nome, cpf)
	case 2:
		// Lendo entrada do usuário
		var origem, destino string
		fmt.Print("Digite a cidade de origem: ")
		fmt.Scanf("%s", &origem)

		fmt.Print("Digite o destino: ")
		fmt.Scanf("%s", &destino)

		valido := VerificarRota(origem, destino)
		if valido {
			caminhos := EncontrarCaminho(origem, destino)
			fmt.Println("Caminhos encontrados:", caminhos)
		} else {
			fmt.Println("Não há rota disponível entre as cidades.")
		}
	case 3:
		// Função de compra ainda não implementada
		fmt.Println("Função de compra ainda não implementada.")
	default:
		fmt.Println("Operação inválida.")
	}

}

func SolicitarDados(conn net.Conn) {
	dados := Dados{
		Request:      ROTAS,
		DadosCompra:  nil,
		DadosUsuario: nil,
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
	for cidade := range rotas {
		if origem == cidade {
			existe_origem = true
		}

		if destino == cidade {
			existe_destino = true
			indice_destino = Buscarindice(cidade)
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
		Nome: nome,
		Cpf:  cpf,
	}
	dados := Dados{
		Request:      CADASTRO,
		DadosCompra:  nil,
		DadosUsuario: &user, // Deve ser um ponteiro
	}

	//Converter dados para JSON
	jsonData, err := json.Marshal(dados)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}

	// // Enviar o JSON ao servidor
	fmt.Println("Enviando dados:", string(jsonData)) // Exibe o JSON como string
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
	// Função de compra ainda não implementada
}

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
	for i, cidadeindice := range indiceCidade { // busca o índice da cidade
		if cidadeindice == cidade {
			return i
		}
	}
	return -1
}
