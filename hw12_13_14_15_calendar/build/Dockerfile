# Собираем в гошке
FROM golang:1.22 as build

ARG SERVICE_NAME
ENV BIN_FILE /opt/${SERVICE_NAME}/${SERVICE_NAME}-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "${LDFLAGS}" \
        -o ${BIN_FILE} cmd/${SERVICE_NAME}/*

# На выходе тонкий образ
FROM alpine:3.9

ARG SERVICE_NAME
LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="$SERVICE_NAME"
LABEL MAINTAINERS="yevhenii.hulkov@gmail.com"

ENV BIN_FILE /opt/${SERVICE_NAME}/${SERVICE_NAME}-app
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/${SERVICE_NAME}/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}

RUN wget \
    https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
    sh

CMD ${BIN_FILE} -config ${CONFIG_FILE}
