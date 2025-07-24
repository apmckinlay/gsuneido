// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
StatusbarControl
	{
	Name:		"Status"
	Xstretch:	1
	New()
		{
		super()
		.Map = Object()
		.Map[NM.CLICK] = "NM_CLICK"
		}

	NM_CLICK()
		{
		if false isnt matches = (.error $ ' ').Match('\s[1-9][0-9]*\s')
			.gotoSyntaxError(matches)
		return 0
		}

	gotoSyntaxError(matches)
		{
		line = .error[matches[0][0] :: matches[0][1]].Trim()
		if line.Number?()
			.Send('SetFirstVisibleLine', Max(Number(line) - 10 /*= shift lines*/, 0))
		}

	ClearError()
		{ .set(.left, '', .right, false) }

	left: 	''
	error:	''	// the center is ONLY for error messages
	right: 	''
	color: 	false
	// To clear, pass ' ' in their places IE: ' \t \t ' will clear all position
	Status(status, invalid = false, valid = false)
		{
		parts = status.Split('\t').Set_default('')

		if false is color = .getColor(invalid, valid)
			.error = ''

		.set(.getLeft(parts[0]),
			parts[1].Blank?() ? .error : parts[1],
			parts[2].Blank?() ? .right  : parts[2],
			color)
		}

	set(.left, .error, .right, .color)
		{
		.Set(left $ '\t' $ error $ '\t' $ right)
		.SetBkColor(.color)
		}

	getColor(invalid, valid)
		{
		if invalid or .color is CLR.ErrorColor and not .error.Blank?()
			return CLR.ErrorColor
		return valid ? CLR.GREEN : false
		}

	getLeft(left)
		{
		if left is '' and false isnt rec = .rec()
			{
			left = .dateStr(rec.lib_modified, rec.lib_committed)
			if .error.Blank?()
				.color = false
			}
		return left
		}

	rec()
		{
		return Query1(SvcTable(.Send(#CurrentTable)).NameQuery(.Send(#CurrentName)))
		}

	dateStr(lib_modified, lib_committed)
		{
		str = lib_modified isnt ''
			? '   Modified: ' $ lib_modified.ShortDateTime()
			: ''
		str $= lib_committed isnt ''
			? '   Committed: ' $ lib_committed.ShortDateTime()
			: ''
		return str
		}
	}
