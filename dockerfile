FROM golang:1.21.5 as builder

RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12

WORKDIR /cmd

COPY . .

RUN swag init -g cmd/main.go --parseDependency --parseInternal -d ./,internal/structures,pkg/handlers && go build cmd/main.go

ENTRYPOINT ["/cmd/main"]