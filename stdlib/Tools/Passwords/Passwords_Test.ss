// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Pwnage()
		{
		cl = Passwords
			{
			Passwords_getPwnedList(unused)
				{
				return "F1ACE0EC3AB1DC9864E65339CB239ACE358:3\r\n" $
					"F1DCFD56F18E1470520A8F27981968F3B90:4\r\n" $
					"F2726B76AA110BCB5D8F4F4BA475EF90599:3\r\n" $
					"F2F322BD56610D4E92FA69389A301B57269:1\r\n" $
					"F30CDFB4A9533B6DA58233FF9679B555D8C:1\r\n" $
					"F30DDADBCC105DCB63F012E1C3D1EB0360C:5\r\n" $
					"F33F8BF412ABA3456BA34B393D9B1A2279A:2\r\n" $
					"F3669FA1C5937091C5B6750A23779A63FA9:2\r\n" $
					"F3775E7DD264FE23D9C1E15360DCA2AB0C9:2\r\n" $
					"F40A222A474262718AEEC4AC7B8F696834F:1\r\n" $
					"F446212A11E757E933EECAA514DB635E9AA:24\r\n" $
					"F4D2341104A2C6955D28A7E94F88D0D09FE:3\r\n" $
					"F187EBB7080BD75AAC9160214E6B1E49F7D:14226\r\n" $
					"F4DEF69DC640227EEAA708C92EC2051F7F4:13\r\n" $
					"F55363341C0B8C28844A782F1031C291E15:5\r\n" $
					"F5B5A9A41E4DE0B1D8C6DFAE37145E40448:1\r\n" $
					"F5C853C88F10777C46893C30EF287C46141:1\r\n" $
					"F5FA4FC835F65E7A8819ED5664FB5F5478C:2\r\n" $
					"F682695505B62C1F11A8D276F53C1778033:7"
				}
			}
		Assert(cl.Passwords_pwnedPassword?('P@ssw0rd1') isnt: "", msg: 'pwned')

		cl = Passwords
			{
			Passwords_getPwnedList(unused)
				{
				return "C538251EA611EBC1286E59048D7A558EEA9:4\r\n" $
					"C5B51F137107AED33ABC8237A49DE224170:1\r\n" $
					"F8D86B5E96C28941D007576E84B124FDB4C:3\r\n" $
					"FAFA56560D8103F3CDBEA1D9051C671F6F0:1\r\n" $
					"FB197B4D833828A72FAA042CCF2D377E0A3:1\r\n" $
					"FB340B774121574CB5D250B4DB1440D55CA:1\r\n" $
					"FB7441F8BFBE3AD5A730976B951EFCD9786:1\r\n" $
					"FB8AC54CACB760A67EAA21748A5D50FC121:24\r\n" $
					"FBF6E51C9C55C43EE6D289555521F650A0E:2\r\n" $
					"FC06AAF703F2299E64ACE678CC0DC34B705:1\r\n" $
					"FCDA46D860715D57EBCCCD9F4152108B35F:16\r\n" $
					"FE1B27820797232DDF5B9296546D1D4D01B:1\r\n" $
					"FE92E7939329AFC83FBCBED302AFAD593ED:61\r\n" $
					"FE99093409383848888D14CD9BD72EDA1BA:1\r\n" $
					"FF5539D9188FA4A9B97B3C9B078FF890D6E:3"
				}
			}
		Assert(cl.
			Passwords_pwnedPassword?('BelieveItOrNotThisIsNotAPwnedPassword') is: '',
				msg: 'not pwned')

		cl = Passwords
			{
			Passwords_getPwnedList(unused) { "" }
			Passwords_logError(unused) { }
			}
		Assert(cl.Passwords_pwnedPassword?('test') is: '', msg: 'empty content')

		cl = Passwords
			{
			Passwords_getPwnedList(unused)
				{ throw "curl: (6) Could not resolve host: api.pwnedpasswords.com" }
			Passwords_logError(unused) { }
			}
		Assert(cl.Passwords_pwnedPassword?('test') is: '', msg: 'connection error')
		}
	Test_passwordComplexity()
		{
		func = Passwords.Passwords_passwordComplexity
		Assert(func('') isnt: "", msg: "empty")
		Assert(func('justgo4') isnt: "", msg: "justgo4")
		Assert(func('text123') isnt: "", msg: "text123")
		Assert(func('2til1ow') isnt: "", msg: "2til1ow")
		Assert(func('p2sswor') isnt: "", msg: "p2sswor")
		Assert(func('PSW0RD1') isnt: "", msg: "PSW0RD1")
		Assert(func('PASSword1!') is: '', msg: "PASSword1!")
		Assert(func('PASSwordd!') is: '', msg: "PASSwordd!")
		Assert(func('justgoingforlength') is: '', msg: "justgoingforlength")
		Assert(func('P@ssw0rd1') is: '', msg: "ssw0rd1")
		Assert(func('!@#$BACK1') is: '', msg: "!@#$BACK1")
		}

	Test_genPasswords()
		{
		func = Passwords.Passwords_genPassword
		pw = func(50, 'Easy to read')
		Assert(pw isSize: 50)
		Assert(pw.Alpha?(), msg: 'easy to read')

		pw = func(50, 'Easy to spell')
		Assert(pw isSize: 50)
		Assert(pw.Tr('0oOiIlL|1"\'`+') is: pw)

		pw = func(14, 'All Characters')
		Assert(pw, isSize: 14)

		pw = func(4, 'Pass Phrase')
		count = 0
		pw.ForEachMatch('[[:upper:]]') { |unused| ++count }
		Assert(count, is: 4)
		}

	Test_generateRetries()
		{
		mock = Mock(Passwords)
		mock.When.genPassword([anyArgs:]).Return('test')
		mock.When.pwnedPassword?([anyArgs:]).Return(false)
		mock.When.validateStrength([anyArgs:]).CallThrough()
		mock.When.lastTryPassword([anyArgs:]).CallThrough()
		mock.Eval(Passwords.Generate, 8 /*=length*/, 'Easy to read')
		mock.Verify.Times(51).validateStrength([anyArgs:])
		}

	Test_lastTryPassword()
		{
		fn = Passwords.Passwords_lastTryPassword

		for size in #(5, 8, 10, 12, 15)
			{
			pw = fn(size)
			Assert(pw isSize: size)
			types = Passwords.Passwords_getUniqueTypes(pw)
			Assert(types.number, msg: 'number')
			Assert(types.symbol, msg: 'symbol')
			}
		}
	Test_GenerateSimple()
		{
		fn = Passwords.GenerateSimple

		Assert(fn(0, 0, 0) is: '')

		Assert(fn() isSize: 15)

		.assertComponentCounts(fn, 7, 0, 0)
		.assertComponentCounts(fn, 7, 7, 0)
		.assertComponentCounts(fn, 7, 0, 7)
		for ..20
			.assertComponentCounts(fn, Random(6), Random(6), Random(6))
		}

	assertComponentCounts(fn, pwLength, symbols, numbers)
		{
		password = fn(pwLength, symbols, numbers)
		symbolCount = numberCounts = letterCount = 0
		for c in password
			if c.Alpha?()
				letterCount++
			else if c.Number?()
				numberCounts++
			else
				symbolCount++
		Assert(password isSize: Max(pwLength, symbols + numbers))
		Assert(symbolCount 	is: symbols)
		Assert(numberCounts is: numbers)
		Assert(letterCount 	is: Max(0, pwLength - symbols - numbers))
		}
	}