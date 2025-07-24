// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// Add to a rotating series of Nlogs files, one file per day,
class
	{
	// hoursPastMidnight - this is to handle processes that run over midnight
	// ensures that the logs still go in the same day.
	Nlogs: 10
	DaysPerLog: 1
	DefaultFolder: ""
	FileSizeLimit: 300_000_000 // ???
	CallClass(name, s, multi_line = false, hoursPastMidnight = 0)
		{
		if Sys.Client?()
			{
			ServerEval(Name(this), name, s, multi_line, hoursPastMidnight)
			return
			}
		name = Paths.ToStd(.DefaultFolder $ name)
		try .Synchronized()
			{
			file = .rotate(name, hoursPastMidnight)
			}
		catch (err)
			{
			SuneidoLog('ERRATIC: (CAUGHT) Rlog rotate failed',
				params: Object(:err, :name, :s), caughtMsg: 'rlog aborted')
			return
			}
		.log(file, s, multi_line)
		}

	rlogStartDate: #20070101
	rotate(name, hoursPastMidnight)
		{
		rotationNum = .getRotationNum(.getCurrDate(), .DaysPerLog, hoursPastMidnight)
		file = .CurrentLog(name, rotationNum)
		.deleteOldFile(name, file, rotationNum)
		return file
		}

	getRotationNum(date, daysPerLog, hoursPastMidnight = 0)
		{
		date = date.Plus(hours: -hoursPastMidnight)
		return (date.MinusDays(.rlogStartDate) / daysPerLog).Floor()
		}

	getCurrDate()
		{
		return Date().Plus(hours: OptContribution('Rlog_HoursOffset', 0))
		}

	GetLastLogFileName(name)
		{
		if Sys.Client?()
			return ServerEval('Rlog.GetLastLogFileName', name)
		return .getLastLog(name)
		}
	getLastLog(name)
		{
		logs = .dirLogFiles(name)
		if logs.Size() is 0
			return false
		logs.Sort!(By(#date))
		logFile = logs.Last().name
		return logFile
		}
	dirLogFiles(name) // extracted for test
		{
		return Dir(name $ '*.log', details:)
		}
	// consider GetLastLogFileName to avoid midnight date change issue
	CurrentLog(name, rotationNum = false, offset = 0) // 0 for current, -1 for previous
		{
		if rotationNum is false
			rotationNum = .getRotationNum(.getCurrDate(), .DaysPerLog)
		n = (rotationNum + offset) % .Nlogs
		return name $ n $ ".log"
		}

	deleteOldFile(name, file, rotationNum)
		{
		cache = Suneido.GetInit(#Rlog, { Object().Set_default(false) })
		if cache[name] isnt rotationNum
			{
			cache[name] = rotationNum
			.deleteFile(file, rotationNum)
			}
		}

	deleteFile(file, rotationNum)
		{
		if FileExists?(file)
			{
			File(file, "r")
				{|f|
				line = f.Readline()
				}
			if line isnt false and .createNewFile?(line, rotationNum)
				DeleteFile(file)
			}
		}

	createNewFile?(line, rotationNum)
		{
		firstFileDate = Date(line.BeforeFirst('.'))
		firstPossibleDate = .firstPossibleRotationDate(rotationNum)
		return firstFileDate < firstPossibleDate
		}

	firstPossibleRotationDate(rotationNum)
		{
		return .rlogStartDate.Plus(days: rotationNum * .DaysPerLog)
		}

	log(file, s, multi_line)
		{
		try
			AddFile(file, .formatMessage(Date(), s, multi_line), limit: .FileSizeLimit)
		catch (unused, "*cannot find|*no such")
			{
			dir = ServerPath().Dir() $ '/' $ file.BeforeLast('/') $ "/"
			CreateDirectories(dir)
			AddFile(file, .formatMessage(Date(), s, multi_line), limit: .FileSizeLimit)
			}
		}

	formatMessage(date, s, multi_line = false)
		{
		s = s.Replace('%m', .MemoryArena)
		s = s.Replace('%s', .sessionId)
		if multi_line
			{
			sepRepeats = 28
			s = Display(date)[1 ..] $ ' ' $ '-+'.Repeat(sepRepeats) $ '\r\n' $ s
			}
		else
			s = Display(date)[1 ..] $ '\t' $ s.Tr('\r\n', ' ')
		return s $ '\r\n'
		}

	MemoryArena(unused = false)
		{
		return (MemoryArena() / 1.Mb()).Round(0) $ 'mb'
		}

	sessionId(unused)
		{
		return Database.SessionId()
		}
	}
