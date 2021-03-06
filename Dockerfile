FROM golang as build
WORKDIR /work
COPY . .
ENV CGO_ENABLED 0
RUN go build -o spas ./app
RUN chmod +x ./spas
RUN echo "spasuser:x:65534:65534:spasuser:/:" > /etc_passwd

FROM scratch as run
WORKDIR /app
COPY --from=build /work/spas .
COPY --from=build /etc_passwd /etc/passwd
USER spasuser
ENTRYPOINT ["/app/spas"]