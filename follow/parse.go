package follow

import (
	"encoding/json"

	"smartconn.cc/tosone/parsing-go/logstash"
)

// Parse 序列化脚本
func Parse(script []byte) (global GlobalSetting, allPages WholePages, err error) {
	var rawPages Pages
	allPages = make(map[string]FollowItem)
	if err = json.Unmarshal(script, &rawPages); err != nil {
		return
	}
	for key, val := range rawPages {
		if key == "globalSetting" {
			if err = json.Unmarshal(val, &global); err != nil {
				return
			}
		} else {
			var page FollowSetting
			if err = json.Unmarshal(val, &page); err != nil {
				return
			}
			var item = FollowItem{}
			item.Follows = make(map[string]FollowActon)
			for k, v := range page {
				if k == "noRecord" && (string(v) == "true" || string(v) == "false") {
					if string(v) == "true" {
						item.NoRecord = true
					} else if string(v) == "false" {
						item.NoRecord = false
					}
				} else if k == "order" {
					var orderList []string
					if err = json.Unmarshal(v, &orderList); err != nil {
						logstash.Error(err.Error())
					}
					item.Order = orderList
				} else if k == "innerOrder" {
					var innerOrderList map[string][]string
					if err = json.Unmarshal(v, &innerOrderList); err != nil {
						logstash.Error(err.Error())
					}
					item.InnerOrder = innerOrderList
				} else {
					var followItem = FollowActon{}
					if err = json.Unmarshal(v, &followItem); err != nil {
						return
					}
					item.Follows[k] = followItem
				}
			}
			allPages[key] = item
		}
	}
	return
}
