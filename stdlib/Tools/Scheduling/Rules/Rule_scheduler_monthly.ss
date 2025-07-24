// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	str = ''
	suffix = " of month"
	oddInputs = #(beginning end middle "2nd" "3rd")
	for (i = 0; i < 28 + oddInputs.Size() - 3; i++)
		if (i < oddInputs.Size())
			str $= oddInputs[i] $ suffix $ ','
		else
			str $= String(i - 1) $ "th" $ suffix $ ','
	return str[.. -1]
	}