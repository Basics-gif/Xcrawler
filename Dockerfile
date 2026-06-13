FROM golang:1.26-bookworm

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
  libnss3 libatk1.0-0 libatk-bridge2.0-0 \
  libcups2 libxkbcommon0 libxcomposite1 \
  libxdamage1 libxrandr2 libgbm1 libasound2 \
  && rm -rf /var/lib/apt/lists/*


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o crawler .

RUN go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps chromium

RUN mkdir -p debug

CMD ["./crawler"]
