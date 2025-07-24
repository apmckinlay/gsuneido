// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Process(errMsg, errorOb, msg, classId, threshold = 10)
		{
		logType = "ERROR: (CAUGHT) "
		caughtMsg = "generic connection error handling; may need attention"
		if .errorMatched?(errMsg, errorOb)
			{
			numberOfConnectionErrs = .CalcNumberOfConnectionErrs(classId, errorOb, errMsg)
			if numberOfConnectionErrs < threshold
				{
				logType = "INFO: "
				caughtMsg = ''
				}
			}
		err = errMsg.Has?('<Message>')
			? errMsg.AfterFirst('<Message>').BeforeFirst('</Message>')
			: errMsg
		.AddToLog(logType $ classId $ ': ' $ msg $ ' - ' $ err, :caughtMsg)
		}

	errorMatched?(errMsg, errorOb)
		{
		for error in errorOb
			if errMsg.Has?(error)
				return true

		return false
		}

	CalcNumberOfConnectionErrs(classId, errorOb, errMsg)
		{
		if .errorMatched?(errMsg, errorOb)
			return QueryCount('suneidolog where
				sulog_timestamp > ' $ Display(Date().Plus(hours: -24)) $
				' and sulog_message =~ ' $ Display(classId))
		else
			return 0
		}

	AddToLog(msg, caughtMsg = '')
		{
		SuneidoLog(msg, :caughtMsg)
		}
	}