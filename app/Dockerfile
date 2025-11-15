# 開発用イメージ
FROM golang:1.25.1-alpine

# 必要なパッケージ
RUN apk add --no-cache git bash curl

WORKDIR /app

# gqlgen や air をインストール
RUN go install github.com/99designs/gqlgen@latest
RUN go install github.com/air-verse/air@latest

# golangci-lint の最新版をインストール
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.6.2
# ホットリロード用設定ファイルコピー（任意）
COPY .air.toml /app/.air.toml

# ソースコードはマウントするので COPY は不要（マウントする場合）
COPY . .
RUN go get -tool github.com/evilmartians/lefthook/v2

EXPOSE 8080

# 開発用に air で起動
CMD ["air"]
