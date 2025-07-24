// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
//	- Comments between '/*' and '* /' must include * at the start of each line to
//	alert this method it is still in a comment in order to avoid line length warnings

class
	{
	CallClass(recordData, minimizeOutput? = false)
		{
		lineWarnings = Object()
		warnings = .findLongLines(recordData, minimizeOutput?, lineWarnings)
		rating = .lineLengthRating(warnings)
		desc = .createLineDescription (warnings, minimizeOutput?)
		return Object (:rating, :warnings, :desc, :lineWarnings)
		}

	findLongLines(recordData, minimizeOutput?, lineWarnings)
		{
		warnings = Object()
		line_number = 1
		for line in recordData.code.Lines()
			{
			if line.Detab().RightTrim().Size() > CheckCode.MaxLineLength and
				not line.Trim().Prefix?('catch') and not line.Trim().Prefix?("/") and
				not line.Trim().Prefix?("*") // (/*, *, //) comments supported
					{
					if minimizeOutput?
						{
						lineWarnings.Add(Object(line_number))
						warnings.AddUnique(Record(name: recordData.lib $ ':' $
							recordData.recordName $ ':' $ line_number $ " - Long line"))
						}
					else
						warnings.AddUnique(Record(name: "Line: " $ line_number))
					}
			++line_number
			}
		return warnings
		}

	createLineDescription(warnings, minimizeOutput?)
		{
		if not warnings.Empty?()
			return 'Some lines are too long'

		if not minimizeOutput?
			return 'All lines are of adequate length'
		return ''
		}

	lineLengthRating(warnings)
		{
		return Max(0, Qc_CalculateCodeRating.MaxRating - warnings.Size())
		}
	}






