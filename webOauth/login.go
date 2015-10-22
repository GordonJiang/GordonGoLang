package main

import(
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"net/http"
	"log"
	"html/template"
)

type authpage struct{
	templ *template.Template
}

func New() authpage{
	auth := authpage{}
	auth.templ =template.Must(template.New("authpage").Parse(templateContent))
	return auth
}

func (a *authpage) ServeHTTP(w http.ResponseWriter, r *http.Request){
	log.Println("In ServerHttp function")
	cookie, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// not authenticated
		log.Println("Didn't find auth cookie. Redirect to google authentication")
		// redirect to google oauth URL
		provider, err := gomniauth.Provider("google")
		if err != nil{
			log.Println("authPage ServerHttp: failed to get provider:", err)
			return
		}
		loginURL, err := provider.GetBeginAuthURL(nil, nil)
		if err !=nil{
			log.Println("authPage ServeHttp: failed to get auth URL:", err)
			return
		}
		log.Println("authPage, redirect to the loginUrl : ", loginURL)
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil{
		log.Println("failed to get cookie: ", err)
	} else{
		// parse the user name from cookie and show in page
		username := objx.MustFromBase64(cookie.Value)["name"]
		log.Println("show the page with name:",username)
		a.templ.Execute(w,username)
	}
}

const templateContent = "<html><head><title>Main page</title></head><body> you are authenticated. Name is {{.}}</body></html>"
