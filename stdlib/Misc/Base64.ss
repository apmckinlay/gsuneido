// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
class
	{
	encodedChunkSize: 4
	decodedChunkSize: 3
	Encode(s)
		{
		return s.MapN(.decodedChunkSize, .encode)
		}
	encode(next3)
		{
		enc = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
		n = (next3[0].Asc() << 16) | (next3[1].Asc() << 8) | next3[2].Asc()
		return enc[n >> 18] $
			enc[(n >> 12) & 0x3f] $
			(next3.Size() > 1 ? enc[(n >> 6) & 0x3f] : '=') $
			(next3.Size() > 2 ? enc[n & 0x3f] : '=')
		}

	EncodeLines(src, eol = '\r\n', linelen = 70)
		{
		src = .Encode(src)
		for (dst = ""; src isnt ""; src = src[linelen..])
			dst $= src[.. linelen] $ eol
		return dst
		}

	Decode(base64)
		{
		base64 = base64.Tr('\r\n').RightTrim('=')
		return base64.MapN(.encodedChunkSize, .decode)
		}
	decode(next4)
		{
		dec = #(y: 50, x: 49, w: 48, v: 47, "9": 61, u: 46, "8": 60, t: 45, "7": 59,
			s: 44, "6": 58, r: 43, "5": 57, q: 42, "4": 56, p: 41, "3": 55, o: 40,
			"2": 54, n: 39, "1": 53, m: 38, "0": 52, "/": 63, l: 37, k: 36, j: 35,
			i: 34, "+": 62, h: 33, g: 32, f: 31, e: 30, d: 29, c: 28, b: 27, a: 26,
			Z: 25, Y: 24, X: 23, W: 22, V: 21, U: 20, T: 19, S: 18, R: 17, Q: 16,
			P: 15, O: 14, N: 13, M: 12, L: 11, K: 10, J: 9, I: 8, H: 7, G: 6, F: 5,
			E: 4, D: 3, C: 2, B: 1, A: 0, z: 51, '': 0)
		n = (dec[next4[0]] << 18) | (dec[next4[1]] << 12) |
			(dec[next4[2]] << 6) | dec[next4[3 /* = 4th encoded char*/]]
		return (n >> 16).Chr() $
			(next4[2] is '' ? '' : ((n >> 8) & 0xff).Chr()) $
			(next4[3 /* = 4th encoded char*/] is '' ? '' : (n & 0xff).Chr())
		}
	}