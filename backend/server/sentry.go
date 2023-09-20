package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) sentry(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		ctx := r.Context()
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}

		options := []sentry.SpanOption{
			sentry.WithOpName("http.server"),
			sentry.ContinueFromRequest(r),
			sentry.WithTransactionSource(sentry.SourceURL),
		}
		// We don't mind getting an existing transaction back so we don't need to
		// check if it is.
		transaction := sentry.StartTransaction(ctx,
			fmt.Sprintf("%s %s", r.Method, r.URL.Path),
			options...,
		)
		defer transaction.Finish()
		r = r.WithContext(transaction.Context())
		hub.Scope().SetRequest(r)
		defer recoverWithSentry(hub, r)
		handler.ServeHTTP(ww, r)

		transaction.Status = httpStatusToSentryStatus(ww.Status())
		rctx := chi.RouteContext(r.Context())
		transaction.Name = rctx.RouteMethod + " " + rctx.RoutePattern()
	})
}

func recoverWithSentry(hub *sentry.Hub, r *http.Request) {
	if err := recover(); err != nil {
		hub.RecoverWithContext(
			context.WithValue(r.Context(), sentry.RequestContextKey, r),
			err,
		)
	}
}

func httpStatusToSentryStatus(status int) sentry.SpanStatus {
	// c.f. https://develop.sentry.dev/sdk/event-payloads/span/

	if status >= 200 && status < 400 {
		return sentry.SpanStatusOK
	}

	switch status {
	case 499:
		return sentry.SpanStatusCanceled
	case 500:
		return sentry.SpanStatusInternalError
	case 400:
		return sentry.SpanStatusInvalidArgument
	case 504:
		return sentry.SpanStatusDeadlineExceeded
	case 404:
		return sentry.SpanStatusNotFound
	case 409:
		return sentry.SpanStatusAlreadyExists
	case 403:
		return sentry.SpanStatusPermissionDenied
	case 429:
		return sentry.SpanStatusResourceExhausted
	case 501:
		return sentry.SpanStatusUnimplemented
	case 503:
		return sentry.SpanStatusUnavailable
	case 401:
		return sentry.SpanStatusUnauthenticated
	default:
		return sentry.SpanStatusUnknown
	}
}
