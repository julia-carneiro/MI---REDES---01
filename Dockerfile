FROM golang:1.23

WORKDIR /app

COPY go.mod . 
COPY /servidor /app/servidor
COPY /cliente /app/cliente

RUN go build -o bin-servidor ./servidor
RUN go build -o bin-cliente ./cliente


ENTRYPOINT [ "/app/bin-servidor" ]  
