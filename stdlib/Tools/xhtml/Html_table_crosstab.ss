// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function (query, rows='', cols='', value='', func='total')
	// pre:	query is a database query
	//		rows is the field whose values will be the rows
	//		cols is the field whose values will be the columns
	//		value is the field whose value will be processed
	//		func is the function to be applied
	//			total, max, min, average, count
	//		if func is count, then value is ignored
	// post:	returns a string like <table><tr><td>...</td></tr></table>
	{
	query = QueryStripSort(query)

	columns = Object()
	if (cols isnt "")
		QueryApply(query $ ' project ' $ cols)
			{|x| columns.Add(x[cols])}
	if (.func is 'count')
		value = ''
	else if (value is "")
		throw "CrossTable " $ func $ " requires Value field"

	if (rows isnt "")
		query $=  ' rename ' $ rows $ ' to rowfield'
	query $= " summarize "
	if (rows isnt "")
		query $= " rowfield, "
	if (cols isnt "")
		query $= cols $ ", "
	query $= func $ " " $ value

	result = '<table border="1">\n'
	result $= '<tr>' $
	result $= Xml('td', rows)
	for (field in columns)
		result $= Xml('td', field)
	result $= '</tr>' $

	QueryApply(query)
		{ |x|
		result $= '<tr>'
		for (field in columns)
			result $= Xml('td', x[field])
		result $= '</tr>\n'
		}
	result $= '</table>'
	return result
	}