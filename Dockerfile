FROM golang:1.17-alpine as builder

WORKDIR /go/src/github.com/jsightapi/jsight-server
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/bin/jsight-server .

FROM scratch
ARG CORS
ARG STATISTICS
ENV JSIGHT_SERVER_CORS=$CORS
ENV JSIGHT_SERVER_STATISTICS=$STATISTICS
COPY --from=builder /go/bin/jsight-server .
EXPOSE 8080
CMD [ "/jsight-server" ]
