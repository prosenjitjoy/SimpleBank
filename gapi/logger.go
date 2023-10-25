package gapi

import (
	"context"
	"net/http"
	"os"
	"time"

	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	startTime := time.Now()
	resp, err = handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	slogAttrs := []slog.Attr{
		slog.String("protocol", "grpc"),
		slog.String("method", info.FullMethod),
		slog.Duration("duration", duration),
		slog.Int("status_code", int(statusCode)),
		slog.String("status_text", statusCode.String()),
	}

	var logger *slog.Logger

	if os.Getenv("ENVIRONMENT") == "dev" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil).WithAttrs(slogAttrs))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil).WithAttrs(slogAttrs))
	}

	if int(statusCode) != 0 {
		logger.Error("received a gRPC request")
	} else {
		logger.Info("received a gRPC request")
	}

	return
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rr *ResponseRecorder) WriteHeader(statusCode int) {
	rr.StatusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func (rr *ResponseRecorder) Write(body []byte) (int, error) {
	rr.Body = body
	return rr.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rr := &ResponseRecorder{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rr, r)
		duration := time.Since(startTime)

		slogAttrs := []slog.Attr{
			slog.String("protocol", "http"),
			slog.String("method", r.Method),
			slog.String("path", r.RequestURI),
			slog.Duration("duration", duration),
			slog.Int("status_code", rr.StatusCode),
			slog.String("status_text", http.StatusText(rr.StatusCode)),
		}

		var logger *slog.Logger

		if os.Getenv("ENVIRONMENT") == "dev" {
			logger = slog.New(slog.NewTextHandler(os.Stdout, nil).WithAttrs(slogAttrs))
		} else {
			logger = slog.New(slog.NewJSONHandler(os.Stdout, nil).WithAttrs(slogAttrs))
		}

		if rr.StatusCode != http.StatusOK {
			logger.Error("received a HTTP request", slog.String("body", string(rr.Body)))
		} else {
			logger.Info("received a HTTP request")
		}
	})
}
