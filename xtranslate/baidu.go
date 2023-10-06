package xtranslate

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"time"
)

func BaiduTranslate(BaiduConfig *BaiduConfigType, from, to, text string) ([]string, string, error) {
	if BaiduConfig == nil || BaiduConfig.Url == "" || BaiduConfig.AppId == "" || BaiduConfig.Key == "" {
		return nil, "", errors.New("百度翻译配置异常")
	}
	// 语言标记转换
	from, err := SafeLangType(from, Baidu)
	to, err = SafeLangType(to, Baidu)

	if err != nil {
		return nil, "", err
	}

	if to == "auto" {
		return nil, "", errors.New("转换后语言不能为auto")
	}

	salt := gtime.Now().UnixMilli()
	signStr := fmt.Sprintf("%s%s%d%s", BaiduConfig.AppId, text, salt, BaiduConfig.Key)
	sign, err := gmd5.EncryptString(signStr)
	// 处理MD5加密失败
	if err != nil {
		return nil, "", err
	}
	// 发起请求
	post, err := g.Client().SetTimeout(time.Duration(BaiduConfig.CurlTimeOut)*time.Millisecond).Post(gctx.New(), BaiduConfig.Url, g.Map{
		"q":     text,
		"from":  from,
		"to":    to,
		"appid": BaiduConfig.AppId,
		"salt":  salt,
		"sign":  sign,
	})
	// 处理请求失败
	if err != nil {
		return nil, "", err
	}
	// 推出函数时关闭链接
	defer func() { _ = post.Close() }()
	// 判断状态码
	if post.StatusCode != 200 {
		return nil, "", errors.New("请求失败")
	}
	// 返回的json解析
	respStr := post.ReadAllString()
	json, err := gjson.DecodeToJson(respStr)
	// 处理json错误
	if err != nil {
		return nil, "", err
	}
	// 判断获取到的数据是否正常
	if json.Get("trans_result").IsEmpty() {
		return nil, "", errors.New(fmt.Sprintf("请求数据异常 账号: %s 返回结果: %s", BaiduConfig.AppId, respStr))
	}
	// 循环获取数据
	var arr []string

	for _, v := range json.Get("trans_result").Maps() {
		arr = append(arr, gvar.New(v["dst"], true).String())
	}

	lang, err := GetYouDaoLang(json.Get("from").String(), Baidu)

	if err != nil {
		return nil, "", err
	}

	return arr, lang, nil
}
