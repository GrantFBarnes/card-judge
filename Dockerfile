FROM registry.hub.docker.com/library/golang
WORKDIR /app
COPY . .
WORKDIR /app/src
RUN go mod download
RUN go build -o /website
EXPOSE 8080
CMD ["/website"]