// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		restrictions = Object()
		Assert(BuildQueryWhere(restrictions) is: '')

		restrictions = Object(Object('' '' ''))
		Assert(BuildQueryWhere(restrictions) is: ' where   ""')

		restrictions = Object(Object('' '' '' built: false))
		Assert(BuildQueryWhere(restrictions) is: ' where   ""')
		}

	Test_preBuilt()
		{
		restrictions = Object(Object('' built: true))
		Assert(BuildQueryWhere(restrictions) is: ' ')

		restrictions = Object(Object('abc < 123' built: true))
		Assert(BuildQueryWhere(restrictions) is: ' abc < 123')
		}

	Test_lists()
		{
		restrictions = Object(Object('bizpartner_city' 'in list'
			Object('Bristol' 'Cochrane' 'Edmonton'))) //in list
		Assert(BuildQueryWhere(restrictions)
			is: ' where bizpartner_city in ("Bristol", "Cochrane", "Edmonton")')

		restrictions = Object(Object('bizpartner_city' 'not in list'
			Object('Bristol' 'Cochrane' 'Edmonton'))) //not in list
		Assert(BuildQueryWhere(restrictions)
			is: ' where bizpartner_city not in ("Bristol", "Cochrane", "Edmonton")')

		// integers (this is where commas between negative values is important
		// in order to prevent suneido from thinking it is an expression to evaluate
		restrictions = Object(Object('number' 'in list' Object(1, -10, -161)))
		Assert(BuildQueryWhere(restrictions) is: ' where number in (1, -10, -161)')

		//in list, with no object defining a list, just a string
		restrictions = Object(Object('bizpartner_city' 'in list' 'list'))
		Assert(BuildQueryWhere(restrictions) is: '')

		//in list, with no object defining a list, blank
		restrictions = Object(Object('bizpartner_city' 'in list' ''))
		Assert(BuildQueryWhere(restrictions) is: '')
		}

	Test_regular_operators()
		{
		//greater than
		restrictions = Object(Object('apcheck_date_report' '>' '#20081016'))
		Assert(BuildQueryWhere(restrictions)
			is: ' where apcheck_date_report > "#20081016"')

		//equals
		restrictions = Object(Object('bizpartner_num_supplier' 'is' 'Revcor Lines'))
		Assert(BuildQueryWhere(restrictions)
			is: ' where bizpartner_num_supplier is "Revcor Lines"')

		//empty
		restrictions = Object(Object('arivc_custpo' 'is' ''))
		Assert(BuildQueryWhere(restrictions) is: ' where arivc_custpo is ""')

		//not empty
		restrictions = Object(Object('arivc_custpo' '>' ''))
		Assert(BuildQueryWhere(restrictions) is: ' where arivc_custpo > ""')

		//contains
		restrictions = Object(Object('bizpartner_num_supplier' '=~' '(?i)(?q)Revcor'))
		Assert(BuildQueryWhere(restrictions)
			is: ' where bizpartner_num_supplier =~ "(?i)(?q)Revcor"')

		//range
		restrictions = Object(Object('bizpartner_city' '>=' 'Bristol'))
		restrictions.Add(Object('bizpartner_city' '<=' 'Edmonton'))
		Assert(BuildQueryWhere(restrictions)
			is: ' where bizpartner_city >= "Bristol" and bizpartner_city <= "Edmonton"')

		//not in range
		restrictions = Object(
			Object('bizpartner_city' 'not in range' 'Bristol' 'Edmonton'))
		Assert(BuildQueryWhere(restrictions)
			is: ' where (bizpartner_city < "Bristol" or bizpartner_city > "Edmonton")')
		}

	Test_multiple_criteria()
		{
		//in list
		restrictions = Object(Object('bizpartner_city' 'in list'
				Object('Bristol' 'Cochrane' 'Edmonton')))
		//greater than
		restrictions.Add(Object('apcheck_date_report' '>' '#20081016'))
		Assert(BuildQueryWhere(restrictions)
			is: ' where bizpartner_city in ("Bristol", "Cochrane", "Edmonton")' $
				' and apcheck_date_report > "#20081016"')

		//greater than
		restrictions = Object(Object('etaorder_order' '>' 1050))
		//in list
		restrictions.Add(Object('etaorder_status' 'in list'
			Object('Preplanned' 'Available' 'Quote')))
		//empty
		restrictions.Add(Object('etaorder_pod_date' 'is' ''))
		//not equal to
		restrictions.Add(Object('bizpartner_num_shipper' 'isnt' #20081016.094026052))
		//not empty
		restrictions.Add(Object('etaequip_num_tractor' '>' ''))
		//equals
		restrictions.Add(Object('etaorder_commodity' 'is' 'Fertilizer'))
		//contains
		restrictions.Add(Object('etaorder_custpo' '=~' '(?i)(?q)test ref'))
		//range
		restrictions.Add(Object('etaorder_invoice' '>=' '11020'))
		//range
		restrictions.Add(Object('etaorder_invoice' '<=' '11024'))
// need to wait for __trial changes to be moved to std to uncomment this
//		Assert(BuildQueryWhere(restrictions)
//			is: ' where Number?(etaorder_order) and etaorder_order > 1050' $
//			' and etaorder_status in ("Preplanned", "Available", "Quote")' $
//			' and etaorder_pod_date is ""' $
//			' and bizpartner_num_shipper isnt #20081016.094026052' $
//			' and etaequip_num_tractor > ""' $
//			' and etaorder_commodity is "Fertilizer"' $
//			' and etaorder_custpo =~ "(?i)(?q)test ref"' $
//			' and etaorder_invoice >= "11020" and etaorder_invoice <= "11024"')

		//equals
		restrictions = Object(Object('user' '=' 'admin'))
		//not in list
		restrictions.Add(Object('etaorder_status' 'not in list'
			Object("Completed" "Overdue")))
		Assert(BuildQueryWhere(restrictions) is: ' where user = "admin"' $
			' and etaorder_status not in ("Completed", "Overdue")')
		}

	Test_not_in_range()
		{
		restrictions = Object(Object('date', 'not in range', #20100101, #20101231))
		Assert(BuildQueryWhere(restrictions)
			is: ' where (date < #20100101 or date > #20101231)')
		}

	Test_callable()
		{
		restrictions = Object(Object('date', 'not in range', #20100101, #20101231))
		fn = BuildQueryWhere(restrictions, build_callable:)
		Assert(fn([date: #20110101]))
		Assert(fn([date: #20101231]) is: false)
		Assert(fn([date: #20101230]) is: false)

		restrictions = Object()
		fn = BuildQueryWhere(restrictions, build_callable:)
		Assert(fn([value: 'any']))

		restrictions = Object(Object('date', 'in list', Object(#20100101, #20101231), ''))
		fn = BuildQueryWhere(restrictions, build_callable:)
		Assert(fn([date: #20100101]))
		Assert(fn([date: #20190101]) is: false)

		restrictions = Object(
			Object('date', 'not in list', Object(#20100101, #20101231), ''))
		fn = BuildQueryWhere(restrictions, build_callable:)
		Assert(fn([date: #20100101]) is: false)
		Assert(fn([date: #20190101]))

		restrictions = Object(Object('ScintillaRichStripHTML(comments)', 'is', ''))
		fn = BuildQueryWhere(restrictions, build_callable:)
		Assert(fn([comments: '']))
		Assert(fn([comments: 'test']) is: false)

		// test throw
		cl = BuildQueryWhere
			{
			BuildQueryWhere_compile(unused)
				{
				throw _compileResult
				}
			BuildQueryWhere_client?()
				{
				return _client
				}
			}

		_client = true
		restrictions = Object(Object('date', 'in list', Object(#20100101, #20101231), ''))
		_compileResult = 'compile error @something'
		Assert({ cl(restrictions, build_callable:) } throws:
			'SHOW: There was a problem with the filter.\n' $
			'This could be caused by invalid filter options or too many In List values')
		_compileResult = 'other errors'
		Assert({ cl(restrictions, build_callable:) } throws: 'other errors')
		}
	}
