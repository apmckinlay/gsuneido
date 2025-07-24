// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// e.g. MshtmlControl(Html_table_query('tables'))
function (query)
	// pre:	query is a database query
	// post:	returns a string like <table><tr><td>...</td></tr></table>
	//		columns are the query columns
	//		one row for each record
	{
	fields = QueryColumns(query)
	result = Xml('tr') { fields.Map({ Xml('th', Heading(it)) }).Join() } $ '\n'
	QueryApply(query)
		{ |x|
		result $= Xml('tr') { fields.Map({ Xml('td', x[it], valign: 'top') }).Join() }
		result $= '\n'
		}
	return Xml('table border="1" cellpadding="3"', result) $ '\n'
	}