FROM node:18-bullseye-slim as build-frontend
WORKDIR /app
COPY ./templates .

RUN yarn install && yarn prod

FROM golang:1.21-bullseye as build

WORKDIR /go/src/app
COPY . .
COPY --from=build-frontend /app /go/src/app/templates

RUN go mod download
RUN go vet .
RUN CGO_ENABLED=0 go build -o /go/bin/app .

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/app /

EXPOSE 8080
CMD ["/app"]
