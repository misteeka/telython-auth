package main

import (
	"bufio"
	"bytes"
	"fmt"
	transport "github.com/misteeka/fasthttp"
	"os"
	"strings"
)

type Response []byte

var (
	SUCCESS        = Response{0}
	NOT_FOUND      = Response{1}
	WRONG_PASSWORD = Response{2}
	EXISTS         = Response{3}
	SERVER_ERROR   = Response{5}
	WRONG_DATA     = Response{6}
	YES            = Response{7}
	NO             = Response{8}
)

func init() {

}

func responseToString(response Response) string {
	if bytes.Equal(response, SUCCESS) {
		return "SUCCESS"
	}
	if bytes.Equal(response, NOT_FOUND) {
		return "NOT_FOUND"
	}
	if bytes.Equal(response, WRONG_PASSWORD) {
		return "WRONG_PASSWORD"
	}
	if bytes.Equal(response, EXISTS) {
		return "EXISTS"
	}
	if bytes.Equal(response, SERVER_ERROR) {
		return "SERVER_ERROR"
	}
	if bytes.Equal(response, WRONG_DATA) {
		return "WRONG_DATA"
	}
	if bytes.Equal(response, YES) {
		return "YES"
	}
	if bytes.Equal(response, NO) {
		return "NO"
	}
	return fmt.Sprintf("%s", response)

}

func get(function string) ([]byte, error) {
	resp, err := transport.Get("http://127.0.0.1:8001/auth/" + function)
	if err != nil {
		return nil, err
	}
	response := resp.Body()
	transport.ReleaseResponse(resp)
	return response, nil
}
func post(function string, json string) ([]byte, error) {
	resp, err := transport.Post("http://127.0.0.1:8001/auth/"+function, []byte(json))
	if err != nil {
		return nil, err
	}
	response := resp.Body()
	transport.ReleaseResponse(resp)
	return response, nil
}
func put(function string, json string) ([]byte, error) {
	resp, err := transport.Put("http://127.0.0.1:8001/auth/"+function, []byte(json))
	if err != nil {
		return nil, err
	}
	response := resp.Body()
	transport.ReleaseResponse(resp)
	return response, nil
}

func SignIn(username string, password string) (Response, error) {
	return put("signIn", fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password))
}
func CheckPassword(username string, password string) (Response, error) {
	return get(fmt.Sprintf("checkPassword?u=%s&p=%s", username, password))
}
func ResetPassword(username string, oldPassword string, newPassword string) (Response, error) {
	return put("resetPassword", fmt.Sprintf(`{"username":"%s", "oldPassword":"%s", "newPassword":"%s"}`, username, oldPassword, newPassword))
}
func RequestSignUpCode(username string, email string) (Response, error) {
	return post("requestSignUpCode", fmt.Sprintf(`{"username":"%s", "email":"%s"}`, username, email))
}
func RequestPasswordRecovery(username string) (Response, error) {
	return put("requestPasswordRecovery", fmt.Sprintf(`{"username":"%s"}`, username))
}
func RecoverPassword(username string, newPassword string, code string) (Response, error) {
	return put("recoverPassword", fmt.Sprintf(`{"username":"%s", "newPassword":"%s", "code":"%s"}`, username, newPassword, code))
}
func SignUp(username string, password string, code string) (Response, error) {
	return post("signUp", fmt.Sprintf(`{"username":"%s", "password":"%s", "code":"%s"}`, username, password, code))
}

func printResponse(resp Response, err error) {
	if err != nil {
		fmt.Println("ERR: " + err.Error())
		return
	}
	fmt.Println(responseToString(resp))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Telyauth Shell")
	fmt.Println("---------------------")
	printResponse(CheckPassword("misteeka", "qazwsx"))
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		text = strings.ReplaceAll(text, "\r", "")
		parts := strings.Split(text, " ")
		if len(parts) < 1 {
			fmt.Println("Wrong command")
			continue
		}
		cmd := parts[0]
		args := parts[1:]
		if strings.Compare("signIn", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			password := args[1]
			printResponse(SignIn(username, password))
		} else if strings.Compare("checkPassword", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			password := args[1]
			printResponse(CheckPassword(username, password))
		} else if strings.Compare("resetPassword", cmd) == 0 {
			if len(args) < 3 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			oldPassword := args[1]
			newPassword := args[2]
			printResponse(ResetPassword(username, oldPassword, newPassword))
		} else if strings.Compare("requestSignUpCode", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			email := args[1]
			printResponse(RequestSignUpCode(username, email))
		} else if strings.Compare("requestPasswordRecovery", cmd) == 0 {
			if len(args) < 1 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			printResponse(RequestPasswordRecovery(username))
		} else if strings.Compare("recoverPassword", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			newPassword := args[1]
			code := args[2]
			printResponse(RecoverPassword(username, newPassword, code))
		} else if strings.Compare("signUp", cmd) == 0 {
			if len(args) < 3 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			password := args[1]
			code := args[2]
			printResponse(SignUp(username, password, code))
		} else {
			fmt.Println("Unknown command.")
		}
	}

}
