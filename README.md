# Kurs Pajak

This repository exists to scrape all kurs posted by Kemenkeu every Wednesday. This kurs is used by beacukai to calculate tax. You can use this API for free.

## Live API

**BASE URL:** `https://api.akhmad.id/pajak`

### GET /kurs

#### Request:
- Parameter: None
- Header:
  - Authorization: None
- Body: None

#### Response:
- Header:
  - `Content-Type: Application/json`
- Body:
  ```json
  {
    "response_code": 200,
    "message": "Success",
    "data": {
      "currencies": [
        {
          "changes": "-32.00",
          "currency": "Dolar Amerika Serikat (USD)",
          "symbol": "USD",
          "value": "14309.00"
        },
        {
          "changes": "31.36",
          "currency": "Dolar Australia (AUD)",
          "symbol": "AUD",
          "value": "10346.84"
        },
      ],
      "updated_at": 1642850483,
      "valid_from": "19 Januari 2022",
      "valid_to": "25 Januari 2022"
    }
  }
  ```
  
## Roadmap
- [ ] Add a way to query kurs by date (just for history data i think?)
- [ ] Save scrapped data to db somewhere
- [x] Add authorization and rate limiter to prevent abuser (done by using API Gateway)
