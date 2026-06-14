package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackwener/aiview/pkg/aiview"
)

// 高级使用示例 - 演示错误处理、数据过滤、批量操作等高级用法
func main() {
	fmt.Println("=== Aiview 高级使用示例 ===")

	// 1. 错误处理模式
	errorHandlingExample()

	// 2. 数据过滤和统计
	dataFilteringExample()

	// 3. 批量获取视频信息
	batchVideoInfoExample()

	// 4. 数据导出为 JSON
	exportJSONExample()

	// 5. 多平台数据聚合
	multiPlatformAggregation()
}

// 错误处理模式 - 演示如何处理各种错误场景
func errorHandlingExample() {
	fmt.Println("--- 错误处理模式 ---")

	client, err := aiview.New("bilibili")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	biliClient, err := client.BilibiliClient()
	if err != nil {
		log.Fatalf("获取 Bilibili 客户端失败: %v", err)
	}

	// 场景 1: 处理不存在的视频
	fmt.Println("\n【处理不存在的视频】")
	_, err = biliClient.GetVideoInfo("BV_NOT_EXIST_12345")
	if err != nil {
		fmt.Printf("✓ 预期错误: %v\n", err)
	}

	// 场景 2: 处理空搜索结果
	fmt.Println("\n【处理空搜索结果】")
	results, err := biliClient.SearchVideo("xxxxxxxxxxxxxxxxxxxx", 1)
	if err != nil {
		log.Printf("搜索失败: %v", err)
	} else if len(results) == 0 {
		fmt.Println("✓ 搜索结果为空")
	} else {
		fmt.Printf("找到 %d 个结果\n", len(results))
	}

	// 场景 3: 重试机制
	fmt.Println("\n【重试机制】")
	var hotVideos []aiview.VideoInfo
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		hotVideos, err = biliClient.GetHotVideos(1, 10)
		if err == nil {
			fmt.Printf("✓ 第 %d 次尝试成功\n", i+1)
			break
		}
		fmt.Printf("第 %d 次尝试失败: %v\n", i+1, err)
		if i < maxRetries-1 {
			time.Sleep(time.Second * 2)
		}
	}
	if err != nil {
		log.Printf("重试 %d 次后仍然失败", maxRetries)
	} else {
		fmt.Printf("成功获取 %d 个热门视频\n", len(hotVideos))
	}
}

// 数据过滤和统计 - 演示如何过滤和统计数据
func dataFilteringExample() {
	fmt.Println("\n--- 数据过滤和统计 ---")

	client, err := aiview.New("bilibili")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	biliClient, err := client.BilibiliClient()
	if err != nil {
		log.Fatalf("获取 Bilibili 客户端失败: %v", err)
	}

	// 获取热门视频
	hotVideos, err := biliClient.GetHotVideos(1, 20)
	if err != nil {
		log.Fatalf("获取热门视频失败: %v", err)
	}

	// 过滤：播放量超过 100 万的视频
	fmt.Println("\n【播放量超过 100 万的视频】")
	popularVideos := filterByPlay量(hotVideos, 1000000)
	fmt.Printf("找到 %d 个视频:\n", len(popularVideos))
	for i, v := range popularVideos {
		if i >= 5 {
			fmt.Printf("... 还有 %d 个视频\n", len(popularVideos)-5)
			break
		}
		fmt.Printf("%d. %s - 播放: %d\n", i+1, v.Title, v.Play)
	}

	// 统计：计算平均播放量
	fmt.Println("\n【统计数据】")
	stats := calculateStats(hotVideos)
	fmt.Printf("总视频数: %d\n", stats.Count)
	fmt.Printf("平均播放量: %.2f\n", stats.AvgPlay)
	fmt.Printf("最高播放量: %d (%s)\n", stats.MaxPlay, stats.MaxPlayTitle)
	fmt.Printf("最低播放量: %d (%s)\n", stats.MinPlay, stats.MinPlayTitle)
}

type VideoStats struct {
	Count        int
	AvgPlay      float64
	MaxPlay      int
	MaxPlayTitle string
	MinPlay      int
	MinPlayTitle string
}

