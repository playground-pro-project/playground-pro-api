FROM golang:1.20-alpine

# create directory folder
RUN mkdir /app

# set working directory
WORKDIR /app

# COPY ./ /app
COPY app/ /app/
COPY features/ /app/features/
COPY utils/ /app/utils/
COPY main.go /app/main.go
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
COPY docker-compose.yaml /app/docker-compose.yaml
COPY redis.conf /app/redis.conf

RUN go mod tidy

# create executable file with name "playgroundpro-api"
RUN go build -o playgroundpro-api

# run executable file
CMD ["./playgroundpro-api"]