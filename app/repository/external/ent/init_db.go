package ent

import (
	"context"
	"fmt"
	"log"

	"sleeve/ent"
	"sleeve/repository/external/utils"

	_ "github.com/lib/pq"
)

// NewDBClient は環境変数から DB 接続情報を読み込み、Ent Client を生成します
func NewDBClient() (*ent.Client, error) {
	var dsn string
	var client *ent.Client
	var err error

	dsn, err = utils.CreateDSN()
	if err != nil {
		return nil, err
	}

	// Ent Client を生成
	client, err = ent.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	return client, nil
}

func main() {
	var client *ent.Client
	var err error

	// DB クライアントを生成
	client, err = NewDBClient()
	if err != nil {
		log.Fatalf("DB接続エラー: %v", err)
	}
	defer client.Close()

	// Ping して接続確認
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	log.Println("DB connected and schema applied!")
}
