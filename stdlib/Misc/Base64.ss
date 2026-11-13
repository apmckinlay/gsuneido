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
		mask = 0x3f
		enc = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
		n = (next3[0].Asc() << 16) | (next3[1].Asc() << 8) | next3[2].Asc()
		return enc[n >> 18] $ enc[(n >> 12) & mask] $
			(next3.Size() > 1 ? enc[(n >> 6) & mask] : '=') $
			(next3.Size() > 2 ? enc[n & mask] : '=')
		}

	// Microsoft Exchange server has issue with line length that is not divisible by 4
	// and could end up with file corrupt
	// should be also consistent with Base64BufferEncodeFile.lineLength
	EncodeLines(src, eol = "\r\n", linelen = 72)
		{
		src = .Encode(src)
		for (dst = ""; src isnt ""; src = src[linelen..])
			dst $= src[..linelen] $ eol
		return dst
		}

	EncodeUrl(s) // rfc 4648, safe in URLs, handled by Decode
		{
		return s.MapN(.decodedChunkSize, .encodeUrl)
		}

	encodeUrl(next3)
		{
		mask = 0x3f
		enc = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
		n = (next3[0].Asc() << 16) | (next3[1].Asc() << 8) | next3[2].Asc()
		return enc[n >> 18] $ enc[(n >> 12) & mask] $
			(next3.Size() > 1 ? enc[(n >> 6) & mask] : "") $
			(next3.Size() > 2 ? enc[n & mask] : "")
		}

	Decode(base64)
		{
		base64 = base64.Tr("\r\n").RightTrim('=')
		return base64.MapN(.encodedChunkSize, .decode)
		}

	decode(next4)
		{
		mask = 0xff
		dec = #("": 0, A: 0, B: 1, C: 2, D: 3, E: 4, F: 5, G: 6, H: 7, I: 8, J: 9, K: 10,
			L: 11, M: 12, N: 13, O: 14, P: 15, Q: 16, R: 17, S: 18, T: 19, U: 20, V: 21,
			W: 22, X: 23, Y: 24, Z: 25, a: 26, b: 27, c: 28, d: 29, e: 30, f: 31, g: 32,
			h: 33, i: 34, j: 35, k: 36, l: 37, m: 38, n: 39, o: 40, p: 41, q: 42, r: 43,
			s: 44, t: 45, u: 46, v: 47, w: 48, x: 49, y: 50, z: 51, '0': 52, '1': 53,
			'2': 54, '3': 55, '4': 56, '5': 57, '6': 58, '7': 59, '8': 60, '9': 61,
			'-': 62, '+': 62, '/': 63, '_': 63)
		n = (dec[next4[0]] << 18) | (dec[next4[1]] << 12) | (dec[next4[2]] << 6) |
			dec[next4[3 /* = 4th encoded char*/]]
		return (n >> 16).Chr() $ (next4[2] is "" ? "" : ((n >> 8) & mask).Chr()) $
			(next4[3 /* = 4th encoded char*/] is "" ? "" : (n & mask).Chr())
		}
	}
