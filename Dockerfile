FROM golang

MAINTAINER JACK
 
# Copy Local Project to Container's Workspace
RUN mkdir -p /go/src/github.com/lichao201626/fullstack/server/
WORKDIR /go/src/github.com/lichao201626/fullstack/server/
# COPY vendor/* /go/src/
COPY . /go/src/github.com/lichao201626/fullstack/server

RUN export GOPROXY=https://goproxy.cn
RUN export GO111MODULE=on

RUN cd /go/src/github.com/lichao201626/fullstack/server

RUN go mod init
RUN go mod download
RUN go mod tidy
RUN go mod vendor

#RUN go get ./...
RUN go build .

# Install all packages here
#RUN go get .
# RUN go build /go/src/github.com/lichao201626/fullstack/server/main.go

# Run when the container starts
# ENTRYPOINT /go/bin/github.com/lichao201626/fullstack/server

 # 配置环境变量
 ENV HOST 0.0.0.0
 ENV PORT 8889

# Service listens on port 8080.
EXPOSE 8889

CMD [ "./server" ]