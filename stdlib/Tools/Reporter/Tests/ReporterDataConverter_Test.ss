// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ToFields()
		{
		for caseOb in .cases
			{
			sf = .basicSelectFields()
			data = caseOb.prompts.DeepCopy()
			.addSummarizeFields(sf, data)
			ReporterDataConverter.ToFields(data, sf)
			Assert(data is: caseOb.fields)

			// already has usingFieldsInSave? flag
			ReporterDataConverter.ToFields(data, sf)
			Assert(data is: caseOb.fields)
			}
		}

	Test_ToPrompts()
		{
		for caseOb in .cases
			{
			sf = .basicSelectFields()
			data = caseOb.fields.DeepCopy()
			ReporterDataConverter.ToPrompts1(data, sf)
			.addSummarizeFields(sf, data)
			ReporterDataConverter.ToPrompts2(data, sf)
			Assert(data is: caseOb.prompts)

			// no usingFieldsInSave? flag
			ReporterDataConverter.ToPrompts2(data, sf)
			Assert(data is: caseOb.prompts)
			}
		}

	basicSelectFields()
		{
		sf = SelectFields()
		sf.AddField('a', 'A')
		sf.AddField('b', 'B')
		sf.AddField('c', 'C')
		return sf
		}

	addSummarizeFields(sf, data)
		{
		mock = Mock(ReporterModel)
		mock.ReporterModel_sf = sf
		mock.ReporterModel_data = data
		mock.When.add_summarize_fields().CallThrough()
		mock.When.addSummaryByFields().CallThrough()
		mock.When.add_summarize_field().CallThrough()
		mock.add_summarize_fields()
		}

	cases: [
		[prompts: [report_name: "A",
			heading2: 'B',
			sort0: 'A', total0: true, show0: true,
			sort1: 'B', Source: "",
			select: [
				checkbox0: true, fieldlist0: 'A', oplist0: "contains", val0: "",
					print0: true, menu_option0: true,
				checkbox1: true, fieldlist1: 'C', menu_option1: true],
			columns: #(#(width: 10, text: 'C'),
				#(width: 22, text: 'B'),
				#(width: 15, text: 'A'),
				#(width: 15, text: 'Unknown')),
			coloptions: #(
				'invalid',
				A: [heading: "A"],
				B: [heading: "BB"],
				C: [heading: "CCC"])],
		fields: [report_name: "A",
			heading2: 'B',
			sort0: 'a', total0: true, show0: true,
			sort1: 'b', Source: "",
			select: [
				checkbox0: true, fieldlist0: 'a', oplist0: "contains", val0: "",
					print0: true, menu_option0: true,
				checkbox1: true, fieldlist1: 'c', menu_option1: true],
			columns: #(#(width: 10, text: 'c'),
				#(width: 22, text: 'b'),
				#(width: 15, text: 'a'),
				#(width: 15, text: 'Unknown')),
			coloptions: #(
				'invalid',
				a: [],
				b: [heading: "BB"],
				c: [heading: "CCC"]),
			usingFieldsInSave?:]
			],
		[prompts: [report_name: "A",
			heading2: 'B',
			summarize_by: 'A,B'
			summarize_field0: 'C', summarize_func0: "average",
			summarize_func1: "count",
			summarize_field2: 'C', summarize_func2: "minimum",
			sort0: 'Minimum C', total0: true, show0: true,
			sort1: 'Count', Source: "",
			select: [
				checkbox0: true, fieldlist0: 'A', oplist0: "contains", val0: "",
					print0: true, menu_option0: true,
				checkbox1: true, fieldlist1: 'Average C', menu_option1: true],
			columns: #('Average C', 'B', 'A', 'Minimum C', 'Count', 'Unknown'),
			coloptions: #(
				A: [heading: "A"],
				B: [heading: "BB"],
				Count: [heading: "Count"],
				'Minimum C': [heading: "Min C"],
				'Average C': [heading: "Avg C"]),
			formulas: #([formula: 'A $ " - b: " $ B', key: "20240716095827",
				type: "Text, single line", calc: "test"]),
				],
		fields: [report_name: "A",
			heading2: 'B',
			summarize_by: 'a,b'
			summarize_field0: 'c', summarize_func0: "average",
			summarize_func1: "count",
			summarize_field2: 'c', summarize_func2: "minimum",
			sort0: 'min_c', total0: true, show0: true,
			sort1: 'count', Source: "",
			select: [
				checkbox0: true, fieldlist0: 'a', oplist0: "contains", val0: "",
					print0: true, menu_option0: true,
				checkbox1: true, fieldlist1: 'average_c', menu_option1: true],
			columns: #('average_c', 'b', 'a', 'min_c', 'count', 'Unknown'),
			coloptions: #(
				a: [],
				b: [heading: "BB"],
				count: [],
				'min_c': [heading: "Min C"],
				'average_c': [heading: "Avg C"]),
			formulas: #([formula: 'FORMULA_#0# $ " - b: " $ FORMULA_#1#',
				key: "20240716095827", type: "Text, single line", calc: "test",
				fields: (a, b)]),
			usingFieldsInSave?:]]]
	}