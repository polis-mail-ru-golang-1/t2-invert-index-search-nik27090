package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"./handlers"
	"./invertIndex"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

type Config struct {
	ServerAddress string
	Direct        string
	Addr          string
	Username      string
	Pass          string
	DB            string
}

type AccessLogger struct {
	ZapLogger *zap.SugaredLogger
}

func (ac *AccessLogger) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		ac.ZapLogger.Info(
			zap.String("URL", r.URL.Path),
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Duration("work_time", time.Since(start)),
		)
	})
}

func main() {
	//config
	conf, err := os.Open("config.json")
	check(err)
	defer conf.Close()
	decoder := json.NewDecoder(conf)
	config := Config{}
	err = decoder.Decode(&config)
	check(err)

	// zap
	zapLogger, err := zap.NewProduction()
	defer zapLogger.Sync()
	check(err)
	zap.ReplaceGlobals(zapLogger)

	zapLogger.Info("server is started",
		zap.String("address", config.ServerAddress),
	)

	AccessLogOut := new(AccessLogger)

	sugar := zapLogger.Sugar().With()
	AccessLogOut.ZapLogger = sugar

	//DB
	db := pg.Connect(&pg.Options{
		Addr:     config.Addr,
		User:     config.Username,
		Password: config.Pass,
		Database: config.DB,
	})
	defer db.Close()

	// server stuff
	siteMux := http.NewServeMux()
	siteHandler := AccessLogOut.accessLogMiddleware(siteMux)

	invertIndex.InIn, invertIndex.SliceFiles = invertIndex.OpenFiles(config.Direct)
	zap.S().Info("InvertIndex built.")

	siteMux.HandleFunc("/search", handlers.SearchPage)
	siteMux.HandleFunc("/", handlers.MainPage)
	siteMux.HandleFunc("/add", handlers.AddPage)
	http.ListenAndServe(config.ServerAddress, siteHandler)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
