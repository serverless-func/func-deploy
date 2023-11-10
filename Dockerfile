FROM golang:alpine as builder

ARG USE_MIRROR

RUN mkdir -p /app

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN if [ "$USE_MIRROR" = "true" ]; then go env -w GOPROXY=https://goproxy.cn,direct; fi

ENV CGO_ENABLED=0

RUN go mod download

COPY . .

RUN GOOS=linux go build -o /bin/app .

FROM bitnami/kubectl:1.27.7
LABEL maintainer="mail@dongfg.com"

ENV TZ=Asia/Shanghai

COPY --from=builder /bin/app /bin/app

ENTRYPOINT ["/bin/app"]