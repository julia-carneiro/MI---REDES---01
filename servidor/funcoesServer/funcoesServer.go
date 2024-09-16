package funcoesServer

import (
	"bufio"
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
	Nome    string `json:"Nome"`
	Cpf     string `json:"Cpf"`
	Origem  string `json:"Origem"`
	Destino string `json:"Destino"`
}

type User struct {
	Nome string `json:"Nome"`
	Cpf  string `json:"Cpf"`
}

type Dados struct {
	Request      Request  `json:"Request"`
	DadosCompra  *Compra  `json:"DadosCompra"`
	DadosUsuario *User    `json:"DadosUsuario"`
}

func ValidarCompra(info Compra) bool {
	var indice_destino int
	var origemEncontrada, destinoEncontrado bool
	origem := info.Origem
	destino := info.Destino

	for cidade := range rotas {
		if cidade == origem {
			origemEncontrada = true
		}
		if cidade == destino {
			indice_destino = Buscarindice(cidade)
			destinoEncontrado = true
		}
	}

	if !origemEncontrada || !destinoEncontrado {
		return false
	}

	if vagas[origem][indice_destino] > 0 {
		vagas[origem][indice_destino] -= 1
		return true
	}

	return false
}

func Buscarindice(cidade string) int {
	for i, cidadeincice := range indiceCidade {
		if cidadeincice == cidade {
			return i
		}
	}
	return -1
}

func Get() ([]byte, error) {
	dados := map[string]interface{}{
		"vagas": vagas,
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
