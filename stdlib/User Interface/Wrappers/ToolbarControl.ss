// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: "Toolbar"
	Ystretch: 0
	New(@buttons)
		{
		.horz = .Vert.HorzEqualHeight
		.commands = .Window.Commands()
		.dropButtons = Object()
		for b in buttons.Values(list:)
			{
			if b is ""
				.horz.Append('EtchedVertLine')
			else if b is ">"
				.horz.Append('Fill')
			else
				{
				if Object?(b)
					{
					if b.GetDefault('drop', false)
						.dropButtons.Add(b[0])
					b = b[0]
					}
				.addtoolbarbutton(b)
				}
			}
		}
	Controls:
		(Vert xstretch: 1 ystretch: 0
			(EtchedLine before: 0 after: 0)
			(HorzEqualHeight ystretch: 1))
	addtoolbarbutton(cmd)
		{
		command = .commands.GetDefault(cmd, Record(bitmap: cmd))
		accel = command.accel
		tip = TranslateLanguage((command.help is '' ? cmd : command.help).Tr("_", " "))
		if cmd is 'Copy' and accel is ''
			accel = 'Ctrl+C'
		if accel > ""
			tip $= " (" $ TranslateLanguage(accel) $ ")"
		if .dropButtons.Has?(cmd)
			tip $= ", right click for more options"
		singleLetter? = command.bitmap.Size() is 1
		name = singleLetter? ? command.bitmap : command.bitmap.Lower()
		image = .checkUnknownIcon(name $ ".emf")
		.horz.Append(Object("EnhancedButton", command: cmd, :image, :tip,
			imagePadding: singleLetter? ? 0.05 : 0.15 /*=padding*/,
			mouseEffect:, ystretch: 0, name: cmd))
		}
	checkUnknownIcon(image)
		{
		if Query1('imagebook', name: image) isnt false or
			IconFont().MapToCharCode(image) isnt false
			return image
		SuneidoLog('ERROR: (CAUGHT) unknown icon: ' $ image, calls:,
			caughtMsg: 'fall back to triangle-warning')
		return 'triangle-warning'
		}
	GetReadOnly()		// read-only not applicable to toolbar
		{ return true }
	EnhancedRButtonDown(rc, source)
		{
		if .dropButtons.Has?(buttonName = source.GetButtonName())
			.Send("Drop_" $ buttonName, rc, :source)
		}
	EnableButton(button, enable)
		{
		button = "On_" $ ToIdentifier(button.Trim())
		buttons = .horz.GetChildren()
		for b in buttons
			if b.Member?(#GetCommand) and b.GetCommand() is button
				{
				b.SetEnabled(enable)
				return true
				}
		return false
		}
	}
