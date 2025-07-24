// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function (prompt, emptyOK = false, description = 'Custom Prompt')
	{
	if prompt is ""
		return emptyOK ? "" : description $ " can not be empty."
	if prompt.Blank?()
		return description $ " must contain more than just blank spaces."
	if prompt !~ "[a-zA-Z]"
		return description $ " must contain at least one alpha character."
	if prompt =~ "[^a-zA-Z0-9 /\?#]"
		return description $ " can not contain special characters."
	if prompt.Size() > (max = 50)
		return description $ ' must be less than ' $ max $ ' characters. \r\n\r\n' $
			'Please consider using "Tooltips" or ' $
			'typing static text in the layouts directly.'
	return ""
	}