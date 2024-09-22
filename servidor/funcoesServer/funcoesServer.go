package funcoesServer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	// "path/filepath"
)

type Request int

const (
	ROTAS Request = iota
	COMPRA
	CADASTRO
	LERCOMPRAS
)

func (s Request) String() string {
	switch s {
	case ROTAS:
		return "ROTAS"
	case COMPRA:
		return "COMPRA"
	case CADASTRO:
		return "CADASTRO"
	case LERCOMPRAS:
		return "LERCOMPRAS"
	}
	return "DESCONHECIDO"
}

type Compra struct {
	Cpf     string   `json:"Cpf"`
	Caminho []string `json:"Caminho"`
}
type User struct {
	Cpf  string `json:"Cpf"`
}

type Dados struct {
	Request      Request `json:"Request"`
	DadosCompra  *Compra `json:"DadosCompra"`
	DadosUsuario *User   `json:"DadosUsuario"`
}

type Rota struct {
	Destino string `json:"Destino"`
	Vagas   int    `json:"Vagas"`
	Peso    int    `json:"Peso"`
}

var rotas map[string][]Rota
var filePathRotas = "/app/dados/rotas.json"
var filePathUsers = "/app/dados/users.json"
var filePathCompras = "/app/dados/compras.json"

func SalvarCompra(compra Compra) error {
	// Ler o conteúdo existente do arquivo de compras
	content, err := os.ReadFile(filePathCompras)
	if err != nil {
		if os.IsNotExist(err) {
			// Se o arquivo não existir, inicialize uma lista vazia de compras
			content = []byte("[]")
		} else {
			return fmt.Errorf("erro ao ler o arquivo de compras: %v", err)
		}
	}

	// Tratar caso o arquivo esteja vazio
	if len(content) == 0 {
		content = []byte("[]")
	}

	// Decodificar o conteúdo do arquivo para uma lista de compras
	var compras []struct {
		Cpf     string     `json:"Cpf"`
		Caminho [][]string `json:"Caminho"`
	}
	err = json.Unmarshal(content, &compras)
	if err != nil {
		return fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	// Verificar se já existe uma entrada para o CPF fornecido
	var usuarioExistente *struct {
		Cpf     string     `json:"Cpf"`
		Caminho [][]string `json:"Caminho"`
	}
	for i, c := range compras {
		if c.Cpf == compra.Cpf {
			usuarioExistente = &compras[i]
			break
		}
	}

	if usuarioExistente != nil {
		// Adicionar a nova rota à lista existente
		usuarioExistente.Caminho = append(usuarioExistente.Caminho, compra.Caminho)
	} else {
		// Adicionar nova compra se não houver entrada para o CPF
		compras = append(compras, struct {
			Cpf     string     `json:"Cpf"`
			Caminho [][]string `json:"Caminho"`
		}{
			Cpf:     compra.Cpf,
			Caminho: [][]string{compra.Caminho},
		})
	}

	// Converter a lista atualizada de volta para JSON
	jsonData, err := json.MarshalIndent(compras, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao converter dados para JSON: %v", err)
	}

	// Escrever o JSON atualizado no arquivo
	err = os.WriteFile(filePathCompras, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("erro ao escrever no arquivo de compras: %v", err)
	}

	return nil
}

func LerCompras(cpf string) ([][]string, error) {
	// Ler o conteúdo do arquivo de compras
	content, err := os.ReadFile(filePathCompras)
	if err != nil {
		if os.IsNotExist(err) {
			// Se o arquivo não existir, retornar uma lista vazia de compras
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao ler o arquivo de compras: %v", err)
	}

	// Tratar caso o arquivo esteja vazio
	if len(content) == 0 {
		return nil, nil
	}

	// Decodificar o conteúdo do arquivo para uma lista de compras
	var compras []struct {
		Cpf     string     `json:"Cpf"`
		Caminho [][]string `json:"Caminho"`
	}
	err = json.Unmarshal(content, &compras)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	// Procurar a compra com o CPF fornecido
	for _, c := range compras {
		if c.Cpf == cpf {
			// Retornar a lista de rotas do usuário encontrado
			return c.Caminho, nil
		}
	}

	// Se o usuário não for encontrado, retornar uma lista vazia
	return nil, nil
}

// Busca e lê o arquivos de rotas
func BuscarArquivosRotas() map[string][]Rota {
	// Defina o caminho do arquivo JSON
	// filePath := `C:\Users\thiag\OneDrive\Documentos\Meus projetos\MI---REDES---01\dados\rotas.json`

	// Abra o arquivo
	file, err := os.Open(filePathRotas)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return nil
	}
	defer file.Close()

	// Criar um mapa para armazenar as rotas
	var rotas map[string][]Rota

	// Decodificar o arquivo JSON para o mapa
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&rotas); err != nil {
		fmt.Println("Erro ao decodificar o JSON:", err)
		return nil
	}
	return rotas
}

