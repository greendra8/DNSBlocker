package custom_user_plugin

import (
    "context"
    "database/sql"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "strings"

    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"
)

// CustomUserPlugin implements the CoreDNS plugin interface and handles user-specific DNS rules.
type CustomUserPlugin struct {
    Next       plugin.Handler
    BaseDomain string
    DB         *sql.DB
}

// ServeDNS implements the plugin.Handler interface.
func (p *CustomUserPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
    log.Printf("CustomUserPlugin: Received query for %s", r.Question[0].Name)

    userID, ok := ctx.Value("userID").(string)
    if !ok {
        log.Printf("CustomUserPlugin: No user ID found in context")
        return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
    }

    log.Printf("CustomUserPlugin: Request for user %s", userID)

    allowed := p.checkUserRules(userID, r.Question[0].Name)

    if !allowed {
        log.Printf("CustomUserPlugin: Query not allowed for user %s", userID)
        m := new(dns.Msg)
        m.SetReply(r)
        
        // Instead of NXDOMAIN, return a specific IP for blocked domains
        m.Answer = []dns.RR{
            &dns.A{
                Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 1},
                A:   net.ParseIP("0.0.0.0"),
            },
        }
        
        w.WriteMsg(m)
        return dns.RcodeSuccess, nil
    }

    log.Printf("CustomUserPlugin: Query allowed for user %s, passing to next plugin", userID)
    
    customWriter := &ttlWriter{ResponseWriter: w, ttl: 1}
    return plugin.NextOrFailure(p.Name(), p.Next, ctx, customWriter, r)
}

// Name implements the plugin.Handler interface.
func (p *CustomUserPlugin) Name() string { return "custom_user_plugin" }

// AddRule adds a blocking rule for a specific user and domain.
func (p *CustomUserPlugin) AddRule(userID, domain string) error {
    // First, check if the user exists
    var exists bool
    err := p.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
    if err != nil {
        log.Printf("Error checking if user exists: %v", err)
        return err
    }

    // If the user doesn't exist, create them
    if !exists {
        _, err = p.DB.Exec("INSERT INTO users (id) VALUES ($1)", userID)
        if err != nil {
            log.Printf("Error creating user %s: %v", userID, err)
            return err
        }
        log.Printf("Created new user: %s", userID)
    }

    // Now add the rule
    domain = strings.TrimSuffix(domain, ".") // Remove trailing dot if present
    _, err = p.DB.Exec("INSERT INTO blocked_sites (user_id, domain) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, domain)
    if err != nil {
        log.Printf("Error adding rule for user %s: %v", userID, err)
        return err
    }
    log.Printf("Added rule for user %s: blocking %s", userID, domain)
    return nil
}

// checkUserRules checks if a query is allowed based on user-specific rules.
func (p *CustomUserPlugin) checkUserRules(userID, qname string) bool {
    qname = strings.TrimSuffix(qname, ".") // Remove trailing dot from qname if present
    qnameParts := strings.Split(qname, ".")

    // Fetch all blocked domains for the user
    rows, err := p.DB.Query("SELECT domain FROM blocked_sites WHERE user_id = $1", userID)
    if err != nil {
        log.Printf("Error fetching rules for user %s: %v", userID, err)
        return true // Allow by default in case of error
    }
    defer rows.Close()

    for rows.Next() {
        var blockedDomain string
        if err := rows.Scan(&blockedDomain); err != nil {
            log.Printf("Error scanning blocked domain: %v", err)
            continue
        }

        blockedParts := strings.Split(blockedDomain, ".")
        if isSubdomain(qnameParts, blockedParts) {
            log.Printf("Blocking access to %s for user %s (matches blocked domain %s)", qname, userID, blockedDomain)
            return false
        }
    }

    return true
}

// isSubdomain checks if the query domain is a subdomain of the blocked domain
func isSubdomain(queryParts, blockedParts []string) bool {
    if len(queryParts) < len(blockedParts) {
        return false
    }
    for i := 1; i <= len(blockedParts); i++ {
        if queryParts[len(queryParts)-i] != blockedParts[len(blockedParts)-i] {
            return false
        }
    }
    return true
}

// ServeHTTP handles DoH requests.
func (p *CustomUserPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

    if r.Method == "OPTIONS" {
        return
    }

    switch {
    case r.Method == http.MethodPost && r.URL.Path == "/add_rule":
        p.handleAddRule(w, r)
    case r.Method == http.MethodPost && r.URL.Path == "/remove_rule":
        p.handleRemoveRule(w, r)
    case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/rules/"):
        p.handleGetRules(w, r)
    default:
        // Existing DoH handling code
        userID := r.URL.Path[1:]
        ctx := context.WithValue(r.Context(), "userID", userID)
        msg, err := dohDecode(r)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        dnsWriter := &dohResponseWriter{w: w}
        _, err = p.ServeDNS(ctx, dnsWriter, msg)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
}

// dohDecode decodes a DNS-over-HTTPS request.
func dohDecode(r *http.Request) (*dns.Msg, error) {
    var msg *dns.Msg
    var err error

    switch r.Method {
    case http.MethodGet:
        dnsBase64 := r.URL.Query().Get("dns")
        dnsWire, err := base64.RawURLEncoding.DecodeString(dnsBase64)
        if err != nil {
            return nil, err
        }
        msg = new(dns.Msg)
        err = msg.Unpack(dnsWire)
    case http.MethodPost:
        dnsWire, err := ioutil.ReadAll(r.Body)
        if err != nil {
            return nil, err
        }
        msg = new(dns.Msg)
        err = msg.Unpack(dnsWire)
    default:
        return nil, fmt.Errorf("unsupported HTTP method: %s", r.Method)
    }

    return msg, err
}

// dohEncode encodes a DNS message for a DNS-over-HTTPS response.
func dohEncode(w http.ResponseWriter, msg *dns.Msg) error {
    buf, err := msg.Pack()
    if err != nil {
        return err
    }

    w.Header().Set("Content-Type", "application/dns-message")
    _, err = w.Write(buf)
    return err
}

// dohResponseWriter adapts http.ResponseWriter to dns.ResponseWriter.
type dohResponseWriter struct {
    w   http.ResponseWriter
    msg *dns.Msg
}

func (d *dohResponseWriter) WriteMsg(msg *dns.Msg) error {
    d.msg = msg
    return dohEncode(d.w, msg)
}

func (d *dohResponseWriter) Write(b []byte) (int, error)     { return d.w.Write(b) }
func (d *dohResponseWriter) Close() error                    { return nil }
func (d *dohResponseWriter) TsigStatus() error               { return nil }
func (d *dohResponseWriter) TsigTimersOnly(bool)             {}
func (d *dohResponseWriter) Hijack()                         {}
func (d *dohResponseWriter) LocalAddr() net.Addr             { return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 443} }
func (d *dohResponseWriter) RemoteAddr() net.Addr            { return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0} }
func (d *dohResponseWriter) WriteRaw([]byte) error           { return nil }
func (d *dohResponseWriter) Flush()                          {}


