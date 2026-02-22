# CodeOwner: @docker_owner

# This is a single-line comment

FROM alpine:3.18

# Install dependencies
RUN apk add --no-cache curl  # Inline comment

ENV APP_PORT=8080  # Inline comment

COPY . /app  # Copy application files

WORKDIR /app

# Run the application
CMD ["./start.sh"]
