# Serverul MCP public (romana-public), pentru mcp.doitmagic.ro.
# Se construiește din rădăcina repo-ului: docker compose -f deploy/mcp/docker-compose.yml up -d --build
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSIUNE=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.versiune=${VERSIUNE}" -o /romana-public ./cmd/romana-public

# distroless: fără shell și fără manager de pachete; are certificatele CA (pentru ACCES=chatgpt)
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /romana-public /romana-public
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/romana-public"]
