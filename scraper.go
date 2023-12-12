package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/labstack/echo/v4"
)

type CastItem struct {
	Actor string   `json:"actor"`
	Roles []string `json:"role"`
}

type MovieData struct {
	Title       string     `json:"title"`
	Rate        float64    `json:"rate"`
	PosterImage string     `json:"posterImage"`
	Duration    string     `json:"duration"`
	Year        string     `json:"year"`
	Genres      []string   `json:"genres"`
	Rank        int        `json:"rank"`
	Directors   []string   `json:"director"`
	Writers     []string   `json:"writers"`
	Cast        []CastItem `json:"cast"`
}

func GetMovieData(h echo.Context) error {

	slug := h.Param("slug")

	url := fmt.Sprintf("https://www.imdb.com/title/%v/", slug)

	response := MovieData{}

	// initialize the Collector
	c := colly.NewCollector()

	// set a valid User-Agent header
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"

	c.OnHTML(".jqlHBQ h1 span", func(e *colly.HTMLElement) {
		response.Title = e.Text
	})

	c.OnHTML(".cMEQkK", func(e *colly.HTMLElement) {
		var rate float64
		var err error
		rate, err = strconv.ParseFloat(e.Text, 64)
		if err != nil {
			rate = 0
		}
		response.Rate = rate
	})

	c.OnHTML(".dEqUUl div.ipc-poster img.ipc-image", func(e *colly.HTMLElement) {
		response.PosterImage = e.Attr("src")
	})

	c.OnHTML(".jqlHBQ ul li:first-child", func(e *colly.HTMLElement) {
		response.Year = e.Text
	})

	c.OnHTML(".jqlHBQ ul li:last-child", func(e *colly.HTMLElement) {
		response.Duration = e.Text
	})

	c.OnHTML(".ktjuZl .ipc-chip-list__scroller a span", func(e *colly.HTMLElement) {
		response.Genres = append(response.Genres, e.Text)
	})

	c.OnHTML(".eWQwwe .fTREEx", func(e *colly.HTMLElement) {
		var rank int
		var err error
		rank, err = strconv.Atoi(e.Text)
		if err != nil {
			rank = 0
		}
		response.Rank = rank
	})

	c.OnHTML(".bHYmJY>li:first-child>div>ul>li>a", func(h *colly.HTMLElement) {
		response.Directors = append(response.Directors, h.Text)
	})

	c.OnHTML(".bHYmJY>li:nth-child(2)>div>ul>li>a", func(h *colly.HTMLElement) {
		response.Writers = append(response.Writers, h.Text)
	})

	c.OnHTML(".hNfYaW .gWwKlt", func(h *colly.HTMLElement) {
		x := h.DOM.Find(".gCQkeh")
		newPerson := CastItem{
			Actor: x.Text(),
			Roles: make([]string, 0),
		}
		y := h.DOM.Find("ul")
		y.Children().Each(func(i int, s *goquery.Selection) {
			newPerson.Roles = append(newPerson.Roles, s.Text())
		})

		response.Cast = append(response.Cast, newPerson)
	})

	c.Visit(url)

	var hasErr bool
	c.OnError(func(r *colly.Response, err error) {
		if err != nil {
			hasErr = true
			// return h.JSON(http.StatusNotFound, "not found")
		}
	})
	if hasErr {
		return h.JSON(http.StatusInternalServerError, "this one's one us")
	}

	if response.Title == "" || response.Rate == 0 {
		return h.JSON(http.StatusNotFound, "not found")
	}

	return h.JSON(http.StatusOK, response)
}
