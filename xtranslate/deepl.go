package xtranslate

import (
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func DeeplTranslate(config *DeeplConfigType, from, to, text string) (result []string, fromLang string, err error) {
	if config == nil || config.Url == "" || config.Key == "" {
		err = errors.New("deepl翻译配置异常")
		return
	}
	ctx := gctx.New()
	// 语言标记转换
	from, err = SafeLangType(from, Deepl)
	if err != nil {
		return
	}
	to, err = SafeLangType(to, Deepl)
	if err != nil {
		return
	}
	if from == "auto" {
		from = ""
	}
	// 调用翻译
	HttpResult, err := g.Client().SetTimeout(time.Duration(config.CurlTimeOut)*time.Millisecond).Header(g.MapStrStr{
		"Authorization": fmt.Sprintf("DeepL-Auth-Key %s", config.Key),
	}).Post(ctx, config.Url, g.Map{
		"text":        text,
		"source_lang": from,
		"target_lang": to,
	})
	// 处理调用接口错误
	if err != nil {
		return
	}
	// 推出函数时关闭链接
	defer func() { _ = HttpResult.Close() }()
	// 判断状态码
	respStr := HttpResult.ReadAllString()
	if HttpResult.StatusCode != 200 {
		err = fmt.Errorf("请求失败 状态码: %d 返回结果: %s", HttpResult.StatusCode, respStr)
		return
	}
	// 返回的json解析
	json, err := gjson.DecodeToJson(respStr)
	if err != nil {
		return
	}
	// 获取源语言
	dsl := json.Get("translations.0.detected_source_language")
	if dsl.IsEmpty() {
		fromLang = from
	} else {
		fromLang = dsl.String()
	}
	// 返回翻译结果
	tr := json.Get("translations.0.text")
	if tr.IsEmpty() {
		err = errors.New("翻译失败请重试 " + respStr)
		return
	} else {
		result = tr.Strings()
	}
	// 将语言种类转换为有道标准
	fromLang, err = GetYouDaoLang(fromLang, Deepl)
	return
}
