# Currency Exchange Rate Parser
This app diligently fetches currency exchange rates from [mig.kz](https://mig.kz), saves them in PostgreSQL and provides REST API endpoints to retrieve the latest rates and their averages over a specified period, so you don’t have to Google it every time.  
---

## Features

- **Web Scraping**: Fetches currency exchange rates (buy and sell) from mig.kz.
- **Database Storage**: Stores parsed data in a PostgreSQL database.
- **REST API**: Provides endpoints to retrieve the latest rates and average rates over a specified period.
- **Cron Job**: Periodically updates the currency data (runs every hour).

---
⚠️ Attention!
No AI was harmed in the making of this project
## Prerequisites

Before running the project, ensure you have the following installed:

- [Go](https://golang.org/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Git](https://git-scm.com/downloads)

---

## API Endpoints

### 1. Get Latest Exchange Rates

**Endpoint:** `GET /currencies/latest`

**Description:** Retrieves the latest exchange rates for all currencies.

#### Example Request:

```bash
curl http://localhost:8080/currencies/latest
```

#### Example Response:

```json
[
    {
        "id": 1,
        "currency_code": "USD",
        "buy_rate": 493.1,
        "sell_rate": 496.5,
        "timestamp": "2025-03-09T13:42:00Z"
    },
    {
        "id": 2,
        "currency_code": "EUR",
        "buy_rate": 532,
        "sell_rate": 538,
        "timestamp": "2025-03-09T13:42:00Z"
    },
    {
        "id": 3,
        "currency_code": "RUB",
        "buy_rate": 5.45,
        "sell_rate": 5.59,
        "timestamp": "2025-03-09T13:42:00Z"
    },
    {
        "id": 4,
        "currency_code": "KGS",
        "buy_rate": 5.49,
        "sell_rate": 6.09,
        "timestamp": "2025-03-09T13:42:00Z"
    },
    {
        "id": 5,
        "currency_code": "GBP",
        "buy_rate": 635,
        "sell_rate": 655,
        "timestamp": "2025-03-09T13:42:00Z"
    },
    {
        "id": 6,
        "currency_code": "CNY",
        "buy_rate": 69.1,
        "sell_rate": 72.1,
        "timestamp": "2025-03-09T13:42:00Z"
    },
    {
        "id": 7,
        "currency_code": "GOLD",
        "buy_rate": 45750,
        "sell_rate": 48450,
        "timestamp": "2025-03-09T13:42:00Z"
    },
]
```

### 2. Get Average Exchange Rates

**Endpoint:** `GET /currencies/average`

**Description:** Retrieves the average buy and sell rates for all currencies over a specified date range.

#### Query Parameters:

- `start_date`: Start date in `DD-MM-YYYY` format.
- `end_date`: End date in `DD-MM-YYYY` format.

#### Example Request:

```bash
curl "http://localhost:8080/currencies/average?start_date=09-03-2025&end_date=10-03-2025"
```

#### Example Response:

```json
{
    "averages": [
        {
            "currency_code": "RUB",
            "average_buy": 5.45,
            "average_sell": 5.59
        },
        {
            "currency_code": "CNY",
            "average_buy": 69.1,
            "average_sell": 72.1
        },
        {
            "currency_code": "EUR",
            "average_buy": 532,
            "average_sell": 538
        },
        {
            "currency_code": "USD",
            "average_buy": 493.1,
            "average_sell": 496.5
        },
        {
            "currency_code": "GOLD",
            "average_buy": 45750,
            "average_sell": 48450
        },
        {
            "currency_code": "GBP",
            "average_buy": 635,
            "average_sell": 655
        },
        {
            "currency_code": "KGS",
            "average_buy": 5.49,
            "average_sell": 6.09
        }
    ],
    "end_date": "09-03-2025",
    "start_date": "01-01-2025"
}
