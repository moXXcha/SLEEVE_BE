# 開発用イメージ
FROM golang:1.25.1-alpine

# 必要なパッケージ
RUN apk add --no-cache git bash

WORKDIR /app

# gqlgen や air をインストール
RUN go install github.com/99designs/gqlgen@latest
RUN go install github.com/air-verse/air@latest

# ホットリロード用設定ファイルコピー（任意）
COPY .air.toml /app/.air.toml

# ソースコードはマウントするので COPY は不要（マウントする場合）
COPY . .

EXPOSE 8080

# 開発用に air で起動
CMD ["air"]
