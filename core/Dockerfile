FROM golang
COPY ./ /core
WORKDIR /core
RUN go build -ldflags "-w -extldflags -static" -tags netgo -installsuffix netgo -o ./core

FROM scratch
COPY --from=0 /core/core .
ENTRYPOINT [ "./core" ]
