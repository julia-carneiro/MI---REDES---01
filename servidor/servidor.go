package main

import (
	"fmt"
	"net"
	"sessao3/servidor/funcoesServer"
)

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

		go funcoesServer.HandleConnection(conn)
	}
}
