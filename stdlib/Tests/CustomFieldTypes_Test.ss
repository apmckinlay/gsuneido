// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		.user = Suneido.User
		}
	Test_compatibility()
		{
		baseType = #(name: 'baseType', base: 'type1', compatible: 't1', customize:)
		nonCompBaseType = #(name: 'nonCompatibleBaseType', base: 'type1', customize:)
		otherType = #(name: 'otherType', base: 'type2', compatible: 't2', customize:)
		compOtherType = #(name: 'compatibleOtherType', base: 'type2',
			compatible: function (_type = 't2') { return type }, customize:)

		_types = Object(baseType, nonCompBaseType, otherType, compOtherType)
		cl = CustomFieldTypes { CustomFieldTypes_getTypes() { return _types } }

		Assert(cl() is: Object(baseType, nonCompBaseType, otherType, compOtherType))
		Assert(cl(filterBy: 'type1') is: Object(baseType))
		Assert(cl(filterBy: 'type2') is: Object(otherType, compOtherType))
		_type = 't3'
		Assert(cl(filterBy: 'type2') is: Object(otherType))
		}

	Test_string_compatibility()
		{
		Suneido.User = 'default'
		types = CustomFieldTypes(filterBy: 'string_custom')
		Assert(types isSize: 2)
		Assert(types.HasIf?({ it.name is "Text, single line" }))
		Assert(types.HasIf?({ it.name is "Text, multi line" }))

		Suneido.User = 'axon'
		types = CustomFieldTypes(filterBy: 'string_custom')
		Assert(types isSize: 1)
		Assert(types.HasIf?({ it.name is "Text, single line" }))
		}

	Teardown()
		{
		Suneido.User = .user
		super.Teardown()
		}
	}
