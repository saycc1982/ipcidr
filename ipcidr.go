package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 国家代码到名称的映射表（已移除通常没有独立IP分配的8个国家）
var countryNames = map[string]string{
	"AD": "安道尔", "AE": "阿拉伯联合酋长国", "AF": "阿富汗", "AG": "安提瓜和巴布达",
	"AI": "安圭拉", "AL": "阿尔巴尼亚", "AM": "亚美尼亚", "AO": "安哥拉",
	"AR": "阿根廷", "AS": "美属萨摩亚", "AT": "奥地利", "AU": "澳大利亚",
	"AW": "阿鲁巴", "AX": "奥兰群岛", "AZ": "阿塞拜疆", "BA": "波斯尼亚和黑塞哥维那",
	"BB": "巴巴多斯", "BD": "孟加拉国", "BE": "比利时", "BF": "布基纳法索",
	"BG": "保加利亚", "BH": "巴林", "BI": "布隆迪", "BJ": "贝宁",
	"BL": "圣巴泰勒米", "BM": "百慕大", "BN": "文莱", "BO": "玻利维亚",
	"BQ": "博奈尔、圣尤斯特歇斯和萨巴", "BR": "巴西", "BS": "巴哈马", "BT": "不丹",
	"BW": "博茨瓦纳", "BY": "白俄罗斯", "BZ": "伯利兹", "CA": "加拿大",
	"CC": "科科斯（基林）群岛", "CD": "刚果民主共和国", "CF": "中非共和国", "CG": "刚果共和国",
	"CH": "瑞士", "CI": "科特迪瓦", "CK": "库克群岛", "CL": "智利",
	"CM": "喀麦隆", "CN": "中国", "CO": "哥伦比亚", "CR": "哥斯达黎加",
	"CU": "古巴", "CV": "佛得角", "CW": "库拉索", "CX": "圣诞岛",
	"CY": "塞浦路斯", "CZ": "捷克共和国", "DE": "德国", "DJ": "吉布提",
	"DK": "丹麦", "DM": "多米尼克", "DO": "多米尼加共和国", "DZ": "阿尔及利亚",
	"EC": "厄瓜多尔", "EE": "爱沙尼亚", "EG": "埃及", "EH": "西撒哈拉",
	"ER": "厄立特里亚", "ES": "西班牙", "ET": "埃塞俄比亚", "FI": "芬兰",
	"FJ": "斐济", "FK": "福克兰群岛", "FM": "密克罗尼西亚联邦", "FO": "法罗群岛",
	"FR": "法国", "GA": "加蓬", "GB": "英国", "GD": "格林纳达",
	"GE": "格鲁吉亚", "GF": "法属圭亚那", "GG": "根西岛", "GH": "加纳",
	"GI": "直布罗陀", "GL": "格陵兰", "GM": "冈比亚", "GN": "几内亚",
	"GP": "瓜德罗普", "GQ": "赤道几内亚", "GR": "希腊", "GS": "南乔治亚和南桑威奇群岛",
	"GT": "危地马拉", "GU": "关岛", "GW": "几内亚比绍", "GY": "圭亚那",
	"HK": "中国香港", "HM": "赫德岛和麦克唐纳群岛", "HN": "洪都拉斯", "HR": "克罗地亚",
	"HT": "海地", "HU": "匈牙利", "ID": "印度尼西亚", "IE": "爱尔兰",
	"IL": "以色列", "IM": "马恩岛", "IN": "印度", "IO": "英属印度洋领地",
	"IQ": "伊拉克", "IR": "伊朗", "IS": "冰岛", "IT": "意大利",
	"JE": "泽西岛", "JM": "牙买加", "JO": "约旦", "JP": "日本",
	"KE": "肯尼亚", "KG": "吉尔吉斯斯坦", "KH": "柬埔寨", "KI": "基里巴斯",
	"KM": "科摩罗", "KN": "圣基茨和尼维斯", "KP": "朝鲜", "KR": "韩国",
	"KW": "科威特", "KY": "开曼群岛", "KZ": "哈萨克斯坦", "LA": "老挝",
	"LB": "黎巴嫩", "LC": "圣卢西亚", "LI": "列支敦士登", "LK": "斯里兰卡",
	"LR": "利比里亚", "LS": "莱索托", "LT": "立陶宛", "LU": "卢森堡",
	"LV": "拉脱维亚", "LY": "利比亚", "MA": "摩洛哥", "MC": "摩纳哥",
	"MD": "摩尔多瓦", "ME": "黑山", "MF": "圣马丁", "MG": "马达加斯加",
	"MH": "马绍尔群岛", "MK": "北马其顿", "ML": "马里", "MM": "缅甸",
	"MN": "蒙古", "MO": "中国澳门", "MP": "北马里亚纳群岛", "MQ": "马提尼克",
	"MR": "毛里塔尼亚", "MS": "蒙特塞拉特", "MT": "马耳他", "MU": "毛里求斯",
	"MV": "马尔代夫", "MW": "马拉维", "MX": "墨西哥", "MY": "马来西亚",
	"MZ": "莫桑比克", "NA": "纳米比亚", "NC": "新喀里多尼亚", "NE": "尼日尔",
	"NF": "诺福克岛", "NG": "尼日利亚", "NI": "尼加拉瓜", "NL": "荷兰",
	"NO": "挪威", "NP": "尼泊尔", "NR": "瑙鲁", "NU": "纽埃",
	"NZ": "新西兰", "OM": "阿曼", "PA": "巴拿马", "PE": "秘鲁",
	"PF": "法属波利尼西亚", "PG": "巴布亚新几内亚", "PH": "菲律宾", "PK": "巴基斯坦",
	"PL": "波兰", "PM": "圣皮埃尔和密克隆", "PN": "皮特凯恩群岛", "PR": "波多黎各",
	"PS": "巴勒斯坦", "PT": "葡萄牙", "PW": "帕劳", "PY": "巴拉圭",
	"QA": "卡塔尔", "RE": "留尼汪", "RO": "罗马尼亚", "RS": "塞尔维亚",
	"RU": "俄罗斯", "RW": "卢旺达", "SA": "沙特阿拉伯", "SB": "所罗门群岛",
	"SC": "塞舌尔", "SD": "苏丹", "SE": "瑞典", "SG": "新加坡",
	"SH": "圣赫勒拿", "SI": "斯洛文尼亚", "SJ": "斯瓦尔巴群岛和扬马延岛", "SK": "斯洛伐克",
	"SL": "塞拉利昂", "SM": "圣马力诺", "SN": "塞内加尔", "SO": "索马里",
	"SR": "苏里南", "SS": "南苏丹", "ST": "圣多美和普林西比", "SV": "萨尔瓦多",
	"SX": "荷属圣马丁", "SY": "叙利亚", "SZ": "斯威士兰", "TC": "特克斯和凯科斯群岛",
	"TD": "乍得", "TF": "法属南部领地", "TG": "多哥", "TH": "泰国",
	"TJ": "塔吉克斯坦", "TK": "托克劳", "TL": "东帝汶", "TM": "土库曼斯坦",
	"TN": "突尼斯", "TO": "汤加", "TR": "土耳其", "TT": "特立尼达和多巴哥",
	"TV": "图瓦卢", "TW": "中国台湾", "TZ": "坦桑尼亚", "UA": "乌克兰",
	"UG": "乌干达", "UM": "美国本土外小岛屿", "US": "美国", "UY": "乌拉圭",
	"UZ": "乌兹别克斯坦", "VA": "梵蒂冈", "VC": "圣文森特和格林纳丁斯", "VE": "委内瑞拉",
	"VG": "英属维尔京群岛", "VI": "美属维尔京群岛", "VN": "越南", "VU": "瓦努阿图",
	"WF": "瓦利斯和富图纳", "WS": "萨摩亚", "YE": "也门", "YT": "马约特",
	"ZA": "南非", "ZM": "赞比亚", "ZW": "津巴布韦",
}

