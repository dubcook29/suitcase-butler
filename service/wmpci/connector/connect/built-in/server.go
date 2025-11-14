package connect_builtin

import (
	"context"

	"github.com/suitcase/butler/wmpci"
)

type BuiltinServer interface {
	WMPService(ctx context.Context, request wmpci.WMPRequest) (wmpci.WMPResponse, error)
	WMPConfig(custom map[string]wmpci.WMPCustom) (bool, error)
	WMPRegist() (BuiltinServer, wmpci.WMPRegistrars)
}
