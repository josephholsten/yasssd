package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/josephholsten/yasssd/store"
)

func init() {
	http.HandleFunc("/", yasssdHandler)
}

// TODO: params: bind port
func main() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func writeErrorResp(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errorBody := map[string]string{"error": msg}
	bs, _ := json.MarshalIndent(errorBody, "", "  ")
	fmt.Fprint(w, string(bs))
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func yasssdHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: why didn't I write this with contexts from the begining?
	// TODO: validate request content-types
	if r.URL.Path == "/register" {
		if r.Body == nil {
			writeErrorResp(w, http.StatusBadRequest, "request body was empty, expected json with username and password fields")
			return
		}

		a, err := store.AccountUnmarshal(r.Body)
		if err != nil {
			writeErrorResp(w, http.StatusBadRequest, err.Error())
			return
		}

		if ok, msg := a.IsValid(); !ok {
			writeErrorResp(w, http.StatusBadRequest, msg)
			return
		}

		// FIXME: ensure account does not exist
		if err := a.Create(); err != nil {
			log.Println(err)
			writeErrorResp(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.URL.Path == "/login" {
		if r.Method != "POST" {
			writeErrorResp(w, http.StatusBadRequest, fmt.Sprintf("unexpected method:%v", r.Method))
			return
		}
		if r.Body == nil {
			writeErrorResp(w, http.StatusBadRequest, "request body was empty, expected json with username and password fields")
			return
		}

		a, err := store.AccountUnmarshal(r.Body)
		if err != nil {
			writeErrorResp(w, http.StatusBadRequest, err.Error())
			return
		}

		if !a.IsAuthenticated() {
			writeErrorResp(w, http.StatusForbidden, "Invalid username and password")
			return
		}

		t, err := a.CreateToken()
		if err != nil {
			log.Println(err)
			writeErrorResp(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		bs, _ := json.Marshal(t)
		w.Write(bs)
		return
	}
	if r.URL.Path == "/files" {
		t := r.Header["X-Token"]
		if t == nil || len(t) != 1 || t[0] == "" {
			writeErrorResp(w, http.StatusBadRequest, "request header X-Token was empty, expected a valid session token")
			return
		}
		a, err := store.FindAccountByToken(t[0]) // Request.Header concatenates header values, only use first
		if err != nil {
			writeErrorResp(w, http.StatusForbidden, err.Error())
			return
		}

		fs, err := a.AllFiles()
		if err != nil {
			log.Println(err)
			writeErrorResp(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		bs, _ := json.Marshal(fs)
		w.Write(bs)
		return
	}

	filesPattern := regexp.MustCompile("/files/(.*)")
	matches := filesPattern.FindStringSubmatch(r.URL.Path)
	if len(matches) != 0 {
		fileID := matches[1]
		t := r.Header["X-Token"]
		if t == nil || len(t) != 1 || t[0] == "" {
			writeErrorResp(w, http.StatusBadRequest, "request header X-Token was empty, expected a valid session token")
			return
		}
		a, err := store.FindAccountByToken(t[0]) // Request.Header concatenates header values, only use first
		if err != nil {
			writeErrorResp(w, http.StatusForbidden, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v/%v", a.Username, fileID)
		return
	}
	http.NotFound(w, r)
}
