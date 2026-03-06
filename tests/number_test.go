// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

var _ = Register("number bitwise operations", `
'0xffffffff | 1', 0xffffffff
'0xffffffff ^ 1', 0xfffffffe
'0xffffffff >> 1', 0x7fffffff
'1 << 31', 0x80000000
'0x80000000 >> 31', 1
`)

var _ = Register("number number.Format", `
'0.Format("###")', "'0'"
'0.Format("###.")', "'0.'"
'0.Format("#.##")', "'.00'"
'1e-5.Format("#.##")', "'.00'"
'123456.Format("###")', "'#'"
'123456.Format("###.###")', "'#'"
'123.Format("###")', "'123'"
'1.Format("###")', "'1'"
'10.Format("###")', "'10'"
'100.Format("###")', "'100'"
'.08.Format("#.##")', "'.08'"
'.18.Format("#.#")', "'.2'"
'.08.Format("#.#")', "'.1'"
'6.789.Format("#.##")', "'6.79'"
'123.Format("##")', "'#'"
'(0-1).Format("#.##")', "'-'"
'(0-12).Format("-####")', "'-12'"
'(0-12).Format("(####)")', "'(12)'"
'123.Format("###,###")', "'123'"
'1234.Format("###,###")', "'1,234'"
'12345.Format("###,###")', "'12,345'"
'123456.Format("###,###")', "'123,456'"
'1e5.Format("###,###")', "'100,000'"
'.8.Format("Foo")', "'#'"
`)

var _ = Register("number number.Round", `
'0.Round(0)', 0
'0.Round(2)', 0
'0.Round(-2)', 0
'1e-9.Round(5)', 0
'1.Round(0)', 1
'12.34.Round(0)', 12
'12.56.Round(0)', 13
'12.34.Round(1)', 12.3
'12.56.Round(1)', 12.6
'12.34.Round(2)', 12.34
'12.56.Round(2)', 12.56
'12.34.Round(3)', 12.34
'12.56.Round(3)', 12.56
'1234.Round(-1)', 1230
'1256.Round(-1)', 1260
'1234.Round(-2)', 1200
'1256.Round(-2)', 1300
'.33333333.Round(2)', .33
'.66666666.Round(2)', .67
'(-.33333333).Round(2)', -.33
'(-.66666666).Round(2)', -.67
'.5.Round(0)', 1
'9999.Round(-4)', 10000
'123456789012.34.Round(2)', 123456789012.34
'123456789012345.Round(-2)', 123456789012300
'125.Round(-1)', 130
'.09.Round(0)', 0
`)

var _ = Register("number number.RoundUp", `
'.33333333.RoundUp(2)', .34
'.66666666.RoundUp(2)', .67
'(-.33333333).RoundUp(2)', -.34
'(-.66666666).RoundUp(2)', -.67
'(-23.00).RoundUp(0)', -23
'(-23.35).RoundUp(0)', -24
'(-23.50).RoundUp(0)', -24
'(-23.67).RoundUp(0)', -24
'(-999.9999999).RoundUp(0)', -1000
'123456789012.11.RoundUp(1)', 123456789012.2
'123456789012311.RoundUp(-2)', 123456789012400
'23.67.RoundUp(0)', 24
'23.50.RoundUp(0)', 24
'23.35.RoundUp(0)', 24
'23.00.RoundUp(0)', 23
'0.RoundUp(0)', 0
'23.50.RoundUp(1)', 23.5
'23.35.RoundUp(1)', 23.4
'23.00.RoundUp(1)', 23.0
'23.50.RoundUp(-1)', 30
'23.35.RoundUp(-1)', 30
'23.00.RoundUp(-1)', 30
'999.9999999.RoundUp(0)', 1000
`)

var _ = Register("number number.RoundDown", `
'.33333333.RoundDown(2)', .33
'.66666666.RoundDown(2)', .66
'(-.33333333).RoundDown(2)', -.33
'(-.66666666).RoundDown(2)', -.66
'(-23.00).RoundDown(0)', -23
'(-23.35).RoundDown(0)', -23
'(-23.50).RoundDown(0)', -23
'(-23.67).RoundDown(0)', -23
'(-999.9999999).RoundDown(0)', -999
'123456789012.99.RoundDown(1)', 123456789012.9
'123456789012399.RoundDown(-2)', 123456789012300
'23.67.RoundDown(0)', 23
'23.50.RoundDown(0)', 23
'23.35.RoundDown(0)', 23
'23.00.RoundDown(0)', 23
'0.RoundDown(0)', 0
'23.50.RoundDown(1)', 23.5
'23.35.RoundDown(1)', 23.3
'23.00.RoundDown(1)', 23.0
'23.50.RoundDown(-1)', 20
'23.35.RoundDown(-1)', 20
'23.00.RoundDown(-1)', 20
'999.9999999.RoundDown(0)', 999
`)

var _ = Register("number number.Frac", `
'0.Frac()', 0
'123.Frac()', 0
'(-123).Frac()', 0
'.123.Frac()', .123
'(-.123).Frac()', -.123
'123.456.Frac()', .456
'(-123.456).Frac()', -.456
'.001.Frac()', .001
'100.002.Frac()', .002
`)

var _ = Register("number number.Int", `
'0.Int()', 0
'123.Int()', 123
'(-123).Int()', -123
'123.456.Int()', 123
'(-123.456).Int()', -123
'.123.Int()', 0
'(-.123).Int()', 0
'1e-5.Int()', 0
'(-1e-5).Int()', 0
`)

var _ = Register("number number.Hex", `
'255.Hex()', "'ff'"
'32768.Hex()', "'8000'"
'65535.Hex()', "'ffff'"
'(0-1).Hex()', "'ffffffffffffffff'"
'4294967295.Hex()', "'ffffffff'"
'0xffffffff.Hex()', "'ffffffff'"
'2864434397.Hex()', "'aabbccdd'"
'0x80000000.Hex()', "'80000000'"
`)
