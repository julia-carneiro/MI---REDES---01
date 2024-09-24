package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sessao3/cliente/funcoesCliente"
	"strings"
	
)

var ADRESS string = "servidor-container:22356"
// type User struct {
// 	Cpf  string `json:"Cpf"`
// }

func main() {
	var cpf_valido = false
	//Roda enquanto o cpf for inválido
	for !cpf_valido{
		//Logar
		fmt.Println("----- Faça login/cadastro-----")
		var conn net.Conn

		var  cpf string
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Digite seu CPF(Apenas números):")
		cpf, _ = reader.ReadString('\n')
		cpf = strings.TrimSpace(cpf)

		
		user := funcoesCliente.User{
			Cpf:  cpf,
		}
		
		cpf_valido = funcoesCliente.ValidarCPF(cpf)
		
		if cpf_valido{
			conn = funcoesCliente.ConectarServidor(ADRESS)
			funcoesCliente.Cadastrar(conn, cpf)
			defer conn.Close()
			funcoesCliente.Menu(ADRESS, user)
			
		}else{
			fmt.Printf("CPF inválido\n")
		}
}
}

