package wmptool

import (
	"fmt"
)

func PrintBasicAssetAttribute(raw_data map[string][]interface{}) (groupId string, assetId string, taskId string, schedulerId string) {
	var data = raw_data

	if v, ok := data["group_id"]; ok {
		for _, vv := range v {
			if vvv, ok := vv.(string); ok {
				groupId = vvv
			}
		}
	}

	if v, ok := data["asset_id"]; ok {
		for _, vv := range v {
			if vvv, ok := vv.(string); ok {
				assetId = vvv
			}
		}
	}

	if v, ok := data["task_id"]; ok {
		for _, vv := range v {
			if vvv, ok := vv.(string); ok {
				taskId = vvv
			}
		}
	}

	if v, ok := data["scheduler_id"]; ok {
		for _, vv := range v {
			if vvv, ok := vv.(string); ok {
				schedulerId = vvv
			}
		}
	}

	return
}

func checkAndAddDataToWMPResponse(data map[string][]interface{}, key string, value []interface{}) {
	if v, ok := data[key]; ok {
		data[key] = removeDuplicates2(append(v, value...))
	} else {
		data[key] = removeDuplicates2(value)
	}
}

func stringItemsToInterfaceItems(list []string) []interface{} {
	if list == nil {
		return nil
	}
	var out []interface{}
	for _, v := range list {
		out = append(out, v)
	}
	return out
}

func removeDuplicates(strings []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, str := range strings {
		if _, ok := seen[str]; !ok && str != "" {
			seen[str] = struct{}{}
			result = append(result, str)
		}
	}

	return result
}

func removeDuplicates2(input []interface{}) []interface{} {
	uniqueMap := make(map[string]struct{})
	var result []interface{}

	for _, item := range input {
		itemStr := fmt.Sprintf("%v", item)
		if _, exists := uniqueMap[itemStr]; !exists {
			uniqueMap[itemStr] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}
