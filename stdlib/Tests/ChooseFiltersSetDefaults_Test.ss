// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		// defaultFilters - none
		Assert(ChooseFiltersSetDefaults(#(), #()) is: #())

		// no defaultFilters but savedFilters
		Assert(ChooseFiltersSetDefaults(false, #(a, b)) is: #(a, b))

		// default: #(a, b, c)		saved: #(d, e) = #(a, b, c, d, e)
		defFitlerFields = #(a, b, c)
		savedFilters = Object(
			#(condition_field: 'd', d: #(operation: 'equals', value: 'fred', value2: ''))
			#(condition_field: 'e', e: #(operation: '', value: '', value2: '')))
		Assert(ChooseFiltersSetDefaults(defFitlerFields, savedFilters)
			is: #((condition_field: 'a' a: #(operation: "", value: "", value2: ""))
				(condition_field: 'b' b: #(operation: "", value: "", value2: ""))
				(condition_field: 'c' c: #(operation: "", value: "", value2: ""))
				(condition_field: 'd' d: #(operation: "equals", value: "fred",value2: ""))
				(condition_field: 'e' e: #(operation: "", value: "", value2: ""))))

		// default: #(a, b, c)		saved: #() = #(a, b, c)
		Assert(ChooseFiltersSetDefaults(defFitlerFields, Object())
			is: #((condition_field: 'a' a: #(operation: "", value: "", value2: ""))
				(condition_field: 'b' b: #(operation: "", value: "", value2: ""))
				(condition_field: 'c' c: #(operation: "", value: "", value2: ""))))

		// default: #(a, b, c)		saved: #(a, c, e) = #(a, b, c, e)
		savedFilters = Object(
			#(condition_field: 'a', a: #(operation: '', value: '', value2: ''))
			#(condition_field: 'c', c: #(operation: '', value: '', value2: ''))
			#(condition_field: 'e', e: #(operation: '', value: '', value2: '')))
		Assert(ChooseFiltersSetDefaults(defFitlerFields, savedFilters)
			is: #((condition_field: 'a' a: #(operation: "", value: "", value2: ""))
				(condition_field: 'b' b: #(operation: "", value: "", value2: ""))
				(condition_field: 'c' c: #(operation: "", value: "", value2: ""))
				#(condition_field: 'e', e: #(operation: '', value: '', value2: ''))))

		// default: #(a, b, c)		saved:#(e, a, c) = #(b, a, c, e)
		savedFilters = Object(
			#(condition_field: 'e', e: #(operation: '', value: '', value2: ''))
			#(condition_field: 'a', a: #(operation: '', value: '', value2: ''))
			#(condition_field: 'c', c: #(operation: '', value: '', value2: '')))
		Assert(ChooseFiltersSetDefaults(defFitlerFields, savedFilters)
			is: #((condition_field: 'e' e: #(operation: "", value: "", value2: ""))
				(condition_field: 'b' b: #(operation: "", value: "", value2: ""))
				(condition_field: 'a' a: #(operation: "", value: "", value2: ""))
				#(condition_field: 'c', c: #(operation: '', value: '', value2: ''))))

		// default: #(a, b, c)		saved:#(a, b, c, a) = #(a, b, c, a)
		savedFilters = Object(
			#(condition_field: 'a', a: #(operation: '', value: '', value2: ''))
			#(condition_field: 'b', b: #(operation: '', value: '', value2: ''))
			#(condition_field: 'c', c: #(operation: '', value: '', value2: ''))
			#(condition_field: 'a', a: #(operation: '', value: '', value2: '')))
		Assert(ChooseFiltersSetDefaults(defFitlerFields, savedFilters)
			is: #((condition_field: 'a' a: #(operation: "", value: "", value2: ""))
				(condition_field: 'b' b: #(operation: "", value: "", value2: ""))
				(condition_field: 'c' c: #(operation: "", value: "", value2: ""))
				#(condition_field: 'a', a: #(operation: '', value: '', value2: ''))))
		}
	}
