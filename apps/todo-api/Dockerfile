FROM golang:1.19-alpine AS build
WORKDIR /apps/todo
COPY apps/todo-api/go.mod apps/todo-api/go.sum ./
COPY libs/common /libs/common
RUN go mod download
COPY apps/todo-api/cmd ./cmd
COPY apps/todo-api/pkg ./pkg
ENV GOARCH=amd64
RUN go build -o /go/bin/api cmd/api/main.go \
    && go build -o /go/bin/migrate cmd/migrate/main.go


FROM alpine:latest  
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
EXPOSE 8080

ENV TZ=Asia/Bangkok
ENV APP_MODE=production

COPY apps/todo-api/migrations /app/migrations
COPY --from=build /go/bin/api /app/api
COPY --from=build /go/bin/migrate /app/migrate


CMD ["/app/api"]  