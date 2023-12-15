package xtranslate

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

var TranslateDebug = func() bool {

	get, err := g.Cfg().Get(gctx.New(), "server.translateDebug", false)

	if err != nil {
		panic(fmt.Sprintf("翻译初始化失败 %s", err))
	}

	return get.Bool()
}()

const (
	YouDao = "YouDao"
	Baidu  = "Baidu"
	Google = "Google"
	Deepl  = "Deepl"
)

// BaiduConfigType 百度的配置类型
type BaiduConfigType struct {
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
	AppId       string `json:"appId"`
	Key         string `json:"key"`
}

// YouDaoConfigType 有道配置类型
type YouDaoConfigType struct {
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
	AppKey      string `json:"appKey"`
	SecKey      string `json:"secKey"`
}

// GoogleConfigType 谷歌配置类型
type GoogleConfigType struct {
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
	Key         string `json:"key"`
}

// DeeplConfigType Deepl配置类型
type DeeplConfigType struct {
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
	Key         string `json:"key"`
}

// 需要初始化的平台配置
var initPlatformConf = []string{
	YouDao,
	Baidu,
	Google,
	Deepl,
}

// InitTranslateBaseConf 初始化翻译基础配置
var InitTranslateBaseConf = func() map[string]map[string]string {
	translate := gfile.GetContents("./translate.json")

	if translate == "" {
		return nil
	}

	json, err := gjson.DecodeToJson(translate)

	if err != nil {
		return nil
	}

	m := make(map[string]map[string]string, 1)

	for s, v := range json.Var().MapStrVar() {
		m[s] = v.MapStrStr()
	}

	return m
}

// BaseTranslateConf 基础翻译配置
var BaseTranslateConf map[string]map[string]string

// BasePlatformTranslateConf 基础平台翻译配置
var BasePlatformTranslateConf map[string][]map[string]*gvar.Var

// InitPlatformConfFunc 初始化平台配置
var InitPlatformConfFunc = initPlatformConfHandler

func InitTranslate() {
	// 初始化基本配置
	BaseTranslateConf = InitTranslateBaseConf()
	// 初始化配置
	BasePlatformTranslateConf = InitPlatformConfFunc()
}

func initPlatformConfHandler() map[string][]map[string]*gvar.Var {
	cfgMap := make(map[string][]map[string]*gvar.Var, 1)

	for _, a := range initPlatformConf {
		// 获取配置
		cfgData, err := getConf(a)
		// 处理配置错误
		if err != nil {
			continue
		}
		// 循环处理配置进行解析
		for _, datum := range cfgData {
			m := make(map[string]*gvar.Var, 1)
			d := datum.MapStrVar()
			// 转换为map数据
			for s, v := range d {
				m[s] = v
			}
			// 将数据写入到配置
			cfgMap[a] = append(cfgMap[a], m)
		}
	}

	return cfgMap
}

func SafeLangType(t, app string) (string, error) {
	if t == "auto" {
		return "auto", nil
	}

	a := BaseTranslateConf[app]

	if a == nil {
		return "", errors.New("没有找到应用")
	}

	l := a[t]

	if l == "" {
		return "", errors.New("不支持的语言类型")
	}

	return l, nil
}

func GetYouDaoLang(lang, app string) (string, error) {
	if lang == "auto" {
		return "auto", nil
	}

	a := BaseTranslateConf[app]

	if a == nil {
		return "", errors.New("没有找到应用")
	}

	for s, s2 := range a {
		if s2 == lang {
			return s, nil
		}
	}

	return lang, nil
}

// TranslationResData 翻译数据
type TranslationResData struct {
	From string `json:"from"` // 原文语言
	To   string `json:"to"`   // 翻译后语言
	Text string `json:"text"` // 原文
}

type TranslationRespData struct {
	Result  []string `json:"result"`  // 翻译后内容
	From    string   `json:"from"`    // 实际原文语言
	FromLen int      `json:"fromLen"` // 翻译原文字数
	ToLen   []int    `json:"toLen"`   // 翻译后文字字数
	To      string   `json:"to"`      // 翻译后语言
	Tool    string   `json:"tool"`    // 翻译工具
}

// Translation 翻译
func Translation(t *TranslationResData, cfg ...map[string][]map[string]*gvar.Var) (*TranslationRespData, error) {
	// 需要一个随机顺序的翻译工具 (利用map无序机制)
	mPlatformConf := make(map[int]string, 1)
	// 将 arr 写入到 map
	for i, s := range initPlatformConf {
		mPlatformConf[i] = s
	}
	// 循环调用进行翻译
	var (
		translate []string
		from      string
		err       error
	)

	for _, s := range mPlatformConf {
		// 获取配置
		var i []map[string]*gvar.Var

		if len(cfg) > 0 {
			i = cfg[0][s]
		} else {
			i = BasePlatformTranslateConf[s]
		}
		// 判断配置不存在的话跳过
		if len(i) <= 0 {
			continue
		}
		// 调用翻译
		for _, cfg := range i {
			switch s {
			case Baidu:
				translate, from, err = BaiduTranslate(&BaiduConfigType{
					CurlTimeOut: cfg["curlTimeOut"].Int(),
					Url:         cfg["url"].String(),
					AppId:       cfg["appId"].String(),
					Key:         cfg["key"].String(),
				}, t.From, t.To, t.Text)
				break
			case YouDao:
				translate, from, err = YouDaoTranslate(&YouDaoConfigType{
					CurlTimeOut: cfg["curlTimeOut"].Int(),
					Url:         cfg["url"].String(),
					AppKey:      cfg["appKey"].String(),
					SecKey:      cfg["secKey"].String(),
				}, t.From, t.To, t.Text)
				break
			case Google:
				translate, from, err = GoogleTranslate(&GoogleConfigType{
					CurlTimeOut: cfg["curlTimeOut"].Int(),
					Url:         cfg["url"].String(),
					Key:         cfg["key"].String(),
				}, t.From, t.To, t.Text)
				break
			case Deepl:
				translate, from, err = DeeplTranslate(&DeeplConfigType{
					CurlTimeOut: cfg["curlTimeOut"].Int(),
					Url:         cfg["url"].String(),
					Key:         cfg["key"].String(),
				}, t.From, t.To, t.Text)
				break
			}
			// 处理翻译失败 翻译失败时跳过本次自动执行下次翻译
			if err != nil {
				continue
			}
			// 获取翻译后字数长度
			var toLen []int

			for _, s2 := range translate {
				toLen = append(toLen, gstr.LenRune(s2))
			}
			// 返回翻译内容
			return &TranslationRespData{
				Result:  translate,
				From:    from,
				FromLen: gstr.LenRune(t.From),
				ToLen:   toLen,
				To:      t.To,
				Tool:    s,
			}, nil
		}
	}
	// 返回结果
	return nil, errors.New("翻译失败")
}

// 获取配置
func getConf(cfgName string) ([]*gvar.Var, error) {
	// 获取百度配置
	cfg, err := g.Cfg().Get(gctx.New(), cfgName)
	// 处理获取配置错误
	if err != nil {
		return nil, err
	}
	// 转换为组数据
	return cfg.Vars(), nil
}
