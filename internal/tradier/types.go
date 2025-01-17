package tradier

type ErrorMessage struct {
	Err error
}

type QuoteUpdate struct {
	Quote Quote
}

type SummaryUpdate struct {
	Summary Summary
}

type Quote struct {
	Type    string  `json:"type"`
	Symbol  string  `json:"symbol"`
	Bid     float64 `json:"bid"`
	BidSize int     `json:"bidsz"`
	BidExch string  `json:"bidexch"`
	BidDate string  `json:"biddate"`
	Ask     float64 `json:"ask"`
	AskSize int     `json:"asksz"`
	AskExch string  `json:"askexch"`
	AskDate string  `json:"askdate"`
}

type Trade struct {
	Type   string `json:"type"`
	Symbol string `json:"symbol"`
	Exch   string `json:"exch"`
	Price  string `json:"price"`
	Size   string `json:"size"`
	Cvol   string `json:"cvol"`
	Date   string `json:"date"`
	Last   string `json:"last"`
}

type Summary struct {
	Type      string `json:"type"`
	Symbol    string `json:"symbol"`
	Open      string `json:"open"`
	High      string `json:"high"`
	Low       string `json:"low"`
	PrevClose string `json:"prevClose"`
}

type Timesale struct {
	Type       string `json:"type"`
	Symbol     string `json:"symbol"`
	Exch       string `json:"exch"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Last       string `json:"last"`
	Size       string `json:"size"`
	Date       string `json:"date"`
	Seq        int    `json:"seq"`
	Flag       string `json:"flag"`
	Cancel     bool   `json:"cancel"`
	Correction bool   `json:"correction"`
	Session    string `json:"session"`
}

type Tradex struct {
	Type   string  `json:"type"`
	Symbol string  `json:"symbol"`
	Exch   string  `json:"exch"`
	Price  float64 `json:"price"`
	Size   int     `json:"size"`
	Cvol   int     `json:"cvol"`
	Date   int64   `json:"date"`
	Last   float64 `json:"last"`
}
