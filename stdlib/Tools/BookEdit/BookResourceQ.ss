// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function (name, imagesOnly? = false, readOnly? = false)
	{
	if name !~ `^/res\>`
		return false
	if imagesOnly?
		pat = "(?i)[.](gif|jpg|png|bmp|emf|wmf|ico|cur)$"
	else if readOnly?
		pat = "(?i)[.](gif|jpg|png|bmp|emf|wmf|ico|cur|ttf|svg|map|afm)$"
	else
		pat = "(?i)[.](gif|jpg|png|bmp|emf|wmf|ico|cur|ttf|svg|map|afm|js|css)$"
	return name =~ pat
	}
