package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"net/http"
	"gophernews"
)

// cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// login handler

func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	redirectTarget := "/"
	if name == "sai" && pass == "sai" {
		// .. check credentials ..
		setSession(name, response)
		redirectTarget = "/internal"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// index page

const indexPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="name">User name</label>
    <input type="text" id="name" name="name">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

// internal page

const internalPage = `
<h1>%s</h1>
<hr>
<small>User: %s</small>
<div>
	Click <a href="%s"> here </a> to access the link
</div>
<form method="post" action="/internal">
    <button type="submit">Next</button>
</form>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	client := gophernews.NewClient()
  topStories, _ := client.GetTop100();
  // currentStory, _ := client.GetMaxItem();
  // story, _ := client.GetStory(8412605) //=> Returns a Story struct
  // comment, _ := client.GetComment(8412767) //=> Returns a Comment struct
  // poll, _ := client.GetPoll(126809) //=> Returns a Poll struct
  // part, _ := client.GetPart(160705)
  // fmt.Println(currentStory)
	i := 0;
	test, _ := client.GetStory(topStories[i])

  // fmt.Println("len:", topStories[150])
  // for i := 0; i < 5; i++ {
  //   // int test
  //   // test = topStories[i];
	// 	fmt.Println(client.GetStory(topStories[i]))
	// }
	userName := getUserName(request)
	if userName != "" {
		// keys := make([]int, len(test))
		// i := 0
		// for k := range test {
		//     keys[i] = k
		//     i++
		// }
		fmt.Println(test)
		fmt.Println(test.Title)
		fmt.Println(test.By)
		fmt.Println(test.Score)
		fmt.Println(test.Time)
		fmt.Println(test.Type)
		fmt.Println(test.URL)
		fmt.Println(test.ID)
		fmt.Fprintf(response, internalPage, test.Title, test.By, test.URL)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// server main method

var router = mux.NewRouter()

func main() {

	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/internal", internalPageHandler)

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
