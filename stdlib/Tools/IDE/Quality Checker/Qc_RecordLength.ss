// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	warningThreshold: 700 //If > this many lines rating goes down
	notificationThreshold: 500 //Alert the user the line count is getting high

	CallClass(recordData, minimizeOutput? = false)
		{
		nlines = recordData.code.Lines().Filter(.countLine?).Count()

		warnings = Object()
		desc = .getDescription(nlines, minimizeOutput?)
		rating = .getRating(nlines)
		size = nlines > .warningThreshold ? 1 : 0
		return Object(:desc, :rating, :warnings, :size)
		}

	countLine?(line)
		{
		line = line.Trim()
		return not (line.Prefix?('//') or line in ('', '{', '}', '[', ']', '(', ')'))
		}

	getRating(recordLength)
		{
		secondaryRatingPenalty = 200
		rating = 5
		if recordLength > .warningThreshold + secondaryRatingPenalty
			rating = 0
		else if recordLength > .warningThreshold
			rating = 3
		return rating
		}

	getDescription(recordLength, minimizeOutput?)
		{
		if minimizeOutput? and recordLength < .notificationThreshold
			return ''

		warnText = recordLength < .warningThreshold and
			recordLength > .notificationThreshold and minimizeOutput?
			? 'Remain under '
			: 'Limit to '

		return  'Record is ' $ recordLength $ ' lines long. ' $ warnText $
			.warningThreshold $ ' lines'
		}
	}
