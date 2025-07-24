// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function (data)
	{
	if not String?(data)
		data = Display(data)
	str = ""
	dataCopy = data
	if not dataCopy.Prefix?('<span style=')
		return dataCopy
	parser = XmlParser('<html>' $ dataCopy.Replace('<br />', '@\r\n@') $ '</html>')
	for child in parser.Children()
		{
		atts = child.Attributes()
		if atts.Member?('style') and atts.style.Has?(' line-through')
			continue
		str $= child.Text().Replace('@\r\n@', '\r\n')
		}
	return str
	}