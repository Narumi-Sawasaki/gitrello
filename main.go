// gitrelloというコマンドを作りたい
// gitrello -h でヘルプが表示される
// gitrello -v でバージョンが表示される

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"os/exec"

	"github.com/hashicorp/go-envparse"
	"github.com/urfave/cli/v2"
)

type ParsedRes map[string]interface{}

// func callGetApi() (ParsedRes, error) {
// 	url := fmt.Sprintf("https://api.trello.com/1/search?query=board:%s%%2029&key=%s&token=%s", env["BOARD_ID"], env["API_KEY"], env["API_TOKEN"])
// 	resp, _ := http.Get(url)
// 	defer resp.Body.Close()
// 	body, _ := io.ReadAll(resp.Body)
// 	var data ParsedRes
// 	if err := json.Unmarshal(body, &data); err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }

type QueryParams map[string]string
type Config map[string]string
type ApiClient interface {
	Get(url string, queryParams QueryParams) (ParsedRes, error)
}
type ApiClientImpl struct {
	Config Config
}

func (client *ApiClientImpl) Get(url string, queryParams QueryParams) (ParsedRes, error) {
	// queryParamsを組み立てる
	keyValues := []string{}
	for key, value := range queryParams {
		keyValues = append(keyValues, key+"="+value)
	}
	urlWithParams := url
	if len(keyValues) > 0 {
		urlWithParams += fmt.Sprintf("?%s", strings.Join(keyValues, "&"))
	}
	resp, _ := http.Get(urlWithParams)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data ParsedRes
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil

}

func fetchCards(client ApiClient, config Config, keyword string) (ParsedRes, error) {
	return client.Get("https://api.trello.com/1/search", QueryParams{
	"query": fmt.Sprintf("board:%s%%20%s", config["BOARD_ID"], keyword),
	"key":   config["API_KEY"],
	"token": config["API_TOKEN"],
})
}

func showCard(_ *cli.Context) error {
	repoRoot, err := getRepoRoot()
	if err != nil {
		return err
	}
	trellorcPath := repoRoot + "/.trellorc"
	file, err := os.Open(trellorcPath)
	if err != nil {
		return err
	}
	env, err := envparse.Parse(bufio.NewReader(file))
	if err != nil {
		return err
	}

	currentTaskName, err := getCurrentTaskName()
	if err != nil {
		return err
	}

	client := &ApiClientImpl{}
	res, err := fetchCards(client, Config(env), currentTaskName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	for _, v := range res["cards"].([]interface{}) {
		card := v.(map[string]interface{})
		urlRegex := fmt.Sprintf("^(https://trello.com/c/.*/%s-).*", currentTaskName)
		if matched, _ := regexp.MatchString(urlRegex, card["url"].(string)); matched {
			fmt.Println(fmt.Sprintf("%s\n%s", card["name"], card["url"]))
			return nil
		}
	}
	return errors.New("Not matched")
}

func openCard(_ *cli.Context) error {
	repoRoot, err := getRepoRoot()
	if err != nil {
		return err
	}
	trellorcPath := repoRoot + "/.trellorc"
	file, err := os.Open(trellorcPath)
	if err != nil {
		return err
	}
	env, err := envparse.Parse(bufio.NewReader(file))
	if err != nil {
		return err
	}

	currentTaskName, err := getCurrentTaskName()
	if err != nil {
		return err
	}

	client := &ApiClientImpl{}
	res, err := fetchCards(client, Config(env), currentTaskName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	for _, v := range res["cards"].([]interface{}) {
		card := v.(map[string]interface{})
		urlRegex := fmt.Sprintf("^(https://trello.com/c/.*/%s-).*", currentTaskName)
		if matched, _ := regexp.MatchString(urlRegex, card["url"].(string)); matched {
			err := exec.Command("open", card["url"].(string)).Start()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return nil
		}
	}
	return errors.New("Not matched")
}

func main() {
	app := &cli.App{
		Usage: "gitrello",
		Commands: []*cli.Command{
			{
				Name:    "c",
				Aliases: []string{"card"},
				Action:  showCard,
				Subcommands: []*cli.Command{
					{
						Name: "open",
						Action: openCard,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
