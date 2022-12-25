# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod ./

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY . .

# compile application
RUN go build -o /godocker
 
# command to be used to execute when the image is used to start a container
CMD [ "/godocker" ]