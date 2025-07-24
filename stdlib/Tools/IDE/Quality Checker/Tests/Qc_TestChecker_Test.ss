// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Continuous_TestFileExistence()
		{
		testClass = "//Suneido ...
		class
			{
			}"

		testInheritance1 = "//Suneido...
		Qc_Coupling
			{
			}"

		testInheritance2 = "//Suneido...
		_Qc_Coupling
			{
			}"

		testInheritance3 = "//Suneido...
		_Qc_Coupling("

		testFunction = "//Suneido...
		function ()
			{
			}"

		testFunction2 = "//Suneido...
		function ()
			{
			}"

		testFunction3 = "//Suneido...
		function     ()
			{
			}"

		testFunction4 = "//Suneido...
		Test
			{
			}"

		testFunction5 = "//Suneido...
		class
			{
			}"

		testContrib1 = "//Suneido...
		#(Multi: (table: 'eta_orders',
			invoiceField: 'etaorder_invoice', clearFields: (etaorder_invoice_on)
			confirm: function (@unused) { return true }
			func: function (sourceRec, ivcRec, t /*unused*/)
				{
				sourceRec.etaorder_previous_invoice = ivcRec.arivc_invoice
				for ob in sourceRec.etaorder_taxcodes_ob
					ob.Delete('rate')
				}
			))"

		.MakeLibraryRecord([name: "Fake1", text: testClass])
		.MakeLibraryRecord([name: "Fake2", text: testInheritance1])
		.MakeLibraryRecord([name: "Fake2_Test", text: "gibberish"])
		.MakeLibraryRecord([name: "Fake3", text: testInheritance2])
		.MakeLibraryRecord([name: "Fake3Test", text: testInheritance3])

		.MakeLibraryRecord([name: "Fake4", text: testFunction])
		.MakeLibraryRecord([name: "Fake5", text: testFunction2])
		.MakeLibraryRecord([name: "Fake5_Test",text: "gibberish"])
		.MakeLibraryRecord([name: "Fake6", text: testFunction3])
		.MakeLibraryRecord([name: "Fake6Test", text: "gibberish"])


		.MakeLibraryRecord([name: "Fake7?", text: testClass])
		.MakeLibraryRecord([name: "Fake7_Test", text: "gibberish"])
		.MakeLibraryRecord([name: "Fake9?", text: testClass])
		.MakeLibraryRecord([name: "Fake10Tests", text: testFunction4])
		.MakeLibraryRecord([name: "Fake11Tests", text: testFunction5])
		.MakeLibraryRecord([name: "Fake1Contrib", text: testContrib1])

		recordData = Record()
		recordData.lib = "Test_lib"
		recordData.recordName = "Fake1"
		recordData.code = testClass
		Assert(Qc_TestChecker(recordData, false) is:
			#(warnings: #(),
			desc: "A test class was not found for this record. Please create one",
			maxRating: 4, size: 1))
		Assert(Qc_TestChecker(recordData, true) is:
			#(warnings: #(),
			desc: "A test class was not found for this record. Please create one",
			maxRating: 4, size: 1))

		recordData.recordName = "Fake2"
		recordData.code = testInheritance1
		Assert(Qc_TestChecker(recordData, false) is:
			#(warnings: #(),
			desc: "A test class was found for this record or is not required",
			maxRating: 5, size: 0))
		Assert(Qc_TestChecker(recordData, true) is:
			#(warnings: #(), desc: "", maxRating: 5, size: 0))

		recordData.recordName = "Fake3"
		recordData.code = testInheritance3
		Assert(Qc_TestChecker(recordData, false) is:
			#(warnings: #(),
			desc: "A test class was found for this record or is not required",
			maxRating: 5, size: 0))
		Assert(Qc_TestChecker(recordData, false) is:
			#(warnings: #(),
			desc: "A test class was found for this record or is not required",
			maxRating: 5, size: 0))

		recordData.recordName = "Fake4"
		recordData.code = testFunction
		Assert(Qc_TestChecker(recordData, false) is:
		#(warnings: #(),
			desc: "A test class was not found for this record. Please create one",
			maxRating: 4, size: 1))

		Assert(Qc_TestChecker(recordData, true) is:
		#(warnings: #(),
			desc: "A test class was not found for this record. Please create one",
			maxRating: 4, size: 1))

		recordData.recordName = "Fake5"
		recordData.code = testFunction2
		Assert(Qc_TestChecker(recordData, false) is:
			#(warnings: #(),
			desc: "A test class was found for this record or is not required",
			maxRating: 5, size: 0))
		Assert(Qc_TestChecker(recordData, true) is:
			#(warnings: #(), desc: "", maxRating: 5, size: 0))

		recordData.recordName = "Fake6"
		recordData.code = testFunction3
		Assert(Qc_TestChecker(recordData, false) is:
			#(warnings: #(),
			desc: "A test class was found for this record or is not required",
			maxRating: 5, size: 0))
		Assert(Qc_TestChecker(recordData, true) is:
			#(warnings: #(), desc: "", maxRating: 5, size: 0))

		recordData.recordName = "Fake2_Test"
		recordData.code = "Test\n { }"
		Assert(Qc_TestChecker(recordData, false) is:
			#(warnings: #(),
			desc: "A test class was found for this record or is not required",
			maxRating: 5, size: 0))
		Assert(Qc_TestChecker(recordData, true) is:
			#(warnings: #(), desc: "", maxRating: 5, size: 0))

		recordData.recordName = "Fake3Test"
		recordData.code = "Test\n {}"
		Assert(Qc_TestChecker(recordData, false) is:
			#(warnings: #(),
			desc: "A test class was found for this record or is not required",
			maxRating: 5, size: 0))
		Assert(Qc_TestChecker(recordData, true) is:
			#(warnings: #(), desc: "", maxRating: 5, size: 0))

		recordData.recordName = "Fake7?"
		recordData.code = testClass
		Assert(Qc_TestChecker(recordData, false) is:
			#(warnings: #(),
			desc: "A test class was found for this record or is not required",
			maxRating: 5, size: 0))
		Assert(Qc_TestChecker(recordData, true) is:
			#(warnings: #(), desc: "", maxRating: 5, size: 0))

		recordData.recordName = "Fake9?"
		recordData.code = testClass
		Assert(Qc_TestChecker(recordData, false) is:
			#(warnings: #(),
			desc: "A test class was not found for this record. Please create one",
			maxRating: 4, size: 1))
		Assert(Qc_TestChecker(recordData, true) is:
			#(warnings: #(),
			desc: "A test class was not found for this record. Please create one",
			maxRating: 4, size: 1))

		recordData.recordName = "Fake10Tests"
		recordData.code = testFunction4
		Assert(Qc_TestChecker(recordData, false) is:
			#(warnings: #(),
			desc: "A test class was found for this record or is not required",
			maxRating: 5, size: 0))
		Assert(Qc_TestChecker(recordData, true) is:
			#(warnings: #(), desc: "", maxRating: 5, size: 0))

		recordData.recordName = "Fake11Tests"
		recordData.code = testFunction5
		Assert(Qc_TestChecker(recordData, false) is:
			#(warnings: #(),
			desc: "A test class was not found for this record. Please create one",
			maxRating: 4, size: 1))
		Assert(Qc_TestChecker(recordData, true) is:
			#(warnings: #(),
			desc: "A test class was not found for this record. Please create one",
			maxRating: 4, size: 1))

		recordData.recordName = "Fake1Contrib"
		recordData.code = testContrib1
		Assert(Qc_TestChecker(recordData, false) is:
			#(warnings: #(),
			desc: "A test class was found for this record or is not required",
			maxRating: 5, size: 0))
		Assert(Qc_TestChecker(recordData, true) is:
			#(warnings: #(), desc: "", maxRating: 5, size: 0))
		}
	}






