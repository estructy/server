---
tags:
    - middleware
    - register
last_reviewed: 2025-08-03
---

# Register Middleware

Under the `api/{version}/middlewares` directory, you should:

1. **Create a new directory** for your middleware. 
```bash
server/
├── api/                  
│   └── v1/
│     ├── handlers/ 
│     ├── middlewares/    
│     │   └── your_middleware/ # New middleware directory
│     │   └── middlewares.go
│     └── routes/
├── cmd/
├── docs/              
├── internal/              
├── docker-compose.yaml 
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

2. **Create a `middlewares.go` file** inside your middleware directory. This file will contain the logic for your middleware.
```go
package name_middleware

import (
	"net/http"
)

type NameMiddleware struct {
    dependency string
}

func NewNameMiddleware(logger *zap.Logger) *NameMiddleware {
	return &NameMiddleware{
		logger: logger,
	}
}

func (m *NameMiddleware) Handle(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Middleware logic goes here
		next.ServeHTTP(w, r)
        // Optional: post-processing logic
	})
}
```

3. **Register your middleware** in the `middlewares.go` file located in the `api/{version}/middlewares` directory. This file is responsible for initializing and registering all middlewares.
```go 
package middleswares

import (
	"net/http"

	logger_middleware "github.com/estructy/server/api/v1/middlewares/logger"
	"go.uber.org/zap"
)

type Middleware func(http.Handler) http.HandlerFunc

type MiddlewareOrchestrator struct {
	Logger func(http.Handler) http.HandlerFunc
    NameMiddleware func(http.Handler) http.HandlerFunc
}

func NewMiddlewareOrchestration(logger *zap.Logger) *MiddlewareOrchestrator {
	loggerMiddleware := logger_middleware.NewLoggerMiddleware(logger)
    nameMiddleware := name_middleware.NewNameMiddleware(logger) // Replace with your middleware

	return &MiddlewareOrchestrator{
		Logger: loggerMiddleware.Handle,
        NameMiddleware: nameMiddleware.Handle, // Register your middleware here
	}
}

```

4. **Update the `main.go` file** in the `cmd/controlriver` directory to include your middleware orchestrator. This is where you initialize your middleware and pass it to the router.

5. **Aply your middleware** in the `api/{version}/routes/routes.go` file. This is where you define your routes and apply the middleware to them.
```go
func health(router *http.ServeMux, middlewares *middleswares.MiddlewareOrchestrator) {
	router.HandleFunc("GET /health", middlewares.Chain(health_handler.GetHealth, middlewares.Logger, middlewares.NameMiddleware)) // Replace with your middleware
}
```

