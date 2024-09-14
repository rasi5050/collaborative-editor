FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o collaborative-editor-backend .

EXPOSE 8080

CMD ["./collaborative-editor-backend"]
