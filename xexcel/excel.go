package xexcel

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/xuri/excelize/v2"
	"math"
)

type Excel struct {
	excelFile *excelize.File
	keySort   [][]string
	rename    []map[string]string
	data      [][]map[string]any
}

// CreateExcelFromGdbResult GF all 查询模式快速转换
func CreateExcelFromGdbResult(data ...gdb.Result) *Excel {
	mapList := make([][]map[string]any, 0)
	for _, item := range data {
		mapList = append(mapList, item.List())
	}
	return CreateExcel(mapList...)
}

// CreateExcel 创建 Excel
func CreateExcel(data ...[]map[string]any) *Excel {
	return &Excel{
		excelFile: excelize.NewFile(),
		data:      data,
	}
}

// ReName 修改标题名称
func (t *Excel) ReName(names []map[string]string) *Excel {
	t.rename = names
	return t
}

// Sort 排序
func (t *Excel) Sort(keys [][]string) *Excel {
	t.keySort = keys
	return t
}

// WriteFile 写入到文件
func (t *Excel) WriteFile(path string) (err error) {
	// 调用处理
	if err = t.handler(); err != nil {
		return
	}
	// 写入和关闭文件
	defer func() {
		if closeErr := t.excelFile.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	if saveErr := t.excelFile.SaveAs(path); saveErr != nil {
		err = saveErr
	}
	return
}

func (t *Excel) handler() (err error) {
	// 循环处理
	for i, item := range t.data {
		if len(item) <= 0 {
			continue
		}
		var viewKey []string
		// 判断当前是否有指定顺序
		if i+1 <= len(t.keySort) {
			viewKey = t.keySort[i]
		} else {
			for k, _ := range item[0] {
				viewKey = append(viewKey, k)
			}
		}
		// 安全获取 label map
		rename := make(map[string]string, 1)
		if i+1 <= len(t.rename) {
			rename = t.rename[i]
		}
		// 获取的当前 sheet
		sheet := fmt.Sprintf("Sheet%d", i+1)
		// 处理不存在的 sheet
		index, newSheetErr := t.excelFile.NewSheet(sheet)
		if newSheetErr != nil {
			err = newSheetErr
			return
		}
		// 生成当前 sheet 所使用的列号
		sn := t.createColumnSerialNumber(item[0])
		// 循环处理内容写入
		if err = t.handleSheet(item, sn, viewKey, rename, sheet); err != nil {
			return err
		}
		t.excelFile.SetActiveSheet(index)
	}
	return
}

func (t *Excel) handleSheet(item []map[string]any, sn, viewKey []string, rename map[string]string, sheet string) (err error) {
	for k, mItem := range item {
		mItemI := 0
		for _, v := range viewKey {
			// 安全获取数据
			var data any
			if mData, ok := mItem[v]; ok {
				data = mData
			} else {
				continue
			}
			// 处理列 label
			if k == 0 {
				var label = v
				if val, ok := rename[v]; ok {
					label = val
				}
				if err = t.excelFile.SetCellValue(sheet, fmt.Sprintf("%s%d", sn[mItemI], k+1), label); err != nil {
					return err
				}
			}
			if mItemI+1 > len(sn) {
				err = errors.New("输入的数据每个成员map成员个数不一致")
				return
			}
			if err = t.excelFile.SetCellValue(sheet, fmt.Sprintf("%s%d", sn[mItemI], k+2), data); err != nil {
				return err
			}
			mItemI += 1
		}
	}
	return
}

// 创建列位置编号
func (t *Excel) createColumnSerialNumber(m map[string]any) (sn []string) {
	sn = make([]string, 0)
	// 计数
	mCount := len(m)
	n := 1
	// 获取生成序列基础值
	s := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	// 获取值
	if mCount <= len(s) {
		sn = s[0:mCount]
	} else {
		for i := 0; i < int(math.Ceil(float64(mCount)/float64(len(s)))); i++ {
			c := -1
			if i > 0 {
				c = i - 1
			}
			for b := 0; b < len(s); b++ {
				if n > mCount {
					break
				}
				if c >= 0 {
					sn = append(sn, fmt.Sprintf("%s%s", s[c], s[b]))
				} else {
					sn = append(sn, s[b])
				}
				n += 1
			}
		}
	}
	return
}
