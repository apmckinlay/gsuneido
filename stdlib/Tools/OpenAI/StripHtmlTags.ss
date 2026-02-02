// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function (html)
	{
	result = ""
	inTag = false
	tag = ""
	blockTags = Object("br", "p", "div", "h1", "h2", "h3", "h4", "h5", "h6",
		"li", "ul", "ol", "table", "tr")

	i = 0
	while i < html.Size()
		{
		c = html[i]
		if not inTag
			{
			if c is '<'
				{
				inTag = true
				tag = ""
				}
			else
				// Only append content if it's not just indentation/formatting whitespace
				{
				// If it's whitespace, only append if last result char is not whitespace
				if " \t\r\n".Has?(c)
					{
					// Only add a single space if previous char is not space or newline
					if result.Size() > 0 and not " \n".Has?(result[result.Size()-1])
						result $= " "
					}
				else
					result $= c
				}
			}
		else
			{
			if c is '>'
				{
				inTag = false
				// Get tag name (up to space, /, or >)
				tn = tag.Trim().Lower()
				isClosing = false
				if tn[0] is '/'
					{
					isClosing = true
					tn = tn[1..]
					}
				tn = tn.Trim(' /')
				// Only add newline for certain tags, and only after closing tag (except br)
				if blockTags.Has?(tn)
					{
					if tn is "br" or isClosing
						result $= "\n"
					}
				}
			else if c isnt '<'
				tag $= c
			}
		i++
		}
	// Condense multiple newlines and strip lines with only whitespace
	return result.Lines().Filter({ not it.Trim().Blank?() }).Join('\n')
	}