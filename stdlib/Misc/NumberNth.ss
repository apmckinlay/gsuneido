// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
function (number)
	{
	special_th_numbers = #(11,12,13)
	if special_th_numbers.Has?(number)
		return number $ 'th'
	digit = number % 10 /* = mod 10 to get last digit as remainder*/
	return number $ #(th, st, nd, rd, th, th, th, th, th, th)[digit]
	}