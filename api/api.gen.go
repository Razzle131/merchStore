//go:build go1.22

// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// AuthRequest defines model for AuthRequest.
type AuthRequest struct {
	// Password Пароль для аутентификации.
	Password string `json:"password"`

	// Username Имя пользователя для аутентификации.
	Username string `json:"username"`
}

// AuthResponse defines model for AuthResponse.
type AuthResponse struct {
	// Token JWT-токен для доступа к защищенным ресурсам.
	Token *string `json:"token,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	// Errors Сообщение об ошибке, описывающее проблему.
	Errors *string `json:"errors,omitempty"`
}

// InfoResponse defines model for InfoResponse.
type InfoResponse struct {
	CoinHistory *struct {
		Received *[]struct {
			// Amount Количество полученных монет.
			Amount *int `json:"amount,omitempty"`

			// FromUser Имя пользователя, который отправил монеты.
			FromUser *string `json:"fromUser,omitempty"`
		} `json:"received,omitempty"`
		Sent *[]struct {
			// Amount Количество отправленных монет.
			Amount *int `json:"amount,omitempty"`

			// ToUser Имя пользователя, которому отправлены монеты.
			ToUser *string `json:"toUser,omitempty"`
		} `json:"sent,omitempty"`
	} `json:"coinHistory,omitempty"`

	// Coins Количество доступных монет.
	Coins     *int `json:"coins,omitempty"`
	Inventory *[]struct {
		// Quantity Количество предметов.
		Quantity *int `json:"quantity,omitempty"`

		// Type Тип предмета.
		Type *string `json:"type,omitempty"`
	} `json:"inventory,omitempty"`
}

// SendCoinRequest defines model for SendCoinRequest.
type SendCoinRequest struct {
	// Amount Количество монет, которые необходимо отправить.
	Amount int `json:"amount"`

	// ToUser Имя пользователя, которому нужно отправить монеты.
	ToUser string `json:"toUser"`
}

// PostApiAuthJSONRequestBody defines body for PostApiAuth for application/json ContentType.
type PostApiAuthJSONRequestBody = AuthRequest

// PostApiSendCoinJSONRequestBody defines body for PostApiSendCoin for application/json ContentType.
type PostApiSendCoinJSONRequestBody = SendCoinRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Аутентификация и получение JWT-токена. При первой аутентификации пользователь создается автоматически.
	// (POST /api/auth)
	PostApiAuth(w http.ResponseWriter, r *http.Request)
	// Купить предмет за монеты.
	// (GET /api/buy/{item})
	GetApiBuyItem(w http.ResponseWriter, r *http.Request, item string)
	// Получить информацию о монетах, инвентаре и истории транзакций.
	// (GET /api/info)
	GetApiInfo(w http.ResponseWriter, r *http.Request)
	// Отправить монеты другому пользователю.
	// (POST /api/sendCoin)
	PostApiSendCoin(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// PostApiAuth operation middleware
func (siw *ServerInterfaceWrapper) PostApiAuth(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostApiAuth(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetApiBuyItem operation middleware
func (siw *ServerInterfaceWrapper) GetApiBuyItem(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "item" -------------
	var item string

	err = runtime.BindStyledParameterWithOptions("simple", "item", r.PathValue("item"), &item, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "item", Err: err})
		return
	}

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetApiBuyItem(w, r, item)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetApiInfo operation middleware
func (siw *ServerInterfaceWrapper) GetApiInfo(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetApiInfo(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostApiSendCoin operation middleware
func (siw *ServerInterfaceWrapper) PostApiSendCoin(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostApiSendCoin(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

// ServeMux is an abstraction of http.ServeMux.
type ServeMux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("POST "+options.BaseURL+"/api/auth", wrapper.PostApiAuth)
	m.HandleFunc("GET "+options.BaseURL+"/api/buy/{item}", wrapper.GetApiBuyItem)
	m.HandleFunc("GET "+options.BaseURL+"/api/info", wrapper.GetApiInfo)
	m.HandleFunc("POST "+options.BaseURL+"/api/sendCoin", wrapper.PostApiSendCoin)

	return m
}
