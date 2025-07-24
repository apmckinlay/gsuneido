// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Init()
		{
		r = new .test_refactor
		Assert(r.Init(Record(text: "class { }")) is: false)
		Assert(r.Msg is: 'Can only convert functions')
		}

	Test_Convert()
		{
		class_text =
'class
	{
	CallClass()
		{
		}
	}'
		function_text =
'function ()
	{
	}'
		Assert(Refactor_Convert_Function_To_Class.Convert(function_text) is: class_text)

		class_text = class_text.Replace('\(\)', '(a, b = 1, c = #())')
		function_text = function_text.Replace('\(\)', '(a, b = 1, c = #())')
		Assert(Refactor_Convert_Function_To_Class.Convert(function_text) is: class_text)
		}

	test_refactor: Refactor_Convert_Function_To_Class
		{
		Msg: false
		RefactorConvert_alert(msg) { .Msg = msg }
		}
	}