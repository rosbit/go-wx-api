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
)

func StartWxApp() {
	wxapi.InitWxAPI(WorkerNum, os.Stdout)

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

