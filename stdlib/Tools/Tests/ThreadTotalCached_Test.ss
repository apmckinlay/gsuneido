// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Suneido.Delete('ThreadTotal_TestThreadTotalCached')

		cl = ThreadTotalCached
			{
			ThreadTotalCached_client?()
				{
				return false
				}
			ThreadTotalCached_getFunc(unused, filters)
				{
				switch (filters.Size())
					{
					case 1:
						return [totals: 1]
					case 2:
						return [totals: 2]
					default:
						return []
					}
				}
			}

		filters = Object([xxtran_date_default: #(value: "", operation: "", value2: ""),
			condition_field: "xxtran_date_default"],
			[bizemp_num_xxtrans: #(value: "", value2: "", operation: ""),
				condition_field: "bizemp_num_xxtrans"],
			[xxtran_dept_desc: #(value: "", operation: "", value2: ""),
				condition_field: "xxtran_dept_desc"],
			[xxtran_type: #(value: "", operation: "", value2: ""),
				condition_field: "xxtran_type"],
			[xxtran_desc: #(value: "", operation: "", value2: ""),
				condition_field: "xxtran_desc"],
			[xxtran_reference: #(value: "", operation: "", value2: ""),
				condition_field: "xxtran_reference"],
			[xxtran_check: #(value: "", operation: "", value2: ""),
				condition_field: "xxtran_check"],
			[xxtran_checkdate: #(value: "", operation: "", value2: ""),
				condition_field: "xxtran_checkdate"])

		filters[1] = [bizemp_num_xxtrans:
			#(value: "#20210414.110739805", value2: "", operation: "equals"),
				condition_field: "bizemp_num_xxtrans"]
		f1 = filters.Copy()
		Assert(cl('test', 'test_table where num is 1', 'TestThreadTotalCached', f1) is:
			[totals: 1])
		Assert(Suneido.ThreadTotal_TestThreadTotalCached.GetMissRate() is: 1)
		Assert(cl('test', 'test_table where num is 1', 'TestThreadTotalCached', f1) is:
			[totals: 1])
		Assert(Suneido.ThreadTotal_TestThreadTotalCached.GetMissRate() is: 0.5)

		filters[4] = [xxtran_desc: #(value: "Regular", operation: "equals", value2: ""),
				condition_field: "xxtran_desc"]
		f2 = filters.Copy()
		Assert(cl('test', 'test_table where num is 2', 'TestThreadTotalCached', f2) is:
			[totals: 2])
		Assert(Suneido.ThreadTotal_TestThreadTotalCached.GetMissRate().Round(2) is: 0.67)
		Assert(cl('test', 'test_table where num is 2', 'TestThreadTotalCached', f2) is:
			[totals: 2])
		Assert(Suneido.ThreadTotal_TestThreadTotalCached.GetMissRate() is: 0.5)

		// test switching filters order
		temp = filters[1]
		filters.Delete(1).Add(temp)
		Assert(cl('test', 'test_table where num is 2', 'TestThreadTotalCached',
			filters) is: [totals: 2])
		Assert(Suneido.ThreadTotal_TestThreadTotalCached.GetMissRate() is: 0.4)
		Assert(cl('test', 'test_table where num is 2', 'TestThreadTotalCached',
			filters) is: [totals: 2])
		Assert(Suneido.ThreadTotal_TestThreadTotalCached.GetMissRate().Round(2) is: 0.33)
		}
	}
