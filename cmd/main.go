package main

import (
	"flag"
	"fmt"
	"github.com/alextreichler/stonks/internal/tradier"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os"
)

// Market Struct stuff

type config struct {
	apiKey string
}

type marketEvents struct {
	symbol string
}

// Bubble tea stuff

type model struct {
	quote   tradier.Quote
	summary tradier.Summary
	client  *tradier.Client
	err     error
}

func main() {

	var cfg config
	var me marketEvents

	flag.StringVar(&me.symbol, "symbol", "", "Stock symbol to monitor")
	flag.StringVar(&cfg.apiKey, "api-key", "", "Tradier API key")

	flag.Parse()
	if me.symbol == "" {
		fmt.Println("Please provide a stock symbol using -symbol flag")
		flag.Usage()
		os.Exit(1)
	}

	if cfg.apiKey == "" {
		cfg.apiKey = os.Getenv("TRADIER_API_KEY")
		if cfg.apiKey == "" {
			fmt.Println("Please provide a Tradier API key using -api-key flag or TRADIER_API_KEY environment variable")
			os.Exit(1)
		}
	}

	client := tradier.NewClient(cfg.apiKey)
	p := tea.NewProgram(initialModel(client, me))
	client.Program = p // Store the program in the client

	if _, err := p.Run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func initialModel(client *tradier.Client, me marketEvents) model {
	return model{
		quote:  tradier.Quote{Symbol: me.symbol},
		client: client,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		func() tea.Msg {
			go func() {
				err := m.client.StreamQuotes(m.quote.Symbol)
				if err != nil {
					m.client.Program.Send(tradier.ErrorMessage{err})
				}
			}()
			return nil
		},
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tradier.QuoteUpdate:
		m.quote = msg.Quote
		return m, nil
	case tradier.SummaryUpdate:
		m.summary = msg.Summary
		return m, nil
	case tradier.ErrorMessage: // Handle potential errors from StreamQuotes
		m.err = msg.Err
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}
	if m.quote.Symbol == "" { // Handle initial state
		return "Loading..."
	}

	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")). // Example color
		Padding(1, 2)                           // Add some padding

	return style.Render(fmt.Sprintf(
		"Symbol: %s\nBid: %.2f (Size: %d)\nAsk: %.2f (Size: %d)\n",
		m.quote.Symbol, m.quote.Bid, m.quote.BidSize, m.quote.Ask, m.quote.AskSize,
	))
}
