#!/bin/bash

# Clone the icapeg repository
git clone https://github.com/egirna/icapeg.git

# Change directory to the icapeg project folder
cd icapeg

# Install dependencies (if any)
go get -v -t -d ./... 

# Build the Go project
go build .

# Run the compiled icapeg binary
./icapeg 
