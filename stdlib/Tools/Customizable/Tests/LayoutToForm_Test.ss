// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		sf = SelectFields()
		sf.AddField('custom_1', 'One')
		sf.AddField('custom_2', 'Two')
		sf.AddField('custom_3', 'Three')
		sf.AddField('custom_4', 'Four')
		sf.AddField('custom_5', 'Five')
		sf.AddField('custom_6', 'Six')
		sf.AddField('custom_7', 'Seven')
		sf.AddField('custom_8', 'Eight')
		sf.AddField('custom_9', 'Nine')
		sf.AddField('custom_10', 'Ten')
		sf.AddField('custom_11', 'Eleven')
		sf.AddField('custom_12', 'Twelve')
		sf.AddField('custom_13', 'Thirteen')
		Assert(LayoutToForm("", sf) is: #(Form))
		Assert(LayoutToForm("One", sf) is: #(Form, (custom_1)))
		Assert(LayoutToForm("One Two", sf)
			is: #(Form, (custom_1), (Static, ' '), (custom_2)))
		Assert(LayoutToForm("One\nTwo", sf)
			is: #(Form, (custom_1, group: 0), nl, (custom_2, group: 0)))
		Assert(LayoutToForm("Title\nOne", sf)
			is: #(Form, (Static Title) nl (custom_1)))

		Assert(LayoutToForm("----Title Here----\nOne Two\n Done", sf),
			is: #(Form, (Static "----Title Here----") nl
				(custom_1)(Static, ' ')(custom_2) nl
				(Static " Done")))
		s = '
abc One   Two
	Three Four
Five'.Trim().Detab()
		expectedForm = #(Form,
			(Static, 'abc '), (custom_1, group: 0), (Static, '   '),
				(custom_2, group: 1), nl,
			(Static, '    '), (custom_3, group: 0), (Static, ' '),
				(custom_4, group: 1), nl
			(custom_5))
		.checkLayoutsPreserved(s, expectedForm, sf)

s = '
One sample text Two
Three 			Four
 Five
  Six
   Seven'.Trim().Detab()
		expectedForm = #(Form,
			(custom_1, group: 0), (Static, ' sample text '), (custom_2, group: 1), nl,
			(custom_3, group: 0), (Static, '           '), (custom_4, group: 1), nl,
			(Static, ' '), (custom_5), nl,
			(Static, '  '), (custom_6), nl,
			(Static, '   '), (custom_7))
		.checkLayoutsPreserved(s, expectedForm, sf)
		s = '
One    Two
 Three Four
 Five'.Trim().Detab()
		expectedForm = #(Form,
			(custom_1), (Static, '    '), (custom_2, group: 1), nl
			(Static, ' '), (custom_3, group: 0), (Static, ' '), (custom_4, group: 1), nl
			(Static, ' '), (custom_5, group: 0))
		.checkLayoutsPreserved(s, expectedForm, sf)
s = '
static text  One Two  Three Four
Five         Six Seven
Eight            Nine Ten
Eleven       Twelve
 Thirteen'.Trim().Detab()
		expectedForm = #(Form,
			(Static, 'static text  '), (custom_1, group: 1), (Static, ' '),
				(custom_2, group: 2), (Static, '  '), (custom_3, group: 3),
				(Static, ' '), (custom_4), nl,
			(custom_5, group: 0), (Static, '         '),
				(custom_6, group: 1), (Static, ' '), (custom_7, group: 2), nl,
			(custom_8, group: 0), (Static, '            '),
				(custom_9, group: 2), (Static, ' '), (custom_10, group: 3), nl,
			(custom_11, group: 0), (Static, '       '), (custom_12, group: 1), nl,
			(Static, ' '), (custom_13))
		.checkLayoutsPreserved(s, expectedForm, sf)
		// test that no duplicate fields added
		Assert(LayoutToForm("One One", sf) is: #(Form, (custom_1) (Static, ' ')))
		}

	checkLayoutsPreserved(s, expectedForm, sf)
		{
		form1 = LayoutToForm(s, sf)
		Assert(form1 is: expectedForm)
		s1 = LayoutToForm.Revert(form1, sf)
		form2 = LayoutToForm(s1, sf)
		s2 = LayoutToForm.Revert(form2, sf)
		Assert(form1 is: form2)
		Assert(s1 is: s2)
		Assert(s2 is: s)
		}

	Test_promptChanged()
		{
		sf = SelectFields()
		sf.AddField('f1', 'field1')
		sf.AddField('f2', 'field2')
		sf.AddField('f3', 'field3')
		sf.AddField('f4', 'field4')
		sf.AddField('f5', 'field5')
s = '
field1 field2
field3 field4
	   field5'.Trim().Detab()
		form = LayoutToForm(s, sf)
		Assert(form is: #(Form, (f1, group: 0), (Static, ' '), (f2, group: 1), nl,
			(f3, group: 0), (Static, ' '), (f4, group: 1), nl,
			(Static, '       '), (f5, group: 1)))
		Assert(LayoutToForm.Revert(form, sf) is: s)
		// rename f1, it gets longer
		sf.Fields.Remove('f1')
		sf.Delete('SelectFields_lookup')
		sf.AddField('f1', 'newfield1')
		renamed = LayoutToForm.Revert(form, sf)
		expected = '
