// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// e.g. MshtmlControl(Html_table(#(((tl da: 1), tr ra: 2) (bl, br), border: 1)))
function (data)
	// pre:	data is a list of lists of values
	//		values can be scalars or objects like #(scalar, attribute: scalar)
	// post:	returns a string like <table><tr><td>...</td></tr></table>
	//		named members of data become table attributes
	//		named members of rows become row (tr) attributes
	{
	attributes = function (data)
		{
		result = ""
		dm = data.Members()
		for (n = data.Size(), i = data.Size(list:); i < n; ++i)
			result $= ' ' $ dm[i] $ '="' $ String(data[dm[i]]) $ '"'
		return result
		}
	result = '<table' $ attributes(data) $ '>\n'
	for row in data.Values(list:)
		{
		result $= '<tr' $ attributes(row) $ '>'
		for val in row.Values(list:)
			result $= Xml('td', val)
		result $= '</tr>\n'
		}
	result $= '</table>'
	return result
	}
