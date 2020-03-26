FROM golang as build
WORKDIR /work
COPY . .
RUN go build -o spas ./app
RUN chmod +x ./spas

FROM ubuntu as run
WORKDIR /app
COPY --from=build /work/spas .
ENTRYPOINT ["/app/spas"]