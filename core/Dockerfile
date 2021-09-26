FROM golang
COPY ./main.go .
RUN go get -d && go build -ldflags "-w -extldflags -static" -tags netgo -installsuffix netgo -o ./core

FROM scratch
COPY --from=0 /go/core .
ENTRYPOINT [ "./core" ]