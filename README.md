# 🐭 HamstersTunnel

HamstersTunnel is a flexible tunneling service that allows clients 
to configure and manage network services with different protocols (HTTP, TCP, UDP). 
It provides an API to manage and create service configurations for tunneling,
enabling dynamic mapping of local services to public ports on a remote server.

## ✨ Features

✅ Reverse Proxy for **TCP** connections  
🔄 Planned support for **UDP** and **HTTP**  
🛠️ **Built-in Service Management** for dynamic port mapping  
⚡ **Fast and lightweight**, ideal for tunneling and proxying  
🔒 **Security-focused**, designed with service isolation  

## 📌 Roadmap

✅ **TCP Reverse Proxy** (Base version)  
🔄 **UDP Support** (In Progress)  
🔄 **HTTP Proxy Support** (Planned)  
📖 **Documentation & Examples**  
🚀 **Performance Optimizations**  

## 🚀 Installation & Usage

Clone the repository:

```sh
git clone https://github.com/typegaro/HamstersTunnel.git
cd HamstersTunnel
```

### 🏗️ Build & Run

Use **make** for easy management:

```sh
# Run the server
test
make server

# Run tests
make test

# Build the server
make build-server

# Clean build files
make clean
```

### 🐳 Docker

To run **HamstersTunnel** using Docker, follow these steps:

1. **Build the Docker image**:  
   First, ensure you've built the Docker image using the following command:

   ```sh
   docker build -t HamstersTunnel .
   ```

2. **Run the Docker container**:  
   You can run the Docker container and map the internal port (8080) to an external port on your host using the `docker run` command:

   ```sh
   docker run -p 8080:8080 HamstersTunnel
   ```

   This will start the server inside the container and expose the internal port `8080` to the external port `8080` on your host. You can access the service by navigating to `http://localhost:8080` in your browser.

3. **Custom port mapping**:  
   If you want to use a different external port, you can modify the command like so:

   ```sh
   docker run -p <external-port>:8080 HamstersTunnel
   ```

   Replace `<external-port>` with the desired port on your host (e.g., `9000:8080`).


## 🏆 Contributing

Contributions are welcome! Feel free to open issues and PRs.



