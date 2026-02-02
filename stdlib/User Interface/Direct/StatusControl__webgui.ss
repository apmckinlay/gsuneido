// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 'Status'
	ComponentName: 'Status'
	ComponentArgs: #()
	Xstretch: 1
	New(text = "")
		{
		.Set(text)
		.status = 'normal'
		}

	statusLimit: 500
	Set(text, normal = false, warn = false, invalid = false)
		{
		.setColor(normal, warn, invalid)
		text = String(text)
		.untranslated = text
		.text = TranslateLanguage(text).Ellipsis(.statusLimit, atEnd:)
		.set()
		}

	prevText: false
	set()
		{
		text = .text is '' ? .defaultMsg : .text
		if .prevText isnt text
		.Act('Set', .prevText = text)
		}

	setColor(normal, warn, invalid)
		{
		if invalid
			.SetValid(false)
		else if warn
			.SetWarning()
		else if normal
			.SetValid()
		}

	defaultMsg: ""
	SetDefaultMsg(text)
		{
		.defaultMsg = TranslateLanguage(text).Ellipsis(.statusLimit, atEnd:)
		.set()
		}

	SetValid(valid = true)
		{
		status = valid is true ? 'normal' : 'error'
		if .status isnt status
			{
			.status = status
			.Act('SetValid', valid)
			}
		}

	SetWarning(warn = true)
		{
		status = warn is true ? 'warn' : 'normal'
		if .status isnt status
			{
			.status = status
			.Act('SetWarning', warn)
			}
		}

	GetValid()
		{
		return .status is 'normal'
		}

	Get()
		{
		return .untranslated
		}

	GetReadOnly() // read-only not applicable to status
		{
		return true
		}
	}