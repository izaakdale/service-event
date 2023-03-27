FROM golang:1.20-alpine as builder
WORKDIR /

COPY . ./
RUN go mod download


RUN go build -o /service-event

FROM alpine
COPY --from=builder /service-event .

EXPOSE 80
CMD [ "/service-event" ]