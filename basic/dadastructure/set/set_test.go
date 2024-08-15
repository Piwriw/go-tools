package utils

import (
	"strings"
	"testing"
)

func TestSet(t *testing.T) {
	set := NewSet()
	sql := "7D6B57BC75255DF2FCF1603BB64EC6AA"
	conid := "1003"
	b1 := QueryBody{SqlID: "7D6B57BC75255DF2FCF1603BB64EC6AA", ConID: "10002"}
	b2 := QueryBody{SqlID: "A86446692DE3CA32943798A09019B333", ConID: "10002"}
	b3 := QueryBody{SqlID: "A86446692DE3CA32943798A09019B333", ConID: "3"}
	b4 := QueryBody{SqlID: "A86446692DE3CA32943798A09019B333", ConID: "3"}
	b5 := QueryBody{SqlID: sql, ConID: conid}

	set.Add(b1, b2, b3, b4, b5)
	set.Add(b5)
	elements, err := set.SqlIDElements()
	if err != nil {
		t.Error(err)
	}
	t.Log(elements)
	var sb strings.Builder
	strArr := elements["10002"]
	// 循环遍历数组
	for i, s := range strArr {
		// 在每个字符串两边加上单引号
		sb.WriteString("'")
		sb.WriteString(s)
		sb.WriteString("'")
		// 如果不是最后一个字符串，则加上逗号和空格
		if i < len(strArr)-1 {
			sb.WriteString(", ")
		}
	}
	t.Log("str", sb.String())

	//sprintf := fmt.Sprintf("SELECT query_sql FROM oceanbase.GV$OB_PLAN_CACHE_PLAN_STAT WHERE TENANT_ID=%d  AND sql_id IN (%s) ORDER BY type", strings.Join(set.StringElements(), ","))
	//fmt.Println(sprintf)
}
