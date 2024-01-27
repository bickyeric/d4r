# -----------------------------------------------------------------------------
# The image for developing d4r in container

FROM golang:1.20.10 AS develop

RUN go install golang.org/x/tools/gopls@latest && \
    go install github.com/cweill/gotests/gotests@latest && \
    go install github.com/fatih/gomodifytags@latest && \
    go install github.com/josharian/impl@latest && \
    go install github.com/haya14busa/goplay/cmd/goplay@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest

# -----------------------------------------------------------------------------
# The base image for building the d4r binary

FROM golang:1.20.10-alpine AS build

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o d4r main.go

FROM alpine
COPY --from=build /app/d4r /bin/d4r
ENTRYPOINT "/bin/d4r"
