package main

import (
    "encoding/json"
    "bufio"
    "fmt"
    "net"
    "strings"
    "sessao3/servidor/funcoesServer"
)

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

type compra struct{
    nome string
    cpf string
    origem string
    destino string
}

// var metodo Request = GET
// fmt.Println(metodo.String())

type Dados struct{
    request Request  //Tipo de requisição
    atributo *string // Caso seja um get, qual da a informação que o  cliente(cidade,rota,vagas)
    dadosCompra *compra       //Caso seja um post, as informações da viagem e do passageiro
    
}


func main() {
    // Escutando na porta 8080
    ln, err := net.Listen("tcp", "172.16.103.4:8080")
    if err != nil {
        fmt.Println("Erro ao iniciar o servidor:", err)
        return
    }
    defer ln.Close()
    fmt.Println("Servidor iniciado na porta 8080...")

    for {
        // Aceitando conexões dos clientes
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Erro ao aceitar conexão:", err)
            continue
        }

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
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
    
    if(dados.request = GET){

    }
    fmt.Println("Mensagem recebida do cliente:", message)

    // Separar a origem e o destino
    parts := strings.Split(strings.TrimSpace(message), ",")
    if len(parts) != 2 {
        conn.Write([]byte("Formato inválido. Use 'origem,destino'.\n"))
        return
    }
    origem := parts[0]
    destino := parts[1]

    // Verificar e responder
    aprovado := funcoesServer.VerificaVagas(origem, destino)
    
    var result string
    if aprovado {
        result = "APROVADA"
    } else {
        result = "RECUSADA"
    }

    // Enviar resposta ao cliente
    conn.Write([]byte("Sua compra foi " + result + "\n"))
}