# Nix Package MCP Server

## How to use?

The MCP server was built using the Go programming language!

**ENDPOINT:** The MCP server default endpoint for communication is `/mcp`.

If you have Nix installed running this server is as simple as:

```bash
# If you're inside the root of the repo
nix run
```

```bash
# If you don't have the repo cloned
nix run github:ElrohirGT/Redes_MCPServer
```

You can also just compile a binary using normal go commands:

```bash
# If you want to just build a binary
go build .
```

```bash
# If you want to just run the program
go run . 
```

This MCP server supports being run with HTTP or with STDIN/STDOUT, please run
the following command to check the available options:

```bash
go run . -h
```

## How to change the server address

By default the server runs on: `0.0.0.0:8080` but you can change this using the
following flag:

```bash
go run . -a <custom host and port>
```

Again you can check all available configuration options by using:

```bash
go run . -h
```

## Available Tools

1. **get_card** - Get detailed information about a specific MTG card by ID
   (name, mana cost, colors, type, rarity, text, image URL).
1. **get_cards** - Get a list of cards from the MTG database.
1. **get_formats** - Get all available MTG formats to play.
