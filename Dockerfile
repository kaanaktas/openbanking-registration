FROM scratch

# Copy the binary file
COPY ./main ./
# Copy aspsp config folder
COPY ./aspsp ./aspsp

EXPOSE 8080

# Command to run the executable
CMD ["./main"]
