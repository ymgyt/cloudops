// Code generated by "stringer -type ResourceType resource.go"; DO NOT EDIT.

package core

import "strconv"

const _ResourceType_name = "InvalidResourceS3ResourceLocalFileResource"

var _ResourceType_index = [...]uint8{0, 15, 25, 42}

func (i ResourceType) String() string {
	if i < 0 || i >= ResourceType(len(_ResourceType_index)-1) {
		return "ResourceType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ResourceType_name[_ResourceType_index[i]:_ResourceType_index[i+1]]
}
