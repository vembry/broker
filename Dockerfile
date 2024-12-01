# NOTE: this docker file came 

# Build the application from source
FROM golang:1.23.2-alpine AS base

# copy broker
WORKDIR /home/app/broker
COPY ./ ./

FROM base AS build

# build broker's binary 
WORKDIR /home/app/broker
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/broker

FROM build AS dev
# add debug tools
# ...

# rebuild the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS release

WORKDIR /

COPY --from=build /home/app/broker/bin/broker /

EXPOSE 4000
