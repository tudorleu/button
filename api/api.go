package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tudorleu/button/models"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var kParamUserId = "userId"
var kParamAmount = "amount"
var kParamEmail = "email"
var kParamFirstName = "firstName"
var kParamLastName = "lastName"

type Context struct {
	Request *http.Request
	Vars    map[string]string
	Params  url.Values
}

type HandlerFuncWithContext func(context Context) (interface{}, error)

// Assembles the params from the URL path and the ones sent in the query string
// or associated form for the HTTP requests in a context to be used by API
// handler functions.
func WithContext(f HandlerFuncWithContext) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		context := Context{
			Request: request,
			Vars:    mux.Vars(request),
			Params:  request.Form,
		}
		apiResponse, err := f(context)
		if err != nil {
			// TODO(tudor): Give out different status codes depending on error.
			response.WriteHeader(500)
			fmt.Fprintf(response, "{\"error\":\"%s\"}", err)
		} else {
			enc := json.NewEncoder(response)
			if err := enc.Encode(apiResponse); err != nil {
				log.Printf("Error while writing response: %s\n", err)
			}
		}
	}
}

func NewUser(context Context) (interface{}, error) {
	email := context.Params.Get(kParamEmail)
	firstName := context.Params.Get(kParamFirstName)
	lastName := context.Params.Get(kParamLastName)

	user, err := models.NewUser(email, firstName, lastName)
	return user, err
}

func GetUser(context Context) (interface{}, error) {
	userIdStr := context.Vars[kParamUserId] // From request path.
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return nil, err
	}
	return models.GetUser(userId)
}

func NewTransfer(context Context) (interface{}, error) {
	userIdStr := context.Vars[kParamUserId] // From request path.
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return nil, err
	}
	amountStr := context.Params.Get(kParamAmount)
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return nil, errors.New("Invalid amount parameter.")
	}

	transfer, err := models.NewTransfer(userId, amount)
	return transfer, err
}

func GetTransfers(context Context) (interface{}, error) {
	userIdStr := context.Vars[kParamUserId] // From request path.
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return nil, err
	}
	return models.GetTransfers(userId)
}
