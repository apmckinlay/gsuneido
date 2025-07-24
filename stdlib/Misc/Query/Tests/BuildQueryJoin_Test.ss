// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table1 = .MakeTable('(num) key(num)')
		table2 = .MakeTable('(num, name, abbrev) key(num)')

		Assert(BuildQueryJoin("", #()) is: '')
		Assert(BuildQueryJoin(table1, #()) is: table1)

		joinob = Object(
			Object(str: " leftjoin by(num) (" $ table2 $ " project num)", fields: #()))
		Assert(BuildQueryJoin(table1, joinob)
			is: table1 $ " leftjoin by(num) (" $ table2 $ " project num)")

		query = table1 $ " extend keep"
		joinob = Object(
			Object(str: " leftjoin by(num) (" $ table2 $ " project num, name, abbrev)",
				fields: #(name, abbrev)))
		Assert(BuildQueryJoin(query, joinob)
			is: "(" $ table1 $ " extend keep) leftjoin by(num) (" $
				table2 $ " project num, name, abbrev)")

		query = table1 $ " extend keep, abbrev"
		joinob = Object(
			Object(str: " leftjoin by(num) (" $ table2 $ " project num, name, abbrev)",
				fields: #(name, abbrev)))
		Assert(BuildQueryJoin(query, joinob)
			is: "(" $ table1 $ " extend keep, abbrev remove abbrev) " $
				"leftjoin by(num) (" $ table2 $ " project num, name, abbrev)")

		table3 = .MakeTable('(num3, name3, abbrev3) key(num3)')
		table4 = .MakeTable('(num4, name4, abbrev4) key(num4)')
		query = table1 $ " extend abbrev4_billto, keep, abbrev, name3"
		joinob = Object(
			Object(str: " leftjoin by(num) (" $ table2 $ " project num, name, abbrev)",
				fields: #(name, abbrev)),
			Object(str: " leftjoin by(num3) (" $ table3 $" project num3, name3, abbrev3)",
				fields: #(name3, abbrev3)),
			Object(str: " leftjoin by(num4) (" $ table4 $ " project num4, name4, abbrev4
				rename num4 to num4_billto,	name4 to name4_billto,
					abbrev4 to abbrev4_billto",
				fields: #(name4_billto, abbrev4_billto)))
		Assert(BuildQueryJoin(query, joinob)
			is: "(" $ query $ " remove abbrev4_billto, abbrev, name3) " $
				"leftjoin by(num) (" $ table2 $ " project num, name, abbrev)" $
				" leftjoin by(num3) (" $ table3 $ " project num3, name3, abbrev3)" $
				" leftjoin by(num4) (" $ table4 $ " project num4, name4, abbrev4
				rename num4 to num4_billto,	name4 to name4_billto,
					abbrev4 to abbrev4_billto")
		}
	Test_vlsortextend()
		{
		table1 = .MakeTable('(user, num, num2, readyto_pay?, rate, date, num3, finalized?)
			key(num)')
		table2 = .MakeTable('(num, name, abbrev) key(num)')
		query = table1 $
			` rename user to bizuser_user_cur, num to num_new,
				num2 to num_prtrans, readyto_pay? to readyto_pay_default,
				rate to rate_default, date to date_default, num3 to num3_tran
			extend prtra_check, bizprotecthis_keyField = 'num_new',
				finalized_default = finalized? is true
			extend linkfield_tran
			where finalized_default isnt true
			/*vl_sort_extend*/
			 leftjoin by(num_prtrans) ` $
				`(` $ table2 $ ` project num, name, abbrev ` $
				`rename num to num_prtrans ` $
				`rename name to name_prtrans ` $
				`rename abbrev to abbrev_prtrans)`
		joinob = Object(
			Object(str: " leftjoin by(num_prtrans) " $
					"(" $ table2 $ " project num, name, abbrev " $
					"rename num to num_prtrans " $
					"rename name to name_prtrans " $
					"rename abbrev to abbrev_prtrans) ",
				fields: #(name_prtrans, abbrev_prtrans)))
		Assert(BuildQueryJoin(query, joinob) is: "(" $ query $ ")")
		}
	}