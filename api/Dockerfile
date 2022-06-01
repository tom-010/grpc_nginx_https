FROM golang AS build

ENV location /app

WORKDIR ${location}/server

ADD ./server.go ./
ADD ./go.mod ${location}/
ADD ./proto ${location}/proto

RUN go get -d ./... 
RUN go install ./...

RUN CGO_ENABLED=0 go build -o /bin/server

# Build stage II : Go binaries are self-contained executables.
FROM scratch
COPY --from=build /bin/server /bin/server

ENTRYPOINT ["/bin/server"]
EXPOSE 50051