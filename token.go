package token

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var Ccu ClientConectUser
var retToken RetornoToken

type ClientConectUser struct {
	ClientConect ClientHttpUser
}

type ClientHttpUser struct {
	Client *http.Client
	Host   string
	User   string
	Auth   string
}

// RetornoToken estrutura do Token do Sistema
type RetornoToken struct {
	Ambiente string `json:"ambiente"`
	ID       string `json:"id"`
	Perfil   string `json:"perfil"`
	Emissao  string `json:"emissao"`
	ExpiraEm string `json:"expiraEm"`
	Token    string `json:"token"`
}

func GetToken(host, user, pass string) (token string, err error) {
	Ccu = NewClientConectUser(host, user, pass)

	return Ccu.GerarToken()

}

func NewClientConectUser(host, user, pass string) ClientConectUser {
	return ClientConectUser{
		ClientConect: NewClientUser(host, user, pass),
	}
}

func NewClientUser(enderecoBase, usuario, senha string) ClientHttpUser {

	client := ClientHttpUser{
		Client: &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}},
		Host:   enderecoBase,
		User:   usuario,
		Auth: fmt.Sprintf(
			"Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", usuario, senha))),
		),
	}
	return client
}

func (c ClientHttpUser) NewRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.Host, endpoint), body)
	if err != nil {
		return nil, fmt.Errorf("erro: cmp: newrequest 1: %s", err)
	}
	req.Header.Add("Authorization", c.Auth)
	return req, nil
}

func (ccu ClientConectUser) GerarToken() (string, error) {
	req, err := ccu.ClientConect.NewRequest("POST", "", nil)
	if err != nil {
		return "", err
	}

	res, err := ccu.ClientConect.Client.Do(req)
	if err != nil {
		//	return Cartao{}, fmt.Errorf("GerarToken 2: %s", err)
		return "", err
	}
	if res.StatusCode != 200 && res.StatusCode != 201 {
		return "", fmt.Errorf("gerartoken: %s", res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	err = json.Unmarshal(b, &retToken)
	if err != nil {
		return "", err
	}

	return retToken.Token, nil
}
