// gitrelloというコマンドを作りたい
// gitrello -h でヘルプが表示される
// gitrello -v でバージョンが表示される

package main

import (
	"encoding/json"
	"fmt"
	"bufio"
	"net/http"
	"io"
	"github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
	"os"
	"regexp"
	"github.com/hashicorp/go-envparse"
)

type ParsedRes map[string]interface{}

func callGetApi() (ParsedRes, error) {
	repoRoot, err := getRepoRoot()
	if err != nil {
		return nil, err
	}
	trellorcPath := repoRoot + "/.trellorc"
	file, err := os.Open(trellorcPath)
	if err != nil {
		return nil, err
	}
	env, err := envparse.Parse(bufio.NewReader(file))
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.trello.com/1/search?query=board:%s%%2029&key=%s&token=%s", env["BOARD_ID"], env["API_KEY"], env["API_TOKEN"])
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data ParsedRes
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// func getRepoRoot() (string, error) {
// 	r, err := git.PlainOpen(".")
// 	if err != nil {
// 		return "", err
// 	}
// 	w, err := r.Worktree()
// 	if err != nil {
// 		return "", err
// 	}
// 	return w.Filesystem.Root(), nil
// }

func getCurrentBranchName() string {
	r, err := git.PlainOpen(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ref, err := r.Head()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return ref.Name().Short()
}

func getCurrentTaskName() string {
	branchName := getCurrentBranchName()

	prefix := "task/"
	if len(branchName) < len(prefix) || branchName[:len(prefix)] != prefix {
		fmt.Fprintln(os.Stderr, "ブランチ名が不正です")
		os.Exit(1)
	}

	return branchName[len(prefix):]
}

func showCard(c *cli.Context) error {
	res, err := callGetApi()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	for _, v := range res["cards"].([]interface{}) {
		card := v.(map[string]interface{})

		urlRegex := fmt.Sprintf("^(https://trello.com/c/.*/%s-).*", getCurrentTaskName())
		if matched, _ := regexp.MatchString(urlRegex, card["url"].(string)); matched {
			fmt.Println(card["name"])
		} else {
			fmt.Println("not matched")
		}

	}
	return nil
}

func main() {
	app := &cli.App{
		Usage: "gitrello",
		Commands: []*cli.Command{
			{
				Name:    "c",
				Aliases: []string{"card"},
				Action:  showCard,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
