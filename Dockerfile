FROM golang:1.17

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . .

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin


CMD ["air"]