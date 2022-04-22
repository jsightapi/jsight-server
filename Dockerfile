FROM golang:1.17-alpine as builder

RUN apk add --no-cache git
ARG GITHUB_TOKEN
# Install jschema library dependency
ARG JSCHEMA_BRANCH
WORKDIR /go/src/j/schema
RUN git clone -b ${JSCHEMA_BRANCH} --depth 1 \
          https://${GITHUB_TOKEN}@github.com/jsightapi/jsight-schema-go-library.git . \
    && git branch --show-current \
    && git show -s

# Install japi library dependency
ARG JAPI_BRANCH
WORKDIR /go/src/j/japi
RUN git clone -b ${JAPI_BRANCH} --depth 1 \
          https://${GITHUB_TOKEN}@github.com/jsightapi/jsight-api-go-library.git . \
    && git branch --show-current \
    && git show -s

# build
WORKDIR /go/src/j/server
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o jsight-server .


FROM alpine
ARG CORS
ARG STATISTICS
ENV JSIGHT_SERVER_CORS=$CORS
ENV JSIGHT_SERVER_STATISTICS=$STATISTICS
COPY --from=builder /go/src/j/server/jsight-server .
EXPOSE 8080
CMD [ "/jsight-server" ]
