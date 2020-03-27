FROM golang as build
WORKDIR /work
COPY . .
ENV CGO_ENABLED 0
RUN go build -o spas ./app
RUN chmod +x ./spas

FROM alpine as run
WORKDIR /app
COPY --from=build /work/spas .
ENTRYPOINT ["/app/spas"]