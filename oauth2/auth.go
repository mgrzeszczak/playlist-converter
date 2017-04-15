package oauth2

import (
	"encoding/base64"
	"fmt"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"runtime"
	"os/exec"
	"net"
	"encoding/json"
	"github.com/franela/goreq"
	"log"
)

const (
	redirect_url = "http://localhost:8080"
	authorize_url_format = "%s?client_id=%s&scope=%s&response_type=code&redirect_uri=%s"

	http_addr = "localhost:8080"
	tcp = "tcp"

	close_html = "resources/close.html"
)

type Credentials struct {
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Config struct {
	Spotify Credentials `json:"spotify"`
	Youtube Credentials `json:"youtube"`
}

type AuthData struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (a *AuthData) GetScopes() []string {
	return strings.Split(a.Scope, " ")
}

type AuthArgs struct {
	ClientId     string
	ClientSecret string
	AuthCodeUrl  string
	AuthTokenUrl string
	Scopes       []string
}

func Authorize(args AuthArgs) (*AuthData, error) {
	link := getAuthorizationLink(args)

	code, err := getAuthCode(link)
	if err != nil {
		return nil, err
	}

	data, err := getAuthData(args, code)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getScope(scopes... string) string {
	val := strings.Join(scopes, " ")
	return url.PathEscape(val)
}

func getAuthorizationLink(args AuthArgs) string {
	scope := getScope(args.Scopes...)
	return fmt.Sprintf(authorize_url_format, args.AuthCodeUrl, args.ClientId, scope, redirect_url)
}

func getAuthCode(authorizeUrl string) (string, error) {
	openUrl(authorizeUrl)
	var err error
	listener, err := net.Listen(tcp, http_addr)
	if err != nil {
		return "", err
	}
	var code string
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		code = request.URL.Query().Get("code")
		var data []byte
		data, err = ioutil.ReadFile(close_html)
		if err != nil {
			return
		}
		wrote := 0
		for wrote < len(data) {
			var count int
			count, err = writer.Write(data[wrote:])
			if err != nil {
				return
			}
			wrote += count
		}
		listener.Close()
	});
	log.Println("Waiting for authorization...")
	http.Serve(listener, nil)
	if err != nil {
		return "", err
	}
	http.DefaultServeMux = http.NewServeMux()
	log.Println("Success")
	return code, nil
}

func getAuthData(args AuthArgs, code string) (*AuthData, error) {
	link := fmt.Sprintf("%s?grant_type=authorization_code&code=%s&redirect_uri=%s", args.AuthTokenUrl, code, redirect_url)
	authorization_value := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", args.ClientId, args.ClientSecret)))


	request := goreq.Request{
		Uri: link,
		Method: "POST",
		ContentType:"application/x-www-form-urlencoded",
	}


	request.AddHeader("Authorization", fmt.Sprintf("Basic %s", authorization_value))

	resp, err := request.Do();
	if err != nil {
		return nil, err
	}

	data := AuthData{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// source:
// http://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
func openUrl(url string) error {
	log.Printf("Browser tab will open soon\nIf it doesn't go to the following link:\n\n%s\n\n",url)
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}