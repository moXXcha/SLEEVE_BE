package utils

import (
	"fmt"
	"os"
)

// CreateDSN は環境変数から PostgreSQL の DSN を構築します
func CreateDSN() (string, error) {
	var dsn string
	var postgres_user string
	var postgres_password string
	var postgres_host string
	var postgres_port string
	var postgres_db string

	// 環境変数から DB 接続情報を取得
	postgres_user = os.Getenv("POSTGRES_USER")
	postgres_password = os.Getenv("POSTGRES_PASSWORD")
	postgres_host = os.Getenv("POSTGRES_HOST")
	postgres_port = os.Getenv("POSTGRES_PORT")
	postgres_db = os.Getenv("POSTGRES_DB")

	// 必須の環境変数チェック
	if postgres_user == "" || postgres_password == "" || postgres_host == "" || postgres_port == "" || postgres_db == "" {
		return "", fmt.Errorf("必要な環境変数が設定されていません: POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST, POSTGRES_PORT, POSTGRES_DB")
	}

	// DSN を構築
	dsn = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		postgres_user,
		postgres_password,
		postgres_host,
		postgres_port,
		postgres_db,
	)

	return dsn, nil
}
