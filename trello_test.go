package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetCurrentBranchName(branchName string) {
	getCurrentBranchName = func() (string, error) { return branchName, nil }
}

func TestGetCurrentTaskName(t *testing.T) {
	t.Run("ブランチ名からタスク番号を取得する", func(t *testing.T) {
		SetCurrentBranchName("task/30")
		got, err := getCurrentTaskName()
		want := "30"
		assert.Equal(t, err, nil)
		assert.Equal(t, got, want)
	})

	t.Run("ブランチ名が'task/'というprefixが無ければエラーを返す", func(t *testing.T) {
		SetCurrentBranchName("30")
		_, err := getCurrentTaskName()
		assert.NotEqual(t, err, nil)
	})
}

type MockApiClient struct {
	Config Config
}

// http.Get自体をmockするとhttp.Getのみをwrapする構造体を作ることになる
// Clientが提供するメソッドはurl・queryParamsの構築、responseのパースまで行う
func (m *MockApiClient) Get(url string, queryParams QueryParams) (ParsedRes, error) {
	return ParsedRes{
		"cards": []int{1, 2, 3},
	}, nil
}

func TestFetchCards(t *testing.T) {
	config := Config{}
	client := &MockApiClient{
		Config: config,
	}

	cards, _ := fetchCards(client, "1")
	assert.Contains(t, cards, "cards")
}
