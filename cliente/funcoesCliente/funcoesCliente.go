package funcoesCliente

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"strings"
)

type Request int

const (
	ROTAS Request = iota
	COMPRA
	CADASTRO
)

type Compra struct {
	Nome    string
	Cpf     string
	Caminho []string
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

// Rota representa uma rota entre duas cidades com uma quantidade de vagas e um peso.
type Rota struct {
	Destino string
	Vagas   int
	Peso    int
}

// Estrutura de dados para o grafo das rotas.
var rotas = map[string][]Rota{
	"São Paulo": {
		{"Salvador", 4, 15},
		{"Recife", 3, 20},
		{"Feira", 1, 15},
	},
	"Salvador": {
		{"São Paulo", 1, 10},
		{"Recife", 2, 25},
		{"Feira", 2, 5},
	},
	"Recife": {
		{"São Paulo", 1, 20},
		{"Salvador", 1, 25},
	},
	"Feira": {
		{"Salvador", 1, 5},
		{"Recife", 2, 10},
	},
	"Manaus": {
		{"São Paulo", 1, 30},
		{"Recife", 1, 40},
	},
}

// Função para realizar a busca em profundidade para encontrar o caminho com o menor peso total.
func buscaProfundidade(cidadeAtual, destino string, visitado map[string]bool, caminhoAtual []string, pesoAtual, menorPeso *int, melhorCaminho *[]string) {
	if cidadeAtual == destino {
		if pesoAtual != nil && menorPeso != nil && *pesoAtual < *menorPeso {
			*menorPeso = *pesoAtual
			*melhorCaminho = append([]string(nil), caminhoAtual...) // Copia o caminho atual para o melhor caminho
		}
		return
	}

	visitado[cidadeAtual] = true
	for _, rota := range rotas[cidadeAtual] {
		if !visitado[rota.Destino] {
			caminhoAtual = append(caminhoAtual, rota.Destino)
			*pesoAtual += rota.Peso
			buscaProfundidade(rota.Destino, destino, visitado, caminhoAtual, pesoAtual, menorPeso, melhorCaminho)
			*pesoAtual -= rota.Peso
			caminhoAtual = caminhoAtual[:len(caminhoAtual)-1]
		}
	}
	visitado[cidadeAtual] = false
}

func menorCaminhoDFS(inicio, fim string) ([]string, int) {
	visitado := make(map[string]bool)
	caminhoAtual := []string{inicio}
	var melhorCaminho []string
	var pesoAtual, menorPeso int
	menorPeso = math.MaxInt32

	buscaProfundidade(inicio, fim, visitado, caminhoAtual, &pesoAtual, &menorPeso, &melhorCaminho)
	return melhorCaminho, menorPeso
}

func Menu(conn net.Conn) {
	var operacao int

	// Lendo entrada do usuário
	fmt.Println("O que deseja fazer?\n1- Cadastrar usuário\n2- Comprar passagens\n3-sla ")
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
		break

	case 2:

		// Lendo entrada do usuário
		var origem, destino string
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Digite a cidade de origem:")
		origem, _ = reader.ReadString('\n')
		origem = strings.TrimSpace(origem)

		fmt.Println("Digite a cidade de destino:")
		destino, _ = reader.ReadString('\n')
		destino = strings.TrimSpace(destino)


		valido := VerificarCidade(origem, destino)
		if valido {
			user := User{
				Nome: "Júlia",
				Cpf:  "093.234.234-23",
			}
			Comprar(conn, user, origem, destino)
		} else {
			fmt.Println("Não há rota disponível entre essas cidades.")
		}
		break
	case 3:
		// Função de compra ainda não implementada
		
		

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

func VerificarCidade(origem string, destino string) bool {
	//var indice_destino int
	existe_origem := false
	existe_destino := false

	// Encontrar os índices das cidades
	for cidade := range rotas {
		if origem == cidade {
			existe_origem = true
		}

		if destino == cidade {
			existe_destino = true
			//indice_destino = Buscarindice(cidade)
		}
	}

	// Verificar se ambas as cidades existem e se há uma rota
	if existe_origem && existe_destino {
		return true
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

func Comprar(conn net.Conn, user User, origem string, destino string) {

	caminho, _ := menorCaminhoDFS(origem, destino)
	fmt.Printf("Rota encontrada - %s a %s: %v", origem, destino, caminho)

	compra := Compra{
		Nome:    user.Nome,
		Cpf:     user.Cpf,
		Caminho: caminho,
	}

	dados := Dados{
		Request:      COMPRA,
		DadosCompra:  &compra,
		DadosUsuario: nil,
	}

	var resposta int
	fmt.Print("Deseja realizar a compra?\n1- Sim\n2- Não\n")
	fmt.Scanf("%d\n", &resposta)
	if resposta == 1 {
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
}
