FROM golang:latest
ENV GOPROXY https://goproxy.cn,direct
RUN go build ./internal/growth_record/run.go
COPY ./internal/growth_record/run /growth_record
EXPOSE 8000
ENTRYPOINT ["/growth_record"]