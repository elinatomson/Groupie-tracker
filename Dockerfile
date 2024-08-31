FROM golang:1.23.0

# metadata for image
LABEL project="Groupie-tracker"
LABEL author="Elina Tomson"
LABEL exercise="https://github.com/01-edu/public/tree/master/subjects/groupie-tracker"

# creating an /app directory within image that will hold the application source files
RUN mkdir /app

# copying everything from the root directory into /app directory
COPY . /app

# specifying that any further commands will be executed inside the /app directory
WORKDIR /app

# run go build to compile the binary executable of the Go program
RUN go build -o main .

EXPOSE 8080

# start command which kicks off created binary executable
CMD ["/app/main"]