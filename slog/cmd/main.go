package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

var appEnv = os.Getenv("APP_ENV")

// https://betterstack.com/community/guides/logging/logging-in-go/

func main() {
	// oldNewusage()

	// creatingCustomLogger()

	// customDefaultLogger()

	// addContextualAttrsToLogRecords()

	// groupingContextualAttrs()

	// creatinAndUsageOfChildLoggers()

	// customizingSlogLevels()

	// creatingCustomLogLevels()

	// customizingSlogHandlers()

	// customizingSlogHandlers()

	// contextPackageWithSlog()

	// errorLoggingWithSlog()

	// hidingSensitiveFields()

	otherBackendUsage()
}

func oldNewusage() {
	log.Print("Info message")
	slog.Info("Info message")

	//	2024/05/12 13:13:39 Info message
	//
	// 2024/05/12 13:13:39 INFO Info message
}

func creatingCustomLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")
	// {"time":"2024-05-12T13:17:58.628961+04:00","level":"INFO","msg":"Info message"}
	// {"time":"2024-05-12T13:17:58.628977+04:00","level":"WARN","msg":"Warn message"}
	// {"time":"2024-05-12T13:17:58.628979+04:00","level":"ERROR","msg":"Error message"}

	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")
	// time=2024-05-12T13:19:46.994+04:00 level=INFO msg="Info message"
	// time=2024-05-12T13:19:46.994+04:00 level=WARN msg="Warn message"
	// time=2024-05-12T13:19:46.994+04:00 level=ERROR msg="Error message"

}

func customDefaultLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)
	slog.Info("Info message")
	log.Println("Hello from old logger")

	//	{"time":"2024-05-12T13:25:42.731817+04:00","level":"INFO","msg":"Info message"}
	//
	// {"time":"2024-05-12T13:25:42.73197+04:00","level":"INFO","msg":"Hello from old logger"}

	handler := slog.NewJSONHandler(os.Stdout, nil)
	oldLogger := slog.NewLogLogger(handler, slog.LevelError)

	_ = http.Server{
		ErrorLog: oldLogger,
	}
}

func addContextualAttrsToLogRecords() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info(
		"incoming request",
		"method", "GET",
		"time_taken_ms", 158,
		"path", "/hello/world?q=search",
		"status", 200,
		"user_agent", "Googlebot",
	)

	logger.Info(
		"incoming request",
		slog.String("method", "GET"),
		slog.Int("time_taken_ms", 158),
		slog.String("path", "/hello/world?q=search"),
		slog.Int("status", 200),
		slog.String("user_agent", "Googlebot"),
	)

	logger.LogAttrs(context.Background(),
		slog.LevelInfo,
		"incoming request",
		slog.String("method", "GET"),
		slog.Int("time_taken_ms", 158),
		slog.String("path", "/hello/world?q=search"),
		slog.Int("status", 200),
		slog.String("user_agent", "Googlebot"),
	)

	// 	{"time":"2024-05-12T13:38:58.544104+04:00","level":"INFO","msg":"incoming request","method":"GET","time_taken_ms":158,"path":"/hello/world?q=search","status":200,"user_agent":"Googlebot"}
	// {"time":"2024-05-12T13:38:58.544229+04:00","level":"INFO","msg":"incoming request","method":"GET","time_taken_ms":158,"path":"/hello/world?q=search","status":200,"user_agent":"Googlebot"}
	// {"time":"2024-05-12T13:38:58.544233+04:00","level":"INFO","msg":"incoming request","method":"GET","time_taken_ms":158,"path":"/hello/world?q=search","status":200,"user_agent":"Googlebot"}

}

func groupingContextualAttrs() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"image uploaded",

		slog.Int("id", 333),
		slog.Group(
			"properties",
			slog.Int("width", 4000),
			slog.Int("height", 3000),
			slog.String("format", "jpg"),
		),
	)

	// {
	// "time":"2024-05-12T13:42:32.461697+04:00",
	// "level":"INFO",
	// "msg":"image uploaded",
	// "id":333,
	// "properties":{
	// 	"width":4000,
	// 	"height":3000,
	// 	"format":"jpg"}
	// }

	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"image uploaded",

		slog.Int("id", 333),
		slog.Group(
			"properties",
			slog.Int("width", 4000),
			slog.Int("height", 3000),
			slog.String("format", "jpg"),
		),
	)
	// time=2024-05-12T13:44:13.927+04:00 level=INFO msg="image uploaded" id=333 properties.width=4000 properties.height=3000 properties.	format=jpg

}

