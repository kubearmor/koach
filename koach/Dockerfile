### Builder

FROM golang:1.18.5-alpine3.16 as builder

RUN apk update
RUN apk add --no-cache make
RUN apk add build-base

WORKDIR /usr/src/koach

COPY koach .

RUN ls
RUN make build-binary

### Copy executable image

FROM alpine:3.14

COPY --from=builder /usr/src/koach/kubearmor-koach /KubeArmor/kubearmor-koach

### Environment variables

ENV DATABASE_FILE_PATH /KubeArmor/koach.db

ENTRYPOINT ["/KubeArmor/kubearmor-koach"]
