// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function (value)
	{
	return not Date?(value)
		? String(value)
		: value is value.NoTime()
			? value.ShortDate()
			: value.ShortDateTime()
	}