# Build executable.
# docker build -t re_web_page_analyzer -f Dockerfile .
# docker run -it -d -p 9088:9088 --env-file=cmd/.env --name re_web_page_analyzer re_web_page_analyzer:latest
FROM golang:1.15 AS build-env

ENV GO111MODULE=on
WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR cmd

RUN go build -o /go/bin/analyzer

# final stage
FROM golang:1.15
WORKDIR /service/cmd
RUN mkdir ../static
COPY --from=build-env /go/bin/analyzer .
COPY --from=build-env src/static/form.html ../static/.
COPY --from=build-env src/static/report.html ../static/.
ENTRYPOINT ["./analyzer"]