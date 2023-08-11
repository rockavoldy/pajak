# Kurs Pajak

This repository exists to scrape all kurs posted by Kemenkeu every Wednesday. This kurs is used by beacukai to calculate tax.

## Run the app
Using docker, you can build and run with
```sh
docker build -t pajak:latest .
docker run -p 8080:8080 -t pajak:latest
```

## Roadmap
- [ ] Add a way to query kurs by date (just for history data i think?)
- [ ] Save scrapped data to db somewhere
- [x] Add authorization and rate limiter to prevent abuser (done by using API Gateway)
