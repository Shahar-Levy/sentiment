package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cdipaolo/sentiment"
)

func main() {

	http.HandleFunc("/", calculateSentiment)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func calculateSentiment(w http.ResponseWriter, r *http.Request) {

	apiKey := os.Getenv("NYT_KEY")
	if apiKey == "" {
		panic("API key empty")
	}

	model, err := sentiment.Restore()
	if err != nil {
		panic(err)
	}

	scores := Scores{}

	for year := 2019; year <= 2019; year++ {
		strYear := fmt.Sprint(year)
		var yearlySentiment float64

		for month := 1; month <= 12; month++ {
			time.Sleep(time.Second * 6)
			strMonth := fmt.Sprint(month)
			url := "https://api.nytimes.com/svc/archive/v1/" + strYear + "/" + strMonth + ".json?api-key=" + apiKey

			resp, err := http.Get(url)
			if err != nil {
				panic(err)
			}

			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			titles := articleTitles{}
			err = json.Unmarshal(bytes, &titles)
			if err != nil {
				panic(err)
			}

			numTitles := len(titles.Response.Docs)
			var monthlyAvg float64

			for _, title := range titles.Response.Docs {

				analysis := model.SentimentAnalysis(title.Headline.Main, sentiment.English)
				monthlyAvg += float64(analysis.Score) / float64(numTitles)
			}
			yearlySentiment += monthlyAvg / 12
			fmt.Println("monthly sentiment for", strMonth, monthlyAvg)

		}
		fmt.Println("yearly sentiment for", strYear, yearlySentiment)
		scores = append(scores, SentimentStruct{
			Year:           strYear,
			SentimentScore: yearlySentiment,
		})
	}

	fmt.Println("scores:", scores)

	response, err := json.Marshal(scores)
	if err != nil {
		panic(err)
	}

	fmt.Println("response:", string(response))

	w.Write(response)
}

type Scores []SentimentStruct

type SentimentStruct struct {
	Year           string
	SentimentScore float64
}

type articleTitles struct {
	Response struct {
		Docs []struct {
			Headline struct {
				Main string `json:"main"`
			} `json:"headline"`
		} `json:"docs"`
	} `json:"response"`
}
