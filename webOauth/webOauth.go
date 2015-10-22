package main

import(
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/signature"
	"github.com/stretchr/objx"
	"net/http"
	"log"
)

func callbackHandler(w http.ResponseWriter, r *http.Request){
	provider, err:= gomniauth.Provider("google")
	if err != nil{
		log.Println("callback handler: failed to get provider:", err)
		return
	}
	
	// see doc https://godoc.org/github.com/stretchr/gomniauth/common#Provider
	// and https://godoc.org/github.com/stretchr/objx#Map
	creds, err:= provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
	if err != nil{
		log.Println("callback hanlder: failed to complete the auth:", err)
		return
	}
	log.Println("callback handler. successfully get credential:", creds)

	user, err := provider.GetUser(creds)
	if err !=nil{
		log.Println("callback handler. Failed to get user name")
		return
	}

	authCookieValue := objx.New(map[string] interface{}{
		"name": user.Name(),
	}).MustBase64()

	http.SetCookie(w, &http.Cookie{
		Name: "auth",
		Value: authCookieValue,
		Path: "/",
	})

	// redirect to normal page
	w.Header()["Location"] = []string{"/"}
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main(){
	// init gomniauth
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		google.New("263631501515-4l08gqagj93r4iqlk1ciuc1hrpbt69i5.apps.googleusercontent.com","oJAYjKrwZlTvLLVSQlWSTKPL","http://localhost:8080/callback"),
	)

	// binding authpage and callback handler, start web server
	http.HandleFunc("/callback", callbackHandler)
	/*
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Hello auth")
	})*/
	auth := New()
	http.Handle("/", &auth)

	log.Println("Starting web server...")
	log.Fatal(http.ListenAndServe(":8080",nil))
}
