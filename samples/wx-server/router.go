/**
 * REST API router
 * Rosbit Xu
 */
package main

import (
	"github.com/urfave/negroni"
	"github.com/gernest/alien"
	"net/http"
	"fmt"
	"os"
	"github.com/rosbit/go-wx-api"
	"github.com/rosbit/go-wx-api/conf"
	"github.com/rosbit/go-wx-api/samples/wx-server/utils"
)

func _registerMessageHandlers() {
	wxapi.RegisterSubscribeEventHandler(subcribeUser)
	wxapi.RegisterRedictHandler(handleRedirect)
	wxapi.RegisterTextMsgHandler(textMsgReceived)
}

func _inits() {
	utils.StartFilesThreads([]string{WelcomeFile, MenuJsonConf})
	wxapi.InitWxAPI(WorkerNum, os.Stdout)
	createMsgHelperEndpoints()
	_registerMessageHandlers()
}

func StartWxApp() {
	_inits()

	api := negroni.New()
	signatureChecker := wxapi.NewWxSignatureChecker(wxconf.WxParams.Token, WxTimeout, []string{WxService})
	api.Use(negroni.HandlerFunc(signatureChecker))
	api.Use(negroni.NewRecovery())
	api.Use(negroni.NewLogger())

	router := alien.New()
	router.Get(WxService,  wxapi.Echo)
	router.Post(WxService, wxapi.Request)
	router.Get(WxRedirect, wxapi.Redirect)
	api.UseHandler(router)

	listenParam := fmt.Sprintf("%s:%d", ListenHost, ListenPort)
	fmt.Printf("%v\n", http.ListenAndServe(listenParam, api))
}

