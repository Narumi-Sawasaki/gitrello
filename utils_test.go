package main

import (
	"os"
	"testing"
	"os/exec"
	"path/filepath"
)

func TestGetRepoRoot(t *testing.T) {
	// テストケースを定義
	tests := []struct {
		name         string
		createGitDir bool
		wantErr      bool
	}{
		{
			name:         "正常ケース: .gitディレクトリが存在する",
			createGitDir: true,
			wantErr:      false,
		},
		{
			name:         "エラーケース: .gitディレクトリが存在しない",
			createGitDir: false,
			wantErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// テスト用の一時ディレクトリを作成
			tempDir := t.TempDir()
			os.Chdir(tempDir)

			// 必要に応じて.gitディレクトリを作成
			if test.createGitDir {
				err := exec.Command("git", "init").Run()
				if err != nil {
					t.Fatalf("エラー: .gitディレクトリの作成に失敗しました: %v", err)
				}
			}

			gitDirPath, err := getRepoRoot()

			// エラーの期待値を確認
			if (err != nil) != test.wantErr {
				t.Errorf("GetGitDirPath() エラー = %v, wantErr %v", err, test.wantErr)
				return
			}

			// 正常ケースの場合、返されたパスが正しいことを確認
			if !test.wantErr {
				expectedPath := filepath.Join("/private", tempDir)
				t.Logf("expectedPath: %v", expectedPath)
				t.Logf("gitDirPath: %v", gitDirPath)
				if gitDirPath != expectedPath {
					t.Errorf("GetGitDirPath() = %v, want %v", gitDirPath, expectedPath)
				}
			}
		})
	}
}

