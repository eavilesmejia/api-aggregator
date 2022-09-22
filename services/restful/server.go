package restful

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/monitor"
    "log"
    "os"
)

type config struct {
    useLogger bool
}

type router interface {
    SetupRoutes(r fiber.Router)
}

var routes = []router{
    newBlockRouter(),
    newTransactionRouter(),
}

// appSetup setups the config and routes, exported for use with tests
func appSetup(cfg ...*config) *fiber.App {
    
    var conf *config
    if len(cfg) > 0 {
        conf = cfg[0]
    } else {
        conf = &config{}
    }
    
    app := fiber.New()
    
    if conf.useLogger {
        app.Use(logger.New())
    }
    
    // Health check route (public)
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.Status(200).Send([]byte("OK"))
    })
    
    // Monitoring CPU
    app.Get("/dashboard", monitor.New())
    
    // Setup private routes
    group := app.Group("/v1")
    for _, route := range routes {
        route.SetupRoutes(group)
    }
    group.Get("/openapi.yaml", func(ctx *fiber.Ctx) error {
        ctx.Set("Content-Type", "text/yaml")
        return ctx.SendFile("./services/restful/openapi.yaml")
    })
    
    return app
}

func Run() {
    app := appSetup(&config{
        useLogger: false,
    })
    
    port := os.Getenv("RESTFUL_PORT")
    fmt.Printf("Starting restful server on port %s\n", port)
    
    // Listen on server 8080 (default) and fatal if error
    log.Fatal(app.Listen(fmt.Sprintf("0.0.0.0:%s", port)))
}
