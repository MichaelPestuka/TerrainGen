FROM golang:1.24.4
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY dla/*.go ./dla/
RUN go build -o /terraingenbackend
CMD ["/terraingenbackend"]