// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_tabControl()
		{
		mock = Mock(TabsControl)
		mock.When.tabControl([anyArgs:]).CallThrough()

		// Test near empty parameters (excluding instance variables)
		controls = Object()
		tabs = mock.tabControl(controls, [close_button:])
		Assert(tabs.themed)
		Assert(tabs.close_button)
		Assert(tabs.buttonTip is: 'Add Tab')
		Assert(tabs.orientation is: 'top')
		Assert(mock.TabsControl_alternativePos is: false)
		Assert(mock.TabsControl_customizableTabs is: #())

		controls.Add(#(Customizable, Tab: Custom))
		tabs = mock.tabControl(controls,
			[addTabButton?:, buttonTip: 'Test This', orientation: 'right'])
		Assert(tabs.close_button is: false)
		Assert(tabs.buttonTip is: 'Test This')
		Assert(mock.TabsControl_alternativePos)
		Assert(mock.TabsControl_customizableTabs is: #(Custom))

		controls.Add(#(Vert,
			#(Form, #(other_controls), Customizable),
			Tab: Manager)
			)
		tabs = mock.tabControl(controls, [addTabButton?: false])
		Assert(mock.TabsControl_customizableTabs is: #(Custom, Manager))

		controls.Add(#(Vert,
			#(Form,
				#(other_controls,
					#(other_controls,
						#(other_controls,
							#(other_controls,
								#(other_controls,
									#(other_controls,
										#(other_controls,
											#(other_controls,
												#(other_controls,
													#(other_controls, Customizable)
				)))))))))),
			Tab: Special)
			)
		tabs = mock.tabControl(controls, [addTabButton?: false, orientation: 'left'])
		Assert(mock.TabsControl_customizableTabs is: #(Custom, Manager, Special))
		Assert(mock.TabsControl_alternativePos is: false)

		controls.Add(#(Vert,
			#(Form,
				#(other_controls,
					#(other_controls)), nl
				#(other_controls,
					Customizable)),
			Tab: Last)
			)
		tabs = mock.tabControl(controls, [orientation: 'bottom'])
		Assert(mock.TabsControl_customizableTabs is: #(Custom, Manager, Special, Last))
		Assert(mock.TabsControl_alternativePos)

		Assert({ mock.tabControl(controls, [orientation: 'other']) }
			throws: 'unhandled switch value')
		}

	Test_Resize()
		{
		mock = Mock(TabsControl)
		mock.When.Resize([anyArgs:]).CallThrough()
		mock.When.resize([anyArgs:]).Do({ })
		mock.TabsControl_resizing = mock.TabsControl_tabs_changing = true
		mock.TabsControl_being_constructed = 1

		mock.Resize(0, 0, 0, 0)
		mock.Verify.Never().resize(0, 0, 0, 0)

		mock.TabsControl_resizing = false
		mock.Resize(0, 0, 0, 0)
		mock.Verify.Never().resize(0, 0, 0, 0)

		mock.TabsControl_tabs_changing = false
		mock.Resize(0, 0, 0, 0)
		mock.Verify.Never().resize(0, 0, 0, 0)

		mock.TabsControl_being_constructed = false
		mock.Resize(1, 1, 1, 1)
		mock.Verify.resize(1, 1, 1, 1)

		mock.TabsControl_resizing = true
		mock.Resize(0, 0, 0, 0)
		mock.Verify.Never().resize(0, 0, 0, 0)

		mock.TabsControl_resizing = false
		mock.TabsControl_tabs_changing = true
		mock.Resize(0, 0, 0, 0)
		mock.Verify.Never().resize(0, 0, 0, 0)

		mock.TabsControl_tabs_changing = false
		mock.Resize(2, 2, 2, 2)
		mock.Verify.resize(2, 2, 2, 2)
		}
	}