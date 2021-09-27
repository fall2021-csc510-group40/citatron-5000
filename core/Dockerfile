FROM golang:1.17
COPY ./ /core
WORKDIR /core
RUN go build -ldflags "-w -extldflags -static" -tags netgo -installsuffix netgo -o ./core

FROM alpine
COPY --from=0 /core/core .
ENTRYPOINT [ "./core" ]
