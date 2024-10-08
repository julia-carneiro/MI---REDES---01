package main

import (
	"fmt"
	"net"
	"sessao3/servidor/funcoesServer"
)

var ADRESS string = "0.0.0.0:22356"

func main() {
	// Escutando na porta 22356
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
		//cria uma gorotine para cada conexão
		go funcoesServer.HandleConnection(conn)
	}
}
