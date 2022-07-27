package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
	"upspin.io/errors"

	"github.com/Owicca/tr/internal/models/util"
	"github.com/owicca/tr/internal/models/logs"
	msessions "github.com/owicca/tr/internal/models/sessions"

	"go.uber.org/zap"
)

var S *Server

// To be ran on server closing
var (
	Undo              func()
	LoggerSync        func() error
	_, filename, _, _ = runtime.Caller(0)
)

func init() {
}

type Server struct {
	mux.Router
	Config   Config
	Conn     *gorm.DB
	Template *Template
	Data     map[string]any
	Errors   *util.Errors
}

func NewServer(
	cfg Config,
	store sessions.Store,
	conn *gorm.DB,
	tmpl *Template,
) *Server {
	if S == nil {
		S = &Server{
			Config:       cfg,
			SessionStore: store.(*sessions.CookieStore),
			Router:       *mux.NewRouter(),
			Conn:         conn,
			Template:     tmpl,
			Data:         map[string]any{},
			Errors:       util.NewErrors(),
		}
	}

	return S
}

// Get hostname and port.
func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%s", s.Config.HttpHost, s.Config.HttpPort)
}

// Server JSON response.
func (s *Server) JSON(w http.ResponseWriter, r *http.Request, status int, data map[string]any) error {
	const op errors.Op = "server.JSON"

	if data != nil {
		data = MergeMaps(s.Data, data)
	}
	w.Header().Set("Content-Type", "application/json")
	if data != nil {
		json.NewEncoder(w).Encode(data)

		if S.Session != nil {
			session := S.Session
			//log.Printf("end => %+v\n", session)
			user_id_str, ok := session.Values["user_id"]

			if err := session.Save(r, w); err != nil {
				logs.LogErr(op, errors.Errorf("Could not save post session on JSON (%s)!", err))
			} else if ok {
				msessions.Update(S.Conn, user_id_str.(int), S.Data)
			}
		}

		return nil
	}
	return fmt.Errorf("No data to return")
}

// Serve a media file.
func (s *Server) MEDIA(w http.ResponseWriter, r *http.Request, status int, media []byte, mediaType string) {
	w.Header().Set("Content-Type", mediaType)
	w.Header().Set("Cache-Control", "max-age=31536000")
	w.WriteHeader(status)
	w.Write(media)
}

// Server a HTML response.
func (s *Server) HTML(w http.ResponseWriter, r *http.Request, status int, htmlView string, data map[string]any) error {
	const op errors.Op = "server.HTML"
	if data != nil {
		data = MergeMaps(data, s.Data)
	}

	if S.Session != nil {
		session := S.Session
		//log.Printf("end => %+v\n", session)
		user_id_str, ok := session.Values["user_id"]

		if err := session.Save(r, w); err != nil {
			logs.LogErr(op, errors.Errorf("Could not save post session on HTML (%s)!", err))
		} else if ok {
			msessions.Update(S.Conn, user_id_str.(int), S.Data)
		}
	}

	return s.Template.Render(w, r, status, htmlView, data)
}

func (s *Server) Redirect(w http.ResponseWriter, r *http.Request, dst string) {
	const op errors.Op = "server.Redirect"

	if S.Session != nil {
		session := S.Session
		//log.Printf("end => %+v\n", session)
		user_id_str, ok := session.Values["user_id"]

		if err := session.Save(r, w); err != nil {
			logs.LogErr(op, errors.Errorf("Could not save post session on redirect (%s)!", err))
		} else if ok {
			msessions.Update(S.Conn, user_id_str.(int), S.Data)
		}
	}

	http.Redirect(w, r, dst, http.StatusSeeOther)
}

func (s *Server) GenerateUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s/", s.Addr(), strings.Trim(endpoint, "/"))
}

// Get config, db, logger
// set up settings and create Server
func (s *Server) Run() {
	addr := s.Addr()
	msg := fmt.Sprintf("Running at %s", addr)
	zap.L().Info(msg, zap.Int64("timestamp", time.Now().Unix()))

	httpServer := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				s.ServeHTTP(w, r)
			}),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	s.ShutdownOnInterrupt(httpServer)
}

func (s *Server) ShutdownOnInterrupt(srv *http.Server) {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			msg := fmt.Sprintf("Shutting down error (%s)", err)
			zap.L().Info(msg, zap.Int64("timestamp", time.Now().Unix()))
		}
		zap.L().Info("Close everything!", zap.Int64("timestamp", time.Now().Unix()))
		defer LoggerSync()
		defer LUndo()
		// s.Conn.Close()
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		zap.L().Info("Could not listen and serve!", zap.Int64("timestamp", time.Now().Unix()))
	}

	<-idleConnsClosed
}

func Setup(configPath string) (Config, sessions.Store, *gorm.DB, *zap.Logger) {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("{\"level\":\"fatal\", \"message\":\"Could not load configuration (%s)\", \"timestamp\": %d", err, time.Now().Unix())
	}
	logger, err := cfg.Logger.Build()
	if err != nil {
		log.Fatalf("{\"level\":\"fatal\", \"message\":\"Can't initialize zap logger (%s)\", \"timestamp\": %d", err, time.Now().Unix())
	}

	conn, err := GetDbConn(cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.DbUser, cfg.DbPassword)
	if err != nil {
		errMsg := fmt.Sprintf("Error while connecting to db (%s)", err)
		logger.Fatal(errMsg, zap.Int64("timestamp", time.Now().Unix()))
	}

	store := sessions.NewCookieStore([]byte(cfg.Sessions.AuthenticationKey), []byte(cfg.Sessions.EncryptionKey))
	store.Options.SameSite = http.SameSiteStrictMode

	return cfg, store, conn, logger
}
