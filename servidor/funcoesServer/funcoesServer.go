package funcoesServer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Request int

const (
	ROTAS Request = iota
	COMPRA
	CADASTRO
)

func (s Request) String() string {
	switch s {
	case ROTAS:
		return "ROTAS"
	case COMPRA:
		return "COMPRA"
	case CADASTRO:
		return "CADASTRO"
	}
	return "DESCONHECIDO"
}

type Compra struct {
	Nome    string   `json:"Nome"`
	Cpf     string   `json:"Cpf"`
	Caminho []string `json:"Caminho"`
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

type Rota struct {
	Destino string `json:"Destino"`
	Vagas   int    `json:"Vagas"`
	Peso    int    `json:"Peso"`
}

var rotas map[string][]Rota

//Busca e lê o arquivos de rotas
func BuscarArquivosRotas() map[string][]Rota {
	// Defina o caminho do arquivo JSON
	filePath := `C:\Users\thiag\OneDrive\Documentos\Meus projetos\MI---REDES---01\dados\rotas.json`

	// Abra o arquivo
	file, err := os.Open(filePath)
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

//Está errada, não funciona ainda
func Get() ([]byte, error) {
	dados := map[string]interface{}{

		"rotas": rotas,
	}

	jsonData, err := json.MarshalIndent(dados, "", "  ")
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return nil, err
	}
	return jsonData, nil
}

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
		info, _ := Get()
		conn.Write(info) // Envia os bytes diretamente
	case COMPRA:
		if dados.DadosCompra == nil {
			conn.Write([]byte("Dados de compra não fornecidos.\n"))
			return
		}

		aprovado := ValidarCompra(*dados.DadosCompra)
		var result string
		if aprovado {
			//Subtrair o numero de vagas nas rotas
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

		conn.Write([]byte("Usuário cadastrado\n"))
	default:
		conn.Write([]byte("Tipo de requisição inválido.\n"))
	}
}