func creatinAndUsageOfChildLoggers() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	buildInfo, _ := debug.ReadBuildInfo()

	logger := slog.New(handler)

	child := logger.With(
		slog.Group("program_info",
			slog.Int("pid", os.Getegid()),
			slog.String("go_version", buildInfo.GoVersion),
		),
	)

	child.Info("image upload successful", slog.String("image_id", "39sjq2"))
	child.Warn("storage is 90% full", slog.String("available space", "900.1 mb"))
	// {"time":"2024-05-12T13:54:58.702026+04:00","level":"INFO","msg":"image upload successful","program_info":{"pid":20,"go_version":"go1.22.2"},"image_id":"39sjq2"}
	// {"time":"2024-05-12T13:54:58.702156+04:00","level":"WARN","msg":"storage is 90% full","program_info":{"pid":20,"go_version":"go1.22.2"},"available space":"900.1 mb"}

	fmt.Println()
	handler = slog.NewJSONHandler(os.Stdout, nil)
	buildInfo, _ = debug.ReadBuildInfo()
	logger = slog.New(handler).WithGroup("program_info")

	child = logger.With(
		slog.Int("pid", os.Getpid()),
		slog.String("go_version", buildInfo.GoVersion),
	)

	child.Warn("storage is 90% full",
		slog.String("available space", "900.1 MB"),
	)

	// {"time":"2024-05-12T13:58:49.394236+04:00","level":"WARN","msg":"storage is 90% full","program_info":{"pid":10667,"go_version":"go1.22.2","available space":"900.1 MB"}}
}

func customizingSlogLevels() {
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(os.Stdout, opts)

	logger := slog.New(handler)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warning message")
	logger.Error("Error message")

	logLevel := &slog.LevelVar{}
	opts = &slog.HandlerOptions{Level: logLevel}
	handler = slog.NewJSONHandler(os.Stdout, opts)

	// ...

	logLevel.Set(slog.LevelDebug)
}

func creatingCustomLogLevels() {
	const (
		LevelTrace = slog.Level(-8)
		LevelFatal = slog.Level(12)
	)

	opts := &slog.HandlerOptions{Level: LevelTrace}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	ctx := context.Background()
	logger.Log(ctx, LevelTrace, "Trace message")
	logger.Log(ctx, LevelFatal, "Fatal level")

	var LevelNames = map[slog.Leveler]string{
		LevelTrace: "TRACE",
		LevelFatal: "FATAL",
	}

	opts = &slog.HandlerOptions{
		Level: LevelTrace,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	}

	logger.Log(ctx, LevelTrace, "Trace message")
	logger.Log(ctx, LevelFatal, "Fatal level")
}

func customizingSlogHandlers() {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	logger.LogAttrs(
		context.Background(),
		slog.LevelDebug,
		"message sent",
		slog.Bool("read", false),
		slog.Int("secs_ago", 12412421124),
	)
	// {"time":"2024-05-12T14:56:21.25665+04:00","level":"DEBUG","source":{"function":"main.customizingSlogHandlers","file":"/Users/dmitriymamykin/Desktop/mygo/slog/cmd/main.go","line":275},"msg":"message sent","read":false,"secs_ago":12412421124}

	opts = *&slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var (
		handler slog.Handler = slog.NewJSONHandler(os.Stdout, &opts)
	)
	if appEnv == "development" {
		handler = slog.NewTextHandler(os.Stdout, &opts)
	}

	logger = slog.New(handler)

	logger.Info("Info message")
	// time=2024-05-12T14:56:21.256+04:00 level=INFO msg="Info message"
}

type ctxKey string

const (
	slogFields ctxKey = "slog_fields"
)

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}

	v := []slog.Attr{}
	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)
}

func contextPackageWithSlog() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.WithValue(context.Background(), "request_id", "24178")

	logger.InfoContext(ctx, "image uploaded", slog.String("image_id", "12241"))
	// {"time":"2024-05-12T15:08:19.904377+04:00","level":"INFO","msg":"image uploaded","image_id":"12241"}

	h := &ContextHandler{slog.NewJSONHandler(os.Stdout, nil)}
	logger = slog.New(h)

	ctx = AppendCtx(context.Background(), slog.String("request_id", "1214125"))
	logger.InfoContext(ctx, "image uploaded", slog.String("image_id", "3114"))
}

func errorLoggingWithSlog() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := errors.New("smth happened")

	logger.ErrorContext(context.Background(), "upload failed", slog.Any("error", err))
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", u.ID),
		slog.String("email", u.Email),
	)
}

func hidingSensitiveFields() {

	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)

	u := &User{
		ID:        "user2312",
		FirstName: "Jan",
		LastName:  "Doe",
		Email:     "jan@example.com",
		Password:  "13fqpjwqdp12",
	}

	logger.Info("info", "user", u)
}

func otherBackendUsage() {
	zapL := zap.Must(zap.NewProduction())
	defer zapL.Sync()

	logger := slog.New(zapslog.NewHandler(zapL.Core(), nil))
	logger.Info(
		"Incoming Request",
		slog.String("method", "GET"),
		slog.String("path", "/api/user"),
		slog.Int("status", 200),
	)
}