// 被移除的8个国家（通常没有独立IP分配）
var removedCountries = map[string]string{
	"AQ": "南极洲",     // 无永久居民，无需IP分配
	"BV": "布韦岛",     // 无人居住的岛屿
	"IO": "英属印度洋领地", // 人口极少，无独立IP需求
	"NR": "瑙鲁",      // 微型国家，IP资源依赖其他国家
	"PN": "皮特凯恩群岛", // 人口不足50人
	"TK": "托克劳",     // 人口约1500人，互联网极不发达
	"UM": "美国本土外小岛屿", // 由美国管理，IP归属美国
	"VA": "梵蒂冈",     // 面积太小，IP资源依赖意大利
}

// RIR数据源配置
var rirConfigs = []struct {
	name       string
	fixedURL   string
	isAfrinic  bool
	maxRetries int
}{
	{name: "APNIC", fixedURL: "https://ftp.apnic.net/stats/apnic/delegated-apnic-latest", isAfrinic: false, maxRetries: 0},
	{name: "ARIN", fixedURL: "https://ftp.arin.net/pub/stats/arin/delegated-arin-extended-latest", isAfrinic: false, maxRetries: 0},
	{name: "RIPE NCC", fixedURL: "https://ftp.ripe.net/ripe/stats/delegated-ripencc-latest", isAfrinic: false, maxRetries: 0},
	{name: "AFRINIC", fixedURL: "https://ftp.afrinic.net/pub/stats/afrinic/", isAfrinic: true, maxRetries: 30},
	{name: "LACNIC", fixedURL: "https://ftp.lacnic.net/pub/stats/lacnic/delegated-lacnic-latest", isAfrinic: false, maxRetries: 0},
}

