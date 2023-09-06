# build stage
FROM golang:1.20-alpine3.18 AS build

COPY . /app

WORKDIR /app/cmd/user-segmentation

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o user-segmentation .

# run stage
FROM alpine:3.18 AS final

COPY --from=build /app/config/config.yml ./config/config.yml

COPY --from=build /app/cmd/user-segmentation/wait-for-postgres.sh \
    /app/cmd/user-segmentation/user-segmentation ./

RUN apk update && \
    apk add --no-cache postgresql15-client && \
    chmod +x wait-for-postgres.sh user-segmentation

EXPOSE 8081
CMD [ "./wait-for-postgres.sh", "postgres", "./user-segmentation" ]