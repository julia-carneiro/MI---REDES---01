package funcoesCliente

import (
	"bufio"
	"encoding/json"
	"fmt"
	"unicode"

	// "math"
	"net"
	"os"
	"sort"
	"strings"
)

type Request int

const (//Tipos de mensagens que podem ser enviadas ao servidor
	ROTAS Request = iota
	COMPRA
	CADASTRO
	LERCOMPRAS
)

type Compra struct {//Estrutura de dados de compra
	Cpf     string
	Caminho []string
}

type User struct {//Estrtura de dados de usuario
	Cpf string `json:"Cpf"`
}

type Dados struct {//Estrutura de dados de mensagem para o servidor
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

//Função para vaalidar entrada do usuário
func ValidarCPF(s string) bool {
	// Verifica se a string tem exatamente 11 caracteres
	if len(s) != 11 {
		return false
	}

	// Verifica se todos os caracteres são números
	for _, char := range s {
		if !unicode.IsDigit(char) {
			return false
		}
	}

	return true
}

//Transforam a todas as letras da string em minusculas menos a primeira letra de cada palavra
func capitalizeWords(s string) string {
	words := strings.Fields(s) // Divide a string em palavras
	for i, word := range words {
		// Converte todas as letras da palavra para minúsculas
		word = strings.ToLower(word)
		// Capitaliza a primeira letra
		runes := []rune(word)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ") // Junta as palavras novamente
}


//Verifica se a vagas em uma determinada rota
func VerificarVagas(caminho []string) bool {
	CompraValida := false
	for i := 0; i < len(caminho); i++ { //percorre as cidades da rota
		if i+1 != len(caminho) { // verifica se a cidade atual não é o destino final
			for j := 0; j < len(rotas[caminho[i]]); j++ { // percorre as cidades que a cidade atual faz rota
				if rotas[caminho[i]][j].Destino == caminho[i+1] { // verifica a rota é a rota desejada
					if rotas[caminho[i]][j].Vagas > 0 {// caso seja a rota desejada verifica se há vagas
						CompraValida = true
					} else {
						CompraValida = false
						return CompraValida//caso uma das rotas não tenha vaga retorna false
					}
				}
			}
		}
	}
	return CompraValida

}

type Caminho struct {//estrtura que guarda caminho
	Cidades []string
	Peso    int
}

// Função para buscar todos os caminhos
func BuscarTodosCaminhos(origem, destino string) []Caminho {
	var caminhos []Caminho
	var caminhoAtual []string
	caminhoAtual = append(caminhoAtual, origem)

	visitarCidades(origem, destino, caminhoAtual, 0, &caminhos)

	// Ordena a lista de caminhos pelo peso total (menor caminho primeiro)
	sort.Slice(caminhos, func(i, j int) bool {
		return caminhos[i].Peso < caminhos[j].Peso
	})

	return caminhos
}

// Função recursiva para visitar cidades e encontrar todos os caminhos
func visitarCidades(origem, destino string, caminhoAtual []string, pesoAtual int, caminhos *[]Caminho) {
	if origem == destino {
		// Adiciona o caminho encontrado à lista de caminhos
		novoCaminho := make([]string, len(caminhoAtual))
		copy(novoCaminho, caminhoAtual)
		*caminhos = append(*caminhos, Caminho{Cidades: novoCaminho, Peso: pesoAtual})
		return
	}

	for _, rota := range rotas[origem] {
		if !contem(caminhoAtual, rota.Destino) { // evita ciclos
			// Continua a busca a partir do próximo destino
			visitarCidades(rota.Destino, destino, append(caminhoAtual, rota.Destino), pesoAtual+rota.Peso, caminhos)
		}
	}
}

// Função auxiliar para verificar se uma cidade já está no caminho atual (para evitar ciclos)
func contem(caminho []string, cidade string) bool {
	for _, c := range caminho {
		if c == cidade {
			return true
		}
	}
	return false
}
// Função para printar as cidades que é possível viajar
func VerCidades(){
	fmt.Print("Cidades possíveis para viajar:\n")
	for cidade,_ := range rotas{
		fmt.Printf("%s\n",cidade)

	}
}

//Menu do sistema
func Menu(ADRESS string, user User) {
	var operacao int
	var conn net.Conn
	var i = true

	for i {

		// Lendo entrada do usuário
		fmt.Println("O que deseja fazer?\n1- Comprar passagens\n2-Ver passagens compradas\n3-Ver cidades disponíveis\n4-Sair ")
		fmt.Scanf("%d\n", &operacao)

		switch operacao {
		case 1:// Fazer compra
			// Lendo entrada do usuário
			var origem, destino string
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Digite a cidade de origem:")
			origem, _ = reader.ReadString('\n')
			origem = strings.TrimSpace(origem)
			origem = capitalizeWords(origem)

			fmt.Println("Digite a cidade de destino:")
			destino, _ = reader.ReadString('\n')
			destino = strings.TrimSpace(destino)
			destino = capitalizeWords(destino)

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
			} else {
				fmt.Println("Não existe rota")
			}

		case 2:// Ver compras anteriores
			conn = ConectarServidor(ADRESS)
			VerPassagensCompradas(conn, user.Cpf)
			defer conn.Close()
		case 3://Ver cidades que é possível viajar
			conn = ConectarServidor(ADRESS)
			BuscarDados(conn)
			defer conn.Close()
			VerCidades()
			
		case 4://Sair do sistema
			i = false
			break
		default:
			fmt.Println("Operação inválida.")

		}
	}

}
//Conecta com o servidor
func ConectarServidor(ADRESS string) net.Conn {
	// Conectando ao servidor na porta 8080
	conn, err := net.Dial("tcp", ADRESS)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return nil
	}

	return conn
}

