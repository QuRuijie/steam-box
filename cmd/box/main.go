package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/YouEclipse/steam-box/pkg/steambox"
	"github.com/google/go-github/github"
)

var (
	filename string
	lines1   []string
	lines2   []string
	err      error
)

func main() {
	appIDList := GetAppList()
	ctx := context.Background()
	gistID := os.Getenv("GIST_ID")
	gistIDRecent := os.Getenv("GIST_ID_RECENT")
	steamID, _ := strconv.ParseUint(os.Getenv("STEAM_ID"), 10, 64)
	box := GetBox()

	fmt.Println(gistID, gistIDRecent)

	updateOption := os.Getenv("UPDATE_OPTION") // options for update: GIST (Gist only), MARKDOWN (README only), GIST_AND_MARKDOWN (Gist and README)
	markdownFile := os.Getenv("MARKDOWN_FILE") // the markdown filename (e.g. MYFILE.md)
	var updateGist, updateMarkdown bool
	if updateOption == "MARKDOWN" {
		updateMarkdown = true
	} else if updateOption == "GIST_AND_MARKDOWN" {
		updateGist = true
		updateMarkdown = true
	} else {
		updateGist = true
	}

	// 更新总游戏时间
	filename = "🎮 Steam Game Time"
	lines1, err = box.GetPlayTime(ctx, steamID, false, appIDList...)
	if err != nil {
		panic("GetPlayTime err:" + err.Error())
	}

	if updateGist {
		gist, err := box.GetGist(ctx, gistID)
		if err != nil {
			panic("GetGist err:" + err.Error())
		}

		f := gist.Files[github.GistFilename(filename)]

		f.Content = github.String(strings.Join(lines1, "\n"))
		gist.Files[github.GistFilename(filename)] = f

		err = box.UpdateGist(ctx, gistID, gist)
		if err != nil {
			panic("UpdateGist err:" + err.Error())
		}
	}

	// 更新近期游戏时间
	filename = "🎮 Steam Game Recent"
	lines2, err = box.GetRecentGames(ctx, steamID, false)
	if err != nil {
		panic("GetRecentGames err:" + err.Error())
	}

	if updateGist {
		gist, err := box.GetGist(ctx, gistIDRecent)
		if err != nil {
			panic("GetGist err:" + err.Error())
		}

		f := gist.Files[github.GistFilename(filename)]

		f.Content = github.String(strings.Join(lines2, "\n"))
		gist.Files[github.GistFilename(filename)] = f

		err = box.UpdateGist(ctx, gistIDRecent, gist)
		if err != nil {
			panic("UpdateGist err:" + err.Error())
		}
	}

	if updateMarkdown && markdownFile != "" {
		title := filename
		if updateGist {
			title = fmt.Sprintf(`#### <a href="https://gist.github.com/%s" target="_blank">%s</a>`, gistID, title)
		}

		content := bytes.NewBuffer(nil)
		content.WriteString(strings.Join(lines1, "\n"))

		err = box.UpdateMarkdown(ctx, title, markdownFile, content.Bytes())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("updating markdown successfully on ", markdownFile)
	}
}

func GetAppList() []uint32 {
	appIDs := os.Getenv("APP_ID")
	appIDList := make([]uint32, 0)

	for _, appID := range strings.Split(appIDs, ",") {
		appid, err := strconv.ParseUint(appID, 10, 32)
		if err != nil {
			continue
		}
		appIDList = append(appIDList, uint32(appid))
	}
	return appIDList
}

func GetBox() *steambox.Box {
	ghToken := os.Getenv("GH_TOKEN")
	ghUsername := os.Getenv("GH_USER")
	steamAPIKey := os.Getenv("STEAM_API_KEY")
	return steambox.NewBox(steamAPIKey, ghUsername, ghToken)
}