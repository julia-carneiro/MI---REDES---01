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
	LERCOMPRAS
)

type Compra struct {
	Cpf     string
	Caminho []string
}

type User struct {
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
var rotas map[string][]Rota

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

func Menu(ADRESS string, user User) {
	var operacao int
	var conn net.Conn
	var i = true

	for i {

		// Lendo entrada do usuário
		fmt.Println("O que deseja fazer?\n1- Comprar passagens\n2-Ver passagens compradas\n3-Sair ")
		fmt.Scanf("%d\n", &operacao)

		switch operacao {
		case 1:
			// Lendo entrada do usuário
			var origem, destino string
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Digite a cidade de origem:")
			origem, _ = reader.ReadString('\n')
			origem = strings.TrimSpace(origem)

			fmt.Println("Digite a cidade de destino:")
			destino, _ = reader.ReadString('\n')
			destino = strings.TrimSpace(destino)

			inicioExiste := false
			fimExiste := false
		
			conn = ConectarServidor(ADRESS)
			BuscarDados(conn)
			defer conn.Close()
			// Verifica se a cidade inicial existe no mapa de rotas
			if _, existe := rotas[origem]; existe {
				inicioExiste = true
			}
		
			// Verifica se a cidade final existe no mapa de rotas
			if _, existe := rotas[destino]; existe {
				fimExiste = true
			}
			// Caso qualquer uma das cidades não exista, não é necessário continuar
			if inicioExiste && fimExiste {

				conn = ConectarServidor(ADRESS)
				Comprar(conn, user, origem, destino)
				defer conn.Close()
			}else{
				fmt.Println("Não existe rota")
			}

		case 2:
			conn = ConectarServidor(ADRESS)
			VerPassagensCompradas(conn, user.Cpf)
			defer conn.Close()
		case 3:
			i = false
			break
		default:
			fmt.Println("Operação inválida.")

		}
	}

}

func ConectarServidor(ADRESS string) net.Conn {
	// Conectando ao servidor na porta 8080
	conn, err := net.Dial("tcp", ADRESS)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return nil
	}

	return conn
}
func BuscarDados(conn net.Conn) {
	dados := Dados{
		Request:      ROTAS,
		DadosCompra:  nil,
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
	// fmt.Println("Resposta do servidor:", string(buffer[:n]))

	// Desserializar o JSON recebido
	err = json.Unmarshal(buffer[:n], &rotas)
	if err != nil {
		fmt.Println("Erro ao converter JSON para estrutura:", err)
		return
	}
	// Exibir os dados convertidos
	// fmt.Println("Dados convertidos:", rotas)
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

// func VerificarCidade(origem string, destino string) bool {
// 	//var indice_destino int
// 	existe_origem := false
// 	existe_destino := false

// 	// Encontrar os índices das cidades
// 	for cidade := range rotas {
// 		if origem == cidade {
// 			existe_origem = true
// 		}

// 		if destino == cidade {
// 			existe_destino = true
// 			//indice_destino = Buscarindice(cidade)
// 		}
// 	}

// 	// Verificar se ambas as cidades existem e se há uma rota
// 	if existe_origem && existe_destino {
// 		return true
// 	}

// 	return false
// }

func Cadastrar(conn net.Conn, cpf string) {
	user := User{
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
	// fmt.Println("Enviando dados:", string(jsonData)) // Exibe o JSON como string
	conn.Write(jsonData)
	conn.Write([]byte("\n")) // Enviar uma nova linha para indicar o fim da mensagem

	// Ler a resposta do servidor
	// buffer := make([]byte, 1024)
	// _, err := conn.Read(buffer)
	// if err != nil {
	// 	fmt.Println("Erro ao ler a resposta do servidor:", err)
	// 	return
	// }
	// fmt.Println("Resposta do servidor:", string(buffer[:n]))
}

func Comprar(conn net.Conn, user User, origem string, destino string) {

	caminho, _ := menorCaminhoDFS(origem, destino)
	if len(caminho) > 0 {
		fmt.Printf("Rota encontrada - %s a %s: %v", origem, destino, caminho)

		compra := Compra{
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

	} else {
		fmt.Printf("Rota não encontrada")
	}
}

func VerPassagensCompradas(conn net.Conn, cpf string) {
	// Preparar a solicitação para ler as compras
	user := User{
		Cpf: cpf,
	}
	dados := Dados{
		Request:      LERCOMPRAS,
		DadosCompra:  nil,
		DadosUsuario: &user, // Deve ser um ponteiro
	}

	// Converter dados para JSON
	jsonData, err := json.Marshal(dados)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}

	// Enviar o JSON ao servidor
	// fmt.Println("Enviando dados:", string(jsonData)) // Exibe o JSON como string
	conn.Write(jsonData)
	conn.Write([]byte("\n")) // Enviar uma nova linha para indicar o fim da mensagem

	// Ler a resposta do servidor
	buffer := make([]byte, 4096) // Aumentar o buffer se necessário
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Erro ao ler a resposta do servidor:", err)
		return
	}

	response := string(buffer[:n])
	// Verificar se a resposta indica que há passagens compradas
	if response != "null" {
		fmt.Println("Passagens Compradas:")
		fmt.Println(response)
	} else {
		fmt.Println("Nenhuma passagem comprada")
	}
}
