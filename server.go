package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/context"

	"github.com/haruyama/golang-goji-sample/controllers"
	"github.com/haruyama/golang-goji-sample/system"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

func main() {
	filename := flag.String("config", "config.toml", "Path to configuration file")

	flag.Parse()
	defer glog.Flush()

	var application = &system.Application{}

	application.Init(filename)
	application.LoadTemplates()

	// Setup static files
	static := web.New()
	publicPath := application.Config.Get("general.public_path").(string)
	static.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir(publicPath))))

	http.Handle("/assets/", static)

	// Apply middleware
	goji.Use(application.ApplyTemplates)
	goji.Use(application.ApplySessions)
	goji.Use(application.ApplyDbMap)
	goji.Use(application.ApplyAuth)
	goji.Use(application.ApplyIsXhr)
	goji.Use(application.ApplyCsrfProtection)
	goji.Use(context.ClearHandler)

	controller := &controllers.MainController{}

	// Couple of files - in the real world you would use nginx to serve them.
	goji.Get("/robots.txt", http.FileServer(http.Dir(publicPath)))
	goji.Get("/favicon.ico", http.FileServer(http.Dir(publicPath+"/images")))

	// Home page
	goji.Get("/", application.Route(controller, "Index"))

	// Sign In routes
	goji.Get("/signin", application.Route(controller, "SignIn"))
	goji.Post("/signin", application.Route(controller, "SignInPost"))

	// Sign Up routes
	goji.Get("/signup", application.Route(controller, "SignUp"))
	goji.Post("/signup", application.Route(controller, "SignUpPost"))

	// KTHXBYE
	goji.Get("/logout", application.Route(controller, "Logout"))

	graceful.PostHook(func() {
		application.Close()
	})
	goji.Serve()
}