//Solicita ao servidor as rotas,pesos e vagas
func BuscarDados(conn net.Conn) {
	dados := Dados{//estrutura da mensagem enviada para o servidor
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
	buffer := make([]byte, 4096)
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

//cadastra o usuário
func Cadastrar(conn net.Conn, cpf string) {
	user := User{
		Cpf: cpf,
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


}
//Função de compra
func Comprar(conn net.Conn, user User, origem string, destino string) {
	var caminho_final []string
	var caminhos []Caminho = BuscarTodosCaminhos(origem, destino)// busca todos os caminho possiveis
	var vagas = false

	//verifica se o menor caminho encontrado possui vagas, caso não tenha procura nos proximos menores caminhos
	for i := 0; i < len(caminhos); i++ {
		vagas = VerificarVagas(caminhos[i].Cidades)
		if vagas {
			caminho_final = caminhos[i].Cidades

			break
		}
	}

	if len(caminho_final) > 0 {
		if vagas {
			fmt.Printf("Rota encontrada - %s a %s: ", origem, destino)//mostra a rota
			for i := 0; i < len(caminho_final); i++ {
				if(i+1 != len(caminho_final)){

					fmt.Printf("%s-> ",caminho_final[i])
				}else{
					fmt.Printf("%s\n", caminho_final[i])
				}
			}
			compra := Compra{
				Cpf:     user.Cpf,
				Caminho: caminho_final,
			}

			dados := Dados{
				Request:      COMPRA,
				DadosCompra:  &compra,
				DadosUsuario: nil,
			}
			//Verifca se o usuario deseja comprar as passagens
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
		} else {
			fmt.Println("Não há vagas para essa rota.")
		}

	} else {

		fmt.Println("Rota não encontrada")

	}
}
//Exibe as passagens compradas
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

	var paths [][]string
    err = json.Unmarshal(buffer[:n], &paths) // Desserializar para lista de listas
    if err != nil {
        fmt.Println("Erro ao desserializar a resposta do servidor:", err)
        return
    }
	// Verificar se a resposta indica que há passagens compradas
	if len(paths) != 0 {
		fmt.Println("Passagens Compradas:")
		for i := 0; i < len(paths); i++ {
			fmt.Printf("Compra %d:\n",i+1)
			for j := 0; j < len(paths[i]); j++ {
				if(j+1 != len(paths[i])){

					fmt.Printf("%s-> ",paths[i][j])
				}else{
					fmt.Printf("%s",paths[i][j])
				}
			}
			fmt.Print("\n\n")
		}
	} else {
		fmt.Println("Nenhuma passagem comprada")
	}
}
