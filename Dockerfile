FROM golang:latest

LABEL creator="Clement Bolin" email="clement.bolin@epitech.eu" description="golang / prometheus image"

WORKDIR /app

RUN export GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 9000

CMD [ "./main" ]
