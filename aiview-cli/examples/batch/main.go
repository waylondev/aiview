package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackwener/aiview/pkg/aiview"
)

// 批量操作示例 - 演示如何高效批量获取和处理数据
func main() {
	fmt.Println("=== Aiview 批量操作示例 ===")

	// 1. 批量获取用户信息
	batchUserExample()

	// 2. 批量搜索并统计
	batchSearchExample()

	// 3. 并发获取视频详情
	concurrentVideoExample()

	// 4. 分页获取所有数据
	paginationExample()
}

// 批量获取用户信息 - 演示如何批量获取多个用户的信息
func batchUserExample() {
	fmt.Println("--- 批量获取用户信息 ---")

	client, err := aiview.New("bilibili")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	biliClient, err := client.BilibiliClient()
	if err != nil {
		log.Fatalf("获取 Bilibili 客户端失败: %v", err)
	}

	// 用户 ID 列表（示例）
	userIDs := []int{
		37737161,
		// 可以添加更多用户 ID
	}

	fmt.Printf("批量获取 %d 个用户信息:\n", len(userIDs))
	var users []*aiview.UserInfo
	for i, uid := range userIDs {
		fmt.Printf("[%d/%d] 获取用户 %d ... ", i+1, len(userIDs), uid)
		user, err := biliClient.GetUserInfo(uid)
		if err != nil {
			fmt.Printf("失败: %v\n", err)
			continue
		}
		fmt.Printf("成功: %s (粉丝: %d)\n", user.Name, user.Fans)
		users = append(users, user)
		time.Sleep(time.Millisecond * 500) // 避免请求过快
	}

	fmt.Printf("\n成功获取 %d/%d 个用户信息\n", len(users), len(userIDs))

	// 统计粉丝总数
	totalFans := 0
	for _, u := range users {
		totalFans += u.Fans
	}
	fmt.Printf("粉丝总数: %d\n", totalFans)
}

// 批量搜索并统计 - 演示如何批量搜索关键词并统计结果
func batchSearchExample() {
	fmt.Println("\n--- 批量搜索并统计 ---")

	client, err := aiview.New("bilibili")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	biliClient, err := client.BilibiliClient()
	if err != nil {
		log.Fatalf("获取 Bilibili 客户端失败: %v", err)
	}

	// 搜索关键词列表
	keywords := []string{
		"Golang",
		"Python",
		"JavaScript",
		"Rust",
	}

	fmt.Printf("批量搜索 %d 个关键词:\n", len(keywords))
	type SearchResult struct {
		Keyword string
		Count   int
	}
	var results []SearchResult

	for i, keyword := range keywords {
		fmt.Printf("[%d/%d] 搜索 '%s' ... ", i+1, len(keywords), keyword)
		searchResults, err := biliClient.SearchVideo(keyword, 1)
		if err != nil {
			fmt.Printf("失败: %v\n", err)
			continue
		}
		fmt.Printf("成功: 找到 %d 个视频\n", len(searchResults))
		results = append(results, SearchResult{
			Keyword: keyword,
			Count:   len(searchResults),
		})
		time.Sleep(time.Millisecond * 500)
	}

	// 统计结果
	fmt.Println("\n搜索结果统计:")
	for _, r := range results {
		fmt.Printf("  %s: %d 个视频\n", r.Keyword, r.Count)
	}
}

// 并发获取视频详情 - 演示如何使用 goroutine 并发获取数据
func concurrentVideoExample() {
	fmt.Println("\n--- 并发获取视频详情 ---")

	client, err := aiview.New("bilibili")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	biliClient, err := client.BilibiliClient()
	if err != nil {
		log.Fatalf("获取 Bilibili 客户端失败: %v", err)
	}

	// 视频 BV 号列表
	bvList := []string{
		"BV1GJ411x7Rq",
		"BV1x54y1e7zf",
		"BV1hK4y1L75e",
		"BV1vK4y1L75e",
		"BV1hK4y1L75f",
	}

	fmt.Printf("并发获取 %d 个视频信息:\n", len(bvList))

	// 使用 goroutine 并发获取
	var wg sync.WaitGroup
	var mu sync.Mutex
	var videoInfos []*aiview.VideoInfo
	var errors []error

	// 限制并发数
	maxConcurrency := 3
	sem := make(chan struct{}, maxConcurrency)

	startTime := time.Now()

	for i, bv := range bvList {
		wg.Add(1)
		go func(index int, bvid string) {
			defer wg.Done()

			// 获取信号量
			sem <- struct{}{}
			defer func() { <-sem }()

			fmt.Printf("[%d/%d] 获取 %s ... ", index+1, len(bvList), bvid)
			info, err := biliClient.GetVideoInfo(bvid)
			if err != nil {
				fmt.Printf("失败: %v\n", err)
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}
			fmt.Printf("成功: %s\n", info.Title)

			mu.Lock()
			videoInfos = append(videoInfos, info)
			mu.Unlock()

			time.Sleep(time.Millisecond * 300) // 模拟网络延迟
		}(i, bv)
	}

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Printf("\n并发获取完成:\n")
	fmt.Printf("  成功: %d/%d\n", len(videoInfos), len(bvList))
	fmt.Printf("  失败: %d\n", len(errors))
	fmt.Printf("  耗时: %v\n", duration)
	fmt.Printf("  平均每个: %v\n", duration/time.Duration(len(bvList)))
}

// 分页获取所有数据 - 演示如何分页获取所有数据
func paginationExample() {
	fmt.Println("\n--- 分页获取所有数据 ---")

	client, err := aiview.New("bilibili")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	biliClient, err := client.BilibiliClient()
	if err != nil {
		log.Fatalf("获取 Bilibili 客户端失败: %v", err)
	}

	// 分页获取热门视频
	fmt.Println("分页获取热门视频 (每页 10 个，共 3 页):")
	var allVideos []aiview.VideoInfo

	for page := 1; page <= 3; page++ {
		fmt.Printf("获取第 %d 页 ... ", page)
		videos, err := biliClient.GetHotVideos(page, 10)
		if err != nil {
			fmt.Printf("失败: %v\n", err)
			continue
		}
		fmt.Printf("成功: 获取 %d 个视频\n", len(videos))
		allVideos = append(allVideos, videos...)
		time.Sleep(time.Millisecond * 500)
	}

	fmt.Printf("\n总共获取 %d 个视频\n", len(allVideos))

	// 按播放量排序
	fmt.Println("\n按播放量排序 (前 10 名):")
	sortByPlayCount(allVideos)
	for i, v := range allVideos {
		if i >= 10 {
			break
		}
		fmt.Printf("%d. %s - 播放: %d\n", i+1, v.Title, v.Play)
	}
}

func sortByPlayCount(videos []aiview.VideoInfo) {
	// 简单的冒泡排序
	n := len(videos)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if videos[j].Play < videos[j+1].Play {
				videos[j], videos[j+1] = videos[j+1], videos[j]
			}
		}
	}
}
