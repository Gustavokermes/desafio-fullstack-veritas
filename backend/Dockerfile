# Dockerfile para o backend Go do Kanban
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o kanban-backend

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/kanban-backend .
COPY --from=builder /app/tasks.json ./tasks.json
EXPOSE 8080
CMD ["./kanban-backend"]
