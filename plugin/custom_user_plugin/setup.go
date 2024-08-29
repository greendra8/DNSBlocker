package custom_user_plugin

import (
    "database/sql"
    "fmt"  // Add this line
    "log"
    "net/http"

    "github.com/coredns/caddy"
    "github.com/coredns/coredns/core/dnsserver"
    "github.com/coredns/coredns/plugin"
    _ "github.com/lib/pq"
)

func init() {
    log.Println("Registering custom_user_plugin")
    plugin.Register("custom_user_plugin", setup)
}

// setup is the function that gets called when the Corefile is parsed.
func setup(c *caddy.Controller) error {
    log.Println("Setting up custom_user_plugin")
    baseDomain := ""
    httpAddr := ":443" // Default HTTPS port
    dbConnString := "" // Database connection string

    // Parse the configuration from the Corefile
    for c.Next() {
        if c.NextArg() {
            baseDomain = c.Val()
        }
        if c.NextArg() {
            httpAddr = c.Val()
        }
        if c.NextArg() {
            dbConnString = c.Val()
        }
    }

    if baseDomain == "" {
        return plugin.Error("custom_user_plugin", c.Err("base domain is required"))
    }

    if dbConnString == "" {
        return plugin.Error("custom_user_plugin", c.Err("database connection string is required"))
    }

    log.Printf("custom_user_plugin: Base domain set to %s", baseDomain)
    log.Printf("custom_user_plugin: HTTP server listening on %s", httpAddr)

    // Initialize database connection
    db, err := sql.Open("postgres", dbConnString)
    if err != nil {
        return plugin.Error("custom_user_plugin", err)
    }

    // Ping to verify connection
    if err := db.Ping(); err != nil {
        return plugin.Error("custom_user_plugin", fmt.Errorf("unable to connect to database: %v", err))
    }
    log.Println("Successfully connected to the database")

    // Initialize the plugin
    p := &CustomUserPlugin{
        BaseDomain: baseDomain,
        DB:         db,
    }

    // Set up HTTP server for DoH
    go func() {
        http.HandleFunc("/", p.ServeHTTP)
        err := http.ListenAndServeTLS(httpAddr, "cert.pem", "key.pem", nil)
        if err != nil {
            log.Fatalf("Failed to start HTTPS server: %v", err)
        }
    }()

    // Add the plugin to CoreDNS
    dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
        p.Next = next
        return p
    })

    return nil
}