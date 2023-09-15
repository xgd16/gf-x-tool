# 用于 [GF](https://goframe.org/pages/viewpage.action?pageId=1114119) 的工具集
<font size="2" color=#006666>西安豆芽科技有限公司 **制**</font>
## Go mod 拉取

```shell
go get -u github.com/xgd16/gf-x-tool
```

## no required module provides package github.com/gogf/gf/v2/os/gmutex 问题处理
```shell
go get github.com/gogf/gf/v2/os/gmutex@v2.5.1
```

> GrayLog日志扩展 xTool > logger.go

```go
package main

// 使用 GrayLog 演示
func main() {
  // 配置 GrayLog 基础配置 host 和 port
  xTool.SetGrayLogConfig("127.0.0.1", 9803)
  // 配置默认日志
  glog.SetDefaultHandler(xTool.SwitchToGraylog("测试日志"))
  // 配置自定义日志
  g.Log("test").SetHandlers(xTool.SwitchToGraylog("测试02"))
  // 至此接入 GrayLog 完成 日志写法按照 GF 规范编写即可
  g.Log().Error(gctx.New(), "测试")
}
```

> 简略文本数据存储 xTool > xdb.go

```go
// 创建一个全局对象
var xdb = xTool.CreateXDB()

xdb.Set("user", "name", "lz") // 写入
xdb.Get("user", "name") // 获取 返回 *gver.Ver 对象
xdb.Del("user", "name") // 删除 存入的数据 支持 xdb.Del("user", "name", "age") 一次对key 下的多个 field 进行删除
xdb.GetGJson() // 获取 *gjson.Json 对象
xdb.GetJsonStr() // 获取所有数据的 json 字符串
```

> HTTP 快速返回 xTool > resp.go

```go
package main

import (
    "error"
    "errors"
    "fmt"
    "github.com/xgd16/gf-x-tool/xTool"
)

func main () {
    // 默认返回自定义 
    xTool.FastResponseOption = &FastResponseOptionType{
        DefaultSuccessCode:       1000,
        DefaultErrorCode:         1001,
        DefaultSuccessStatusCode: 200,
        DefaultErrorStatusCode:   500,
        MsgName:                  "msg",
        CodeName:                 "code",
        DataName:                 "data",
    }
    // 根据状态 获取 对应 message 实现
    xTool.GetMessageToStatus = func(int status) string {
        return map[int]string{
            1002: "参数错误",
            1001: "请求失败",
        }[status]
    }
    // xTool.FastResp(r).ErrorStatus().SetStatus(1002).Resp() msg 返回 参数错误
}

func test(r *ghttp.Request) {
  // .FastResp() 介绍
  // 参数 1 *ghttp.Request
  // 参数 2 支持类型 error|bool 当传入 err 时会自动判断是否有错误如果有 则触发返回并且自动记录日志 如果是 bool 型 true 会触发返回
  // 参数 3 类型 bool true 返回 SuccessStatus 状态数据 false 返回 ErrorStatus 数据
  // .Resp() 介绍
  // 参数 1 控制返回的 msg
  // 参数 2 控制返回的 code
  // 参数 3 控制返回的 data
    // 返回成功 {"code": 1000, "msg": "SUCCESS !!!", "data": []}
  xTool.FastResp(r).SuccessStatus().Resp()
  // 返回 data {"code": 1000, "msg": "SUCCESS !!!", "data": {"a": 1}}
  xTool.FastResp(r).SetData(g.Map{
    "a": 1,
  }).Resp(); // 默认 SuccessStatus
  // 也可
  xTool.FastResp(r).Resp(nil, nil, g.Map{
    "a": 1,
  });
  // 返回失败 {"code": 1001, "msg": "ERROR !!!", "data": []}
    xTool.FastResp(r).ErrorStatus().Resp()
  // 智能返回 状态判断
    err := errors.New("测试")
    xTool.FastResp(r, err != nil).Resp() // 注意 当第二项参数 结果为 true 时 返回 SuccessStatus
    xTool.FastResp(r, err).Resp() // 当 第二个参数为 error 类型时会自动判断是否存在 error 如果存在状态 ErrorStatus 并自动记录到日志    
    xTool.FastResp(r, err, true).Resp() // 此操作时返回 SuccessStatus 参数为 false 时为 ErrorStatus 第三个参数 控制结果状态 优先级高于 自动状态     // 自动事务结束    tx, _ := g.DB().Begin()
    xTool.FastResp(r, err).TxRollBack(tx).Resp() // 当第二个参数 为 true 时自动触发 RollBack    // 返回时调用
    xTool.FastResp(r, err).Callback(func (t *xTool.FastResponse) {
        fmt.Println("test")
    }).Resp() // test 输出 会在 返回前触发
}
```

> Excel 创建写入

```go
xExcel.CreateExcel([]map[string]any{
        {
            "a": "测试1",
        },
        {
            "a": "测试2",
        },
    }, []map[string]any{
        {
            "a": "测试3",
        },
        {
            "a": "测试4",
        },
}).ReName([]map[string]string{{"a": "名称"}, {"a": "昵称"}}).WriteFile(fmt.Sprintf("./test-%s.xlsx", gtime.Now().Format("Y_m_d_H_i_s")))
```



> 在线翻译支持 translate
>
> >百度 *BaiduTranslate*
> >
> >谷歌 *GoogleTranslate*
> >
> >有道 *YouDaoTranslate*
> >
> >Deepl *DeeplTranslate*

```go
package main

import (
    "github.com/xgd16/gf-x-tool/translate"
)

func main () {
    // 配置获取方式自定义
    translate.InitPlatformConfFunc = func() map[string][]map[string]*gvar.Var {}
    // 扩展提供的统一翻译封装
    translate.Translation(&translate.TranslationResData{
        From: "auto",
        To: "en",
        Text: "测试",
    })
}
```
