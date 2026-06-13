package main

import (
	"fmt"
	"log"

	"github.com/jackwener/aiview/pkg/aiview"
)

// 基础使用示例 - 演示如何使用 Aiview Go Library
func main() {
	fmt.Println("=== Aiview 基础使用示例 ===")

	// 1. Bilibili 示例
	bilibiliExample()

	// 2. 抖音示例
	douyinExample()

	// 3. 小红书示例
	xiaohongshuExample()
}

func bilibiliExample() {
	fmt.Println("--- Bilibili 平台 ---")

	// 创建客户端
	client, err := aiview.New("bilibili")
	if err != nil {
		log.Printf("创建 Bilibili 客户端失败: %v", err)
		return
	}

	fmt.Printf("平台: %s\n", client.PlatformName())

	// 获取 Bilibili 专用客户端
	biliClient, err := client.BilibiliClient()
	if err != nil {
		log.Printf("获取 Bilibili 客户端失败: %v", err)
		return
	}

	// 示例 1: 获取热门视频
	fmt.Println("\n【获取热门视频】")
	hotVideos, err := biliClient.GetHotVideos(1, 5)
	if err != nil {
		log.Printf("获取热门视频失败: %v", err)
	} else {
		for i, video := range hotVideos {
			fmt.Printf("%d. %s - %s\n", i+1, video.Title, video.Author)
			fmt.Printf("   播放: %d, 点赞: %d\n", video.Play, video.Like)
		}
	}

	// 示例 2: 获取视频详情
	fmt.Println("\n【获取视频详情】")
	videoInfo, err := biliClient.GetVideoInfo("BV1GJ411x7Rq")
	if err != nil {
		log.Printf("获取视频详情失败: %v", err)
	} else {
		fmt.Printf("标题: %s\n", videoInfo.Title)
		fmt.Printf("作者: %s (ID: %s)\n", videoInfo.Author, videoInfo.AuthorID)
		fmt.Printf("播放: %d, 弹幕: %d, 点赞: %d\n", videoInfo.Play, videoInfo.Danmaku, videoInfo.Like)
		fmt.Printf("链接: %s\n", videoInfo.URL)
	}

	// 示例 3: 获取用户信息
	fmt.Println("\n【获取用户信息】")
	userInfo, err := biliClient.GetUserInfo(37737161)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
	} else {
		fmt.Printf("昵称: %s\n", userInfo.Name)
		fmt.Printf("签名: %s\n", userInfo.Sign)
		fmt.Printf("粉丝: %d, 关注: %d\n", userInfo.Fans, userInfo.Follow)
	}

	// 示例 4: 搜索视频
	fmt.Println("\n【搜索视频】")
	searchResults, err := biliClient.SearchVideo("Golang", 1)
	if err != nil {
		log.Printf("搜索视频失败: %v", err)
	} else {
		fmt.Printf("找到 %d 个结果:\n", len(searchResults))
		for i, item := range searchResults {
			if i >= 5 {
				break
			}
			fmt.Printf("%d. %s - %s\n", i+1, item.Title, item.Author)
		}
	}
}

func douyinExample() {
	fmt.Println("\n--- 抖音平台 ---")

	// 创建客户端
	client, err := aiview.New("douyin")
	if err != nil {
		log.Printf("创建抖音客户端失败: %v", err)
		return
	}

	fmt.Printf("平台: %s\n", client.PlatformName())

	// 获取抖音专用客户端
	dyClient, err := client.DouyinClient()
	if err != nil {
		log.Printf("获取抖音客户端失败: %v", err)
		return
	}

	// 示例 1: 获取热搜
	fmt.Println("\n【获取热搜】")
	hotSearch, err := dyClient.GetHotSearch()
	if err != nil {
		log.Printf("获取热搜失败: %v", err)
	} else {
		for i, item := range hotSearch {
			if i >= 10 {
				break
			}
			fmt.Printf("%d. %s (热度: %d)\n", i+1, item.Keyword, item.HotValue)
		}
	}

	// 示例 2: 搜索视频
	fmt.Println("\n【搜索视频】")
	searchResults, err := dyClient.Search("AI技术", 1, 5)
	if err != nil {
		log.Printf("搜索视频失败: %v", err)
	} else {
		fmt.Printf("找到 %d 个结果:\n", len(searchResults))
		for i, item := range searchResults {
			fmt.Printf("%d. %s - %s\n", i+1, item.Title, item.Author)
		}
	}
}

func xiaohongshuExample() {
	fmt.Println("\n--- 小红书平台 ---")

	// 创建客户端
	client, err := aiview.New("xiaohongshu")
	if err != nil {
		log.Printf("创建小红书客户端失败: %v", err)
		return
	}

	fmt.Printf("平台: %s\n", client.PlatformName())

	// 获取小红书专用客户端
	xhsClient, err := client.XiaohongshuClient()
	if err != nil {
		log.Printf("获取小红书客户端失败: %v", err)
		return
	}

	// 示例 1: 获取热门笔记
	fmt.Println("\n【获取热门笔记】")
	hotNotes, err := xhsClient.GetHotNotes()
	if err != nil {
		log.Printf("获取热门笔记失败: %v", err)
	} else {
		for i, item := range hotNotes {
			if i >= 10 {
				break
			}
			fmt.Printf("%d. %s (热度: %d)\n", i+1, item.Keyword, item.HotValue)
		}
	}

	// 示例 2: 搜索笔记
	fmt.Println("\n【搜索笔记】")
	searchResults, err := xhsClient.SearchNotes("旅行攻略", 1)
	if err != nil {
		log.Printf("搜索笔记失败: %v", err)
	} else {
		fmt.Printf("找到 %d 个结果:\n", len(searchResults))
		for i, item := range searchResults {
			if i >= 5 {
				break
			}
			fmt.Printf("%d. %s - %s\n", i+1, item.Title, item.Author)
		}
	}
}
