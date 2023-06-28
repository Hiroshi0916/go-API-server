# Goの公式イメージをベースにする
FROM golang:1.20

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係をコピー
COPY go.mod ./
COPY go.sum ./

# 依存関係をダウンロードします。
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド
RUN go build -o main .

# ポートのエクスポート
EXPOSE 8000

# 実行
CMD ["/app/main"]
