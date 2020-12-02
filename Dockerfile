FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Copy the binary file
COPY ./main ./
# Copy aspsp config folder
COPY ./aspsp ./aspsp

EXPOSE 8080

# Command to run the executable
CMD ["./main"]
