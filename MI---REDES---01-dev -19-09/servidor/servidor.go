package main

import (
	"fmt"
	"net"
	"sessao3/servidor/funcoesServer"
)

var ADRESS string = "localhost:22355"

func main() {
	// Escutando na porta 8080
	ln, err := net.Listen("tcp", ADRESS)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer ln.Close()
	fmt.Printf("Servidor iniciado em: %s", ADRESS)

	for {
		// Aceitando conexões dos clientes
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}

		go funcoesServer.HandleConnection(conn)
	}
}
