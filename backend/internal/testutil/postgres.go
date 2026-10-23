// Package testutil provides a real, throwaway PostgreSQL server for tests.
//
// It uses embedded-postgres, which downloads official Postgres binaries on first
// use (cached in ~/.embedded-postgres-go) so tests need neither Docker nor a
// local install. Run `go test -short ./...` to skip database-backed tests.
package testutil

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync/atomic"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/wayne-stays/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	user     = "wayne_stays"
	password = "test"
)

var (
	server  *embeddedpostgres.EmbeddedPostgres
	port    uint32
	dbCount atomic.Int32
)

// StartPostgres starts a Postgres server for the current test binary. Call it
// from TestMain and run the returned stop function after m.Run.
func StartPostgres() (stop func(), err error) {
	free, err := freePort()
	if err != nil {
		return nil, err
	}
	runtimePath, err := os.MkdirTemp("", "wayne-stays-pg-")
	if err != nil {
		return nil, err
	}

	port = free
	server = embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Version(embeddedpostgres.V16).
		Port(port).
		Username(user).
		Password(password).
		RuntimePath(runtimePath).
		Logger(io.Discard))

	if err := server.Start(); err != nil {
		os.RemoveAll(runtimePath)
		server = nil
		return nil, fmt.Errorf("failed to start embedded postgres: %w", err)
	}

	return func() {
		server.Stop()
		os.RemoveAll(runtimePath)
	}, nil
}

// NewDatabase creates an empty database for one test and returns a connection
// to it. Tests are skipped when the server was not started (e.g. -short).
func NewDatabase(t *testing.T) *gorm.DB {
	t.Helper()
	if server == nil {
		t.Skip("embedded postgres not running (use go test without -short)")
	}

	name := fmt.Sprintf("test_%d", dbCount.Add(1))
	admin := open(t, "postgres")
	if err := admin.Exec("CREATE DATABASE " + name).Error; err != nil {
		t.Fatalf("failed to create database %s: %v", name, err)
	}

	return open(t, name)
}

func open(t *testing.T, dbName string) *gorm.DB {
	t.Helper()
	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     fmt.Sprint(port),
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  "disable",
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to %s: %v", dbName, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}
	t.Cleanup(func() { sqlDB.Close() })

	return db
}

func freePort() (uint32, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return uint32(l.Addr().(*net.TCPAddr).Port), nil
}
