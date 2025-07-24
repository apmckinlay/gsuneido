// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		parseSchema = Refactor_Table_Class.Refactor_Table_Class_parseSchema

		buildfkstr = Refactor_Table_Class.Refactor_Table_Class_buildForeignKeyStr
		Assert(buildfkstr(Object()), is: '', msg: 'no keys')

		ob = ('#' $ .test1).SafeEval()
		result = parseSchema(ob.schema)
		Assert(result.keys is: #(etaequip_num), msg: 'keys')
		Assert(result.indexes is: #(etatest_TS, 'one,two'), msg: 'indexes')
		Assert(result.foreignKeys.eta_orders is: #(etaorder_num, cascade:),
			msg: 'cascade foreign key')
		Assert(result.foreignKeys.biz_employees is: 'bizemp_num', msg: 'foreign key')

		str = buildfkstr(result.foreignKeys)
		sbe = "ForeignKeys: (biz_employees: ((from: bizemp_num))
	eta_orders: ((from: etaorder_num, cascade:))
	)"
		Assert(str is: sbe, msg: 'test 1')

		ob = ('#' $ .test2).SafeEval()
		result = parseSchema(ob.schema)
		Assert(result.keys equalsSet: #(etaol_name, etaol_num), msg: 'keys')
		Assert(result.indexes is: #(etaol_status), msg: 'indexes')
		Assert(result.foreignKeys.eta_orders is: 'etaorder_num',
			msg: 'second foreign key')

		str = buildfkstr(result.foreignKeys)
		sbe = "ForeignKeys: (eta_orders: ((from: etaorder_num))\r\n\t)"
		Assert(str is: sbe, msg: 'test 2')
		}
	test1: `(Tables, 'ensures', table: 'eta_order_driving_hours',
		schema: "(etaorder_num, etaorder_driving_hours, etaequip_num, bizemp_num)
		key(etaorder_num) in eta_orders cascade
		key(etaequip_num)
		index(bizemp_num) in biz_employees
		index(etatest_TS)
		index(one, two)
		key(key)")`
	test2: `(Tables, 'ensures', table: 'eta_order_lanes',
		schema: '(etaol_num, etaol_name, etaol_status,
			etaorder_num, etaol_days, etaol_weekstartday,
			Etaol_order, Etaol_order_bizpartner_num)
		key(etaol_num)
		key(etaol_name)
		index(etaol_status)
		index(etaorder_num) in eta_orders')`
	}