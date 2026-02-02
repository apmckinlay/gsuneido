// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cl: ScintillaAddonsCommandManager
		{
		ScintillaAddonsCommandManager_keyPressed?(@unused)
			{
			if not Object?(_keyPressedResults)
				return _keyPressedResults
			return _keyPressedResults.PopFirst()
			}
		ScintillaAddonsCommandManager_logDuplicateLevels(option)
			{
			_logs.Add(option)
			}
		}
	Test_main()
		{
		inst = new .cl
		_logs = Object()
		_keyPressedResults = false

		inst.Set(#())
		Assert(inst.ScintillaAddonsCommandManager_commands isSize: 0)
		Assert(commandOb = inst.BuildCommandOb('A') is: #('a'))
		Assert(inst.GetMethod(commandOb) is: false)
		Assert(commandOb = inst.BuildCommandOb(false) is: #('false'))
		Assert(inst.GetMethod(commandOb) is: false)
		Assert(_logs isSize: 0)

		inst.Set(#('Command0_NoAccels', 'Command1_NoAccels'))
		Assert(inst.ScintillaAddonsCommandManager_commands isSize: 0)
		Assert(commandOb = inst.BuildCommandOb('') is: #(''))
		Assert(inst.GetMethod(commandOb) is: false)
		Assert(_logs isSize: 0)

		inst.Set(#(
			'Command0_NoAccels',
			'Command1\tF1',
			'Command2\tF2',
			'Command3\tF3',
			'Command4\tF4'))
		Assert(inst.ScintillaAddonsCommandManager_commands
			equalsSet: #(
				#(method: 'On_Command1', command: #('f1')),
				#(method: 'On_Command2', command: #('f2')),
				#(method: 'On_Command3', command: #('f3')),
				#(method: 'On_Command4', command: #('f4'))))
		Assert(commandOb = inst.BuildCommandOb('F1') is: #('f1'))
		Assert(inst.GetMethod(commandOb) is: 'On_Command1')
		Assert(commandOb = inst.BuildCommandOb('F2') is: #('f2'))
		Assert(inst.GetMethod(commandOb) is: 'On_Command2')
		Assert(commandOb = inst.BuildCommandOb('F3') is: #('f3'))
		Assert(inst.GetMethod(commandOb) is: 'On_Command3')
		Assert(commandOb = inst.BuildCommandOb('F4') is: #('f4'))
		Assert(inst.GetMethod(commandOb) is: 'On_Command4')
		Assert(commandOb = inst.BuildCommandOb('F5') is: #('f5'))
		Assert(inst.GetMethod(commandOb) is: false)
		Assert(commandOb = inst.BuildCommandOb(false) is: #('false'))
		Assert(inst.GetMethod(commandOb) is: false)
		Assert(_logs isSize: 0)
		inst.Destroy()
		Assert(commandOb = inst.BuildCommandOb('F1') is: #('f1'))
		Assert(inst.GetMethod(commandOb) is: false)
		}

	Test_main_WithKeyPressed()
		{
		inst = new .cl
		_logs = Object()

		// Some duplicate commands are processed, resulting in duplicate logging
		inst.Set(#(
			'Command1\tShift+Ctrl+F1',
			'Command2\tShift+Ctrl+F1',
			'Command3\tF2',
			'Command4\tF2',
			'Command5\tCtrl+C'))
		Assert(inst.ScintillaAddonsCommandManager_commands
			equalsSet: #(
				#(method: 'On_Command1', command: #('ctrl', 'f1', 'shift')),
				#(method: 'On_Command3', command: #('f2')),
				#(method: 'On_Command5', command: #('c', 'ctrl'))))
		_keyPressedResults = Object(
			/*shift: */ true, 	/*ctrl: */ true, 	/*alt: */ false, // On_Command1
			/*shift: */ false, 	/*ctrl: */ false, 	/*alt: */ false, // On_Command3
			/*shift: */ false, 	/*ctrl: */ true, 	/*alt: */ false, // On_Command5
			/*shift: */ false, 	/*ctrl: */ false, 	/*alt: */ true,  // Not defined
			/*shift: */ false, 	/*ctrl: */ true, 	/*alt: */ false, // F1 no shift
			/*shift: */ false, 	/*ctrl: */ false, 	/*alt: */ false, // F2 post destroy
			)
		Assert(commandOb = inst.BuildCommandOb('F1') is: #('ctrl', 'f1', 'shift'))
		Assert(inst.GetMethod(commandOb) is: 'On_Command1')
		Assert(commandOb = inst.BuildCommandOb('F2') is: #('f2'))
		Assert(inst.GetMethod(commandOb) is: 'On_Command3')
		Assert(commandOb = inst.BuildCommandOb('C') is: #('c', 'ctrl'))
		Assert(inst.GetMethod(commandOb) is: 'On_Command5')
		Assert(commandOb = inst.BuildCommandOb('E') is: #('alt', 'e'))
		Assert(inst.GetMethod(commandOb) is: false)
		Assert(_logs isSize: 2)
		Assert(_logs[0] has: 'Command2\tShift+Ctrl+F1')
		Assert(_logs[1] has: 'Command4\tF2')
		Assert(commandOb = inst.BuildCommandOb('F1') is: #('ctrl', 'f1'))
		Assert(inst.GetMethod(commandOb) is: false)
		inst.Destroy()
		Assert(commandOb = inst.BuildCommandOb('F2') is: #('f2'))
		Assert(inst.GetMethod(commandOb) is: false)
		}
	}