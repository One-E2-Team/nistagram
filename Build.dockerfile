FROM golang AS build
COPY . src/nistagram/
ARG ARG_TARGET_MS_GO
ENV CGO_ENABLED=0
#https://stackoverflow.com/questions/36279253/go-compiled-binary-wont-run-in-an-alpine-docker-container-on-ubuntu-host
RUN cd src/nistagram && go mod download && go mod verify && go build -o exec ${ARG_TARGET_MS_GO}/main.go


FROM alpine AS image
COPY --from=build /go/src/nistagram/exec /nistagram/microservice/exec
EXPOSE 8080
ENTRYPOINT ["/bin/sh", "-c" , "cd nistagram/microservice && ./exec"]