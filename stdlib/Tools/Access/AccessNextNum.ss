// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Addon
	{
	New(parent, .nextNum)
		{
		super(parent, false)
		if .nextNum isnt false
			{
			stringVal = .nextNum.GetDefault('stringVal', false)
			.nextNumClass = stringVal is true ? GetNextNumString : GetNextNum
			}
		}

	nextNum_num: false
	SetData(rec, newrec)
		{
		if .nextNum is false
			return

		// put back any previous next num (if not saved)
		if .nextNum_num isnt false
			.PutBack()

		if newrec is true
			{
			.nextnum(rec)
			.schedule_nextnum_renew(rec[.nextNum.field])
			}
		}

	nextNumAttempts: 100
	nextnum(rec, skipLog = false)
		{
		// keep trying until we get a nextnum that hasn't been used
		i = 0
		do
			{
			nextnum = .nextNumClass.Reserve(.nextNum.table, .nextNum.table_field)
			if not .isDuplicate(nextnum)
				break
			.nextNumClass.Confirm(nextnum, .nextNum.table, .nextNum.table_field)
			} while ++i < .nextNumAttempts
		if i > 0 and not skipLog
			.logNextNumWarning('AccessControl.nextnum - had to skip ' $ i $
				' numbers to get: ' $ nextnum)
		.nextNum_num = rec[.nextNum.field] = nextnum
		}

	logNextNumWarning(msg)
		{
		for c in Contributions('LogNextNumWarning')
			c(msg)
		}

	isDuplicate(nextnum)
		{
		return IsDuplicate(.GetBaseQuery(), .nextNum.field, nextnum)
		}

	Confirm() // called from successful Save
		{
		if .nextNum_num is false
			return
		if .nextNum_num isnt .GetData()[.nextNum.field]
			.nextNumClass.PutBack(.nextNum_num, .nextNum.table, .nextNum.table_field)
		else
			.nextNumClass.Confirm(.nextNum_num, .nextNum.table, .nextNum.table_field)
		.nextNum_num = false
		.kill_nextnum_timer()
		}

	PutBack()
		{
		if .nextNum_num is false
			return
		.nextNumClass.PutBack(.nextNum_num, .nextNum.table, .nextNum.table_field)
		.nextNum_num = false
		.kill_nextnum_timer()
		}

	nextnum_timer: false
	schedule_nextnum_renew(nextnum)
		{
		seconds = (.nextNumClass.ReserveSeconds / 3).Int() /*= 100s, to check 3 times*/
		.nextnum_timer = .Delay(seconds.SecondsInMs(), .Renew)
		if .nextNum_num isnt false and .nextNum_num isnt nextnum
			.PutBack()
		}

	renewAttempts: 10
	Renew()
		{
		.kill_nextnum_timer()
		if .nextNum_num is false or not .EditMode?() or not .NewRecord?()
			return true
		origNum = .nextNum_num
		count = .attemptRenew()
		// the SuneidoLog entered in the "if condition" is to deal with the scenario where
		// the system could only end up with a duplicate next number even after max tries
		// searching for an unused one;
		// see suggestion 30325
		if count is .renewAttempts and .isDuplicate(.nextNum_num)
			{
			.logDuplicateNextNum()
			.nextNum_num = false
			return 'Please fix your next ' $ Prompt(.nextNum.field) $ ' before saving'
			}
		.schedule_nextnum_renew(.GetRecordControl().FindControl(.nextNum.field).Get())
		.renewNextNumChanged(origNum)
		return true
		}

	attemptRenew()
		{
		count = 0
		while ++count < .renewAttempts and false is .renew()
			.nextnum(.GetData(), skipLog:)
		if count > 1
			.logNextNumWarning('renew_next_num while looped ' $ count $ ' times')
		return count
		}

	renew()
		{
		return .nextNumClass.
			Renew(.nextNum_num, .nextNum.table, .nextNum.table_field, skipLog:)
		}

	renewNextNumChanged(origNum)
		{
		// .nextNum_num can be set to false IF the user manually changes the field
		// controlling the reserved number. User will be informed via lower level
		// duplicate handling if required instead.
		if .nextNum_num is false or origNum is .nextNum_num
			return
		.AlertWarn('Next Number', 'Another user has taken ' $ origNum $
			'. You have been assigned ' $ .nextNum_num)
		SuneidoLog('WARNING: Next Number started at: ' $ origNum $ ', skipped to: ' $
			.nextNum_num)
		}

	logDuplicateNextNum()
		{
		SuneidoLog('INFO: Next Number ' $ Display(.nextNum_num) $
			' is still a duplicate after maximum tries (' $ .nextNum.table $ ')', calls:)
		}

	kill_nextnum_timer()
		{
		if .nextnum_timer is false
			return
		.nextnum_timer.Kill()
		.nextnum_timer = false
		}
	}
