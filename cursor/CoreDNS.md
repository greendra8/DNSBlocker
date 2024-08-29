# CoreDNS Manual

This document provides a comprehensive overview of CoreDNS, its installation, configuration, and plugin development.

## What is CoreDNS?

CoreDNS is a DNS server written in Go. It distinguishes itself from other DNS servers through its flexibility and plugin-based architecture. Almost all functionality is handled by plugins, which can operate independently or in conjunction to perform DNS functions.

CoreDNS defines a "DNS function" as software that adheres to the CoreDNS Plugin API. This functionality can vary widely, encompassing plugins that augment functionality without directly generating responses (e.g., metrics, cache) and those that do generate responses. These response-generating plugins can interact with various sources like Kubernetes for service discovery, files, or databases.

The default CoreDNS installation includes around 30 plugins, with numerous external plugins available for compilation to expand its capabilities.

## Installation

### Binaries

Pre-compiled binaries for various operating systems are provided for each CoreDNS release.

### Docker

Every release is also pushed as Docker images to the public Docker Hub for the CoreDNS organization. These images typically consist of a minimal base (scratch), CoreDNS, and TLS certificates for DoT, DoH, and gRPC.

### Source

To compile CoreDNS from source, a working Go setup is required. CoreDNS utilizes Go modules for dependency management. Detailed compilation instructions are maintained within the CoreDNS source code.

### Testing

After obtaining a `coredns` binary, the `-plugins` flag lists all compiled plugins. Without a Corefile (see Configuration), CoreDNS loads the `whoami` plugin, which responds with the client's IP address and port. A simple test involves starting CoreDNS on port 1053 and using `dig` to send a query:

bash $ ./coredns -dns.port=1053
In another terminal:

bash $ dig @localhost -p 1053 a whoami.example.org
## Plugins

CoreDNS operates by running Servers, each defined by the zones it serves and its listening port. Each Server possesses its own Plugin Chain.

When processing a query, CoreDNS follows these steps:

1. **Server Selection:** If multiple Servers listen on the queried port, the Server with the most specific zone (longest suffix match) is chosen.

2. **Plugin Chain Execution:** The query is routed through the Server's configured Plugin Chain in a static order defined in `plugin.cfg`.

3. **Plugin Processing:** Each plugin examines the query and decides whether to process it. Possible outcomes include:
    - **Query Processed:** The plugin generates a response and sends it to the client.
    - **Query Not Processed:** The plugin passes the query to the next plugin in the chain.
    - **Query Processed with Fallthrough:** The plugin processes the query but allows subsequent plugins to also handle it.
    - **Query Processed with a Hint:** The plugin processes the query and passes it on, providing a hint for later inspection.

### Unregistered Plugins

Certain plugins don't handle DNS data directly but influence CoreDNS behavior. Examples include:

- `bind`: Controls interface binding.
- `root`: Sets the root directory for plugin file access.
- `health`: Enables an HTTP health check endpoint.
- `ready`: Supports readiness reporting.

### Anatomy of Plugins

Plugins typically consist of Setup, Registration, and Handler components.

- **Setup:** Parses configuration and plugin-specific directives.
- **Handler:** Contains the query processing logic.
- **Registration:** Registers the plugin within CoreDNS during compilation.

### Plugin Documentation

Each plugin has a README file detailing its configuration, examples, and usage considerations. These READMEs are available at [https://coredns.io/plugins](https://coredns.io/plugins) and are also compiled into manual pages.

## Configuration

CoreDNS configuration involves selecting plugins for compilation (for custom builds) and using the Corefile to define server behavior.

### Corefile

The Corefile is the primary configuration file for CoreDNS. It consists of one or more Server Blocks, each specifying zones and a chain of plugins.

**Server Blocks:**

[zone...]:[port] { # Plugins defined here }
- `zone...`: One or more zones the server is authoritative for.
- `port`: Optional port number (defaults to 53).

**Plugins:**

Plugins are listed within Server Blocks. Simple plugins can be added by name:

. { chaos }
Plugins with more configuration options use Plugin Blocks:

. { plugin { # Plugin Block } }
**Example Corefile:**

coredns.io:5300 { file db.coredns.io }

example.io:53 { log errors file db.example.io }

example.net:53 { file db.example.net }

.:53 { kubernetes forward . 8.8.8.8 log errors cache }
### Environment Variables

CoreDNS supports environment variable substitution in the Corefile using the syntax `{$ENV_VAR}` or `{%ENV_VAR%}`.

### Importing Other Files

The `import` plugin allows including configuration from other files.

### Reusable Snippets

Snippets are named blocks of configuration that can be reused with the `import` plugin:

(snip) { prometheus log errors }

. { whoami import snip }
## Specific Setups

This section provides examples of common CoreDNS configurations.

### Authoritative Serving From Files

This setup uses the `file` plugin to serve zone data from a file:

**db.example.org:**

$ORIGIN example.org. @ 3600 IN SOA sns.dns.icann.org. noc.dns.icann.org. ( 2017042745 ; serial 7200 ; refresh (2 hours) 3600 ; retry (1 hour) 1209600 ; expire (2 weeks) 3600 ; minimum (1 hour) )

    3600 IN NS a.iana-servers.net.
    3600 IN NS b.iana-servers.net.

www IN A 127.0.0.1 IN AAAA ::1
**Corefile:**

example.org { file db.example.org log }
### Forwarding

The `forward` plugin forwards queries to a recursor:

**Corefile:**

. { forward . 8.8.8.8 9.9.9.9 log }
### Recursive Resolver

The `unbound` plugin (external) provides recursive resolution functionality:

**Corefile:**

. { unbound cache log }
## Writing Plugins

Plugins are the core of CoreDNS's functionality. Developing a plugin involves creating several Go files:

- `setup.go`: Handles configuration parsing.
- `<plugin_name>.go`: Contains the query processing logic.
- `README.md`: Documents the plugin's usage and configuration.

The `example` plugin serves as a minimal example for plugin development.

### How Plugins Are Called

CoreDNS invokes a plugin's `ServeDNS` method, which receives a `context.Context`, a `dns.ResponseWriter`, and a `*dns.Msg` representing the client request. The method returns a response code and an error.

### Logging From a Plugin

The `log` package provides logging capabilities within a plugin.

### Metrics

Plugins should export metrics with the namespace `plugin.Namespace` (`coredns`) and the subsystem set to the plugin's name.

### Documentation

Plugin documentation should follow a Unix manual page style and include sections for Name, Description, Syntax, Examples, and optionally See Also and Bugs.

### Style

Adhere to the Unix manual page style guidelines for consistency.

### Example Domain Names

Use `example.org` or `example.net` in examples and tests.

### Fallthrough

The `fallthrough` directive allows a plugin to handle only a subset of names within a zone and pass others to the next plugin in the chain.