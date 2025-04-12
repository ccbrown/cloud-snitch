package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func AddRequestToContextMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), requestContextKey, r))
		h.ServeHTTP(w, r)
	})
}

// This middleware sets the Cache-Control header to no-cache.
func NoCachingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		h.ServeHTTP(w, r)
	})
}

type endOfRequestLogFieldsContextKeyType int

var endOfRequestLogFieldsContextKey endOfRequestLogFieldsContextKeyType

// This middleware creates the unauthenticated session and adds logging before and after each
// request.
func LoggingMiddleware(api *API) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			beginTime := time.Now()

			sess := api.app.NewAnonymousSession()

			w = &statusCodeRecorder{
				ResponseWriter: w,
			}

			requestId := model.NewId("req").String()
			sess = sess.WithLogFields(zap.String("request_id", requestId))
			if remote := api.httpRequestIPAddress(r); remote != "" {
				sess = sess.WithLogFields(zap.String("remote", remote))
			}

			var endOfRequestLogFields []zap.Field

			defer func() {
				statusCode := w.(*statusCodeRecorder).StatusCode
				if statusCode == 0 {
					statusCode = 200
				}

				duration := time.Since(beginTime)
				logger := sess.Logger().With(
					zap.Int64("duration_ms", int64(duration/time.Millisecond)),
					zap.Int("status_code", statusCode),
				).With(endOfRequestLogFields...)
				logger.Info(r.Method + " " + r.URL.RequestURI())
			}()

			r = r.WithContext(context.WithValue(context.WithValue(r.Context(), sessionContextKey, sess), endOfRequestLogFieldsContextKey, endOfRequestLogFields))
			r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
			h.ServeHTTP(w, r)
		})
	}
}

// This middleware recovers from panics and makes sure all exposed errors have been sanitized.
func ErrorSanitizationMiddleware(f apispec.StrictHandlerFunc, operationID string) apispec.StrictHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, args any) (ret any, retErr error) {
		sess := ctxSession(ctx)

		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("%v: %s", r, debug.Stack())
				retErr = sess.SanitizedError(err)
			}
		}()

		ret, err := f(ctx, w, r, args)
		return ret, sess.SanitizedError(err)
	}
}

// This middleware adds the operation id to the session and end-of-request logs.
func LogOperationIdMiddleware(f apispec.StrictHandlerFunc, operationID string) apispec.StrictHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, args any) (ret any, retErr error) {
		ctx = context.WithValue(ctx, endOfRequestLogFieldsContextKey, append([]zap.Field{
			zap.String("operation_id", operationID),
		}, ctx.Value(endOfRequestLogFieldsContextKey).([]zap.Field)...))
		sess := ctxSession(ctx).WithLogFields(
			zap.String("operation_id", operationID),
		)
		ctx = context.WithValue(ctx, sessionContextKey, sess)
		return f(ctx, w, r, args)
	}
}

// This middleware authenticates the user.
func AuthMiddleware(f apispec.StrictHandlerFunc, operationID string) apispec.StrictHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, args interface{}) (ret interface{}, retErr error) {
		if auth := r.Header.Get("Authorization"); auth != "" {
			sess := ctxSession(ctx)

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || parts[0] != "token" {
				return nil, app.AuthenticationError{}
			} else if token, err := base64.RawURLEncoding.DecodeString(parts[1]); err != nil {
				return nil, app.AuthenticationError{}
			} else {
				newSess, err := sess.WithUserAccessToken(ctx, token)
				if err != nil {
					return nil, err
				} else if newSess == nil {
					return nil, app.AuthenticationError{}
				}
				if newSess.User() != nil {
					// Make sure we log the user id in the request log too.
					ctx = context.WithValue(ctx, endOfRequestLogFieldsContextKey, append([]zap.Field{
						zap.String("user_id", newSess.User().Id.String()),
					}, ctx.Value(endOfRequestLogFieldsContextKey).([]zap.Field)...))
				}
				ctx = context.WithValue(ctx, sessionContextKey, newSess)
			}
		}

		return f(ctx, w, r, args)
	}
}
