# 用于 [GF](https://goframe.org/pages/viewpage.action?pageId=1114119) 的工具集

## go mod
    go get -u github.com/xgd16/gf-x-tool
## 如果不能拉取,需要配置此代理
    go env -w GOPRIVATE="github.com"

### 1. 快速返回
    // 返回成功
    xTool.CreateFastResponse(r).SuccessStatus().Response()
    // 返回失败
	xTool.CreateFastResponse(r).ErrorStatus().Response()

    // 可以直接在快速返回内执行判断
    xTool.CreateFastResponse(r, err != nil).ErrorStatus().Response("发送失败请稍后再试")
    
        