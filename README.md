# DNS Manager

DNS Manager is a custom DNS management system built on top of CoreDNS. It allows users to manage their own DNS blocking rules through a simple web interface.

## Features

- User-specific DNS blocking rules
- Web interface for managing blocking rules
- DNS-over-HTTPS (DoH) support
- Subdomain blocking
- Low-latency updates to blocking rules

## Components

1. **CoreDNS Server**: A customized CoreDNS server with our custom plugin.
2. **Custom CoreDNS Plugin**: Handles user-specific DNS rules and DoH requests.
3. **Web Frontend**: A SvelteKit application for managing blocking rules.
4. **PostgreSQL Database**: Stores user information and blocking rules.

## Setup

### Prerequisites

- Go 1.16 or later
- Node.js and npm
- PostgreSQL
- CoreDNS source code

### Installation

1. Set up the PostgreSQL database (see `setup_and_run.md` for details).
2. Build the custom CoreDNS server with our plugin.
3. Set up and run the SvelteKit frontend application.

For detailed setup instructions, please refer to `setup_and_run.md`.

## Usage

1. Users set their DNS server to the address of your CoreDNS server.
2. Users access the web interface to manage their blocking rules.
3. The CoreDNS server applies user-specific rules to DNS queries.

## Development

- The custom CoreDNS plugin is located in `plugin/custom_user_plugin/`.
- The SvelteKit frontend is in the `dns-manager-frontend/` directory.

## Known Issues

- DNS caching can delay the effect of newly added blocking rules.

## Future Improvements

- Implement more aggressive caching mitigation strategies.
- Add user authentication for the web interface.
- Expand blocking capabilities (e.g., time-based rules, category-based blocking).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.
