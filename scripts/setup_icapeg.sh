#!/bin/bash

# Check for Go installation
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Installing Go..."

    # Install Go (adjust the installation method based on your OS)
    if [[ "$(uname -s)" == "Darwin" ]]; then # macOS
        brew install go
    elif [[ "$(uname -s)" == "Linux" ]]; then # Linux
        sudo apt-get update
        sudo apt-get install golang-go
    else
        echo "Unsupported operating system. Please install Go manually."
        exit 1
    fi
fi

# Check Go version
go_version=$(go version | awk '{print $3}')
if [[ "${go_version:2}" < "1.16" ]]; then
    echo "Go version is too old. Installing Go 1.16 or later..."

    # Install Go (adjust the installation method based on your OS)
    if [[ "$(uname -s)" == "Darwin" ]]; then # macOS
        brew install go@1.16
    elif [[ "$(uname -s)" == "Linux" ]]; then # Linux
        # You might need to add a PPA or download the Go binaries directly for older versions
        echo "Installing older Go versions on Linux might require additional steps. Please refer to the Go installation documentation."
        exit 1
    else
        echo "Unsupported operating system. Please install Go 1.16 or later manually."
        exit 1
    fi
fi

# Check if a proxy server is configured (simplified check for environment variables)
if [[ -z "$http_proxy" && -z "$https_proxy" ]]; then
    echo "Proxy server environment variables are not set. Configuring proxy..."

    # Set proxy environment variables (replace with your actual proxy server details)
    export http_proxy="http://your_proxy_server:proxy_port"
    export https_proxy="http://your_proxy_server:proxy_port"

    echo "Proxy server configured. Please ensure your proxy server is running."
fi

# Clone the ICAPeg repository
if [ ! -d "icapeg" ]; then
    git clone https://github.com/egirna/icapeg.git
else
    echo "icapeg repository already exists. Skipping cloning."
fi

# Navigate to the ICAPeg directory
cd icapeg

# Build ICAPeg binary
go build .

# Get port number from user input
read -p "Enter the port number for ICAPeg (default is 1344): " port_number
port_number=${port_number:-1344} # Use default 1344 if no input is provided

# Execute ICAPeg with the specified port
./icapeg -port $port_number 
