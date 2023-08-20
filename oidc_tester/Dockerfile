# syntax=docker/dockerfile:1

FROM golang:1.20.4-alpine3.17

## We create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /oidc_tester
## We copy everything in the root directory
## into our /app directory
ADD . /oidc_tester
## We specify that we now wish to execute 
## any further commands inside our /app
## directory
WORKDIR /oidc_tester/src
## we run go build to compile the binary
## executable of our Go program
RUN go build -o main ./main

RUN chmod 755 -R /oidc_tester

EXPOSE 9000
## Our start command which kicks off
## our newly created binary executable
CMD ["main/main"]