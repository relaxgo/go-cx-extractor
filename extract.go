package extractor

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

const (
	blocksWidth  int = 4
	FREQUENT_URL int = 30
)

var (
	MaxThreshold int = 100

// = regexp.MustCompile("<[aA]\\s+[Hh][Rr][Ee][Ff]="+
// 	"[\"|\']?([^>\"\' ]+)[\"|\']?\\s*[^>]*>([^>]+)</a>(\\s*.{0,"+
// 		FREQUENT_URL+"}\\s*<a\\s+href=[\"|\']?([^>\"\' ]+)[\"|\']?\\s*[^>]*>"+
// 		"([^>]+)</[aA]>){2,100}")
)

func preProcess(source string) string {
	source = regexp.MustCompile("(?is)<!DOCTYPE.*?>").ReplaceAllString(source, "")
	source = regexp.MustCompile("(?is)<!--.*?-->").ReplaceAllString(source, "")              // remove html comment
	source = regexp.MustCompile("(?is)<script.*?>.*?</script>").ReplaceAllString(source, "") // remove javascript
	source = regexp.MustCompile("(?is)<style.*?>.*?</style>").ReplaceAllString(source, "")   // remove css
	source = regexp.MustCompile("&.{2,5};|&#.{2,5};").ReplaceAllString(source, " ")          // remove special char

	//剔除连续成片的超链接文本（认为是，广告或噪音）,超链接多藏于span中
	// source = regexp.MustCompile("<[sS][pP][aA][nN].*?>").ReplaceAllString(source, "")
	// source = regexp.MustCompile("</[sS][pP][aA][nN]>").ReplaceAllString(source, "")

	// int len = len(source)
	// for source=linksReg.ReplaceAllString(source, "");len!=len(source){
	// 	len = len(source)
	// }

	//防止html中在<>中包括大于号的判断
	// source = regexp.MustCompile("<[^>'\"]*['\"].*['\"].*?>").ReplaceAllString(source, "")
	source = regexp.MustCompile("\r\n").ReplaceAllString(source, "\n")
	source = regexp.MustCompile("(?is)<.*?>").ReplaceAllString(source, "")
	source = regexp.MustCompile("(?is)<.*?>").ReplaceAllString(source, "")

	return source
}

func ExtactTitle(source string) string {
	t := regexp.MustCompile("(?is)<title>(.*?)</title>").FindString(source)
	return regexp.MustCompile("</?title>").ReplaceAllString(t, "")
}

func ExtractText(source string) string {
	var text string = ""
	source = preProcess(source)
	// fmt.Println(source)
	lines := strings.Split(source, "\n")

	if len(lines) <= blocksWidth {
		return ""
	}

	blocks, threshold := getIndexDistribution(lines)

	// for _, line := range lines {
	// 	fmt.Println(line)
	// }

	// for _, bb := range blocks {
	// 	fmt.Println(bb)
	// }

	fmt.Println("threshold", threshold)
	// fmt.Println(blocks)

	// threshold = 5
	isStart := false
	for i, block := range blocks {
		//块断开
		if block == 0 {
			isStart = false
			continue
		}

		if block > threshold && i < len(blocks)-3 {
			if blocks[i+1] > 0 && blocks[i+2] > 0 && blocks[i+3] > 0 {
				isStart = true
			}
		}

		if isStart {
			text = text + "\n" + lines[i]
		}
	}
	return text
}

/*
param lines
ruturn  indexDistribution,threshold
*/
func getIndexDistribution(lines []string) ([]int, int) {
	blocks := len(lines) - blocksWidth + 1
	blockLengths := make([]int, blocks, blocks)

	cBlockLen, cLineLen, sum, empty := 0, 0, 0, 0

	reg := regexp.MustCompile("\\s+")
	for i, line := range lines {
		line = reg.ReplaceAllString(line, "")
		lines[i] = line
		cLineLen = len(line)

		//计算初始行块长度
		if i < (blocksWidth - 1) {
			cBlockLen += cLineLen
			continue
		}
		// 计算每个行块，加上当前块，去掉最早的一块
		index := i - blocksWidth + 1
		cBlockLen = cBlockLen + cLineLen

		blockLengths[index] = cBlockLen
		//去掉过时的行
		cBlockLen = cBlockLen - len(lines[index])

		if cLineLen == 0 {
			empty++
		}
		sum += cLineLen
	}

	threshold := int(1.2 * math.Min(float64(MaxThreshold), float64(sum/(blocks-empty))))

	return blockLengths, threshold
}
