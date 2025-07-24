// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
function (before, after)
	// pre: before and after are strings
	// where the first four bytes are from IntToStr
	// and the rest is a binary fraction stored in a string
	{
	if after is false
		return IntToStr(StrToInt(before) + 1)  // before + one

	// return (before + after) / 2
	numBytes = 4
	wholeSize = 256
	b = before[numBytes ..]
	a = before[.. numBytes] is after[.. numBytes]
		? after[numBytes ..]
		: '\xff'.Repeat(b.Size() + 1)
	ob = Object()
	n = Max(b.Size(), a.Size())
	// first add them
	carry = 0
	for (i = n - 1; i >= 0; --i)
		{
		ob[i] = b[i].Asc() + a[i].Asc() + carry
		if ob[i] >= wholeSize
			{
			carry = 1
			ob[i] -= wholeSize
			}
		else
			carry = 0
		}
	// then divide by 2 (like right shift)
	for (i = 0; i < n; ++i)
		{
		ob[i] += carry * wholeSize
		carry = ob[i] % 2
		ob[i] = (ob[i] / 2).Int()
		}
	if carry is 1
		ob.Add(0x80)
	// then convert back to binary string
	s = ''
	for n in ob
		s $= n.Chr()
	return before[.. numBytes] $ s
	}