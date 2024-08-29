package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)

var cidades = []string{"São Paulo", "Salvador", "Recife"}
var vagas = [][]int{
    {0, 2, 3},
    {1, 0, 1},
    {1, 1, 0},
}

// Função para verificar vagas e atualizá-las
func verifica_vagas(origem string, destino string) bool {
    var indice_origem, indice_destino int
    var origemEncontrada, destinoEncontrado bool

    // Encontrar os índices de origem e destino
    for i, cidade := range cidades {
        if cidade == origem {
            indice_origem = i
            origemEncontrada = true
        }
        if cidade == destino {
            indice_destino = i
            destinoEncontrado = true
        }
    }

    // Se qualquer uma das cidades não for encontrada
    if !origemEncontrada || !destinoEncontrado {
        return false
    }

    // Verificar se há vagas e atualizar
    if vagas[indice_origem][indice_destino] > 0 {
        vagas[indice_origem][indice_destino] -= 1
        return true
    }

    return false
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
    aprovado := verifica_vagas(origem, destino)
    var result string
    if aprovado {
        result = "APROVADA"
    } else {
        result = "RECUSADA"
    }

    // Enviar resposta ao cliente
    conn.Write([]byte("Sua compra foi " + result + "\n"))
}
