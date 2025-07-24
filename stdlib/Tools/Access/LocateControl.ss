// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
/*
keys is an object where the members are prompts and the values are fields
*/
Controller
	{
	Name: 'Locate'

	New(.query, .keys, columns = false, excludeSelect = #(), sortFields = #(),
		option = false, optionalRestrictions = #(), .customizeQueryCols = false,
		.startLast = false)
		{
		super(controls = .layout(:columns, :excludeSelect, :sortFields,
			:option, :optionalRestrictions))
		if (controls[0] is 'Static')
			{
			.Send('Locate?', false)
			return
			}
		.Send('Locate?', true)
		.locate = .Horz.Locate
		.locateby = .Horz.LocateBy
		.NewValue(.locateby.Get(), .locateby)

		.prevproc = SetWindowProc(.locate.Field.Hwnd, GWL.WNDPROC, .Enterproc)
		}
	prevproc: false
	Enterproc(hwnd, msg, wparam, lparam)
		{
		_hwnd = .WindowHwnd()
		if msg is WM.GETDLGCODE
			.locate.Field.HideBalloonTip()
		if (msg is WM.GETDLGCODE and
			false isnt (m = MSG(lparam)) and
			m.wParam is VK.RETURN and
			(m.message is WM.CHAR or m.message is WM.KEYDOWN))
			return DLGC.WANTALLKEYS
		if (msg is WM.CHAR and wparam is VK.RETURN)
			{
			result = CallWindowProc(.prevproc, hwnd, msg, wparam, lparam)
			.Send("On_Go")
			return result
			}
		return CallWindowProc(.prevproc, hwnd, msg, wparam, lparam)
		}

	layout(columns, excludeSelect, sortFields, option, optionalRestrictions)
		{
		if .keys in ( #(), #(""), #("": "") )
			return #(Static '')

		availableColumns = QuerySelectColumns(.query)
		.useLowerKeys(availableColumns)

		// have to do substr so key on locateby doesn't get too long
		// doesn't use TruncateKey since this does not handle Tr (just replace)
		if option is false
			option = .query.Tr('\r\n', ' ')[.. 50] /*= max size*/
		.save_settings_key = "locateby:" $ option

		firstKeyField = .keys.Values().Sort!()[0]
		keyField = UserSettings.Get(.save_settings_key)
		if not .keys.Has?(keyField)
			keyField = firstKeyField
		keyPrompt = .keys.Find(keyField)

		columns = columns is false
			? availableColumns
			: columns.Copy().Append(Customizable.GetPermissableFields(.query))

		abbrevField = .getAbbrevField(keyField, columns)

		// make sure key is first column
		columns.Remove(keyField).Add(keyField, at: 0)
		return Object('Horz',
			#(Skip 6),
			#(Static 'Locate'),
			#(Skip 4),
			Object('LocateKey', .query, keyField, keys: .keys.Copy().Merge(sortFields),
				width: 10, name: 'Locate',	:columns, :excludeSelect,
				saveInfoName: 'Locate: ' $ option $ ' ' $ firstKeyField, :abbrevField,
				:optionalRestrictions, customizeQueryCols: .customizeQueryCols,
				status: 'use Ctrl+L to jump to this field', startLast: .startLast),
			#(Skip 4),
			#(Static 'by'),
			#(Skip 4),
			(.keys.Size() is 1
				? Object('Static', keyPrompt, name: 'LocateBy')
				: Object('ChooseList', .keys.Members(), set: keyPrompt, name: 'LocateBy'))
			#(Skip 6),
			#(Button Go)
			)
		}
	useLowerKeys(qcols) // pass in qcols to simplify testing
		{
		for prompt in .keys.Members().Copy()
			{
			field = .keys[prompt]
			if qcols.Has?(field $ '_lower!') // should be indexed but no easy check
				{
				.keys.Delete(prompt)
				.keys[prompt $ '*'] = field $ '_lower!'
				}
			}
		}

	getAbbrevField(keyField, columns)
		{
		origKeyField = keyField.RemoveSuffix('_lower!')
		abbrevField = origKeyField.Replace("(_num|_name|_abbrev)$", "_abbrev")
		if abbrevField in (keyField, origKeyField) or not columns.Has?(abbrevField)
			abbrevField = false
		return abbrevField
		}

	On_Go()
		{
		.Send('On_Go')
		}
	current_key: false
	NewValue(value, source)
		{
		if source is .locateby and .keys.Member?(value) and .current_key isnt value
			{
			.current_key = value
			newkey = .keys[.current_key].Split(',')[0]
			.locate.ChangeKey(newkey)
			.locate.FieldReturn()
			if .locate.Get() isnt ''
				.Send('On_Go')
			}
		}
	ListRecordSelected(record, locateby, source)
		{
		if (source is .locate)
			{
			if .keys.Member?(locateby) // only set when locateby is a key
				{
				.locateby.Set(locateby)
				.NewValue(locateby, .locateby)
				}
			keyvalue = ''
			keyfields = .keys[.current_key].Split(',')
			for (k in keyfields)
				keyvalue $= ',' $ record[KeyDisplayField(k)]
			keyvalue = keyvalue[1 ..]
			.locate.Set(keyvalue)

			.Send('On_Go')
			}
		}
	Get()
		{
		return Object(locate: .locate.Get() locateby: .locateby.Get())
		}
	EditHwnd()
		{
		return .locate.EditHwnd()
		}
	locate: false
	Valid?()
		{
		return .locate is false ? true : .locate.Valid?()
		}
	Pos()
		{
		if .locate is false
			return false
		GetClientRect(.EditHwnd(), r = [])
		ClientToScreen(.EditHwnd(), pos = [x: r.left, y: r.top])
		return pos
		}
	BalloonTip(msg)
		{
		.locate.Field.ShowBalloonTip(msg)
		}
	SelectAll()
		{
		if .locate isnt false
			.locate.SelectAll()
		}

	locateby: false
	Destroy()
		{
		if (.locateby isnt false and .keys.Member?(by = .locateby.Get()))
			UserSettings.Put(.save_settings_key, .keys[by])
		if .prevproc isnt false
			{
			SetWindowProc(.locate.Field.Hwnd, GWL.WNDPROC, .prevproc)
			ClearCallback(.Enterproc)
			}
		super.Destroy()
		}
	}
