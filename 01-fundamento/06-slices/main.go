// Pacote main com exemplos de uso de slices em Go.
//
// O que é um slice?
//
// Um slice é uma "janela" sobre um array. Diferente do array, que tem tamanho
// fixo ([5]int), o slice ([]int) pode crescer. Por baixo dos panos, um slice
// é uma pequena estrutura com três campos:
//
//	ponteiro -> aponta para um elemento do array subjacente (underlying array)
//	len      -> quantos elementos o slice enxerga agora
//	cap      -> quantos elementos cabem a partir do ponteiro até o fim do array
//
// Entender esses três campos explica praticamente todo o comportamento dos
// slices: fatiamento, append, compartilhamento de memória etc.
//
// Execute com: go run main.go
package main

import (
	"fmt"
	"slices"
)

func main() {
	exemploBasico()
	exemploDeclaracao()
	exemploNilVsVazio()
	exemploFatiamento()
	exemploAppendECrescimento()
	exemploCompartilhamento()
	exemploCopy()
	exemploTresIndices()
	exemploRemoverEInserir()
	exemploIteracao()
	exemploMultidimensional()
	exemploPassandoParaFuncao()
	exemploPacoteSlices()
}

// imprime mostra len, cap e o conteúdo de um slice. É usada em todos os exemplos.
func imprime(nome string, s []int) {
	fmt.Printf("%-12s len=%d cap=%d %v\n", nome, len(s), cap(s), s)
}

// titulo imprime um cabeçalho para separar os exemplos na saída.
func titulo(t string) {
	fmt.Printf("\n===== %s =====\n", t)
}

// exemploBasico é o exemplo original da aula: len e cap ao fatiar.
func exemploBasico() {
	titulo("1. Básico: len e cap")

	slice := []int{10, 20, 30, 40, 50}

	imprime("slice", slice)         // len=5 cap=5
	imprime("slice[:0]", slice[:0]) // len=0 cap=5 -> vazio, mas o array continua lá
	imprime("slice[:4]", slice[:4]) // len=4 cap=5 -> corta o final, cap não muda
	imprime("slice[2:]", slice[2:]) // len=3 cap=3 -> corta o início, cap DIMINUI

	// Quando o append ultrapassa a capacidade, Go cria um array novo
	// (normalmente com o dobro do tamanho) e copia os elementos.
	slice = append(slice, 60)
	imprime("após append", slice) // len=6 cap=10
}

// exemploDeclaracao mostra as formas de criar um slice.
func exemploDeclaracao() {
	titulo("2. Formas de declarar")

	// Literal: já cria com valores. len = cap = quantidade de elementos.
	literal := []int{1, 2, 3}
	imprime("literal", literal)

	// make(tipo, len): cria com len elementos zerados.
	comLen := make([]int, 3)
	imprime("make(3)", comLen) // [0 0 0]

	// make(tipo, len, cap): reserva espaço extra. Útil quando você já sabe
	// quantos elementos vai adicionar, pois evita realocações no append.
	comCap := make([]int, 0, 10)
	imprime("make(0,10)", comCap) // [] mas com cap=10

	// A partir de um array: o slice aponta para o mesmo array.
	array := [5]int{1, 2, 3, 4, 5}
	doArray := array[1:4]
	imprime("array[1:4]", doArray) // [2 3 4] cap=4
}

// exemploNilVsVazio mostra a diferença entre slice nil e slice vazio.
func exemploNilVsVazio() {
	titulo("3. Nil vs vazio")

	var nilSlice []int // valor zero de um slice: nil
	vazio := []int{}   // não é nil, mas não tem elementos

	fmt.Println("nilSlice == nil:", nilSlice == nil) // true
	fmt.Println("vazio == nil:   ", vazio == nil)    // false

	// Na prática os dois se comportam igual: len=0, dá para fazer range e append.
	imprime("nilSlice", nilSlice)
	imprime("vazio", vazio)

	nilSlice = append(nilSlice, 1) // append em slice nil funciona normalmente
	imprime("nil+append", nilSlice)

	// Dica: para verificar se está vazio, use len(s) == 0, e não s == nil.
}

// exemploFatiamento mostra a sintaxe s[inicio:fim].
func exemploFatiamento() {
	titulo("4. Fatiamento s[inicio:fim]")

	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	// O índice "inicio" é incluído e o "fim" é EXCLUÍDO.
	imprime("s[2:5]", s[2:5]) // [2 3 4]
	imprime("s[:3]", s[:3])   // [0 1 2]    -> do começo até o índice 3 (exclusivo)
	imprime("s[7:]", s[7:])   // [7 8 9]    -> do índice 7 até o fim
	imprime("s[:]", s[:])     // tudo

	// Último elemento: não existe índice negativo em Go.
	fmt.Println("último:", s[len(s)-1])

	// Acessar fora do len causa panic: s[10] -> index out of range.
	// Mas é possível fatiar além do len, até o cap:
	curto := s[:2]
	imprime("curto", curto)         // [0 1] cap=10
	imprime("curto[:5]", curto[:5]) // [0 1 2 3 4] -> "recupera" elementos escondidos
}