// 显示帮助信息
func printHelp() {
	fmt.Println("IP地址范围生成器与国家代码查询工具")
	fmt.Println("使用方法:")
	fmt.Println("  查看帮助:")
	fmt.Println("    go run ipcidr.go -h")
	fmt.Println("  查询国家代码:")
	fmt.Println("    查看单个国家: go run ipcidr.go -name CN")
	fmt.Println("    查看所有国家(默认每行5个): go run ipcidr.go -name all")
	fmt.Println("    自定义每行显示数量: go run ipcidr.go -name all 3")
	fmt.Println("  更新IP数据:")
	fmt.Println("    go run ipcidr.go -update")
	fmt.Println("\n参数说明:")
	fmt.Println("  -h        显示帮助信息")
	fmt.Println("  -name     国家代码查询 (支持单个代码或 'all')")
	fmt.Println("  -update   下载并生成最新IP国家文件")
}

// 打印底部签名信息
func printSignature() {
	fmt.Printf(`

___________________________________

如有建议或 BUG 反馈, 欢迎联络作者
作 者: ED
联 络: https://t.me/hongkongisp
服务器推荐: 顺安云 https://www.say.cc
`)
}

// 获取所有国家代码并排序
func getAllCountryCodes() []string {
	codes := make([]string, 0, len(countryNames))
	for code := range countryNames {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes
}

// 显示国家代码信息，包括已移除的国家列表
func displayCountryInfo(nameParam string, lineCount int) {
	allCodes := getAllCountryCodes()
	targetCode := strings.ToUpper(nameParam)
	
	// 显示单个国家
	if targetCode != "" && targetCode != "ALL" {
		// 检查是否是已移除的国家
		if name, exists := removedCountries[targetCode]; exists {
			fmt.Printf("%s %s (已移除，无独立IP分配)\n", targetCode, name)
			return
		}
		
		name, exists := countryNames[targetCode]
		if exists {
			fmt.Printf("%s %s\n", targetCode, name)
		} else {
			fmt.Printf("未找到国家代码: %s\n", targetCode)
		}
		return
	}
	
	// 显示所有国家，自定义每行数量
	fmt.Printf("所有国家/地区代码 (共 %d 个，每行 %d 个):\n", len(allCodes), lineCount)
	count := 0
	for _, code := range allCodes {
		name := countryNames[code]
		// 格式化输出，保持对齐
		fmt.Printf("%-3s %-18s", code, name)
		count++
		if count%lineCount == 0 {
			fmt.Println()
		}
	}
	// 确保最后一行换行
	if count%lineCount != 0 {
		fmt.Println()
	}
	
	// 显示已移除的国家
	fmt.Printf("\n已移除的国家/地区 (共 %d 个，无独立IP分配):\n", len(removedCountries))
	removedCodes := make([]string, 0, len(removedCountries))
	for code := range removedCountries {
		removedCodes = append(removedCodes, code)
	}
	sort.Strings(removedCodes)
	
	for i, code := range removedCodes {
		fmt.Printf("%-3s %-18s", code, removedCountries[code])
		if (i+1)%5 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

// 下载文件并保存到临时目录
func downloadFile(url string, tempDir string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("下载失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
	}

	parts := strings.Split(url, "/")
	filename := parts[len(parts)-1]
	filePath := filepath.Join(tempDir, filename)

	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer out.Close()

	var reader io.Reader
	if strings.HasSuffix(filename, ".gz") {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("gzip解压失败: %v", err)
		}
		defer gzReader.Close()
		reader = gzReader
	} else {
		reader = resp.Body
	}

	_, err = io.Copy(out, reader)
	if err != nil {
		return "", fmt.Errorf("写入文件失败: %v", err)
	}

	return filePath, nil
}

// 生成AFRINIC的URL并尝试下载
func downloadAfrinicFile(tempDir string, maxRetries int) (string, error) {
	currentDate := time.Now()
	
	for i := 0; i <= maxRetries; i++ {
		year := currentDate.Format("2006")
		dateStr := currentDate.Format("20060102")
		url := fmt.Sprintf("%s%s/delegated-afrinic-extended-%s", 
			rirConfigs[3].fixedURL, year, dateStr)
		
		fmt.Printf("尝试下载AFRINIC数据: %s (尝试 %d/%d)\n", url, i+1, maxRetries+1)
		
		filePath, err := downloadFile(url, tempDir)
		if err == nil {
			return filePath, nil
		}
		
		if strings.Contains(err.Error(), "HTTP状态码错误") {
			currentDate = currentDate.AddDate(0, 0, -1)
			continue
		}
		
		return "", err
	}
	
	fallbackURL := "http://ftp.afrinic.net/pub/stats/afrinic/RIR-Statistics-Exchange-Format.txt"
	fmt.Printf("所有日期尝试失败，尝试备选URL: %s\n", fallbackURL)
	return downloadFile(fallbackURL, tempDir)
}

// 提取所有国家代码，同时过滤已移除的国家
func extractCountryCodes(files []string) ([]string, error) {
	countryCodes := make(map[string]bool)
	countryRegex := regexp.MustCompile(`^[A-Z]{2}$`)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.Split(line, "|")
			if len(parts) >= 2 {
				cc := parts[1]
				// 检查是否是有效的国家代码且未被移除
				if countryRegex.MatchString(cc) && countryNames[cc] != "" {
					countryCodes[cc] = true
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	ccList := make([]string, 0, len(countryCodes))
	for cc := range countryCodes {
		ccList = append(ccList, cc)
	}
	sort.Strings(ccList)

	return ccList, nil
}

// 计算IPv4的CIDR前缀长度
func calculateIPv4Prefix(size uint64) int {
	if size == 0 {
		return 32
	}
	prefix := 32 - int(math.Log2(float64(size)))
	if prefix < 0 {
		return 0
	}
	return prefix
}

// 处理国家的IP地址范围
func processCountryIP(cc string, files []string, dataDir string, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	// 检查是否是已移除的国家，如果是则跳过处理
	if removedCountries[cc] != "" {
		return
	}

	countryDir := filepath.Join(dataDir, strings.ToLower(cc))
	if err := os.MkdirAll(countryDir, 0755); err != nil {
		errChan <- fmt.Errorf("创建目录 %s 失败: %v", countryDir, err)
		return
	}

	ipv4File, err := os.Create(filepath.Join(countryDir, "ipv4.txt"))
	if err != nil {
		errChan <- fmt.Errorf("创建IPv4文件失败: %v", err)
		return
	}
	defer ipv4File.Close()

	ipv6File, err := os.Create(filepath.Join(countryDir, "ipv6.txt"))
	if err != nil {
		errChan <- fmt.Errorf("创建IPv6文件失败: %v", err)
		return
	}
	defer ipv6File.Close()

	currentTime := time.Now().Format(time.RFC3339)
	// 写入带#号的最后更新时间
	if _, err := ipv4File.WriteString(fmt.Sprintf("# last updated: %s\n", currentTime)); err != nil {
		errChan <- err
		return
	}
	if _, err := ipv6File.WriteString(fmt.Sprintf("# last updated: %s\n", currentTime)); err != nil {
		errChan <- err
		return
	}

	// 在最后更新时间下方添加带#号的作者签名信息
	signature := `# 如有建议或 BUG 反馈, 欢迎联络作者
# 作 者: ED
# 联 络: https://t.me/hongkongisp
# 服务器推荐: 顺安云 https://www.say.cc
`
	if _, err := ipv4File.WriteString(signature); err != nil {
		errChan <- err
		return
	}
	if _, err := ipv6File.WriteString(signature); err != nil {
		errChan <- err
		return
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			errChan <- err
			return
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.Split(line, "|")
			if len(parts) < 5 || parts[1] != cc {
				continue
			}

			ipType := parts[2]
			ipStart := parts[3]
			sizeStr := parts[4]

			size, err := strconv.ParseUint(sizeStr, 10, 64)
			if err != nil {
				continue
			}

			switch ipType {
			case "ipv4":
				prefix := calculateIPv4Prefix(size)
				if _, err := ipv4File.WriteString(fmt.Sprintf("%s/%d\n", ipStart, prefix)); err != nil {
					f.Close()
					errChan <- err
					return
				}
			case "ipv6":
				if _, err := ipv6File.WriteString(fmt.Sprintf("%s/%d\n", ipStart, size)); err != nil {
					f.Close()
					errChan <- err
					return
				}
			}
		}

		f.Close()
		if err := scanner.Err(); err != nil {
			errChan <- err
			return
		}
	}
}

// 执行IP数据更新
func updateIPData() {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "ipdata-")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建数据目录
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("创建数据目录失败: %v", err)
	}

	// 下载所有RIR数据文件
	var downloadedFiles []string
	fmt.Println("开始下载RIR数据文件...")
	
	for _, config := range rirConfigs {
		var filePath string
		var err error
		
		if config.isAfrinic {
			filePath, err = downloadAfrinicFile(tempDir, config.maxRetries)
		} else {
			fmt.Printf("下载 %s 数据...\n", config.name)
			filePath, err = downloadFile(config.fixedURL, tempDir)
		}
		
		if err != nil {
			log.Printf("警告: 下载 %s 失败: %v", config.name, err)
			if config.isAfrinic {
				log.Printf("已忽略AFRINIC数据，继续处理其他RIR...")
			}
			continue
		}
		
		downloadedFiles = append(downloadedFiles, filePath)
		fmt.Printf("成功下载 %s 数据\n", config.name)
	}

	if len(downloadedFiles) == 0 {
		log.Fatal("没有成功下载任何数据文件，无法继续")
	}

	// 提取所有国家代码
	fmt.Println("提取国家代码...")
	countryCodes, err := extractCountryCodes(downloadedFiles)
	if err != nil {
		log.Fatalf("提取国家代码失败: %v", err)
	}

	// 处理每个国家的IP地址
	var wg sync.WaitGroup
	errChan := make(chan error, len(countryCodes))
	maxWorkers := 10
	semaphore := make(chan struct{}, maxWorkers)

	fmt.Println("开始处理IP地址范围...")
	for i, cc := range countryCodes {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(cc string, index int, total int) {
			defer func() { <-semaphore }()
			fmt.Printf("处理 %s (%d/%d)\n", cc, index+1, total)
			processCountryIP(cc, downloadedFiles, dataDir, &wg, errChan)
		}(cc, i, len(countryCodes))
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		fmt.Printf("处理过程中出现 %d 个错误:\n", len(errors))
		for _, e := range errors {
			fmt.Printf(" - %v\n", e)
		}
	}

	fmt.Println("所有操作已完成")
	fmt.Printf("结果已保存到 %s 目录\n", dataDir)
	fmt.Printf("共生成 %d 个国家/地区的IP资料夹\n", len(countryCodes))
}

func main() {
	// 解析命令行参数
	helpFlag := flag.Bool("h", false, "显示帮助信息")
	nameParam := flag.String("name", "", "国家代码查询 (支持单个代码或 'all')")
	updateFlag := flag.Bool("update", false, "下载并生成最新IP国家文件")
	flag.Parse()

	// 处理帮助请求
	if *helpFlag {
		printHelp()
		printSignature()
		return
	}

	// 处理国家代码查询
	if *nameParam != "" {
		// 解析自定义每行显示数量（默认5个）
		lineCount := 5
		args := flag.Args()
		if len(args) > 0 {
			if num, err := strconv.Atoi(args[0]); err == nil && num > 0 {
				lineCount = num
			} else {
				fmt.Printf("无效的行数参数: %s，使用默认值 %d\n", args[0], lineCount)
			}
		}
		displayCountryInfo(*nameParam, lineCount)
		printSignature()
		return
	}

	// 处理IP数据更新
	if *updateFlag {
		updateIPData()
		printSignature()
		return
	}

	// 未提供任何参数时显示帮助
	printHelp()
	printSignature()
}
    
