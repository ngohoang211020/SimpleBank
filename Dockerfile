# Use the official Golang image with version 1.20 on Alpine Linux 3.16
#Step 1: Build Stage
FROM golang:1.20-alpine3.16 as builder
#Sets the working directory inside the Docker image to /app. This is the directory where subsequent commands will be executed.
WORKDIR /app
# Copies all files and directories from the current directory (on the host machine) to the /app directory inside the Docker image
COPY . .
# Build the application with the output binary named "main" from the main.go file
RUN go build -o main main.go
#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
#_______________________________________________________________________________
#Step 2: Run Stage
FROM alpine:3.13
WORKDIR /app
#The COPY --from=builder command copies the "main" binary from the previous build stage into the current stage.
COPY --from=builder /app/main .
#COPY --from=builder /app/migrate.linux-amd64 ./migrate

COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration
ENV GIN_MODE=release
ENV PORT=8000
# Informs Docker that the container will listen on port 8080 at runtime.
#However, this does not actually publish the port; it serves as documentation for users of the image
EXPOSE $PORT
# Specifies the default command to run when the container starts.
#it's running the binary executable "main" located in the /app directory inside the container.
#This means that when you run a container from this image, it will automatically execute the "main" binary as the main process.
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]