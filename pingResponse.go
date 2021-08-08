package main

import (
	"github.com/Tnze/go-mc/chat"
	_ "github.com/Tnze/go-mc/data/lang/en-us"
	"github.com/google/uuid"
)

type Status struct {
	Description chat.Message

	Players struct {
		Max    int
		Online int
		Sample []struct {
			ID   uuid.UUID
			Name string
		}
	}
	Version struct {
		Name     string
		Protocol int
	}
}

type IpData struct {
	IP                string  `json:"ip"`
	Success           bool    `json:"success"`
	Type              string  `json:"type"`
	Continent         string  `json:"continent"`
	ContinentCode     string  `json:"continent_code"`
	Country           string  `json:"country"`
	CountryCode       string  `json:"country_code"`
	CountryFlag       string  `json:"country_flag"`
	CountryCapital    string  `json:"country_capital"`
	CountryPhone      string  `json:"country_phone"`
	CountryNeighbours string  `json:"country_neighbours"`
	Region            string  `json:"region"`
	City              string  `json:"city"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Asn               string  `json:"asn"`
	Org               string  `json:"org"`
	Isp               string  `json:"isp"`
	Timezone          string  `json:"timezone"`
	TimezoneName      string  `json:"timezone_name"`
	TimezoneDstOffset int     `json:"timezone_dstOffset"`
	TimezoneGmtOffset int     `json:"timezone_gmtOffset"`
	TimezoneGmt       string  `json:"timezone_gmt"`
	Currency          string  `json:"currency"`
	CurrencyCode      string  `json:"currency_code"`
	CurrencySymbol    string  `json:"currency_symbol"`
	CurrencyRates     float64 `json:"currency_rates"`
	CurrencyPlural    string  `json:"currency_plural"`
	CompletedRequests int     `json:"completed_requests"`
}
