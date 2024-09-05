package funcoesServer

var cidades = []string{"São Paulo", "Salvador", "Recife"}
var vagas = [][]int{
    {0, 2, 3},
    {1, 0, 1},
    {1, 1, 0},
}

var rotas = [][]int{
    {1,2,3},{1,2,3},
}

// Função para verificar vagas e atualizá-las
func VerificaVagas(origem string, destino string) bool {
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

func Get(){

}
