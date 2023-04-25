FROM golang:1.18 as builder

# Set some information.
# For hosting yourself modify these fields. 
LABEL deployedBy="Kjetil Indrehus"
LABEL maintainer="kjetikin@stud.ntnu.no"
LABEL stage=builder

# Copy all files to the app folder.
# The reason why we copy all is beacause all files are relavent to the deployment
COPY . /go/src/app/

# Set up execution environment in container's GOPATH
WORKDIR /go/src/app/cmd

# Run go mod tidy to ensure accuracy of go.mod file
RUN go mod tidy

# Download packages
RUN go mod download 

# Compile binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o server

# Move the executable one level up so that all the paths are correct.
# See the compose file for more information about where the resources are set.
RUN mv server ../server

# Change the working directory back for running the executable. 
# Now on the same level
WORKDIR /go/src/app

# Indicate port on which server listens
EXPOSE 8080

# Instantiate binary
CMD ["./server"]

