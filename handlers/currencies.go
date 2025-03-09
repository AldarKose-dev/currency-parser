package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"mig_parser/models"

)

func GetLatestCurrencies(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`
			WITH latest_timestamps AS (
				SELECT currency_code, MAX(timestamp) as max_timestamp
				FROM currencies
				GROUP BY currency_code
			)
			SELECT c.id, c.currency_code, c.buy_rate, c.sell_rate, c.timestamp
			FROM currencies c
			JOIN latest_timestamps lt 
				ON c.currency_code = lt.currency_code AND c.timestamp = lt.max_timestamp
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve currencies"})
			return
		}
		defer rows.Close()

		var currencies []models.Currency
		for rows.Next() {
			var currency models.Currency
			if err := rows.Scan(&currency.ID, &currency.CurrencyCode, &currency.BuyRate, &currency.SellRate, &currency.Timestamp); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan currencies"})
				return
			}
			currencies = append(currencies, currency)
		}

		c.JSON(http.StatusOK, currencies)
	}
}

func GetAverageCurrencies(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		// Validate dates
		if startDate == "" || endDate == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Both start_date and end_date are required"})
			return
		}

		// Parse dates in DD-MM-YYYY format
		start, err := time.Parse("02-01-2006", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use DD-MM-YYYY"})
			return
		}

		end, err := time.Parse("02-01-2006", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use DD-MM-YYYY"})
			return
		}

		// Add one day to end date to include the entire end day
		end = end.AddDate(0, 0, 1)

		// Query to get average rates
		rows, err := db.Query(`
			SELECT 
				currency_code, 
				AVG(buy_rate) as avg_buy,
				AVG(sell_rate) as avg_sell
			FROM currencies
			WHERE timestamp >= $1 AND timestamp < $2
			GROUP BY currency_code
		`, start, end)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve average rates"})
			return
		}
		defer rows.Close()

		var averages []models.AverageCurrency
		for rows.Next() {
			var avg models.AverageCurrency
			if err := rows.Scan(&avg.CurrencyCode, &avg.AverageBuy, &avg.AverageSell); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan average rates"})
				return
			}
			averages = append(averages, avg)
		}

		c.JSON(http.StatusOK, gin.H{
			"start_date": startDate,
			"end_date":   endDate,
			"averages":   averages,
		})
	}
}