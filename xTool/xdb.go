package xTool

import (
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gfile"
	"sync"
	"sync/atomic"
)

// CreateXDB 创建数据库对象
func CreateXDB() *XDB {
	return (&XDB{
		DBPath: "./XDB",
		DBFile: "xdb.json",
	}).init()
}

type XDB struct {
	DBPath    string `json:"dbPath"`
	DBFile    string `json:"dbFile"`
	cache     map[string]map[string]any
	cacheJson *gjson.Json
	writeNum  atomic.Int64
	mutex     *sync.Mutex
}

func (x *XDB) init() *XDB {

	content := gfile.GetContents(x.getDbFilePath())

	json, _ := gjson.DecodeToJson(content)

	m := make(map[string]map[string]any)

	for k, v := range json.Map() {
		m[k] = v.(map[string]any)
	}

	x.cache = m
	x.cacheJson = json
	x.mutex = new(sync.Mutex)

	return x
}

// Get 获取指定数据
func (x *XDB) Get(key, field string) *gvar.Var {
	x.mutex.Lock()
	defer x.mutex.Unlock()

	return x.cacheJson.Get(fmt.Sprintf("%s.%s", key, field))
}

// GetGJson 获取 gf json 操作对象
func (x *XDB) GetGJson() *gjson.Json {
	x.mutex.Lock()
	defer x.mutex.Unlock()

	return x.cacheJson
}

// GetJsonStr 获取存储的json字符串
func (x *XDB) GetJsonStr() string {
	x.mutex.Lock()
	defer x.mutex.Unlock()

	return x.cacheJson.MustToJsonString()
}

func (x *XDB) SetList(key, field string, data any) error {
	x.writeNum.Add(1)

	x.mutex.Lock()

	x.writeNum.Add(-1)

	defer x.mutex.Unlock()
	// 判断数据是否为空
	if x.cache == nil {
		x.cache = make(map[string]map[string]any)
	}

	if x.cache[key] == nil {
		x.cache = map[string]map[string]any{key: {field: []any{data}}}
	} else {
		a := x.cache[key][field].([]any)

		a = append(a, data)

		x.cache[key][field] = a
	}

	return x.save()
}

// Set 新增操作
func (x *XDB) Set(key, field string, data any) error {
	x.writeNum.Add(1)

	x.mutex.Lock()

	x.writeNum.Add(-1)

	defer x.mutex.Unlock()
	// 判断数据是否为空
	if x.cache == nil {
		x.cache = make(map[string]map[string]any)
	}

	if x.cache[key] == nil {
		x.cache = map[string]map[string]any{key: {field: data}}
	} else {
		x.cache[key][field] = data
	}

	return x.save()
}

// Del 删除操作
func (x *XDB) Del(key string, field ...string) error {
	x.writeNum.Add(1)

	x.mutex.Lock()
	x.writeNum.Add(-1)
	defer x.mutex.Unlock()

	if len(field) <= 0 {
		delete(x.cache, key)
	} else {
		for _, v := range field {
			delete(x.cache[key], v)
		}
	}

	return x.save()
}

// DelList 删除列表操作
func (x *XDB) DelList(key, field string, index ...int) error {
	x.writeNum.Add(1)

	x.mutex.Lock()
	x.writeNum.Add(-1)
	defer x.mutex.Unlock()

	if len(field) <= 0 {
		delete(x.cache, key)
	} else {
		for _, v := range index {
			a := x.cache[key][field].([]any)

			aArr := garray.NewArrayFrom(a)

			aArr.Remove(v)

			x.cache[key][field] = aArr.Interfaces()
		}
	}

	return x.save()
}

// 执行存储操作
func (x *XDB) save() (err error) {
	// 判断是否需要执行写入操作 (防止 再多线程操作下 高频进行 io 写入操作)
	if x.writeNum.Load() > 0 {
		return nil
	}
	// 将数据转换为json数据
	json, err := gjson.EncodeString(x.cache)
	if err != nil {
		return
	}
	// 将数据写入到文件
	if err := gfile.PutContents(x.getDbFilePath(), json); err != nil {
		return err
	}
	// 将新数据进行解析
	jsonData, err := gjson.DecodeToJson(json)
	if err != nil {
		return err
	}

	x.cacheJson = jsonData

	return nil
}

// 获取库文件路径
func (x *XDB) getDbFilePath() string {
	return fmt.Sprintf("%s/%s", x.DBPath, x.DBFile)
}
