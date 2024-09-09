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
	"São Paulo":{0, 2, 3},
	"Salvador":{1, 0, 1},
	"Recife":{1, 1, 0},
}

var rotas = map[string][]int{
	"São Paulo":{1, 1, 1},
	"Salvador":{1, 1, 1},
	"Recife":{1, 1, 1},
}

type Request int

const (
	GET Request = iota
	POST
)

func (s Request) String() string {
	switch s {
	case GET:

	case POST:
		return "POST"

	}
	return ""
}
type CidadeInfo struct {
	Nome  string `json:"nome"`
	Vagas []int  `json:"vagas"`
	Rotas []int  `json:"rotas"`
}

type Compra struct {
	nome    string
	cpf     string
	origem  string
	destino string
}

// var metodo Request = GET
// fmt.Println(metodo.String())

type Dados struct {
	request     Request //Tipo de requisição
	dadosCompra *Compra //Caso seja um post, as informações da viagem e do passageiro

}

// Função para verificar vagas e atualizá-las
func ValidarCompra(info Compra ) bool {
	var indice_destino int
	var origemEncontrada, destinoEncontrado bool
	var origem string = info.origem
	var destino string =  info.destino
	// Encontrar os índices de origem e destino
	for cidade,_ := range rotas {
		if cidade == origem {
			
			
			origemEncontrada = true		
			
		}
		if cidade == destino {
			indice_destino = Buscarindice(cidade)
			destinoEncontrado = true
		}
	}

	// Se qualquer uma das cidades não for encontrada
	if !origemEncontrada || !destinoEncontrado {
		return false
	}

	// Verificar se há vagas e atualizar
	if vagas[origem][indice_destino] > 0 {
		vagas[origem][indice_destino] -= 1
		return true
	}

	return false
}

//Retorna o indice da cidade
func Buscarindice(cidade string)int{
	for i, cidadeincice := range indiceCidade{//busca o indice da cidade
		if cidadeincice == cidade{
			return i
		}
	}
	return -1
}

//Encontrar todos os caminhos possíveis
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


func Get() ([]byte, error) {
	dados := map[string]interface{}{
		"vagas":   vagas,
		"rotas":   rotas,
	}

	// Converter para JSON
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

	// Leitura da mensagem do cliente
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler a mensagem:", err)
		return
	}

	var dados Dados // variável que tem todas as informações que o cliente mandou
	err = json.Unmarshal([]byte(message), &dados)
	if err != nil {
		conn.Write([]byte("Erro no formato dos dados enviados. Esperado JSON.\n"))
		return
	}

	fmt.Println("Mensagem recebida do cliente:", message)

	if dados.request == GET { //Verifica o tipo de requisição
		info, _ := Get()
		conn.Write([]byte(info))

	} else if dados.request == POST {

		// // Separar a origem e o destino
		// parts := strings.Split(strings.TrimSpace(message), ",")
		// if len(parts) != 2 {
		// 	conn.Write([]byte("Formato inválido. Use 'origem,destino'.\n"))
		// 	return
		// }
		// origem := parts[0]
		// destino := parts[1]

		// Verificar e responder
		aprovado := ValidarCompra(*dados.dadosCompra)//envia as informações da compra

		var result string
		if aprovado {
			result = "APROVADA"
		} else {
			result = "RECUSADA"
		}

		// Enviar resposta ao cliente
		conn.Write([]byte("Sua compra foi " + result + "\n"))
	}

}
