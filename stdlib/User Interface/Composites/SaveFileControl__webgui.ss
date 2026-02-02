// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
Controller
	{
	New(.title = 'Save', width = 30, .filter = "", .ext = "", file = "",
		.flags = false, status = '', .mandatory = false)
		{
		super(.layout(:width, :status, :mandatory))
		.field = .FindControl('Field')
		.Set(file)
		.Send('Data')
		}

	layout(width, status, mandatory)
		{
		return Object('Field', :width, :status, :mandatory,
			cue: 'Please enter file name only')
		}

	SetFilter(filter)
		{ .filter = filter	}
	SetDefExt(ext)
		{ .ext = ext }

	Valid?()
		{
		if .GetReadOnly()
			return true

		fileName = .Get()
		if not fileName.Blank?() and not CheckDirectory.ValidFileName?(fileName)
			return false

		if not .FileNameOnly?(fileName)
			return false

		return super.Valid?()
		}

	FileNameOnly?(fileName)
		{
		return Paths.Basename(fileName) is fileName
		}

	Set(value)
		{
		if value isnt ""
			.InitialDir = value.BeforeLast(Paths.Basename(value))
		// if is gived only the path, displayed value is ""
		value = Paths.Basename(value) is "" ? "" : value
		.field.Set(value)
		}

	NewValue(value)
		{
		.Send("NewValue", value)
		}

	Get()
		{
		return .field.Get()
		}

	Dirty?(dirty = "")
		{
		return .field.Dirty?(dirty)
		}

	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}