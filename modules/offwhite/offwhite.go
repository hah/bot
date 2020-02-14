package offwhite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hah/bot/utils"
)

const (
	baseURL   = "https://www.off---white.com/en-it/"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36'"
)

var (
	client *utils.Client
	b      bag
)

// Item is the product from the site
type Item struct {
	Name, Color, Size, URL string
}

type bag struct {
	ID string `json:"bagId"`
}

type searchResult struct {
	Name     interface{} `json:"name"`
	Products struct {
		Entries []struct {
			ID               int     `json:"id"`
			ShortDescription string  `json:"shortDescription"`
			Price            float64 `json:"price"`
			FormattedPrice   string  `json:"formattedPrice"`
			Gender           int     `json:"gender"`
			Slug             string  `json:"slug"`
		} `json:"entries"`
		Number     int `json:"number"`
		TotalItems int `json:"totalItems"`
	} `json:"products"`
}

type productDetail struct {
	Result struct {
		Colors []struct {
			Color struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"color"`
		} `json:"colors"`
		// 	Variants         []struct {
		// 		Attributes []struct {
		// 			Type        int    `json:"type"`
		// 			Value       string `json:"value"`
		// 			Description string `json:"description"`
		// 		} `json:"attributes"`
		// 		AvailableAt []int `json:"availableAt"`
		// 		MerchantID  int   `json:"merchantId"`
		// } `json:"result"`
		Sizes []struct {
			SizeID            string `json:"sizeId"`
			SizeDescription   string `json:"sizeDescription"`
			Scale             string `json:"scale"`
			ScaleAbbreviation string `json:"scaleAbbreviation"`
			IsOneSize         bool   `json:"isOneSize"`
			// Variants          []struct {
			// 	MerchantID                    int      `json:"merchantId"`
			// 	FormattedPrice                string   `json:"formattedPrice"`
			// 	FormattedPriceWithoutDiscount string   `json:"formattedPriceWithoutDiscount"`
			// 	Quantity                      int      `json:"quantity"`
			// 	Barcodes                      []string `json:"barcodes"`
			// } `json:"variants"`
		} `json:"sizes"`
	}
}

func init() {
	client = utils.CreateClient()
	client.Header = &http.Header{}
	client.Header.Set("User-Agent", userAgent)
	client.Header.Set("Authority", "www.off---white.com")
	client.Header.Set("Accept", "application/json, text/plain, */*'")
	client.Header.Set("Referer", "https://www.off---white.com/")

	// getting the bag id
	var url bytes.Buffer
	url.WriteString(baseURL)
	url.WriteString("api/users/me")
	response := client.Perform(http.MethodGet, url.String(), nil)
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	err = json.Unmarshal([]byte(content), &b)
	if err != nil {
		log.Fatal(err)
	}

}

// Search - looking for the product
func (i Item) Search() {
	var url bytes.Buffer
	url.WriteString(baseURL)
	url.WriteString("/api/listing?query=")
	url.WriteString(strings.ReplaceAll(i.Name, " ", "%20"))

	response := client.Perform(http.MethodGet, url.String(), nil)
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	var result searchResult
	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("looking for exact match.")
	for _, entry := range result.Products.Entries {
		if entry.ShortDescription == i.Name {

			var url bytes.Buffer
			url.WriteString(baseURL)
			url.WriteString("/api/products/")
			url.WriteString(strconv.Itoa(entry.ID))
			response := client.Perform(http.MethodGet, url.String(), nil)
			content, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()
			var detail productDetail
			err = json.Unmarshal([]byte(content), &detail)
			if err != nil {
				log.Fatal(err)
			}
			for _, c := range detail.Result.Colors {
				if c.Color.Name == i.Color {
					fmt.Println("match found")
					fmt.Println(entry)
					break
				}
			}
			break
		}
	}

}

// Fetch - fetching product data
func (i Item) Fetch() {

}
