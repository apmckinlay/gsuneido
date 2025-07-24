// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getAndRemoveExpiredChanges()
		{
		c = CustomizableOnServer
			{
			CustomizableOnServer_now() { return _now }
			}
		fn = c.CustomizableOnServer_getAndRemoveExpiredChanges

		_now = #20000101
		allChanges = Object(custom_fields: #(), custom_screen: #())
		Assert(fn(allChanges) is: #(custom_fields: #(), custom_screen: #()))

		allChanges = Object(
			custom_fields: Object(
				Object(key: 'hello', asof: #20000101)
				),
			custom_screen: Object(
				Object(field: 'world', asof: #20000101)
				))
		Assert(fn(allChanges) is: #(custom_fields: #('hello'), custom_screen: #('world')))
		Assert(allChanges is:  Object(
			custom_fields: Object(
				Object(key: 'hello', asof: #20000101)
				),
			custom_screen: Object(
				Object(field: 'world', asof: #20000101)
				)))

		_now = #20000101.1230
		allChanges = Object(
			custom_fields: Object(
				Object(key: 'hello', asof: #20000101.1222)
				Object(key: 'hello2', asof: #20000101.1225)
				Object(key: 'hello3', asof: #20000101.1227)
				Object(key: 'hello3', asof: #20000101.1228)
				Object(key: 'hello4', asof: #20000101.1228)
				),
			custom_screen: Object(
				Object(field: 'world', asof: #20000101.1222)
				Object(field: 'world2', asof: #20000101.1225)
				Object(field: 'world3', asof: #20000101.1227)
				Object(field: 'world3', asof: #20000101.1228)
				Object(field: 'world4', asof: #20000101.1228)
				))
		Assert(fn(allChanges)
			is: #(custom_fields: #('hello3', 'hello4'),
				custom_screen: #('world3', 'world4')))
		Assert(allChanges is: Object(
			custom_fields: Object(
				Object(key: 'hello3', asof: #20000101.1227)
				Object(key: 'hello3', asof: #20000101.1228)
				Object(key: 'hello4', asof: #20000101.1228)
				),
			custom_screen: Object(
				Object(field: 'world3', asof: #20000101.1227)
				Object(field: 'world3', asof: #20000101.1228)
				Object(field: 'world4', asof: #20000101.1228)
				)))

		_now = #20000101.1300
		Assert(fn(allChanges) is: #(custom_fields: #(), custom_screen: #()))
		Assert(allChanges is: #(custom_fields: #(), custom_screen: #()))
		}
	}