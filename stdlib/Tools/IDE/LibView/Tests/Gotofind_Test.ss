// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_FindItems()
		{
		cl = Gotofind
			{ Gotofind_libraryTables() { return #(stdlib, Testlib, otherlib) } }
		m = cl.FindItems
		Assert(m(#DefinitionFormat, true) is: #(DefinitionFormat))
		Assert(m(#DefinitionFormat)
			equalsSet: #(DefinitionFormat, '<lib>_DefinitionFormat'))

		Assert(m(#DefinitionControl, true) is: #(DefinitionControl))
		Assert(m(#DefinitionControl)
			equalsSet: #(DefinitionControl, '<lib>_DefinitionControl'))

		Assert(m(#Trigger_Definition) is: #(Trigger_Definition))

		Assert(m(#Field_definition, true) is: #(Field_definition))
		Assert(m(#Field_definition)
			equalsSet: #(Field_definition, Rule_definition, Rule_definition__protect))

		Assert(m(#Rule_definition, true) is: #(Rule_definition))
		Assert(m(#Rule_definition)
			equalsSet: #(Field_definition, Rule_definition, Rule_definition__protect))

		Assert(m(#Definition, true) is: #(Definition))
		Assert(m(#Definition)
			equalsSet: #(Definition, Field_definition, Rule_definition,
				Rule_definition__protect, Trigger_Definition, DefinitionFormat,
				DefinitionControl, DefinitionComponent, '<lib>_Definition',
				Table_Definition))

		Assert(m(#Stdlib_Definition, true) is: #(Stdlib_Definition))
		Assert(m(#Stdlib_Definition)
			equalsSet: #(Stdlib_Definition, Field_definition, Rule_definition,
				Rule_definition__protect, Trigger_Definition, DefinitionFormat,
				DefinitionControl, DefinitionComponent, '<lib>_Definition',
				Table_Definition))
		}

	Test_onlyExactMatches?()
		{
		m = Gotofind.Gotofind_onlyExactMatches?
		Assert(m(#A_Average_Definition) is: false, msg: 'A_Average_Definition')
		Assert(m(#Trigger_Definition), msg: 'Trigger_Definition')
		Assert(m(#Field_definition), msg: 'Field_definition')
		Assert(m(#Rule_definition), msg: 'Rule_definition')
		Assert(m(#Rule_Control), msg: 'Rule_Control')
		Assert(m(#DefinitionControl) is: false, msg: 'DefinitionControl')
		Assert(m(#DefinitionFormat) is: false, msg: 'DefinitionFormat')
		Assert(m(#DefinitionOther) is: false, msg: 'DefinitionOther')
		}

	Test_addProtectDefinitions?()
		{
		m = Gotofind.Gotofind_addProtectDefinitions?
		Assert(m(#Not_a_Rule_or_Field, false) is: false, msg: 'Not_a_Rule_or_Field false')
		Assert(m(#Not_a_Rule_or_Field, true) is: false, msg: 'Not_a_Rule_or_Field true')

		Assert(m(#Rule_definition, false), msg: 'Rule_definition false')
		Assert(m(#RuleDefinition, false) is: false, msg: 'RuleDefinition false')
		Assert(m(#Rule_definition, true) is: false, msg: 'Rule_definition true')

		Assert(m(#Field_definition, false), msg: 'Field_definition false')
		Assert(m(#FieldDefinition, false) is: false, msg: 'FieldDefinition false')
		Assert(m(#Field_definition, true) is: false, msg: 'Field_definition true')
		}

	Test_libraryPrefix?()
		{
		cl = Gotofind
			{ Gotofind_libraryTables() { return #(stdlib, Testlib, otherlib) } }
		m = cl.Gotofind_libraryPrefix?
		Assert(m(#Notlib_prefix) is: false, msg: 'Notlib_prefix')
		Assert(m(#Stdlib_prefix), msg: 'Stdlib_prefix')
		Assert(m(#Testlib_prefix), msg: 'Testlib_prefix')
		Assert(m(#Otherlib_prefix), msg: 'Otherlib_prefix')
		}
	}