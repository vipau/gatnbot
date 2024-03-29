package fakernews_mod

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/mb-14/gomarkov"
)

const (
	hnBaseURL        = "https://hacker-news.firebaseio.com/v0/"
	hnTopStoriesPath = "topstories.json"
	hnStoryItemPath  = "item/"
)

type hnStory struct {
	Title string `json:"title"`
}

// TrainModel refreshes model with the latest 500 HN stories
func TrainModel() {
	chain, err := buildModel()
	if err != nil {
		fmt.Println(err)
		return
	}
	saveModel(chain)
}

// GenerateNews generates fake HN story (public method)
func GenerateNews() string {
	chain, err := loadModel()
	if err != nil {
		fmt.Println(err)
		return "A bug just happened! Please report this to gatnbot devs"
	}
	return generateHNStory(chain)
}

// build the markov model from scratch
func buildModel() (*gomarkov.Chain, error) {
	stories, err := fetchHNTopStories()
	if err != nil {
		return nil, err
	}
	chain := gomarkov.NewChain(1)
	var wg sync.WaitGroup
	wg.Add(len(stories))
	fmt.Println("Adding HN story titles to markov chain...")
	for _, storyID := range stories {
		go func(storyID int) {
			defer wg.Done()
			story, err := fetchHNStory(storyID)
			if err != nil {
				fmt.Println(err)
				return
			}
			chain.Add(strings.Split(story.Title, " "))
		}(storyID)
	}
	wg.Wait()
	return chain, nil
}

// load markov model from file
func loadModel() (*gomarkov.Chain, error) {
	var chain gomarkov.Chain
	data, err := os.ReadFile("model.json")
	if err != nil {
		return &chain, err
	}
	err = json.Unmarshal(data, &chain)
	if err != nil {
		return &chain, err
	}
	return &chain, nil
}

// fetch top 500 articles from HN
func fetchHNTopStories() ([]int, error) {
	fmt.Println("Fetching HN top stories...")
	resp, err := http.Get(fmt.Sprintf("%s%s", hnBaseURL, hnTopStoriesPath))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var stories []int
	err = json.Unmarshal(body, &stories)
	return stories, err
}

// fetch single HN article
func fetchHNStory(storyID int) (hnStory, error) {
	var story hnStory
	resp, err := http.Get(fmt.Sprintf("%s%s%d.json", hnBaseURL, hnStoryItemPath, storyID))
	if err != nil {
		return story, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return story, err
	}
	err = json.Unmarshal(body, &story)
	return story, err
}

// save markov model to file
func saveModel(chain *gomarkov.Chain) {
	jsonObj, _ := json.Marshal(chain)
	err := os.WriteFile("model.json", jsonObj, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

// generate fake HN story
func generateHNStory(chain *gomarkov.Chain) string {
	tokens := []string{gomarkov.StartToken}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, _ := chain.Generate(tokens[(len(tokens) - 1):])
		tokens = append(tokens, next)
	}
	return strings.Join(tokens[1:len(tokens)-1], " ")
}
