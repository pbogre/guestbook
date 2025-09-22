package main

import (
	"flag"
        "os"
        "strconv"
)

type Config struct {
    Title               string
    UseRateLimit        bool
    RateLimit           float64
    BurstLimit          int
    EntriesPerPage      int
    Port                int
}

var GuestbookConfig Config

func loadArguments() {
    title := flag.String("title", getenvDefault("GB_TITLE", "Guestbook"), "Title displayed at the top of the webpage")
    useRateLimit := flag.Bool("use-rate-limit", mustParseBool(getenvDefault("GB_USE_RATELIMIT", "true")), "Whether or not to use ratelimiting")
    rateLimit := flag.Float64("rate-limit", mustParseFloat(getenvDefault("GB_RATELIMIT", "0.2")), "Rate limit of requests per second")
    burstLimit := flag.Int("burst-limit", mustParseInt(getenvDefault("GB_BURSTLIMIT", "2")), "Maximum permitted burst of requests")
    entriesPerPage := flag.Int("entries-per-page", mustParseInt(getenvDefault("GB_ENTRIES_PER_PAGE", "10")), "Number of entries displayed per page")
    port := flag.Int("port", mustParseInt(getenvDefault("PORT", "8080")), "Port number to run Guestbook on")

    flag.Parse()

    GuestbookConfig.Title = *title
    GuestbookConfig.UseRateLimit = *useRateLimit
    GuestbookConfig.RateLimit = *rateLimit
    GuestbookConfig.BurstLimit = *burstLimit
    GuestbookConfig.EntriesPerPage = *entriesPerPage
    GuestbookConfig.Port = *port
}

func getenvDefault(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}

func mustParseInt(s string) int {
    v, _ := strconv.Atoi(s)
    return v
}

func mustParseFloat(s string) float64 {
    v, _ := strconv.ParseFloat(s, 64)
    return v
}

func mustParseBool(s string) bool {
    v, _ := strconv.ParseBool(s)
    return v
}