// exemploAppendECrescimento mostra como a capacidade cresce com o append.
func exemploAppendECrescimento() {
	titulo("5. Append e crescimento da capacidade")

	var s []int
	capAnterior := cap(s)
	for i := 0; i < 20; i++ {
		s = append(s, i)
		// Só imprime quando houve realocação (cap mudou).
		if cap(s) != capAnterior {
			fmt.Printf("len=%-2d cap mudou de %-2d para %d\n", len(s), capAnterior, cap(s))
			capAnterior = cap(s)
		}
	}

	// append aceita vários valores de uma vez...
	s2 := append([]int{}, 1, 2, 3)
	// ...e também outro slice, usando "..." para "espalhar" os elementos.
	s2 = append(s2, []int{4, 5}...)
	imprime("s2", s2)

	// IMPORTANTE: sempre atribua o retorno do append (s = append(s, x)),
	// pois ele pode devolver um slice apontando para um array novo.
}

// exemploCompartilhamento mostra que slices podem compartilhar o mesmo array.
// Essa é a principal "pegadinha" dos slices.
func exemploCompartilhamento() {
	titulo("6. Compartilhamento de memória")

	original := []int{1, 2, 3, 4, 5}
	parte := original[1:3] // [2 3], mas aponta para o MESMO array

	parte[0] = 99                 // altera também o original!
	imprime("original", original) // [1 99 3 4 5]
	imprime("parte", parte)       // [99 3]

	// Pegadinha com append: "parte" tem cap=4, então o append NÃO realoca
	// e sobrescreve o elemento original[3].
	parte = append(parte, 100)
	imprime("original", original) // [1 99 3 100 5]  <- o 4 sumiu!
	imprime("parte", parte)

	// Quando o append ultrapassa a capacidade, um array novo é criado
	// e a ligação entre os dois slices é quebrada.
	parte = append(parte, 200, 300)
	parte[0] = -1
	imprime("original", original) // não é mais afetado
	imprime("parte", parte)
}

// exemploCopy mostra como criar uma cópia independente com copy.
func exemploCopy() {
	titulo("7. Copy: cópia independente")

	origem := []int{1, 2, 3, 4, 5}

	// copy(destino, origem) copia min(len(destino), len(origem)) elementos
	// e retorna quantos foram copiados. O destino precisa ter len suficiente.
	destino := make([]int, len(origem))
	n := copy(destino, origem)
	fmt.Println("copiados:", n)

	destino[0] = 999 // não afeta a origem
	imprime("origem", origem)
	imprime("destino", destino)

	// Destino menor: copia só o que cabe.
	pequeno := make([]int, 2)
	copy(pequeno, origem)
	imprime("pequeno", pequeno) // [1 2]

	// Alternativa moderna (Go 1.21+): slices.Clone
	clone := slices.Clone(origem)
	imprime("clone", clone)
}

// exemploTresIndices mostra a "full slice expression" s[inicio:fim:max],
// que limita a capacidade do novo slice.
func exemploTresIndices() {
	titulo("8. Fatiamento com 3 índices s[i:j:k]")

	original := []int{1, 2, 3, 4, 5}

	// cap = k - i. Aqui: cap = 3 - 1 = 2, igual ao len.
	limitado := original[1:3:3]
	imprime("limitado", limitado) // [2 3] cap=2

	// Como não sobra capacidade, o append é obrigado a criar um array novo,
	// protegendo o "original" da pegadinha do exemplo 6.
	limitado = append(limitado, 100)
	imprime("original", original) // [1 2 3 4 5] intacto
	imprime("limitado", limitado)
}

