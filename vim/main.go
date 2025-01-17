package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zgs225/youdao"
)

func init() {
	log.SetPrefix("[i] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

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

	// 针对返回的数据做处理
	if r.Translation != nil {
		// 只展示第一个翻译结果
		(Str)(newLine((*r.Translation)[0], newLineLimit)).End()
	}

}

type Str string

func (r Str) End() {
	fmt.Print(r)
	os.Exit(0)
}

// newLine 针对 string 进行换行
func newLine(translation string, newLineLimit int) string {

	// 删除翻译结尾的"。" 避免后续做多余的处理
	if strings.HasSuffix(translation, "。") {
		translation = translation[:len(translation)-3]
	}

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
	// 这里的msg 换行需要处理一下的 因为vim里面的Tab键和换行符并不是\t \n
	// 详细的 在 vim
	//		:%!xxd  查看文件十六进制
	// 		set list 查看不可见字符
	//		:help digraph-tablk  => 0x09 和0x0a
	NL := 10
	s := strings.Join(res, string(NL))
	return s
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