func filterByPlay量(videos []aiview.VideoInfo, threshold int) []aiview.VideoInfo {
	var filtered []aiview.VideoInfo
	for _, v := range videos {
		if v.Play >= threshold {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func calculateStats(videos []aiview.VideoInfo) VideoStats {
	if len(videos) == 0 {
		return VideoStats{}
	}

	stats := VideoStats{
		Count:   len(videos),
		MaxPlay: videos[0].Play,
		MinPlay: videos[0].Play,
	}

	totalPlay := 0
	for _, v := range videos {
		totalPlay += v.Play
		if v.Play > stats.MaxPlay {
			stats.MaxPlay = v.Play
			stats.MaxPlayTitle = v.Title
		}
		if v.Play < stats.MinPlay {
			stats.MinPlay = v.Play
			stats.MinPlayTitle = v.Title
		}
	}
	stats.AvgPlay = float64(totalPlay) / float64(stats.Count)

	return stats
}

// 批量获取视频信息 - 演示如何批量处理多个视频
func batchVideoInfoExample() {
	fmt.Println("\n--- 批量获取视频信息 ---")

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
	}

	fmt.Println("批量获取视频信息:")
	var videoInfos []*aiview.VideoInfo
	for i, bv := range bvList {
		fmt.Printf("[%d/%d] 获取 %s ... ", i+1, len(bvList), bv)
		info, err := biliClient.GetVideoInfo(bv)
		if err != nil {
			fmt.Printf("失败: %v\n", err)
			continue
		}
		fmt.Printf("成功: %s\n", info.Title)
		videoInfos = append(videoInfos, info)
		time.Sleep(time.Millisecond * 500) // 避免请求过快
	}

	fmt.Printf("\n成功获取 %d/%d 个视频信息\n", len(videoInfos), len(bvList))
}

// 数据导出为 JSON - 演示如何将数据导出到文件
func exportJSONExample() {
	fmt.Println("\n--- 数据导出为 JSON ---")

	client, err := aiview.New("bilibili")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	biliClient, err := client.BilibiliClient()
	if err != nil {
		log.Fatalf("获取 Bilibili 客户端失败: %v", err)
	}

	// 获取热门视频
	hotVideos, err := biliClient.GetHotVideos(1, 10)
	if err != nil {
		log.Fatalf("获取热门视频失败: %v", err)
	}

	// 导出为 JSON
	filename := "hot_videos.json"
	err = exportToJSON(hotVideos, filename)
	if err != nil {
		log.Printf("导出失败: %v", err)
	} else {
		fmt.Printf("✓ 成功导出 %d 个视频到 %s\n", len(hotVideos), filename)
	}
}

func exportToJSON(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return aiverr.APIError("export", fmt.Sprintf("编码 JSON 失败: %v", err))
	}

	return nil
}

// 多平台数据聚合 - 演示如何从多个平台获取数据并聚合
func multiPlatformAggregation() {
	fmt.Println("\n--- 多平台数据聚合 ---")

	type PlatformHot struct {
		Platform string
		Items    []aiview.HotItem
	}

	var allHot []PlatformHot

	// 1. 获取抖音热搜
	fmt.Println("获取抖音热搜...")
	dyClient, err := aiview.New("douyin")
	if err != nil {
		log.Printf("创建抖音客户端失败: %v", err)
	} else {
		dy, err := dyClient.DouyinClient()
		if err == nil {
			hotSearch, err := dy.GetHotSearch()
			if err == nil {
				allHot = append(allHot, PlatformHot{
					Platform: "抖音",
					Items:    hotSearch,
				})
				fmt.Printf("✓ 获取到 %d 条热搜\n", len(hotSearch))
			}
		}
	}

	// 2. 获取小红书热门
	fmt.Println("获取小红书热门...")
	xhsClient, err := aiview.New("xiaohongshu")
	if err != nil {
		log.Printf("创建小红书客户端失败: %v", err)
	} else {
		xhs, err := xhsClient.XiaohongshuClient()
		if err == nil {
			hotNotes, err := xhs.GetHotNotes()
			if err == nil {
				allHot = append(allHot, PlatformHot{
					Platform: "小红书",
					Items:    hotNotes,
				})
				fmt.Printf("✓ 获取到 %d 条热门\n", len(hotNotes))
			}
		}
	}

	// 3. 聚合展示
	fmt.Println("\n【热门内容聚合】")
	for _, platform := range allHot {
		fmt.Printf("\n%s 热门 (共 %d 条):\n", platform.Platform, len(platform.Items))
		for i, item := range platform.Items {
			if i >= 3 {
				fmt.Printf("... 还有 %d 条\n", len(platform.Items)-3)
				break
			}
			fmt.Printf("  %d. %s", i+1, item.Keyword)
			if item.HotValue > 0 {
				fmt.Printf(" (热度: %d)", item.HotValue)
			}
			fmt.Println()
		}
	}
}