// ttlWriter is a custom dns.ResponseWriter that modifies the TTL of the response
type ttlWriter struct {
    dns.ResponseWriter
    ttl uint32
}

// WriteMsg implements the dns.ResponseWriter interface
func (w *ttlWriter) WriteMsg(res *dns.Msg) error {
    // Modify the TTL of all answer RRs
    for _, rr := range res.Answer {
        rr.Header().Ttl = w.ttl
    }
    return w.ResponseWriter.WriteMsg(res)
}

// Add these methods to your CustomUserPlugin struct

func (p *CustomUserPlugin) RemoveRule(userID string, ruleID int) error {
    _, err := p.DB.Exec("DELETE FROM blocked_sites WHERE user_id = $1 AND id = $2", userID, ruleID)
    if err != nil {
        log.Printf("Error removing rule for user %s: %v", userID, err)
        return err
    }
    log.Printf("Removed rule ID %d for user %s", ruleID, userID)
    return nil
}

// Add this new function to handle rule addition
func (p *CustomUserPlugin) handleAddRule(w http.ResponseWriter, r *http.Request) {
    var req RuleRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := p.AddRule(req.UserID, req.Domain); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "rule added"})
}

// Add this struct for JSON parsing
type RuleRequest struct {
    UserID string `json:"user_id"`
    Domain string `json:"domain"`
}

func (p *CustomUserPlugin) handleRemoveRule(w http.ResponseWriter, r *http.Request) {
    var req struct {
        UserID string `json:"user_id"`
        RuleID int    `json:"rule_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := p.RemoveRule(req.UserID, req.RuleID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "rule removed"})
}

func (p *CustomUserPlugin) handleGetRules(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Path[len("/rules/"):]
    if userID == "" {
        log.Printf("Error: Empty user ID")
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }
    log.Printf("Fetching rules for user: %s", userID)
    
    rules, err := p.GetRules(userID)
    if err != nil {
        log.Printf("Error getting rules for user %s: %v", userID, err)
        http.Error(w, fmt.Sprintf("Error getting rules: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(rules); err != nil {
        log.Printf("Error encoding rules for user %s: %v", userID, err)
        http.Error(w, fmt.Sprintf("Error encoding rules: %v", err), http.StatusInternalServerError)
        return
    }
    
    log.Printf("Successfully fetched %d rules for user %s", len(rules), userID)
}

func (p *CustomUserPlugin) GetRules(userID string) ([]map[string]interface{}, error) {
    log.Printf("Querying database for rules of user: %s", userID)
    
    rows, err := p.DB.Query("SELECT id, domain FROM blocked_sites WHERE user_id = $1", userID)
    if err != nil {
        log.Printf("Database query error for user %s: %v", userID, err)
        return nil, fmt.Errorf("database query error: %v", err)
    }
    defer rows.Close()

    var rules []map[string]interface{}
    for rows.Next() {
        var id int
        var domain string
        if err := rows.Scan(&id, &domain); err != nil {
            log.Printf("Error scanning row for user %s: %v", userID, err)
            return nil, fmt.Errorf("error scanning row: %v", err)
        }
        rules = append(rules, map[string]interface{}{
            "id":     id,
            "domain": domain,
        })
    }
    
    if err := rows.Err(); err != nil {
        log.Printf("Error after scanning all rows for user %s: %v", userID, err)
        return nil, fmt.Errorf("error after scanning all rows: %v", err)
    }
    
    log.Printf("Retrieved %d rules for user %s", len(rules), userID)
    return rules, nil
}

