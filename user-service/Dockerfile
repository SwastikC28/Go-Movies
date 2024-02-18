FROM golang as builder

ENV GO111MODULE=on
RUN mkdir /build
WORKDIR /build

COPY . .

RUN apt-get update && apt-get install wget binutils xz-utils -y && apt-get clean
RUN wget -q https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz
RUN tar xf upx-3.96-amd64_linux.tar.xz && mv upx-3.96-amd64_linux/upx .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -ldflags "-s -w" -ldflags '-extldflags "-static"' -o main ./cmd
# RUN strip main && ./upx --best --lzma main

# FROM scratch
FROM bash
COPY --from=builder /build/main /app/
WORKDIR /app

CMD ["./main"]


