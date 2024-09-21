package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sessao3/cliente/funcoesCliente"
)

var ADRESS string = "0.0.0.0:22355"
// type User struct {
// 	Cpf  string `json:"Cpf"`
// }

func main() {
	var cpf_valido = true

	for cpf_valido{
		//Logar
		fmt.Println("----- Faça login/cadastro-----")
		var conn net.Conn

		var  cpf string
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Digite seu CPF:")
		cpf, _ = reader.ReadString('\n')
		cpf = strings.TrimSpace(cpf)

		conn = funcoesCliente.ConectarServidor(ADRESS)
		funcoesCliente.Cadastrar(conn, cpf)
		defer conn.Close()

		user := funcoesCliente.User{
			Cpf:  cpf,
		}

		if user.Cpf != ""{
			funcoesCliente.Menu(ADRESS, user)
			cpf_valido = false
		}else{
			fmt.Printf("CPF não fornecido\n")
		}
}
}

