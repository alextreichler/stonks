package tradier

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func (c *Client) StreamQuotes(symbol string) error {
	// Get streaming session
	sessionID, err := c.createStreamingSession()
	if err != nil {
		return fmt.Errorf("error creating streaming session: %w", err)
	}
	//log.Printf("Session ID: %s", sessionID)

	// Create websocket connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	wsURL := "wss://ws.tradier.com/v1/markets/events"
	conn, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
		HTTPHeader: http.Header{
			"Authorization": []string{"Bearer " + c.apiKey},
			"SessionId":     []string{sessionID},
		},
	})
	if err != nil {
		return fmt.Errorf("error connecting to websocket: %w", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing connection")

	// Subscribe to quotes
	subscribeMsg := map[string]interface{}{
		"symbols":   []string{symbol},
		"type":      "quotes",
		"sessionid": sessionID,
	}

	if err := wsjson.Write(ctx, conn, subscribeMsg); err != nil {
		return fmt.Errorf("error subscribing to quotes: %w", err)
	}

	for {
		var data map[string]interface{}
		err := wsjson.Read(ctx, conn, &data)
		if err != nil {
			return fmt.Errorf("error reading message: %w", err)
		}

		if eventType, ok := data["type"].(string); ok {
			switch eventType {
			case "quote":
				var quote Quote
				quoteJSON, err := json.Marshal(data) // Convert map to JSON
				if err != nil {
					log.Printf("Error marshalling quote: %v", err)
					continue
				}
				if err := json.Unmarshal(quoteJSON, &quote); err != nil {
					log.Printf("Error unmarshalling quote: %v", err)
					continue
				}
				c.Program.Send(QuoteUpdate{Quote: quote})

				//log.Printf("Quote: %+v", quote)
				//				fmt.Printf("\rQuote: %s --- %v:%v", quote.Symbol, quote.Ask, quote.Bid)

			case "trade":
				var trade Trade
				tradeJSON, err := json.Marshal(data)
				if err != nil {
					log.Printf("Error marshalling trade: %v", err)
					continue
				}
				if err := json.Unmarshal(tradeJSON, &trade); err != nil {
					log.Printf("Error unmarshalling trade: %v", err)
					continue
				}
				//				fmt.Printf("\rTrade: %+v", trade)

			case "timesale":
				var timesale Timesale
				timesaleJSON, err := json.Marshal(data)
				if err != nil {
					log.Printf("Error marshalling timesale: %v", err)
					continue
				}
				if err := json.Unmarshal(timesaleJSON, &timesale); err != nil {
					log.Printf("Error unmarshalling timesale: %v", err)
					continue
				}
				//				fmt.Printf("srade: %+v", timesale)

			case "summary":
				var summary Summary
				summaryJSON, err := json.Marshal(data)
				if err != nil {
					log.Printf("Error unmarshalling summary: %v", err)
					continue
				}
				if err := json.Unmarshal(summaryJSON, &summary); err != nil {
					log.Printf("Error unmarshalling summary: %v", err)
				}
				//				fmt.Printf("\rsummary: %v+", summary)
				c.Program.Send(SummaryUpdate{Summary: summary})

			// ... handle other event types (summary, timesale, tradex)
			default:
				log.Printf("Unknown event type: %s", eventType)
			}
		} else {
			log.Printf("Invalid message format: %+v", data)
		}
	}
}
