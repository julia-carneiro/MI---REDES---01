package funcoesCliente

var cidades = []string{"São Paulo", "Salvador", "Recife"}
var rotas = [][]int{
	{0, 1, 1},
	{1, 0, 1},
	{1, 1, 0},
}

func VerificarRota(origem string, destino string) bool {
	var indice_origem, indice_destino int
	existe_origem := false
	existe_destino := false

	// Encontrar os índices das cidades
	for i, cidade := range cidades {
		if origem == cidade {
			existe_origem = true
			indice_origem = i
		}

		if destino == cidade {
			existe_destino = true
			indice_destino = i
		}
	}

	// Verificar se ambas as cidades existem e se há uma rota
	if existe_origem && existe_destino {
		return rotas[indice_origem][indice_destino] == 1
	}

	return false
}
