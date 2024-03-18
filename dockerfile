FROM golang:1.22

RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12

WORKDIR /cmd

COPY . .

RUN swag init -g cmd/main.go --parseDependency --parseInternal -d ./,internal/structures,pkg/handlers
RUN go build cmd/main.go

ENTRYPOINT ["/cmd/main"]