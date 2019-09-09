package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

var _ = builtin2("WideCharToMultiByte(string, cp = 1252)",
	func(s, c Value) Value {
		// UTF-16 to 1252 or UTF-8
		utf16 := ToStr(s)
		if strings.HasSuffix(utf16, "\x00\x00") {
			utf16 = utf16[:len(utf16)-2]
		}
		utf16decode := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
		utf8, err := utf16decode.String(utf16)
		if err != nil {
			panic("WideCharToMultiByte " + err.Error())
		}
		if codePage(c) == cp_utf8 {
			return SuStr(utf8)
		}
		encoder := charmap.Windows1252.NewEncoder()
		s1252, err := encoder.String(utf8)
		if err != nil {
			panic("WideCharToMultiByte 1252 " + err.Error())
		}
		return SuStr(s1252)
	})

var _ = builtin2("MultiByteToWideChar(string, cp = 1252)",
	func(str, cp Value) Value {
		// 1252 or UTF-8 to UTF-16
		s := ToStr(str)
		if codePage(cp) == 1252 {
			decoder := charmap.Windows1252.NewDecoder()
			utf8, err := decoder.String(s)
			if err != nil {
				panic("MultiByteToWideChar 1252 " + err.Error())
			}
			s = utf8
		}
		// s is now UTF-8
		utf16encode := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
		utf16, err := utf16encode.String(s)
		if err != nil {
			panic("MultiByteToWideChar " + err.Error())
		}
		return SuStr(utf16 + "\x00\x00")
	})

const (
	acp        = 0
	thread_acp = 3
	cp_utf8    = 65001
)

func codePage(cp Value) int {
	switch ToInt(cp) {
	case acp, thread_acp, 1252:
		return 1252
	case cp_utf8:
		return cp_utf8
	default:
		panic("invalid code page")
	}
}
