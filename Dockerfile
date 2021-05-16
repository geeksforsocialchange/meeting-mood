# syntax=docker/dockerfile:1
FROM golang:1.16-alpine
WORKDIR /go/src/github.com/wheresalice/meeting-mood/
COPY . .
#RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o meeting-mood .

FROM scratch
COPY --from=0 /go/src/github.com/wheresalice/meeting-mood/meeting-mood /
COPY Procfile /
CMD ["/meeting-mood"]