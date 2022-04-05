package main

import (
	"bytes"
	"fmt"
	"main/cfg"
	"main/database"
	"main/log"
	"main/response"
	"math/rand"
	"strconv"
	"time"
)

func requestEmailCode(username string) response.Response {
	email, found, err := database.SingleUsers.GetString(username, "email")
	if err != nil {
		return response.INTERNAL_SERVER_ERROR
	}
	if !found {
		return response.NOT_FOUND
	}
	fmt.Println(email)
	// sendEmail(email, "", "")
	return response.SUCCESS
}
func isEmailExists(email string) response.Response {
	_, found, err := database.Users.GetString("email", email, "password")
	if err != nil {
		return response.INTERNAL_SERVER_ERROR
	}
	if !found {
		return response.NO
	} else {
		return response.YES
	}
}
func isUsernameExists(username string) response.Response {
	_, found, err := database.SingleUsers.GetString(username, "password")
	if err != nil {
		return response.INTERNAL_SERVER_ERROR
	}
	if !found {
		return response.NO
	} else {
		return response.YES
	}
}

func signIn(username string, password string, ip string) response.Response {
	checkResponse := checkPassword(username, password)
	if bytes.Compare(checkResponse, response.SUCCESS) == 0 {
		err := database.SingleUsers.Set(username, []string{"last_ip", "last_login"}, []interface{}{ip, time.Now().UnixMicro()})
		if err != nil {
			return response.INTERNAL_SERVER_ERROR
		}
		return response.SUCCESS
	} else {
		return checkResponse
	}
}
func checkPassword(username string, password string) response.Response {
	val, found, err := database.SingleUsers.GetString(username, "password")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return response.INTERNAL_SERVER_ERROR
	}
	if !found {
		return response.AUTHORIZATION_FAILED
	}
	if val != password {
		return response.AUTHORIZATION_FAILED
	}
	return response.SUCCESS
}
func resetPassword(username string, oldPassword string, newPassword string) response.Response {
	checkResponse := checkPassword(username, oldPassword)
	if bytes.Compare(checkResponse, response.SUCCESS) == 0 {
		err := database.SingleUsers.SingleSet(username, "password", newPassword)
		if err != nil {
			return response.INTERNAL_SERVER_ERROR
		}
		return response.SUCCESS
	} else {
		return checkResponse
	}
}
func requestSignUpCode(username string, email string, ip string) response.Response {
	isEmailExists := isEmailExists(email)
	if bytes.Equal(isEmailExists, response.YES) {
		return response.ALREADY_EXISTS
	} else if bytes.Equal(isEmailExists, response.INTERNAL_SERVER_ERROR) {
		return response.INTERNAL_SERVER_ERROR
	}
	code := rand.Intn(999999)
	go func() {
		err := sendEmail(email, "Confirm email", fmt.Sprintf(`<h2>Telython registration</h2>
	<div>
		<div>Hello, %s.</div>
		<div>Use code <b>%d</b> to confirm the email for registration.</div>
		<div>Enter this code in the registration form in your app.</div>
		<div>If you did not request a registration, please ignore this message.</div>
	</div>`, username, code))
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	err := database.PendingEmailConfirmations.Put(username, []string{"email", "code", "timestamp"}, []interface{}{email, code, time.Now().UnixMicro()})
	if err != nil {
		return response.INTERNAL_SERVER_ERROR
	}
	return response.SUCCESS
}
func requestPasswordRecovery(username string) response.Response {
	code := strconv.FormatInt(int64(rand.Intn(999999)), 10)
	email, found, err := database.SingleUsers.GetString(username, "email")
	if err != nil {
		return response.INTERNAL_SERVER_ERROR
	}
	if !found {
		return response.NOT_FOUND
	}
	err = database.EmailCodes.Put(username, []string{"code"}, []interface{}{code})
	if err != nil {
		return response.INTERNAL_SERVER_ERROR
	}
	go func() {
		err := sendEmail(email, "Code from Telython", fmt.Sprintf(`
		<h2>Email confirmation</h2>
		<div>Use code <b>%s</b> to confirm the email for password recovery.</div>
		<div>Enter this code in the registration form in your app.	</div>
		<div>If you did not request a password recovery, please ignore this message.</div>
	`, code))
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	return response.SUCCESS
}
func recoverPassword(username string, code string, newPassword string) response.Response {
	savedCode, found, err := database.EmailCodes.GetString(username, "code")
	if err != nil {
		return response.INTERNAL_SERVER_ERROR
	}
	if !found {
		return response.NOT_FOUND
	}
	if savedCode == code {
		err := database.SingleUsers.SingleSet(username, "password", newPassword)
		if err != nil {
			return response.INTERNAL_SERVER_ERROR
		}
		return response.SUCCESS
	} else {
		return response.AUTHORIZATION_FAILED
	}
}
func signUp(username string, password string, code string, ip string) response.Response {
	isUsernameExists := isUsernameExists(username)
	if bytes.Equal(isUsernameExists, response.YES) {
		return response.ALREADY_EXISTS
	} else if bytes.Equal(isUsernameExists, response.INTERNAL_SERVER_ERROR) {
		return response.INTERNAL_SERVER_ERROR
	}
	email, found, err := database.PendingEmailConfirmations.GetString(username, "email")
	if err != nil {
		return response.INTERNAL_SERVER_ERROR
	}
	if !found {
		return response.NOT_FOUND
	}
	isEmailExists := isEmailExists(email)
	if bytes.Equal(isEmailExists, response.YES) {
		return response.ALREADY_EXISTS
	} else if bytes.Equal(isEmailExists, response.INTERNAL_SERVER_ERROR) {
		return response.INTERNAL_SERVER_ERROR
	}
	err = database.SingleUsers.Put(username, []string{"password", "email", "reg_ip", "last_ip", "reg_date", "last_login"}, []interface{}{password, email, ip, ip, time.Now().UnixMicro(), time.Now().UnixMicro()})
	if err != nil {
		return response.INTERNAL_SERVER_ERROR
	}
	go func() {
		database.PendingEmailConfirmations.Remove(username)
		err = sendEmail(email, "You was register on Telython!", fmt.Sprintf(`<h2>Telython registration</h2>
	<div>
		<div>Hello, %s.</div>
		<div>Your was registered on Telython.</div>
		<div>Enjoy messaging!</div>
	</div>`, username))
		if err != nil {
			log.ErrorLogger.Println(err.Error())
		}
	}()
	return response.SUCCESS
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	log.InfoLogger.Println("Starting...")
	log.InfoLogger.Println("Config loading")
	panicIfError(cfg.LoadConfig())
	log.InfoLogger.Println("Database start")
	panicIfError(database.InitDatabase())

	log.InfoLogger.Println("Gomail start")
	initMailClient()
	// TestMail()

	log.InfoLogger.Println("TLS initialization")
	initTLS()

	panicIfError(initFiber())
	ListenAndServe() // Blocking

	log.InfoLogger.Println("Shutdown...")
	log.InfoLogger.Println("Goodbye!")
}
