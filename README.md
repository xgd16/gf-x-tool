# 用于 GF 的工具集
    
### 1. 快速返回
    // 返回成功
    tool.CreateFastResponse(r).SuccessStatus().Response()
    // 返回失败
	tool.CreateFastResponse(r).ErrorStatus().Response()
        