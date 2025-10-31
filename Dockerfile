# Use an existing image as a base
FROM uozi/debian-base-slim:latest

# Copy the compiled binary
COPY burn-api /app/burn-api

# Define the entrypoint
CMD ["/app/burn-api", "-config", "/config/app.ini"]
