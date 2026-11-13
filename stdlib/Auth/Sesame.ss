// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Gen(s)
		{
		log = Object()
		a = #Test
		b = #Temp
		r = ""
		bytes = ""
		h = Sha256()
		t = a $ b
		size = 32
		for c in s
			switch c
				{
			case 'r':
				r = Random()
				Print(:r)
			case 's':
				Random.Seed(Adler32(bytes))
				Print("seed:", Adler32(bytes), "<-", bytes.ToHex())
			case 'b':
				bytes = Random.Bytes(size)
			case 'h':
				h.Update(bytes)
				Print(hash: bytes.ToHex())
			case 'T':
				Print(#create, t, ':', b $ a $ r, '=', bytes.ToHex())
				x = Object()
				x[b $ a $ r] = bytes
				log[t] = x
			case 'L':
				t = b $ a $ r
			case 'H':
				bytes = h.Value()
			case ' ':
				}
		return log
		}

	Verfy(log)
		{
		h = Sha256()
		bytes = log.TestTemp.TempTest
		h.Update(bytes)
		Print(hash: bytes.ToHex())
		Random.Seed(Adler32(bytes))
		Print("seed:", Adler32(bytes), "<-", bytes.ToHex())
		for i in ..2
			{
			r1 = Random()
			r2 = Random()
			bytes = log["TempTest" $ r1]["TempTest" $ r2]
			h.Update(bytes)
			Print(hash: bytes.ToHex())
			}
		r1 = Random()
		r2 = Random()
		bytes = log["TempTest" $ r1]["TempTest" $ r2]
		return bytes is h.Value()
		}
	}
