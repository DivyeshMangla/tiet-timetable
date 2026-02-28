# Build frontend
FROM node:22-alpine AS frontend

WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install

COPY frontend/ .
RUN npm run build

# Build backend
FROM golang:1.25-alpine AS backend

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app ./cmd/app

# Final image
FROM gcr.io/distroless/static-debian12

WORKDIR /app

USER nonroot:nonroot

COPY --from=backend /app/app .
COPY --from=backend /app/frontend/dist ./frontend/dist

EXPOSE 8080

ENTRYPOINT ["./app"]