// exemploRemoverEInserir mostra operações comuns de remoção e inserção.
func exemploRemoverEInserir() {
	titulo("9. Remover e inserir elementos")

	// Remover o elemento no índice i, mantendo a ordem:
	// junta tudo antes de i com tudo depois de i.
	s := []int{10, 20, 30, 40, 50}
	i := 2
	s = append(s[:i], s[i+1:]...)
	imprime("remove i=2", s) // [10 20 40 50]

	// Remover sem manter a ordem (mais rápido, O(1)):
	// troca pelo último e encolhe o slice.
	s2 := []int{10, 20, 30, 40, 50}
	s2[1] = s2[len(s2)-1]
	s2 = s2[:len(s2)-1]
	imprime("remove rápido", s2) // [10 50 30 40]

	// Inserir no índice i:
	s3 := []int{1, 2, 4, 5}
	s3 = append(s3[:2], append([]int{3}, s3[2:]...)...)
	imprime("insere i=2", s3) // [1 2 3 4 5]

	// Com o pacote slices (Go 1.21+) fica bem mais legível:
	s4 := []int{1, 2, 4, 5}
	s4 = slices.Insert(s4, 2, 3)
	imprime("slices.Insert", s4)
	s4 = slices.Delete(s4, 0, 1) // remove o intervalo [0, 1)
	imprime("slices.Delete", s4)

	// Filtrar (ex.: manter só os pares) reaproveitando o mesmo array:
	nums := []int{1, 2, 3, 4, 5, 6}
	pares := nums[:0] // len=0, mas usa o array de nums
	for _, n := range nums {
		if n%2 == 0 {
			pares = append(pares, n)
		}
	}
	imprime("pares", pares) // [2 4 6]
}

// exemploIteracao mostra as formas de percorrer um slice.
func exemploIteracao() {
	titulo("10. Iteração")

	frutas := []string{"maçã", "banana", "uva"}

	// range devolve índice e uma CÓPIA do valor.
	for i, fruta := range frutas {
		fmt.Printf("%d: %s\n", i, fruta)
	}

	// Ignorando o índice com _
	for _, fruta := range frutas {
		fmt.Print(fruta, " ")
	}
	fmt.Println()

	// Como o valor é uma cópia, para alterar o slice use o índice:
	nums := []int{1, 2, 3}
	for _, n := range nums {
		n *= 10 // não altera nums
		_ = n
	}
	imprime("sem efeito", nums)

	for i := range nums {
		nums[i] *= 10 // altera nums
	}
	imprime("com índice", nums)
}

// exemploMultidimensional mostra slice de slices (matriz).
func exemploMultidimensional() {
	titulo("11. Slice de slices (matriz)")

	// Cada linha é um slice independente, e podem ter tamanhos diferentes.
	matriz := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	for _, linha := range matriz {
		fmt.Println(linha)
	}
	fmt.Println("matriz[1][2] =", matriz[1][2]) // 6

	// Criando uma matriz linhas x colunas com make:
	linhas, colunas := 2, 4
	grade := make([][]int, linhas)
	for i := range grade {
		grade[i] = make([]int, colunas)
	}
	grade[0][3] = 1
	fmt.Println("grade:", grade)
}

// dobra recebe um slice e altera seus elementos. Como o slice carrega um
// ponteiro para o array, a função altera os dados de quem chamou.
func dobra(s []int) {
	for i := range s {
		s[i] *= 2
	}
}

// adiciona faz append dentro da função. O len do slice de quem chamou NÃO muda,
// porque o slice (ponteiro, len, cap) é passado por valor. Por isso é preciso
// retornar o novo slice.
func adiciona(s []int, v int) []int {
	return append(s, v)
}

// exemploPassandoParaFuncao mostra o comportamento de slices como parâmetro.
func exemploPassandoParaFuncao() {
	titulo("12. Slices em funções")

	nums := []int{1, 2, 3}
	dobra(nums)
	imprime("após dobra", nums) // [2 4 6] -> alterado

	adiciona(nums, 99)           // retorno ignorado
	imprime("sem retorno", nums) // [2 4 6] -> continua igual

	nums = adiciona(nums, 99)
	imprime("com retorno", nums) // [2 4 6 99]
}

// exemploPacoteSlices mostra funções úteis do pacote padrão "slices" (Go 1.21+).
func exemploPacoteSlices() {
	titulo("13. Pacote slices")

	s := []int{5, 2, 8, 1, 9, 3}

	fmt.Println("Contains 8:", slices.Contains(s, 8))
	fmt.Println("Index de 1:", slices.Index(s, 1))
	fmt.Println("Max:", slices.Max(s), "Min:", slices.Min(s))

	ordenado := slices.Clone(s)
	slices.Sort(ordenado) // ordena no próprio slice
	imprime("Sort", ordenado)

	// BinarySearch só funciona em slices ordenados.
	pos, achou := slices.BinarySearch(ordenado, 8)
	fmt.Println("BinarySearch 8:", pos, achou)

	slices.Reverse(ordenado)
	imprime("Reverse", ordenado)

	fmt.Println("Equal:", slices.Equal([]int{1, 2}, []int{1, 2}))

	// Compact remove duplicados CONSECUTIVOS (ordene antes para remover todos).
	dup := []int{1, 1, 2, 3, 3, 3, 4}
	imprime("Compact", slices.Compact(dup))
}
