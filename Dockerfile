# ビルドステージ
FROM golang:1.19-alpine AS builder

WORKDIR /app

# 依存関係ファイルをコピーし、ダウンロードする
COPY go.mod go.sum ./
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 実行ステージ
FROM alpine:latest

# ca-certificates をインストール
RUN apk --no-cache add ca-certificates

# 非 root ユーザーを作成
RUN adduser -D appuser
USER appuser

WORKDIR /app

# ビルドステージからビルドされた実行ファイルをコピー
COPY --from=builder /app/main .

CMD ["./main"]