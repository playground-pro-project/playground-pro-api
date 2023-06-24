FROM golang:1.19-alpine

# create directory folder
RUN mkdir /app

# set working directory
WORKDIR /app

COPY ./ /app

RUN go mod tidy

# create executable file with name "playgroundpro-api"
RUN go build -o playgroundpro-api

# run executable file
CMD ["./playgroundpro-api"]