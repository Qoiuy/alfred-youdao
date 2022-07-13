package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zgs225/youdao"
)

const (
	APPID     = ""
	APPSECRET = ""
	MAX_LEN   = 255

	UPDATECMD = "alfred-youdao:update"
)

func init() {
	log.SetPrefix("[i] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

const golangCodeComment = "//" // go代码注释关键字 替换使用

func main() {

	appID := os.Getenv("zhiyun_id")
	appKey := os.Getenv("zhiyun_key")

	if appID == "" || appKey == "" {
		(Str)("错误: 请设置有道API").End()
	}

	client := &youdao.Client{
		AppID:     appID,
		AppSecret: appKey,
	}
	agent := newAgent(client)

	q, from, to, lang, _ := parseArgs(os.Args)
	if lang {
		if err := agent.Client.SetFrom(from); err != nil {
			(Str)(fmt.Sprintf("错误: 源语言不支持[%s]", from)).End()
		}
		if err := agent.Client.SetTo(to); err != nil {
			(Str)(fmt.Sprintf("错误: 目标语言不支持[%s]", to)).End()
		}
	}

	if len(q) == 0 {
		(Str)(`未输入查询内容`).End()
	}

	r, err := agent.Query(q)
	if err != nil {
		panic(err)
	}

	// 针对返回的数据做处理 t.Translation
	for i, translation := range *r.Translation {
		(*r.Translation)[i] = newLine(translation, newLineLimit)
	}

	if r.Translation != nil {
		// 只展示第一个翻译结果
		(Str)((*r.Translation)[0]).End()
	}

}

type Str string

func (r Str) End() {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(r); err != nil {
		panic(err)
	}
	fmt.Print(b.String())
	os.Exit(0)
}

// newLine 针对 string 进行换行
func newLine(translation string, newLineLimit int) string {
	periods := strings.Split(translation, "。")
	res := make([]string, 0)
	for i := 0; i < len(periods); {
		// 默认合并一句话
		mergeNum, mergeString := merge(&periods, i, newLineLimit)
		if mergeNum == 0 {
			i++
			continue
		}

		res = append(res, mergeString)
		i = i + mergeNum
	}
	return strings.Join(res, "\n")
}

const newLineLimit int = 50

// merge 选取字符串数组里面 start开始的字符串 进行合并
// 长度不超过newLineLimit限制
// 如果有异常 返回 0 ,""
// 没有异常返回 合并使用的字符串数 和 拼接好的字符串内容
func merge(periods *[]string, start int, newLineLimit int) (int, string) {
	// 计算需要合并 字符串数组中 几个字符串
	wordCount := 0
	if start < 0 || start >= len(*periods) {
		return wordCount, ""
	}

	wordNumber := 0
	for i := start; i < len(*periods); i++ {
		if wordNumber > newLineLimit {
			break
		}
		wordNumber += len((*periods)[i])
		wordCount++
	}

	return wordCount, "// " + strings.Join((*periods)[start:start+wordCount], " ")

}
