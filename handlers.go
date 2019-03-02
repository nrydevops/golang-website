package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func handleIndex(res http.ResponseWriter, req *http.Request) {

	if  GetUserFromRequest(req) != nil {

		http.Redirect(res, req, "/chat/", 302)
		return;
	}

	if req.Method == http.MethodGet {

		tpl.ExecuteTemplate(res, "index.html", nil)
	}
}

func handleChat(res http.ResponseWriter, req *http.Request) {

	if err := RefreshSession(&res, req); err != nil {

		tpl.ExecuteTemplate(res, "index.html", nil)
		return;
	}

	if req.Method == http.MethodGet {

		res.Header().Add("cache-control", "no-cache, private")
		tpl.ExecuteTemplate(res, "chat.html", nil)

	}
}

func handleNews(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {

		user := GetUserFromRequest(req)

		if err := RefreshSession(&res, req); err != nil {

			fmt.Println(err)
			tpl.ExecuteTemplate(res, "index.html", nil)
			return;
		}

		if user != nil {

			tpl.ExecuteTemplate(
				res,
				"news.html",
				NewsData{GetArticles(), user.Admin})

		} else {

			http.Error(res, "Unknown user", http.StatusUnauthorized)
			return
		}

	} else if req.Method == http.MethodPost {

		user := GetUserFromRequest(req)

		if user == nil || !user.Admin {

			http.Error(res, "Unknown user", http.StatusUnauthorized)
			return
		}

		req.ParseForm()

		importance, _ := strconv.Atoi(req.FormValue("importance"))

		article := Article{
			user.Username,
			req.FormValue("header"),
			req.FormValue("text"),
			importance,
			time.Now().Format("2006-01-02 15:04:05"),
		}

		_ = WriteArticleToDB(article)

		tpl.ExecuteTemplate(
			res,
			"news.html",
			NewsData{GetArticles(), user.Admin})
	}
}

func handleFiles(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {

		user := GetUserFromRequest(req)

		if err := RefreshSession(&res, req); err != nil {

			fmt.Println(err)
			tpl.ExecuteTemplate(res, "index.html", nil)
			return;
		}

		if user != nil {

			tpl.ExecuteTemplate(
				res,
				"files.html",
				FilesData{GetFiles(), user.Admin})

		} else {

			http.Error(res, "Unknown user", http.StatusUnauthorized)
			return
		}

	} else if req.Method == http.MethodPost {

		user := GetUserFromRequest(req)

		if user == nil || !user.Admin {

			http.Error(res, "Unknown user", http.StatusUnauthorized)
			return
		}

		//req.ParseForm()
		req.ParseMultipartForm(32 << 20)

		multipartFile, handler, err := req.FormFile("uploadfile")
		if err != nil {

			return
		}
		defer multipartFile.Close()

		file, err := os.Create("./website/storage/" + handler.Filename)
		if err != nil {

			return
		}
		defer file.Close()

		w := bufio.NewWriter(file)
		w.ReadFrom(multipartFile)

		fileInfo := File{
			user.Username,
			req.FormValue("description"),
			"/file/" + handler.Filename,
			time.Now().Format("2006-01-02 15:04:05"),
		}

		_ = WriteFileToDB(fileInfo)

		tpl.ExecuteTemplate(
			res,
			"files.html",
			FilesData{GetFiles(), user.Admin})
	}

}

func handleLogin(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		req.ParseForm()

		session_token, err := SignIn(req.FormValue("username"), req.FormValue("password"))

		if err != nil {

			res.WriteHeader(http.StatusUnauthorized)
			tpl.ExecuteTemplate(res, "index.html", true)

		} else {

			http.SetCookie(res, &http.Cookie{
				Name: "Token",
				Value: session_token,
				Path: "/",
				Expires: time.Now().Add(300 * time.Second),
			})
			http.Redirect(res, req, "/", 302)
		}
	}
}

func handleLogout(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		req.ParseForm()

		cookie, err := req.Cookie("Token")
		if err != nil {

			return;
		}

		LogOut(cookie.Value)

		http.SetCookie(res, &http.Cookie{
			Name: "Token",
			Value: "",
			Path: "/",
			Expires: time.Now().Add(300 * time.Second),
		})
		http.Redirect(res, req, "/", 302)
	}
}

func handleOptions(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {

		if err := RefreshSession(&res, req); err != nil {

			tpl.ExecuteTemplate(res, "index.html", nil)
			return;
		}

		tpl.ExecuteTemplate(res, "options.html", nil)
	}
}

func handleChangePassword(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		req.ParseForm()

		token, err := req.Cookie("Token")
		if err != nil {

			http.Error(res, "Unknown user", http.StatusUnauthorized)
			return
		}

		user := GetUserByToken(token.Value)
		if user == nil {

			http.Error(res, "Unknown user", http.StatusUnauthorized)
			return
		}

		if user.Password !=  req.FormValue("old_password") {

			tpl.ExecuteTemplate(
				res,
				"options.html",
				1)
			return;
		}

		newPassword := req.FormValue("new_password")

		err = ChangeUserInDB(user.Username, newPassword, user.Admin)

		if err != nil {

			tpl.ExecuteTemplate(
				res,
				"options.html",
				1)
			return;
		}

		user.Password = newPassword

		tpl.ExecuteTemplate(res, "options.html", 2)
	}
}