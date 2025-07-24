// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New()
		{
		super()
		.spinner = .FindControl('pass_length')
		.options = .FindControl('pass_options')
		.generateField = .FindControl('gen_pass')
		.generateField.SetReadOnly(true)
		.Data.AddObserver(.Record_changed)
		}

	Controls()
		{
		return Object('Record',
			  Object('Vert'
				#('Field', '', width: 22, name: 'gen_pass')
				Object('Horz'
					#('Button', 'Generate Password'),
					#('Button', 'Use')
					'Skip',
					Object('Pair',
						#('Static', 'Length'),
						Object('Spinner', width: 1
							rangefrom: .minPw['Pass Phrase'],
							rangeto: .minPw['All Characters'],
							set: .minPw['Pass Phrase'],
							name: 'pass_length')
						),
					),
				#('RadioButtons',
					'Pass Phrase - word combo',
					'Easy to read - no numbers, no symbols',
					'Easy to spell - no ambiguity (o,O,0,1,I,l,|)',
					'All Characters', skip: 8, name: 'pass_options')
				))
		}

	Record_changed(member)
		{
		if member is 'pass_options'
			.OptionChanged()
		}

	On_Generate_Password()
		{
		if '' is pwLength = .spinner.Get()
			return ''
		try // this needs to run from an un-authorized client.
			BookLog('PasswordGenerator: Generate - Option: ' $
				.options.Get().BeforeFirst(' -'))
		catch(unused, 'not authorized')
			{ }
		option = .options.Get().BeforeFirst(' -')
		.generateField.Set(Passwords.Generate(pwLength, option))
		}

	On_Use()
		{
		if '' is pw = .generateField.Get()
			return
		try // this needs to run from an un-authorized client.
			BookLog('PasswordGenerator: Use password')
		catch(unused, 'not authorized')
			{ }
		ClipboardWriteString(pw)
		InfoWindowControl('Generated password copied to clipboard',
			titleSize: 0, marginSize: 7, autoClose: 1)
		.Send('Use_Generated_Password', pw)
		}

	maxPw: 25
	minPw: #('Easy to read': 8,
		'Easy to spell': 8,
		'All Characters': 6,
		'Pass Phrase': 3)
	OptionChanged()
		{
		option = .options.Get().BeforeFirst(' -')
		if option is 'Pass Phrase' /// Allowing 6 words as their max when using passPhase
			.spinner.OverrideRanges(.minPw[option], .minPw['All Characters'], setLow:)
		else
			.spinner.OverrideRanges(.minPw[option], .maxPw)
		}
	}