# Kurs Pajak

This repository exists to scrape all kurs posted by Kemenkeu every Wednesday, and serve the API in this [Kalkulator pajak](https://pajak.akmd.dev) to find a rough estimate how much import tax you need to pay for your items.

## Run the app
Using docker, you can build and run with
```sh
docker build -t pajak:latest .
docker run -p 8080:8080 -t pajak:latest
```

## Roadmap
- [x] Add frontend to the root endpoint
- [x] ~~Add CORS~~ Not needed; since the API and frontend is on the same project