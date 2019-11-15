package builtin

import (
	"net"
	"strconv"

	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("ServerIP()", func() Value {
	host, port, _ := net.SplitHostPort(options.NetAddr)
	if host == "" && port != "" {
		host = "127.0.0.1"
	}
	return SuStr(host)
})

var _ = builtin0("ServerPort()", func() Value {
	_, port, _ := net.SplitHostPort(options.NetAddr)
	if port == "" {
		return EmptyStr
	}
	n, _ := strconv.Atoi(port)
	return IntVal(n)
})

var _ = builtin0("Server?()", func() Value {
	return False
})
