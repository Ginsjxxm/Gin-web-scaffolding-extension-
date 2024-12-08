# Gin-web-scaffolding-extension
---

# Gin Framework Optimization with Zap Logger and Configuration Management

This repository demonstrates how to optimize the Gin framework by integrating:
- **Customizable Zap logger** with advanced features like file rotation.
- **Dynamic configuration loading** using `viper`.
- **Custom middleware** for logging and panic recovery.

## Features
1. **Dynamic Configuration**: Configuration is loaded from a YAML file and supports hot-reloading with Viper.
2. **Enhanced Logging**: Zap logger is integrated to provide structured and efficient logging.
3. **Custom Middlewares**:
    - `GinLogger`: Enhanced request logging.
    - `GinRecovery`: Graceful panic recovery with optional stack trace.

---

## Configuration File (`config.yaml`)

Example of a configuration file to be used with this project:
```yaml
app:
  name: "web_app"
  mode: "dev"
  port: 8081
  version: "v0.1.4"
log:
  level: "debug"
  filename: "web_app.log"
  max_size: 200
  max_age: 30
  max_backups: 7
mysql:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "root"
  dbname: "sql_demo"
  max_open: 200
  max_idle: 50
redis:
  host: 10.188.188.226
  port: 6379
  db: 0
```

---

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-repo/gin-optimization.git
   cd gin-optimization
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

---

## Usage

### 1. Initialize Configuration and Logger
Ensure the configuration file (`config.yaml`) is present in the project root.

```go
package main

import (
	"log"
	"os"

	"BlueBell/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// Load configuration
	if err := initConfig(); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Initialize logger
	if err := logger.Init(&logger.Config{
		Filename:   viper.GetString("log.filename"),
		MaxSize:    viper.GetInt("log.max_size"),
		MaxBackups: viper.GetInt("log.max_backups"),
		MaxAge:     viper.GetInt("log.max_age"),
	}); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zap.L().Sync()

	// Initialize Gin
	gin.SetMode(viper.GetString("app.mode"))
	r := gin.New()

	// Add middlewares
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// Define routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, Gin!"})
	})

	// Run server
	port := viper.GetString("app.port")
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initConfig() error {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
```

### 2. Run the Server
```bash
go run main.go
```

---

## Logger Customization

The logger package uses `zap` for logging. The following components can be customized:
- **Log Levels** (`debug`, `info`, `warn`, `error`, etc.).
- **Log Rotation** using `lumberjack` for size-based file rotation.

Example:
```yaml
log:
  level: "info"
  filename: "app.log"
  max_size: 100
  max_age: 7
  max_backups: 3
```

---

## Middleware Overview

### GinLogger
Logs details about each request, including:
- HTTP method and path.
- Query parameters.
- Client IP address.
- User agent.
- Response status and time taken.

### GinRecovery
Handles panics gracefully, preventing server crashes:
- Logs error details and stack trace (if enabled).
- Returns a 500 Internal Server Error.

---

## Folder Structure

```
.
├── config.yaml         # Configuration file
├── main.go             # Main application entry point
├── logger/             # Custom logger and middleware
│   ├── logger.go       # Zap logger integration
│   └── middleware.go   # GinLogger and GinRecovery middlewares
├── go.mod              # Go module file
└── README.md           # Documentation
```

---

## Dependencies
- [Gin](https://github.com/gin-gonic/gin): Web framework.
- [Viper](https://github.com/spf13/viper): Configuration management.
- [Zap](https://github.com/uber-go/zap): High-performance structured logging.
- [Lumberjack](https://github.com/natefinch/lumberjack): Log file rotation.

---

## License
This project is licensed under the MIT License.