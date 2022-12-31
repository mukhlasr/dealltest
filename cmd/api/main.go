package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simpleblog/auth"
	"simpleblog/httphandler"
	"simpleblog/psqldb"
	"strings"
	"syscall"

	_ "github.com/lib/pq"
)

func main() {
	passwordGenerator := func() (string, error) {
		return auth.GeneratePassword(4, rand.Reader)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db := setupDB()

	jwtPrivKey := getJWTPrivKey(fromEnvWithDefault("AUTH_JWT_PRIVATE_KEY", httphandler.DefaultJWTPrivKey))
	startServer(ctx, &http.Server{
		Addr: getListenAddress(),
		Handler: RouteHTTPHandler(
			RouterDeps{
				UserStorage:       setupUserStorage(db),
				PostStorage:       setupPostStorage(db),
				PasswordHash:      sha256.New(),
				PasswordGenerator: passwordGenerator,
				TokenGenerator:    auth.ES256TokenGenerator(jwtPrivKey),
				TokenParser:       auth.ES256TokenParser(&jwtPrivKey.PublicKey),
			},
		),
	})
}

func startServer(ctx context.Context, srv *http.Server) {
	shutdownChan := make(chan struct{})

	var shutdownErr error
	go func() {
		<-ctx.Done()
		shutdownErr = srv.Shutdown(context.Background())
		close(shutdownChan)

	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Println("failed to start server:", err)
		return
	}

	<-shutdownChan
	if shutdownErr != nil {
		log.Println("failed to shutdown http server:", shutdownErr)
	}
}

func setupDB() *sql.DB {
	db, err := sql.Open("postgres", "")
	if err != nil {
		log.Fatalln("failed to open db:", err)
	}

	_, err = db.Exec(psqldb.DBSchema)
	if err != nil {
		log.Fatalln("failed to setup database:", err)
	}
	return db
}

func setupUserStorage(db *sql.DB) *psqldb.UserStorage {
	psqlStorage := (*psqldb.UserStorage)(db)
	return psqlStorage
}

func setupPostStorage(db *sql.DB) *psqldb.PostStorage {
	_, err := db.Exec(psqldb.DBSchema)
	if err != nil {
		log.Fatalln("failed to setup database:", err)
	}
	psqlStorage := (*psqldb.PostStorage)(db)
	return psqlStorage
}

func getJWTPrivKey(pubKeyString string) *ecdsa.PrivateKey {
	privkey, err := auth.LoadX509ECDSAKey(strings.NewReader(pubKeyString))
	if err != nil {
		log.Fatalln("failed to generate ecdsakey:", err)
	}
	return privkey
}

func getListenAddress() string {
	addr := fromEnvWithDefault("AUTH_LISTEN_ADDR", "0.0.0.0")
	port := fromEnvWithDefault("AUTH_LISTEN_PORT", "8888")
	return addr + ":" + port
}

func fromEnvWithDefault(envName, defaultValue string) string {
	res := os.Getenv(envName)
	if res == "" {
		return defaultValue
	}
	return res
}