func CadastrarUsuario(novoUsuario User) error {
	// Ler o conteúdo existente do arquivo
	content, err := os.ReadFile(filePathUsers)
	if err != nil {
		if os.IsNotExist(err) {
			// Se o arquivo não existir, inicialize uma lista vazia de usuários
			content = []byte("[]")
		} else {
			return fmt.Errorf("erro ao ler o arquivo de usuários: %v", err)
		}
	}

	// Tratar caso o arquivo esteja vazio
	if len(content) == 0 {
		content = []byte("[]")
	}

	// Decodificar o conteúdo do arquivo para uma lista de usuários
	var users []User
	err = json.Unmarshal(content, &users)
	if err != nil {
		return fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	// Verificar se o CPF já está cadastrado
	for _, u := range users {
	 	if u.Cpf == novoUsuario.Cpf {
	 		fmt.Printf("Logando com CPF %s", novoUsuario.Cpf)
			return nil
	 	}
    }

	// Adicionar o novo usuário à lista
	users = append(users, novoUsuario)

	// Converter a lista atualizada de volta para JSON
	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao converter dados para JSON: %v", err)
	}

	// Escrever o JSON atualizado no arquivo
	err = os.WriteFile(filePathUsers, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("erro ao escrever no arquivo de usuários: %v", err)
	}

	return nil
}

func AtualizarVagas(info Compra) {
	for i := 0; i < len(info.Caminho); i++ { //percorre as cidades da rota
		if i+1 != len(info.Caminho) { // verifica se a cidade atual não é o destino final
			for j := 0; j < len(rotas[info.Caminho[i]]); j++ { // percorre as cidades que a cidade atual faz rota
				if rotas[info.Caminho[i]][j].Destino == info.Caminho[i+1] { // verifica se a rota é a rota desejada
					if rotas[info.Caminho[i]][j].Vagas > 0 { // caso seja a rota desejada verifica se há vagas
						rotas[info.Caminho[i]][j].Vagas -= 1 // diminue uma vaga no trecho atual
					}
				}
			}
		}

	}

	// Converter dados para JSON
	jsonData, err := json.MarshalIndent(rotas, "", "  ")
	if err != nil {
		fmt.Println("Erro ao converter dados para JSON:", err)
		return
	}

	// Abrir ou criar o arquivo para sobrescrever
	file, err := os.Create(filePathRotas)
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return
	}
	defer file.Close()

	// Escrever JSON no arquivo
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Erro ao escrever no arquivo:", err)
		return
	}
}

