package offwhite

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hah/bot/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	baseURL   = "https://www.off---white.com/en-it"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36'"
)

var (
	client *utils.Client
	b      bag
)

// Item - desired item details
type Item struct {
	Name, Color, Size, URL string
}

// Product - matched product from the site
type Product struct {
	ID, Scale, Size string
}

type bag struct {
	ID string `json:"bagId"`
}

type atcpayload struct {
	MerchantID       int    `json:"merchantId"`
	ProductID        string `json:"productId"`
	Quantity         int    `json:"quantity"`
	Scale            int    `json:"scale"`
	Size             int    `json:"size"`
	CustomAttributes string `json:"customAttributes"`
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
	} `json:"result"`
	Sizes []struct {
		SizeID            string `json:"sizeId"`
		SizeDescription   string `json:"sizeDescription"`
		Scale             string `json:"scale"`
		ScaleAbbreviation string `json:"scaleAbbreviation"`
		IsOneSize         bool   `json:"isOneSize"`
	} `json:"sizes"`
}

// for later
// login
// type ErrResponse struct {
//     Success      bool        `json:"success"`
//     ErrorCode    json.Number `json:"errorCode"`
//     ErrorMessage string      `json:"errorMessage"`
// }

// type TaskAccount struct {
//     Username   string `json:"username"`
//     Password   string `json:"password"`
//     RememberMe bool   `json:"rememberMe"`
// }

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMilli})
	client = utils.CreateClient()
	client.Header = &http.Header{}
	client.Header.Set("User-Agent", userAgent)
	client.Header.Set("Content-Type", "application/json")
	client.Header.Set("Accept", "application/json, text/plain, */*'")
	client.Header.Set("Referer", "https://www.off---white.com/")

	// getting the bag id
	var url bytes.Buffer
	url.WriteString(baseURL)
	url.WriteString("api/users/me")
	response := client.Perform(http.MethodGet, url.String(), nil)
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer response.Body.Close()

	err = json.Unmarshal([]byte(content), &b)
	if err != nil {
		log.Fatal().Err(err)
	}

}

func getDetail(pid string) productDetail {
	var url bytes.Buffer
	url.WriteString(baseURL)
	url.WriteString("/api/products/")
	url.WriteString(pid)
	response := client.Perform(http.MethodGet, url.String(), nil)
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal().Err(err)
	}
	utils.CheckContentType(response.Header)
	defer response.Body.Close()
	var detail productDetail
	err = json.Unmarshal([]byte(content), &detail)
	if err != nil {
		log.Fatal().Err(err)
	}

	return detail
}

// Search - looking for the product
func (i Item) Search() Product {
	log.Info().Msg("No URL provided trying to find a match")
	var url bytes.Buffer
	url.WriteString(baseURL)
	url.WriteString("/api/listing?query=")
	url.WriteString(strings.ReplaceAll(i.Name, " ", "%20"))

	response := client.Perform(http.MethodGet, url.String(), nil)
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal().Err(err)
	}
	utils.CheckContentType(response.Header)
	defer response.Body.Close()
	var result searchResult
	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msg("Looking for an exact match")
	var matched Product

	for _, entry := range result.Products.Entries {
		if entry.ShortDescription == i.Name {
			detail := getDetail(strconv.Itoa(entry.ID))
			for _, c := range detail.Result.Colors {
				if c.Color.Name == i.Color {
					log.Info().Msg("Color matched")
					matched.ID = strconv.Itoa(entry.ID)
					break
				}
			}

			for _, s := range detail.Sizes {
				if s.SizeDescription == i.Size {
					log.Info().Msg("Size matched")
					matched.Size = s.SizeID
					matched.Scale = s.Scale
					break
				}
			}

			break
		}
	}
	return matched
}

// Fetch - fetching product data
func (i Item) Fetch() Product {
	log.Info().Msg("Fetching product details")
	sl := strings.Split(i.URL, "-")
	pid := sl[len(sl)-1]
	fetched := getDetail(pid)

	var matched Product
	matched.ID = pid
	for _, s := range fetched.Sizes {
		if s.SizeDescription == i.Size {
			log.Info().Msg("Size matched")
			matched.Size = s.SizeID
			matched.Scale = s.Scale
			break
		}
	}
	return matched
}

// ATC - Add to cart
func (p Product) ATC() {
	var url bytes.Buffer
	url.WriteString(baseURL)
	url.WriteString("/api/bags/")
	url.WriteString(b.ID)
	url.WriteString("/items")
	scale, err := strconv.Atoi(p.Scale)
	if err != nil {
		log.Fatal().Err(err)

	}
	sizeid, err := strconv.Atoi(p.Size)
	if err != nil {
		log.Fatal().Err(err)
	}

	payload := atcpayload{
		MerchantID:       12572,
		ProductID:        p.ID,
		Quantity:         1,
		Scale:            scale,
		Size:             sizeid,
		CustomAttributes: "",
	}
	jsonvalue, err := json.Marshal(payload)
	if err != nil {
		log.Fatal().Err(err)
	}
	response := client.Perform(http.MethodPost, url.String(), bytes.NewBuffer(jsonvalue))
	// bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal().Err(err)
	}
	utils.CheckContentType(response.Header)
	defer response.Body.Close()
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)
}
