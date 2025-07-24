// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.

// not working - incomplete

Refactor
	{
	Name: 'Extract Function'
	Desc: 'Convert the selection to a function in a separate library record'
	Controls: (Vert
		(Pair (Static 'Name') (Field font: '@mono' name: function_name))
		(Skip 4)
		(Pair (Static 'Inputs')
			(Field font: '@mono' readonly:, xstretch: 1 name: 'inputs'))
		(Skip 4)
		(Pair (Static 'Outputs')
			(Field font: '@mono' readonly:, xstretch: 1 name: 'outputs'))
		(Skip 4)
		(Static 'Call')
		(DisplayCode, ymin: 70, name: 'call')
		(Skip 4)
		(Static 'Function')
		(DisplayCode, readonly:, ymin: 140, ystretch: 3, name: 'method')
		xmin: 600
		)
	Init(data)
		{
		if .validate(data) is false
			return false
		return true
		}
	validate(data)
		{
		if data.select.cpMin >= data.select.cpMax
			{
			AlertError(.Name,
				'Please select the code you want to extract into a global function')
			return false
			}
		.selection = data.text[data.select.cpMin :: data.select.cpMax - data.select.cpMin]
		return true
		}

	Errors(data)
		{
		if data.function_name !~ .name_pat
			return "Invalid global name"
		return ""
		}
	name_pat: '^[[:upper:]][_[:alpha:][:digit:]]*[?!]?$'

	Process(data /*unused*/)
		{
Inspect(.selection)
		return true
		}
	}