newfield1 field2
field3    field4
		  field5'.Trim().Detab()
		Assert(renamed is: expected)

		// f1 prompt gets shorter
		sf.Fields.Remove('f1')
		sf.Delete('SelectFields_lookup')
		sf.AddField('f1', 'f1')
		renamed = LayoutToForm.Revert(form, sf)
		expected = '
f1     field2
field3 field4
	   field5'.Trim().Detab()
		Assert(renamed is: expected)

		// inlcuding more static text
		str = '
f1      field2
abcdefg field3
   aaa  field4
  field5'.Trim().Detab()
		form = LayoutToForm(str, sf)
		Assert(form is: #(Form, (f1), (Static, '      '), (f2, group: 0), nl,
			(Static, 'abcdefg '), (f3, group: 0), nl,
			(Static, '   aaa  '), (f4, group: 0), nl,
			(Static, '  '), (f5)))
		// f1 prompt gets longer, whitespace maintained
		sf.Fields.Remove('f1')
		sf.Delete('SelectFields_lookup')
		sf.AddField('f1', 'ff1')
		renamed = LayoutToForm.Revert(form, sf)
		expected = '
ff1      field2
abcdefg  field3
   aaa   field4
  field5'.Trim().Detab()
		Assert(renamed is: expected)
		}

	Test_deleteFields()
		{
		sf = SelectFields()
		sf.AddField('f1', 'field1')
		sf.AddField('f2', 'field2')
		sf.AddField('f3', 'field3')
		sf.AddField('f4', 'field4')
		sf.AddField('f5', 'field5')
s = '
field1 field2
field3 field4
	   field5'.Trim().Detab()
		form = LayoutToForm(s, sf)
		Assert(form is: #(Form, (f1, group: 0), (Static, ' '), (f2, group: 1), nl,
			(f3, group: 0), (Static, ' '), (f4, group: 1), nl,
			(Static, '       '), (f5, group: 1)))
		Assert(LayoutToForm.Revert(form, sf) is: s)

		// standard field deleted
		sf.Fields.Remove('f1')
		sf.Delete('SelectFields_lookup')
		expected = '
