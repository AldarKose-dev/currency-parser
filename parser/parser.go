package parser

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mig_parser/models"

	"github.com/PuerkitoBio/goquery"
)

func ParseCurrencies(db *sql.DB) {
	log.Println("Starting currency parsing...")

	// Make HTTP request to mig.kz
	resp, err := http.Get("https://mig.kz")
	if err != nil {
		log.Printf("Failed to fetch webpage: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		log.Printf("Bad response status: %d", resp.StatusCode)
		return
	}

	// Parse HTML using goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Failed to parse HTML: %v", err)
		return
	}

	// Extract timestamp from the page
	var timestamp time.Time
	doc.Find(".informer h4.text-center").Each(func(i int, s *goquery.Selection) {
		dateText := strings.TrimSpace(s.Text())
		if strings.Contains(dateText, "на") {
			dateStr := strings.SplitN(dateText, "на", 2)[1]
			dateStr = strings.TrimSpace(dateStr)
			dateStr = replaceRussianMonths(dateStr)
			parsedTime, err := time.Parse("02 January 2006 15:04", dateStr)
			if err != nil {
				log.Printf("Error parsing date: %v", err)
				return
			}
			timestamp = parsedTime
		}
	})

	// Fallback to current time if timestamp not parsed
	if timestamp.IsZero() {
		timestamp = time.Now().UTC()
		log.Println("Using current time as timestamp")
	}

	// Find currency rows
	currencies := make([]models.Currency, 0)
	doc.Find(".informer table tr").Each(func(i int, s *goquery.Selection) {
		// Skip rows without currency code
		currencyCode := s.Find("td.currency").Text()
		if currencyCode == "" {
			return
		}

		// Extract buy and sell rates
		buyStr := s.Find("td.buy").Text()
		sellStr := s.Find("td.sell").Text()

		// Parse rates
		buyRate, err := strconv.ParseFloat(buyStr, 64)
		if err != nil {
			log.Printf("Error parsing buy rate for %s: %v", currencyCode, err)
			return
		}

		sellRate, err := strconv.ParseFloat(sellStr, 64)
		if err != nil {
			log.Printf("Error parsing sell rate for %s: %v", currencyCode, err)
			return
		}

		currencies = append(currencies, models.Currency{
			CurrencyCode: strings.TrimSpace(currencyCode),
			BuyRate:      buyRate,
			SellRate:     sellRate,
			Timestamp:    timestamp,
		})
	})

	// Save to database
	for _, currency := range currencies {
		_, err := db.Exec(
			"INSERT INTO currencies (currency_code, buy_rate, sell_rate, timestamp) VALUES ($1, $2, $3, $4)",
			currency.CurrencyCode, currency.BuyRate, currency.SellRate, currency.Timestamp,
		)
		if err != nil {
			log.Printf("Failed to save currency %s: %v", currency.CurrencyCode, err)
		}
	}

	log.Printf("Parsed and saved %d currencies", len(currencies))
}

func replaceRussianMonths(dateStr string) string {
	months := map[string]string{
		"января":   "January",
		"февраля":  "February",
		"марта":    "March",
		"апреля":   "April",
		"мая":      "May",
		"июня":     "June",
		"июля":     "July",
		"августа":  "August",
		"сентября": "September",
		"октября":  "October",
		"ноября":   "November",
		"декабря":  "December",
	}

	for ru, en := range months {
		dateStr = strings.Replace(dateStr, ru, en, 1)
	}
	return dateStr
}
