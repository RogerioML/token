package main

func main() {
	estrutura, erro := GetToken("https://apihom.correios.com.br/token/v1/autentica", "17811","gogogo")

	println(estrutura, erro)
}
