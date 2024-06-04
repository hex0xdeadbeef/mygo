package main

// https://betterstack.com/community/guides/logging/go/zap/

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// create()

	// createDevelopmentLogger()

	// devToProdSwitching()

	// settingGlobalLogger()

	// lowLevelLogger()

	// sugaredLogger()

	// conversionsBetweenSugarAndDefaultLogger()

	// sugaredLoggerMethodsNotes()

	// customLoggerUsage()

	// addingContextToLogs()

	// parentAndChildLoggers()

	// errorsLogging()

	// configuringDPANICAndPANICToErrorLevel()

	// logSampling()

	// hidingSensitiveDetails()
}

func create() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	logger.Info("Hello from Zap logger!")
}

func createDevelopmentLogger() {
	logger := zap.Must(zap.NewDevelopment())
	defer logger.Sync()

	logger.Info("Hello from Zap Development logger!")
}

func devToProdSwitching() {
	var (
		logger *zap.Logger
	)

	if os.Getenv("APP_ENV") == "development" {
		logger = zap.Must(zap.NewDevelopment())
	} else {
		logger = zap.Must(zap.NewProduction())
	}

	logger.Info(fmt.Sprintf("Hello from Zap %q logger!", os.Getenv("APP_ENV")))
}

func settingGlobalLogger() {
	var (
		initGlobalLogger = func() func() {
			return zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
		}
	)
	defer initGlobalLogger()

	zap.L().Info("Hello from Global Zap Logger")
}

func lowLevelLogger() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	logger.Info(
		"user logged in",
		zap.String("username", "johndoe"),
		zap.Int("userid", 12345),
		zap.String("provider", "google"),
	)
}

func sugaredLogger() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Info("Hello from Zap Sugar Logger!")
	// {"level":"info","ts":1715518934.1648672,"caller":"cmd/main.go:84","msg":"Hello from Zap Sugar Logger!"}

	sugar.Infoln(
		"Hello from Zap Sugar Logger!",
	)
	// {"level":"info","ts":1715518977.9734578,"caller":"cmd/main.go:87","msg":"Hello from Zap Sugar Logger!"}

	sugar.Infof("Hello from Zap logger! The time is %s", time.Now().Format("03:04 AM"))
	// {"level":"info","ts":1715519061.745663,"caller":"cmd/main.go:93","msg":"Hello from Zap logger! The time is 05:04 AM"}

	sugar.Infow(
		"user logged in",
		"username", "johndoe",
		"userid", 12345,
		zap.String("provider", "google"),
	)
	// {"level":"info","ts":1715519153.376725,"caller":"cmd/main.go:96","msg":"user logged in","username":"johndoe","userid":12345,"provider":"google"}

}

func conversionsBetweenSugarAndDefaultLogger() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	sugaredLogger := logger.Sugar()
	defaultDesugaredLogger := sugaredLogger.Desugar()

	defaultDesugaredLogger.Info("Hello from desugared logger!")
}

func sugaredLoggerMethodsNotes() {
	defer func() {
		if v := recover(); v != nil {
			fmt.Printf("%#[1]v\n\n%[1]T\n", v)
		}
	}()
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	sugaredLogger.Infow("user logged in", 1234, "userID")
}

func createCustomLoggerWithCfg() *zap.Logger {
	encodedCfg := zap.NewProductionEncoderConfig()
	encodedCfg.TimeKey = "timestamp"
	encodedCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     true, // changed
		DisableStacktrace: true, // changed
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encodedCfg,
		OutputPaths:       []string{"stdout"}, // changed
		InitialFields: map[string]interface{}{
			"pid": os.Getgid(),
			"work_dir": func() string {
				if wd, err := os.Getwd(); err != nil {
					return wd
				}
				return "not found"
			}(),
		},
	}

	return zap.Must(config.Build())
}

func createCustomLoggerWithNew() *zap.Logger {

	/*
		Creates write-syncers
	*/
	stdout := zapcore.AddSync(os.Stdout)
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename: "logs/app.log",
		// Data for log-rotating
		MaxSize:    16, // MB
		MaxBackups: 3,
		MaxAge:     7, // Days
		Compress:   false,
	})

	/*
		Create Atomic Level
	*/
	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	/*
		Creates and defines the production/development cfgs
	*/
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	/*
		Defines and creates decoders
	*/
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	/*
		Defines and creates cores. Logger will write out into the stdout/file
	*/
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	return zap.New(core)
}

