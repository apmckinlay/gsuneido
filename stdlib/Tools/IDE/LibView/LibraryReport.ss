// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
QueryFormat
	{
	CallClass(lib)
		{
		return Object('Params', this,
			title: "Print Library"
			name: "print_library"
			Params:
				Object('Vert',
					Object('Pair', #('Static', 'Library')
						Object('ChooseList', Libraries(), name: "lib"))
					#(Pair (Static Paging)
						(ChooseList (Continuous Normal Duplex), name: 'paging')))
			SetParams: Object(:lib)
			header: false
			)
		}
	Query()
		{
		return _report.Params.lib $ " where group is -1 sort name"
		}
	Header()
		{
		return Object('PageHead', _report.Params.lib)
		}
	BeforeOutput(data)
		{
		return Object('Library', data.name, data.text)
		}
	Output: false
	AfterOutput(data/*unused*/)
		{
		switch (_report.Params.paging)
			{
		case 'Continuous' :
			return #(Vskip)
		case 'Normal', '' :
			return 'pg'
		case 'Duplex' :
			return "pgo"
			}
		}
	}