// Verifica se há vagas nas rotas que o usuário deseja comprar
// Depois é necessario subtrair o número de vagas
func ValidarCompra(info Compra) bool {
	rotas = BuscarArquivosRotas()
	CompraValida := false
	for i := 0; i < len(info.Caminho); i++ { //percorre as cidades da rota
		if i+1 != len(info.Caminho) { // verifica se a cidade atual não é o destino final
			for j := 0; j < len(rotas[info.Caminho[i]]); j++ { // percorre as cidades que a cidade atual faz rota
				if rotas[info.Caminho[i]][j].Destino == info.Caminho[i+1] { // verifica a rota é a rota desejada
					if rotas[info.Caminho[i]][j].Vagas > 0 { // caso seja a rota desejada verifica se há vagas
						CompraValida = true
					} else {
						CompraValida = false
					}
				}
			}
		}
	}
	return CompraValida

}

// func Get() ([]byte, error) {
// 	rotas := BuscarArquivosRotas()

// 	jsonData, err := json.MarshalIndent(rotas, "", "  ")
// 	if err != nil {
// 		fmt.Println("Erro ao converter para JSON:", err)
// 		return nil, err
// 	}
// 	return jsonData, nil
// }

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Cliente conectado:", conn.RemoteAddr())

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler a mensagem:", err)
		return
	}

	var dados Dados
	err = json.Unmarshal([]byte(message), &dados)
	if err != nil {
		conn.Write([]byte("Erro no formato dos dados enviados. Esperado JSON.\n"))
		return
	}

	fmt.Println("Mensagem recebida do cliente:", dados)

	switch dados.Request {
	case ROTAS:
		rotas := BuscarArquivosRotas()

		jsonData, err := json.MarshalIndent(rotas, "", "  ")
		if err != nil {
			fmt.Println("Erro ao converter para JSON:", err)

		}
		conn.Write(jsonData)     // Envia os bytes diretamente
		conn.Write([]byte("\n")) // Enviar uma nova linha para indicar o fim da mensagem

	case COMPRA:
		if dados.DadosCompra == nil {
			conn.Write([]byte("Dados de compra não fornecidos.\n"))
			return
		}

		aprovado := ValidarCompra(*dados.DadosCompra)
		var result string
		if aprovado {
			//Subtrair o numero de vagas nas rotas
			AtualizarVagas(*dados.DadosCompra)

			// Salvar a compra no arquivo "compras.json"
			err := SalvarCompra(*dados.DadosCompra)
			if err != nil {
				conn.Write([]byte(fmt.Sprintf("Erro ao salvar a compra: %v\n", err)))
				return
			}

			result = "APROVADA"
		} else {
			result = "RECUSADA"
		}

		conn.Write([]byte("Sua compra foi " + result + "\n"))
	case CADASTRO:

		if dados.DadosUsuario == nil {
			conn.Write([]byte("Dados de usuário não fornecidos.\n"))
			return
		}

		// Tentar cadastrar o usuário
		err := CadastrarUsuario(*dados.DadosUsuario)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("Erro ao cadastrar usuário: %v\n", err)))
			return
		}

		// Confirmar sucesso
		conn.Write([]byte("Operação realizada com sucesso.\n"))

	case LERCOMPRAS:
		if dados.DadosUsuario == nil {
			conn.Write([]byte("Dados de usuário não fornecidos.\n"))
			return
		}

		// Supondo que DadosUsuario contém o CPF do usuário
		cpf := dados.DadosUsuario.Cpf
		if cpf == "" {
			conn.Write([]byte("CPF não fornecido.\n"))
			return
		}

		// Ler as compras do usuário
		compras, err := LerCompras(cpf)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("Erro ao ler as compras: %v\n", err)))
			return
		}

		// Converter as compras para JSON
		jsonData, err := json.MarshalIndent(compras, "", "  ")
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("Erro ao converter dados para JSON: %v\n", err)))
			return
		}

		// Enviar os dados para o cliente
		conn.Write(jsonData)
		conn.Write([]byte("\n")) // Enviar uma nova linha para indicar o fim da mensagem

	default:
		conn.Write([]byte("Tipo de requisição inválido.\n"))
	}
}
