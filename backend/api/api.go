package api

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/app"
)

type sessionContextKeyType int

var sessionContextKey sessionContextKeyType

func ctxSession(ctx context.Context) *app.Session {
	return ctx.Value(sessionContextKey).(*app.Session)
}

type requestContextKeyType int

var requestContextKey requestContextKeyType

func ctxRequest(ctx context.Context) *http.Request {
	ret, _ := ctx.Value(requestContextKey).(*http.Request)
	return ret
}

type statusCodeRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (r *statusCodeRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (api *API) httpRequestIPAddress(r *http.Request) string {
	addr := r.RemoteAddr
	proxyCount := api.config.ProxyCount
	if api.config.ProxySecret == "" || subtle.ConstantTimeCompare([]byte(api.config.ProxySecret), []byte(r.Header.Get("Proxy-Secret"))) == 1 {
		if h := r.Header.Get("X-Forwarded-For"); h != "" && proxyCount > 0 {
			if clients := strings.Split(h, ","); proxyCount >= len(clients) {
				addr = clients[0]
			} else {
				addr = clients[len(clients)-proxyCount-1]
			}
		}
	}
	addr = strings.TrimSpace(addr)
	if host, _, err := net.SplitHostPort(addr); err == nil && host != "" {
		return host
	}
	return addr
}

type API struct {
	app     *app.App
	handler http.Handler
	config  *Config
}

func WriteError(w http.ResponseWriter, r apispec.ErrorResponse, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(r)
}

func New(a *app.App, config Config) *API {
	ret := &API{
		app:    a,
		config: &config,
	}

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "HEAD", "PATCH", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:         10 * 60,
	})

	ret.handler = cors.Handler(apispec.HandlerWithOptions(apispec.NewStrictHandlerWithOptions(ret, []apispec.StrictMiddlewareFunc{
		// We put error sanitization both before and after auth to make sure errors always have the
		// most detailed log fields.
		ErrorSanitizationMiddleware,
		AuthMiddleware,
		ErrorSanitizationMiddleware,
		LogOperationIdMiddleware,
	}, apispec.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			WriteError(w, apispec.ErrorResponse{
				Message: err.Error(),
			}, http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			switch err := err.(type) {
			case app.NotFoundError:
				WriteError(w, apispec.ErrorResponse{
					Message: err.Error(),
				}, http.StatusNotFound)
			case app.AuthorizationError:
				WriteError(w, apispec.ErrorResponse{
					Message: err.Error(),
				}, http.StatusForbidden)
			case app.AuthenticationError:
				WriteError(w, apispec.ErrorResponse{
					Message: err.Error(),
				}, http.StatusUnauthorized)
			case app.InternalError:
				WriteError(w, apispec.ErrorResponse{
					Message: err.Error(),
				}, http.StatusInternalServerError)
			case app.UserFacingError:
				WriteError(w, apispec.ErrorResponse{
					Message: err.Error(),
				}, http.StatusBadRequest)
			default:
				// This should never happen.
				zap.L().Error("An unknown error has occurred: " + err.Error())
				WriteError(w, apispec.ErrorResponse{
					Message: "An unknown error has occurred.",
				}, http.StatusInternalServerError)
			}
		},
	}), apispec.GorillaServerOptions{
		BaseRouter: mux.NewRouter().UseEncodedPath(),
		Middlewares: []apispec.MiddlewareFunc{
			AddRequestToContextMiddleware,
			NoCachingMiddleware,
			LoggingMiddleware(ret),
		},
	}))

	return ret
}

func (*API) GetHealthCheck(ctx context.Context, request apispec.GetHealthCheckRequestObject) (apispec.GetHealthCheckResponseObject, error) {
	return apispec.GetHealthCheck200TextResponse("ok"), nil
}

func (*API) ContactUs(ctx context.Context, request apispec.ContactUsRequestObject) (apispec.ContactUsResponseObject, error) {
	sess := ctxSession(ctx)
	if err := sess.ContactUs(ctx, app.ContactUsInput{
		Name:         request.Body.Name,
		EmailAddress: request.Body.EmailAddress,
		Subject:      request.Body.Subject,
		Message:      request.Body.Message,
	}); err != nil {
		return nil, err
	}
	return apispec.ContactUs200Response{}, nil
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.handler.ServeHTTP(w, r)
}

func nilIfEmpty[T comparable](v T) *T {
	var empty T
	if v == empty {
		return nil
	}
	return &v
}

func emptyIfNil[T any](v *T) T {
	if v == nil {
		var empty T
		return empty
	}
	return *v
}

func pointer[T any](v T) *T { return &v }

func reshape[T any](m any) (T, error) {
	var ret T
	buf, err := json.Marshal(m)
	if err != nil {
		return ret, err
	}
	if err := json.Unmarshal(buf, &ret); err != nil {
		return ret, err
	}
	return ret, nil
}

func mapSlice[T any, U any](slice []T, f func(T) U) []U {
	ret := make([]U, len(slice))
	for i, v := range slice {
		ret[i] = f(v)
	}
	return ret
}
