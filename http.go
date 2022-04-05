package main

import (
	"crypto/tls"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fastjson"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"main/cfg"
	"main/log"
	"main/response"
	"net"
)

var (
	tlsConfig   *tls.Config
	app         *fiber.App
	tlsListener net.Listener
)

func initFiber() error {
	log.InfoLogger.Println("Fiber initialization")
	app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	requestHandler()
	tlsPortListener, err := net.Listen(cfg.GetString("networkName"), "0.0.0.0:443")
	if err != nil {
		return err
	}

	tlsListener = tls.NewListener(tlsPortListener, tlsConfig)
	return nil
}

func initTLS() {
	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("oppositemc.com"),
		Cache:      autocert.DirCache("./certs"),
	}
	tlsConfig = &tls.Config{
		GetCertificate: certManager.GetCertificate,
		NextProtos: []string{
			"http/1.1", acme.ALPNProto,
		},
		InsecureSkipVerify: false,
	}
}

func ListenAndServe() {
	// log("Start HTTPS redirect")
	// go runHTTPSRedirect()
	log.InfoLogger.Println("Fiber server start")

	err := app.Listen(":8001")
	// err := app.Listener(tlsListener)
	panicIfError(err)
}

func runHTTPSRedirect() {

}

func Respond(ctx *fiber.Ctx, resp response.Response) error {
	_, err := ctx.Write(resp)
	return err
}

func Deserialize(jsonBytes []byte) (data *fastjson.Value, err error) {
	var p fastjson.Parser
	data, err = p.ParseBytes(jsonBytes)
	return
}

func requestHandler() {
	app.Put("auth/signIn", func(ctx *fiber.Ctx) error {
		data, err := Deserialize(ctx.Body())
		if err != nil {
			return Respond(ctx, response.INVALID_REQUEST)
		}
		username := string(data.GetStringBytes("username"))
		password := string(data.GetStringBytes("password"))
		return Respond(ctx, signIn(username, password, ctx.IP()))
	})
	app.Get("auth/checkPassword", func(ctx *fiber.Ctx) error {
		username := ctx.FormValue("u")
		password := ctx.FormValue("p")
		return Respond(ctx, checkPassword(username, password))
	})
	app.Put("auth/resetPassword", func(ctx *fiber.Ctx) error {
		data, err := Deserialize(ctx.Body())
		if err != nil {
			return Respond(ctx, response.INVALID_REQUEST)
		}
		username := string(data.GetStringBytes("username"))
		oldPassword := string(data.GetStringBytes("oldPassword"))
		newPassword := string(data.GetStringBytes("newPassword"))
		return Respond(ctx, resetPassword(username, oldPassword, newPassword))
	})
	app.Post("auth/requestSignUpCode", func(ctx *fiber.Ctx) error {
		data, err := Deserialize(ctx.Body())
		if err != nil {
			return Respond(ctx, response.INVALID_REQUEST)
		}
		username := string(data.GetStringBytes("username"))
		email := string(data.GetStringBytes("email"))
		return Respond(ctx, requestSignUpCode(username, email, ctx.IP()))
	})
	app.Put("auth/requestPasswordRecovery", func(ctx *fiber.Ctx) error {
		data, err := Deserialize(ctx.Body())
		if err != nil {
			return Respond(ctx, response.INVALID_REQUEST)
		}
		username := string(data.GetStringBytes("username"))
		return Respond(ctx, requestPasswordRecovery(username))
	})
	app.Put("auth/recoverPassword", func(ctx *fiber.Ctx) error {
		data, err := Deserialize(ctx.Body())
		if err != nil {
			return Respond(ctx, response.INVALID_REQUEST)
		}
		username := string(data.GetStringBytes("username"))
		code := string(data.GetStringBytes("code"))
		newPassword := string(data.GetStringBytes("newPassword"))
		return Respond(ctx, recoverPassword(username, code, newPassword))
	})
	app.Post("auth/signUp", func(ctx *fiber.Ctx) error {
		data, err := Deserialize(ctx.Body())
		if err != nil {
			return Respond(ctx, response.INVALID_REQUEST)
		}
		username := string(data.GetStringBytes("username"))
		password := string(data.GetStringBytes("password"))
		code := string(data.GetStringBytes("code"))
		return Respond(ctx, signUp(username, password, code, ctx.IP()))
	})
}
