// Code generated by "stringer -type=ItemType"; DO NOT EDIT.

package kman

import "strconv"

const _ItemType_name = "ItemTypeTopicItemTypeTerm"

var _ItemType_index = [...]uint8{0, 13, 25}

func (i ItemType) String() string {
	if i < 0 || i >= ItemType(len(_ItemType_index)-1) {
		return "ItemType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ItemType_name[_ItemType_index[i]:_ItemType_index[i+1]]
}
