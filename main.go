package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/ElrohirGT/Redes_RemoteMCPServer/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	var transport string
	flag.StringVar(&transport, "t", "http", "Transport type (stdio or http)")
	var addr string
	flag.StringVar(&addr, "a", ":8080", "Address where the server will listen")
	flag.Parse()

	// Create a new MCP server
	s := server.NewMCPServer(
		"MTG Cards MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add tool
	get_cards_tool := mcp.NewTool("get_cards",
		mcp.WithDescription("Obtains a list of cards from the Magic The Gathering (MTG) database. For each card it gets: id, name, manacost, colors, type, rarity, text and image url."),
	)
	get_card_tool := mcp.NewTool("get_card",
		mcp.WithDescription("Obtains information about a specific card from Magic The Gathering (MTG). Specifically: id, name, manacost, colors, type, rarity, text, image url and text description."),
		mcp.WithString("id", mcp.Description("The id of the card to search for more information.")),
	)
	get_game_formats_tool := mcp.NewTool("get_formats",
		mcp.WithDescription("Get's all the available formats to play Magic the Gathering (MTG)"),
	)

	// Add tool handler
	s.AddTool(get_cards_tool, get_cards)
	s.AddTool(get_card_tool, get_card)
	s.AddTool(get_game_formats_tool, get_game_formats)

	if transport == "http" {
		serverCtx, cancelServerCtx := context.WithCancel(context.Background())
		defer cancelServerCtx()
		serv := server.NewStreamableHTTPServer(
			s,
			server.WithStateLess(true),
		)
		log.Println("HTTP server listening on", addr)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := serv.Start(addr); err != nil {
				log.Printf("Server error: %v", err)
			}
			log.Println("Server execution ended!")
		}()

		var stopChan = make(chan os.Signal, 2)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		<-stopChan // wait for SIGINT
		log.Println("Shutting down server...")
		err := serv.Shutdown(serverCtx)
		if err != nil {
			log.Panic(err)
		}
		wg.Wait()
		log.Println("Server shutdown!")
	} else {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}

func get_card(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request to `get_card`\n%#v", request)
	id, err := request.RequireString("id")
	if err != nil {
		log.Println("No ID provided!")
		return mcp.NewToolResultError(err.Error()), err
	}

	result, err, shouldTerminate := tools.GetCardCore(ctx, id)
	if err != nil {
		if shouldTerminate {
			return mcp.NewToolResultError(err.Error()), err
		} else {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	log.Println("Transforming request from JSON into single string...")
	b := strings.Builder{}
	cardToStringBuilder(&b, result.Card)
	resultText := b.String()
	log.Println("Final output:\n", resultText)

	return mcp.NewToolResultText(resultText), nil
}

func cardToStringBuilder(b *strings.Builder, card tools.MTGCard) {
	b.WriteString(card.Id)
	b.WriteString("\n - Name: ")
	b.WriteString(card.Name)
	b.WriteString("\n - Mana Cost: ")
	b.WriteString(card.ManaCost)
	b.WriteString("\n - Colors: ")
	for _, c := range card.Colors {
		b.WriteString(c)
		b.WriteString(",")
	}
	b.WriteString("\n - Rarity: ")
	b.WriteString(card.Rarity)
	b.WriteString("\n - Type: ")
	b.WriteString(card.Type)

	if card.Text != "" {
		b.WriteString("\n - Text: ")
		b.WriteString(card.Text)
	}
}

func get_cards(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request to `get_cards`\n%#v", request)

	log.Println("Trying to get cards from MTG...")
	result, err, shouldTerminate := tools.GetCardsCore(ctx)
	if err != nil {
		if shouldTerminate {
			return mcp.NewToolResultError(err.Error()), err
		} else {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	log.Println("Transforming request from JSON into single string...")
	b := strings.Builder{}
	for _, card := range limit(result.Cards, 10) {
		card.Text = ""
		cardToStringBuilder(&b, card)
		b.WriteRune('\n')
	}

	resultText := b.String()
	log.Println("Final output:\n", resultText)

	return mcp.NewToolResultText(resultText), nil
}

func limit[T any](arr []T, limit int) []T {
	if len(arr) > limit {
		return arr[:limit]
	}
	return arr
}

func get_game_formats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request to `get_game_formats`\n%#v", request)

	log.Println("Trying to get game formats from MTG...")
	result, err, shouldTerminate := tools.GetGameFormats(ctx)
	if err != nil {
		if shouldTerminate {
			return mcp.NewToolResultError(err.Error()), err
		} else {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	log.Println("Transforming request from JSON into single string...")
	b := strings.Builder{}
	for _, format := range result.Formats {
		b.WriteString(format)
		b.WriteRune('\n')
	}

	resultText := b.String()
	log.Println("Final output:\n", resultText)

	return mcp.NewToolResultText(resultText), nil
}
