// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		.dialog('1 of 4')
		.dialog('2 of 4', closeButton?:)
		.dialog('3 of 4', resizable:)
		.dialog('4 of 4', resizable:, closeButton?:)
		}

	dialog(title, resizable = false, closeButton? = false)
		{
		title = 'DialogTester ' $ title
		ctrl = [#Vert,
			[#Static,
				(resizable ? 'SHOULD' : 'should NOT') $ " be resizable"],
			[#Static,
				(closeButton? ? 'SHOULD' : 'should NOT') $ " have 'X' close button"],
			#Skip
			#(Static "Initial focus should be on the field (typing should work)")
			#(Static "Switching to a different window and back shouldn't hide dialog")
			#(Static "Switching to a different window and back shouldn't keep focus")
			#(Static "ESC or Cancel button should close dialog and return false")
			#(Static "OK button should be highlighted as the default"),
			#(Static "Tabbing through controls should not lose default button highlight"),
			[#Static resizable ? "" : "ENTER should trigger default button"],
			#Skip
			resizable ? 'Editor' : 'Field',
			#Skip,
			#OkCancel]
		result = Dialog(0, [.controller, ctrl], :closeButton?, :title)
		Print('DialogTester', :result)
		}
	controller: Controller {
		On_OK()
			{
			.Window.Result('OK pressed')
			}
		}
	}