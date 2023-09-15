package translate

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/xgd16/gf-x-tool/xTool"
	"math"
	"time"
)

// YouDaoTranslate 有道翻译
func YouDaoTranslate(YouDaoConfig *YouDaoConfigType, from, to, text string) ([]string, string, error) {
	if YouDaoConfig == nil || YouDaoConfig.AppKey == "" || YouDaoConfig.Url == "" || YouDaoConfig.SecKey == "" {
		return nil, "", errors.New("有道翻译配置异常")
	}
	// 语言标记转换
	from, err := SafeLangType(from, YouDao)
	// to, err = SafeLangType(to, YouDao)

	if err != nil {
		return nil, "", err
	}

	if to == "auto" {
		return nil, "", errors.New("转换后语言不能为auto")
	}

	truncate := func(s string) string {
		l := gstr.LenRune(s)

		if l <= 20 {
			return s
		}

		return fmt.Sprintf("%s%d%s", gstr.SubStrRune(s, 0, 10), l, gstr.SubStrRune(s, l-10, l))
	}

	salt := gtime.Now().UnixMilli()
	curTime := int(math.Round(float64(salt / 1000)))
	signStr := fmt.Sprintf("%s%s%d%d%s", YouDaoConfig.AppKey, truncate(text), salt, curTime, YouDaoConfig.SecKey)
	sign := xTool.Sha256(signStr)

	post, err := g.Client().SetTimeout(time.Duration(YouDaoConfig.CurlTimeOut)*time.Millisecond).Post(gctx.New(), YouDaoConfig.Url, g.Map{
		"q":        text,
		"appKey":   YouDaoConfig.AppKey,
		"salt":     salt,
		"from":     from,
		"to":       to,
		"sign":     sign,
		"signType": "v3",
		"curtime":  curTime,
	})

	if err != nil {
		return nil, "", err
	}

	defer func() { _ = post.Close() }()

	if post.StatusCode != 200 {
		return nil, "", errors.New("请求失败")
	}

	postResp := post.ReadAllString()

	json, err := gjson.DecodeToJson(postResp)

	if TranslateDebug {
		g.Log().Info(gctx.New(), YouDao, signStr, g.Map{
			"q":        text,
			"appKey":   YouDaoConfig.AppKey,
			"salt":     salt,
			"from":     from,
			"to":       to,
			"sign":     sign,
			"signType": "v3",
			"curtime":  curTime,
		}, "输出: ", postResp)
	}

	if err != nil {
		return nil, "", err
	}

	if json.Get("errorCode").Int() != 0 {
		return nil, "", errors.New(fmt.Sprintf("请求失败errorCode: %d err: %s", json.Get("errorCode").Int(), postResp))
	}
	// 获取 from
	returnFrom := gstr.Split(json.Get("l").String(), "2")[0]
	return json.Get("translation").Strings(), returnFrom, nil
}
