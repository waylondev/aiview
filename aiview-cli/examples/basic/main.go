package main

import (
	"fmt"
	"log"

	"github.com/jackwener/aiview/pkg/aiview"
)

func main() {
	// 创建 Bilibili 客户端
	client, err := aiview.New("bilibili")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Platform: %s\n", client.PlatformName())

	// 获取 Bilibili 专用客户端
	biliClient, err := client.BilibiliClient()
	if err != nil {
		log.Fatal(err)
	}

	// 获取热门视频
	hotVideos, err := biliClient.GetHotVideos(1, 10)
	if err != nil {
		log.Printf("GetHotVideos error: %v", err)
	} else {
		fmt.Println("\n=== 热门视频 ===")
		for i, video := range hotVideos {
			fmt.Printf("%d. %s - %s (播放: %d)\n", i+1, video.Title, video.Author, video.Play)
		}
	}

	// 获取视频详情
	videoInfo, err := biliClient.GetVideoInfo("BV1GJ411x7h7")
	if err != nil {
		log.Printf("GetVideoInfo error: %v", err)
	} else {
		fmt.Printf("\n=== 视频详情 ===\n")
		fmt.Printf("标题: %s\n", videoInfo.Title)
		fmt.Printf("作者: %s\n", videoInfo.Author)
		fmt.Printf("播放: %d, 点赞: %d, 投币: %d\n", videoInfo.Play, videoInfo.Like, videoInfo.Coin)
	}

	// 搜索视频
	searchResults, err := biliClient.SearchVideo("Golang", 1)
	if err != nil {
		log.Printf("SearchVideo error: %v", err)
	} else {
		fmt.Println("\n=== 搜索结果 ===")
		for i, item := range searchResults {
			if i >= 5 {
				break
			}
			fmt.Printf("%d. %s - %s\n", i+1, item.Title, item.Author)
		}
	}

	// 抖音示例
	fmt.Println("\n=== 抖音示例 ===")
	douyinClient, err := aiview.New("douyin")
	if err != nil {
		log.Printf("Douyin client error: %v", err)
	} else {
		dyClient, err := douyinClient.DouyinClient()
		if err == nil {
			hotSearch, err := dyClient.GetHotSearch()
			if err == nil {
				fmt.Println("抖音热搜:")
				for i, item := range hotSearch {
					if i >= 5 {
						break
					}
					fmt.Printf("%d. %s (热度: %d)\n", i+1, item.Keyword, item.HotValue)
				}
			}
		}
	}

	// 小红书示例
	fmt.Println("\n=== 小红书示例 ===")
	xhsClient, err := aiview.New("xiaohongshu")
	if err != nil {
		log.Printf("Xiaohongshu client error: %v", err)
	} else {
		xhs, err := xhsClient.XiaohongshuClient()
		if err == nil {
			hotNotes, err := xhs.GetHotNotes()
			if err == nil {
				fmt.Println("小红书热门:")
				for i, item := range hotNotes {
					if i >= 5 {
						break
					}
					fmt.Printf("%d. %s\n", i+1, item.Keyword)
				}
			}
		}
	}
}
