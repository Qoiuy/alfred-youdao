package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zgs225/alfred-youdao/alfred"
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
	log.Println(os.Args)

	items := alfred.NewResult()

	appID := os.Getenv("zhiyun_id")
	appKey := os.Getenv("zhiyun_key")

	if appID == "" || appKey == "" {
		items.Append(&alfred.ResultElement{
			Valid:    false,
			Title:    "错误: 请设置有道API",
			Subtitle: "有道词典",
		})
		items.End()
	}

	client := &youdao.Client{
		AppID:     appID,
		AppSecret: appKey,
	}
	agent := newAgent(client)

	q, from, to, lang, isGolangComment := parseArgs(os.Args)
	if lang {
		if err := agent.Client.SetFrom(from); err != nil {
			items.Append(&alfred.ResultElement{
				Valid:    true,
				Title:    fmt.Sprintf("错误: 源语言不支持[%s]", from),
				Subtitle: `有道词典`,
			})
			items.End()
		}
		if err := agent.Client.SetTo(to); err != nil {
			items.Append(&alfred.ResultElement{
				Valid:    true,
				Title:    fmt.Sprintf("错误: 目标语言不支持[%s]", to),
				Subtitle: `有道词典`,
			})
			items.End()
		}
	}

	if len(q) == 0 {
		items.Append(&alfred.ResultElement{
			Valid:    true,
			Title:    "有道词典",
			Subtitle: `查看"..."的解释或翻译`,
		})
		items.End()
	}

	// 普通的小查询限制字符数 golang注释翻译不限制
	if !isGolangComment && len(q) > 255 {
		items.Append(&alfred.ResultElement{
			Valid:    false,
			Title:    "错误: 最大查询字符数为255",
			Subtitle: q,
		})
		items.End()
	}

	r, err := agent.Query(q)
	if err != nil {
		panic(err)
	}

	// 针对返回的数据做处理 t.Translation
	for i, translation := range *r.Translation {
		(*r.Translation)[i] = newLine(translation, newLineLimit)
	}

	mod := map[string]*alfred.ModElement{
		alfred.Mods_Shift: &alfred.ModElement{
			Valid:    true,
			Arg:      toYoudaoDictUrl(q),
			Subtitle: "回车键打开词典网页",
		},
	}
	if r.Basic != nil {
		phonetic := joinPhonetic(r.Basic.Phonetic, r.Basic.UkPhonetic, r.Basic.UsPhonetic)
		for _, title := range r.Basic.Explains {
			mod2 := copyModElementMap(mod)
			mod2[alfred.Mods_Cmd] = &alfred.ModElement{
				Valid:    true,
				Arg:      wordsToSayCmdOption(title, r),
				Subtitle: "发音",
			}
			item := alfred.ResultElement{
				Valid:    true,
				Title:    title,
				Subtitle: phonetic,
				Arg:      title,
				Mods:     mod2,
			}
			items.Append(&item)
		}
	}

	if r.Translation != nil {
		title := strings.Join(*r.Translation, "; ")
		mod2 := copyModElementMap(mod)
		mod2[alfred.Mods_Cmd] = &alfred.ModElement{
			Valid:    true,
			Arg:      wordsToSayCmdOption(title, r),
			Subtitle: "发音",
		}
		item := alfred.ResultElement{
			Valid:    true,
			Title:    title,
			Subtitle: "翻译结果",
			Arg:      title,
			Mods:     mod2,
		}
		items.Append(&item)
	}

	if r.Web != nil {
		items.Append(&alfred.ResultElement{
			Valid:    true,
			Title:    "网络释义",
			Subtitle: "有道词典 for Alfred",
		})

		for _, elem := range *r.Web {
			mod2 := copyModElementMap(mod)
			mod2[alfred.Mods_Cmd] = &alfred.ModElement{
				Valid:    true,
				Arg:      wordsToSayCmdOption(elem.Key, r),
				Subtitle: "发音",
			}
			items.Append(&alfred.ResultElement{
				Valid:    true,
				Title:    elem.Key,
				Subtitle: strings.Join(elem.Value, "; "),
				Arg:      elem.Key,
				Mods:     mod,
			})
		}
	}

	if agent.Dirty {
		if err := agent.Cache.SaveFile(CACHE_FILE); err != nil {
			log.Println(err)
		}
	}

	items.End()
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

	return wordCount, strings.Join((*periods)[start:start+wordCount], " ")

}
