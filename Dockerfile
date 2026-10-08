FROM golang:1.22 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o bootstrap .

FROM public.ecr.aws/lambda/provided:al2023

COPY --from=builder /src/bootstrap ${LAMBDA_TASK_ROOT}/bootstrap

CMD ["bootstrap"]
