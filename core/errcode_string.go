// Code generated by "stringer -type ErrCode errors.go"; DO NOT EDIT.

package core

import "strconv"

const _ErrCode_name = "OKInvalidParamConflictTimeoutInternalExternalNotFoundUnauthorizedUnauthenticatedRateLimitNotImplementedYetCanceledUndefined"

var _ErrCode_index = [...]uint8{0, 2, 14, 22, 29, 37, 45, 53, 65, 80, 89, 106, 114, 123}

func (i ErrCode) String() string {
	if i < 0 || i >= ErrCode(len(_ErrCode_index)-1) {
		return "ErrCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrCode_name[_ErrCode_index[i]:_ErrCode_index[i+1]]
}
