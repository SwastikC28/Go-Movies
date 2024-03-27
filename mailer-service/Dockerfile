# First State- Create Binary File

FROM golang as builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -ldflags "-s -w" -ldflags '-extldflags "-static"' -o main ./cmd


# Stage 2: Create a minimal production image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage to the production image
COPY --from=builder /build/main .

COPY ./internal/templates .


# Command to run the executable
CMD ["./main"]