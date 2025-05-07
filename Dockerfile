FROM golang:1.24 as builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0

COPY / /work
WORKDIR /work

RUN go build -o hec-mock main.go
RUN strip -s hec-mock

FROM scratch
COPY --from=builder /work/hec-mock /hec-mock

USER 999
ENTRYPOINT ["/hec-mock"]

EXPOSE 8080
