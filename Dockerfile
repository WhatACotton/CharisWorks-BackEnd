# ベースとなるDockerイメージ指定
FROM golang:latest
# コンテナログイン時のディレクトリ指定
WORKDIR /appq
# ホストのファイルをコンテナの作業ディレクトリに移行
COPY . /appq/

EXPOSE 5000