func customLoggerUsage() {
	customLoggerWithCfg := createCustomLoggerWithCfg()
	defer customLoggerWithCfg.Sync()
	customLoggerWithCfg.Info("Hello from Zap custom logger!")

	customLoggerWithNew := createCustomLoggerWithNew()
	defer customLoggerWithNew.Sync()
	customLoggerWithNew.Info("Hello from Zap custom logger!")
}

func addingContextToLogs() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	logger.Warn(
		"User account is nearing the storage limit",
		zap.String("username", "john.doe"),
		zap.Float64("storageUsed", 4.5),
		zap.Float64("storageLimit", 5),
	)
}

func parentAndChildLoggers() {
	parentLogger := zap.Must(zap.NewProduction())
	defer parentLogger.Sync()

	buildinfo, _ := debug.ReadBuildInfo()

	childLogger := parentLogger.With(
		zap.String("go_version", buildinfo.GoVersion),
		zap.Int("pid", os.Getpid()),

		zap.String("service", "userService"),
		zap.String("requestID", "abc123"),
	)
	defer childLogger.Sync()

	childLogger.Info(
		"User registration successful",
		zap.String("username", "john.doe"),
		zap.String("email", "john@example.com"),
	)

	childLogger.Info("redirectering user to admin dashboard")
}

func errorsLogging() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	logger.Error(
		"Failed to perform an op",
		zap.String("op", "someOp"),
		zap.Error(errors.New("smth happened")),
		zap.Int("retryAttempts", 3),
		zap.String("user", "john.doe"),
	)

	// Calls the os.Exit(1)
	logger.Fatal(
		"Smth went terribly wrong",
		zap.String("context", "main"),
		zap.Int("code", 500),
		zap.Error(errors.New("An error occured")),
	)

	// Calls panic(...)
	logger.Panic(
		"program is destroyed, it will panic",
		zap.String("context", "main"),
		zap.Int("code", 500),
		zap.Error(errors.New("An error occured")),
	)

}

func configuringDPANICAndPANICToErrorLevel() {
	var (
		lowerCaseLevelEncoder = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			if level == zap.PanicLevel || level == zap.DPanicLevel {
				enc.AppendString("error")
				return
			}

			zapcore.LowercaseColorLevelEncoder(level, enc)
		}
	)

	stdout := zapcore.AddSync(os.Stdout)
	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	productionCfg.EncodeLevel = lowerCaseLevelEncoder

	jsonEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewCore(jsonEncoder, stdout, level)

	logger := zap.New(core)
	defer logger.Sync()

	// Panics
	logger.DPanic(
		"this was never supposed to happen",
	)
}

func logSampling() {
	var (
		lowerCaseLevelEncoder = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			if level == zap.PanicLevel || level == zap.DPanicLevel {
				enc.AppendString("error")
				return
			}

			zapcore.LowercaseColorLevelEncoder(level, enc)
		}
	)

	stdout := zapcore.AddSync(os.Stdout)

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey, productionCfg.EncodeTime = "timestamp", zapcore.ISO8601TimeEncoder
	productionCfg.EncodeLevel = lowerCaseLevelEncoder
	productionCfg.StacktraceKey = "stack"

	jsonEncoder := zapcore.NewJSONEncoder(productionCfg)

	jsonOutCore := zapcore.NewCore(jsonEncoder, stdout, level)

	samplingOutCore := zapcore.NewSamplerWithOptions(
		jsonOutCore,
		time.Millisecond, // interval
		3,                // log first 3 entries
		0,                // thereafter log zero entries with the interval
	)

	logger := zap.New(samplingOutCore)
	defer logger.Sync()

	for i := 0; i < 10; i++ {
		logger.Info("info message")
		logger.Warn("warn message")
	}
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *User) String() string {
	return u.ID
}

type SensitiveFieldEncoder struct {
	zapcore.Encoder
	cfg zapcore.EncoderConfig
}

func hidingSensitiveDetails() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	user := &User{
		ID:    "1245",
		Name:  "John Doe",
		Email: "john@example.com",
	}

	logger.Info("user login", zap.Any("user", user))
}
