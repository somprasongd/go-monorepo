FROM golang:1.19-alpine AS build
WORKDIR /apps/auth
COPY apps/auth-api/go.mod apps/auth-api/go.sum ./
COPY libs/common /libs/common
RUN go mod download
COPY apps/auth-api/config ./config
COPY apps/auth-api/cmd ./cmd
COPY apps/auth-api/pkg ./pkg
ENV GOARCH=amd64
RUN go build -o /go/bin/api cmd/api/main.go \
    && go build -o /go/bin/migrate cmd/migrate/main.go


FROM alpine:latest  
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
EXPOSE 8080

ENV TZ=Asia/Bangkok
ENV APP_MODE=production

COPY apps/auth-api/config /app/config
COPY apps/auth-api/migrations /app/migrations
COPY --from=build /go/bin/api /app/api
COPY --from=build /go/bin/migrate /app/migrate

CMD ["/app/api"]  