# Use a Docker multi-stage build to create a lean production image.
# Use the offical Golang image to create a build artifact.
FROM golang as builder

# Copy local code to the container image.
WORKDIR /src/app
COPY go.mod ./
RUN go mod download
COPY . .

# Build the outyet command inside the container.
# RUN CGO_ENABLED=0 GOOS=linux go build -v -o hello
RUN cd cmd && CGO_ENABLED=0 go build -o /inventory-app -ldflags="-w -s" . 

# Copy the binary to the production image from the builder stage.
FROM alpine
COPY --from=builder /inventory-app app/inventory-app
COPY tmpl/* /tmpl/
COPY key/* /key/
EXPOSE 8080
ENTRYPOINT ["/app/inventory-app"]