???    field2
field3 field4
	   field5'.Trim().Detab()
		removed = LayoutToForm.Revert(form, sf)
		Assert(removed is: expected)
		}

	Test_limit_height()
		{
		.testControl('(Editor)', height: 4)
		.testControl('(ScintillaAddonsEditor)', height: 2)
		.testControl('(ScintillaRichWordAddons)', height: 4)

		.testControl('(Editor, height: 8)', height: 4)
		.testControl('(ScintillaAddonsEditor, height: 8)', height: 4)
		.testControl('(ScintillaRichWordAddons, height: 8)', height: 4)

		.testControl('(EditorControl, height: 8)', height: 4)
		.testControl('(ScintillaAddonsEditorControl, height: 8)', height: 4)
		.testControl('(ScintillaRichWordAddonsControl, height: 8)', height: 4)

		.testControl('(Editor, height: 2)', height: 2)
		.testControl('(ScintillaAddonsEditor, height: 2)', height: 2)
		.testControl('(ScintillaRichWordAddons, height: 2)', height: 2)

		.testControl('(ScintillaRichWordAddons, ystretch: 0, height: 7)', height: 4)
		}

	testControl(ctrl, height)
		{
		editor = .TempTableName()
		prompt = .TempTableName()
		.MakeLibraryRecord([name: "Field_" $ editor, text: `Field_string
			{
			Prompt: '` $ prompt $ `'
			Control: ` $ ctrl $ `
			}`])
		form = Object(#Form, Object(editor))
		LayoutToForm.EnsureEditorHeightLimit(editor, form[1])
		Assert(form is: Object(#Form, Object(editor, :height, ystretch: 0)))
		}

	Test_Revert()
		{
		revert = LayoutToForm.Revert
		sf = SelectFields()
		sf.AddField('t1', 'Test1')
		sf.AddField('t2', 'Test2')
		sf.AddField('t3', 'Testing3')
		sf.AddField('t4', 'Test4')
		Assert(revert(#(), sf) is: '')

		Assert(revert(#(Form, (t1)), sf) is: 'Test1')
		Assert(revert(#(Form, (t1), (Static, ' '), (t2)), sf)
			is: 'Test1 Test2')
		Assert(revert(#(Form, (t1), (Static, ' '), (t2, group: 0) nl,
			(Static, '      '), (t3, group: 0)), sf) is: 'Test1 Test2\r\n      Testing3')
		Assert(revert(#(Form, (t1), (Static, '       '), (t2, group: 0) nl,
			(Static 'hello world '), (t3, group: 0)), sf)
				is: 'Test1       Test2\r\nhello world Testing3')
		Assert(revert(#(Form, (t1, group: 0), (t2, group: 1) nl,
			(Static 'hello world '), (t3, group: 1)), sf)
				is: 'Test1       Test2\r\nhello world Testing3')

		form = #(Form, (t1, group: 0), (Static, ' '), (t2), (Static, " hello world"), nl,
			(Static, '  '), (t3), (t4, group: 0))
		Assert(revert(form, sf) is: 'Test1 Test2 hello world\r\nTest4  Testing3')

		form = #(Form, (t1, group: 0), (Static, ' '), (t2), (Static, " hello world"), nl,
			(Static, '  '), (t3), (t4, group: 0), nl,
			(not_a_real_field, group: 0))
		Assert(revert(form, sf) is: 'Test1 Test2 hello world\r\nTest4  Testing3\r\n???')

		.MakeLibraryRecord([name: "Field_not_a_real_field",
			text: `Field_string { Prompt: 'Not Real' }`])
		Assert(revert(form, sf) is: 'Test1 Test2 hello world\r\nTest4  Testing3' $
			'\r\nNot Real')
		}

	Test_onlyCustomFields()
		{
				sf = SelectFields()
		sf.AddField('f1', 'field1')
		sf.AddField('f2', 'field2')
		sf.AddField('f3', 'field3')
		sf.AddField('f4', 'field4')
		sf.AddField('f5', 'field5')
		sf.AddField('custom_1', 'One')
s = '
field1 field2
field3 field4
	   field5
One'.Trim().Detab()
		form = LayoutToForm(s, sf)
		Assert(form is: #(Form, (f1, group: 0), (Static, ' '), (f2, group: 1), nl,
			(f3, group: 0), (Static, ' '), (f4, group: 1), nl,
			(Static, '       '), (f5, group: 1), nl,
			(custom_1, group: 0)))
		Assert(LayoutToForm.Revert(form, sf) is: s)

		form = LayoutToForm(s, sf, onlyCustomFields?:)
		Assert(form is: #(Form,
			(Static, 'field1 field2'), nl,
			(Static, 'field3 field4'), nl,
			(Static, '       field5'), nl,
			(custom_1)), msg: 'only custom fields')
		Assert(LayoutToForm.Revert(form, sf) is: s, msg: 'revert only custom fields')
		}
	}