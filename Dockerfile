# NOTE: this docker file came 

# Build the application from source
FROM golang:1.23.2-alpine AS build-stage

# copy broker
WORKDIR /home/app/broker
COPY ./ ./

# build broker's binary 
WORKDIR /home/app/broker
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/broker

# rebuild the application binary into a lean image
# FROM gcr.io/distroless/base-debian11 AS build-release-stage
FROM scratch AS build-release-stage

WORKDIR /

COPY --from=build-stage /home/app/broker/bin/broker /broker

EXPOSE 4000
