// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function (text)
	{
	if text is ''
		return ''
	scan = ScannerWithContext(text) // skips comments and whitespace
	token = scan.Next()
	if token is 'class' or token is 'function' or
		token is 'struct' or token is 'dll' or token is 'callback'
		return token
	if token is '#' and scan.Ahead() is '('
		return 'object'
	if token is '[' or (token is '#' and scan.Ahead() is '{')
		return 'record'
	if scan.Type() is #NUMBER or
		(token is '-' and scan.AheadType() is #NUMBER)
		return 'number'
	if scan.Type() is #IDENTIFIER and
		token =~ "^_?[[:upper:]]" and scan.Ahead() is '{'
		return 'class'
	if scan.Type() is #STRING or scan.Type() is #IDENTIFIER
		return 'string'
	return false
	}