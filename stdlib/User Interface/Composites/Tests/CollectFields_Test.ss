// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(CollectFields(#()) isSize: 0)
		Assert(CollectFields(#(#(), #(), #(#()))) isSize: 0)
		layout0 = #(Vert, #(notAField), #(NotAControl), #(Vert, notAField, #(notAField)))
		Assert(CollectFields(layout0) isSize: 0)

		.MakeLibraryRecord([name: showHide = .TempName(),
			text: `class
				{
				CallClass()
					{ return true }
				Hide()
					{ return false}
				}` ])
		layout1 = Object(#Form,
			#(sulog_timestamp, group: 0),
				#(sulog_message, group: 1),
				Object(#ShowHide, #sulog_params, showHide, group: 3), #nl,
			Object(#Horz, #sulog_option
				#sulog_need_to_send?,
				Object(#ShowHide, #sulog_calls,
					{ Global(showHide).Hide() }), group: 0), #nl,
			#(Pair, sulog_user, sulog_session_id, group: 0))
		fields = CollectFields(Object(layout0, layout1))
		Assert(fields isSize: 7)
		Assert(fields has: #sulog_timestamp)
		Assert(fields has: #sulog_message)
		Assert(fields has: #sulog_params)
		Assert(fields has: #sulog_option)
		Assert(fields has: #sulog_need_to_send?)
		Assert(fields has: #sulog_user)
		Assert(fields has: #sulog_session_id)
		Assert(fields hasnt: #sulog_calls)

		layout2 = Object(#Vert,
			Object(#Horz, #sulog_message,
				Object(#Vert, #sulog_user, #(sulog_need_to_send?),
					Object(#Horz, #(Field, name: sulog_timestamp),
						Object(#Vert, #sulog_locals,
							Object(#Horz,
								Object(#ShowHideControl, #sulog_calls, showHide)))))))
		fields = CollectFields(Object(layout0, layout1, layout2, Container))
		Assert(fields isSize: 9)
		Assert(fields has: #sulog_timestamp)
		Assert(fields has: #sulog_message)
		Assert(fields has: #sulog_params)
		Assert(fields has: #sulog_option)
		Assert(fields has: #sulog_need_to_send?)
		Assert(fields has: #sulog_user)
		Assert(fields has: #sulog_session_id)
		Assert(fields has: #sulog_calls)
		Assert(fields has: #sulog_locals)

		prefix = .TempName().Lower() $ '_'
		address1 = .addressField(prefix, #Address1)
		address2 = .addressField(prefix, #Address2)
		city = .addressField(prefix, #City)
		layout3 = Object(#Vert, Object(#Address, prefix, extra_field1: #sulog_calls))
		fields = CollectFields(Object(layout0, layout1, Container, layout3))
		Assert(fields isSize: 11)
		Assert(fields has: #sulog_timestamp)
		Assert(fields has: #sulog_message)
		Assert(fields has: #sulog_params)
		Assert(fields has: #sulog_option)
		Assert(fields has: #sulog_need_to_send?)
		Assert(fields has: #sulog_user)
		Assert(fields has: #sulog_session_id)
		Assert(fields has: #sulog_calls)
		Assert(fields has: address1)
		Assert(fields has: address2)
		Assert(fields has: city)

		layout = Object(layout0, layout1, Container, layout3, Hide?:)
		Assert(CollectFields(Object(layout)) isSize: 0)

		layout1.Hide? = layout3.Hide? = true
		controls = Object(#AccessControl,
			Object(#Tabs, layout0, layout1),
			layout3,
			#(sulog_session_id)
			#(sulog_timestamp, Hide?:))
		Assert(CollectFields(controls) is: #(sulog_session_id))

		prefix = .TempName().Lower() $ '_'
		.MakeLibraryRecord([name: `Field_` $ name = prefix $ 'name',
			text: `Field_string { Prompt: 'Test Name' }`])
		.MakeLibraryRecord([name: `Field_` $ abbrev = prefix $ 'abbrev',
			text: `Field_string { Prompt: 'Test Abbrev' }`])
		Assert(CollectFields(Object(Object(#Vert, Object(#NameAbbrev, prefix))))
			equalsSet: Object(name, abbrev))
		}

	addressField(prefix, prompt)
		{
		.MakeLibraryRecord([name: `Field_` $ field = prefix $ prompt.Lower(),
			text: `Field_string { Prompt: '` $ prompt $ `' }`])
		return field
		}

	Test_skipTab?()
		{
		_user = 'admin'
		c = CollectFields
			{
			CollectFields_user()
				{
				return _user
				}
			}
		fn = c.CollectFields_skipTab?
		Assert(fn(#(Tab: 'test')) is: false)
		Assert(fn(#(Tab: 'Hidden Options')))
		Assert(fn(#(Tab: 'Developer Options')))
		Assert(fn(#(Tab: 'test2', Hide?:)))
		Assert(fn(#(Tab: 'test2', Hide?: false)) is: false)

		_user = 'default'
		Assert(fn(#(Tab: 'test')) is: false)
		Assert(fn(#(Tab: 'Hidden Options')) is: false)
		Assert(fn(#(Tab: 'Developer Options')) is: false)
		Assert(fn(#(Tab: 'test2', Hide?:)))
		Assert(fn(#(Tab: 'test2', Hide?: false)) is: false)

		_user = 'axon'
		Assert(fn(#(Tab: 'test')) is: false)
		Assert(fn(#(Tab: 'Hidden Options')) is: false)
		Assert(fn(#(Tab: 'Developer Options')))
		Assert(fn(#(Tab: 'test2', Hide?:)))
		Assert(fn(#(Tab: 'test2', Hide?: false)) is: false)
		}

	Test_collect_with_path()
		{
		field = .TempTableName()
		.MakeLibraryRecord([name: "Field_" $ field,
			text: `Field_config_cols
				{
				Prompt: "Config Columns"
				}` ])
		layout = Object('Tabs',
			Object(#Form,
				#(sulog_timestamp, group: 0),
				#(sulog_message, group: 1)
				#(Button 'Hello World') #nl
				field, Tab: 'Setting'
				))
		fields = CollectFields(layout, path?:)

		Assert(fields has: #((type: "Tab", section: "Setting")))
		Assert(fields has: #((type: "Tab", section: "Setting"),
			(name: "sulog_timestamp", type: "Field", section: "Timestamp")))
		Assert(fields has: #((type: "Tab", section: "Setting"),
			(name: "sulog_message", type: "Field", section: "Message")))
		Assert(fields has: #((type: "Tab", section: "Setting"),
			(name: "Hello_World", type: "Button", section: "Hello World")))

		Assert(fields has: Object(Object(type: "Tab", section: "Setting"),
			Object(name: field, type: "Field", section: "Config Columns")))

		Assert(fields has: Object(Object(type: "Tab", section: "Setting"),
			Object(name: field, type: "ResetButton", section: "Config Columns > Reset")))

		Assert(fields isSize: 6)
		}
	}
