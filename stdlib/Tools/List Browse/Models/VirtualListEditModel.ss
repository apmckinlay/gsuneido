// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	LockKeyField: false
	New(lockFields = #(), .ProtectField = false, .ValidField = false,
		.warningField = false, keys = false)
		{
		.changedRecords = Object().Set_default(Object())
		if not .Editable?()
			return
		if lockFields is #()
			.LockKeyField = ShortestKey(keys)
		else
			{
			Assert(lockFields isSize: 1)
			.LockKeyField = lockFields[0]
			}
		Assert(not .LockKeyField.Has?(',') and .LockKeyField isnt ''
			msg: 'VirtualList - cannot detect key field for locking record')
		}

	lockHandler: false
	getInitLockHandler()
		{
		if .lockHandler is false
			.lockHandler = LockManagerHandler()
		return .lockHandler
		}

	Editable?()
		{
		return .ProtectField isnt false
		}

	AddChanges(rec, col, val, invalidCols = false)
		{
		if false is tmpChanges = .changedRecords.FindOne({ rec is it.rec })
			.changedRecords.Add(tmpChanges = Object(:rec, changes: []))
		changes = tmpChanges.changes
		changes[col] = val
		if invalidCols isnt false
			changes.invalidCols = invalidCols
		}

	AddInvalidCol(rec, col)
		{
		.initInvalidCols(rec).AddUnique(col)
		}

	initInvalidCols(rec)
		{
		if false is c = .changedRecords.FindOne({ rec is it.rec })
			.changedRecords.Add(c = Object(:rec, changes: []))
		changes = c.changes
		if not changes.Member?('invalidCols')
			changes.invalidCols = Object()
		return changes.invalidCols
		}

	RemoveInvalidCol(rec, col)
		{
		if false is i = .changedRecords.FindIf({ rec is it.rec })
			return
		c = .changedRecords[i]
		changes = .changedRecords[i].changes
		if not changes.Member?('invalidCols')
			return
		changes.invalidCols.Remove(col)
		if not .recordChanged?(c)
			.changedRecords.Delete(i)
		}

	ClearMemberChange(rec, member)
		{
		if false is i = .changedRecords.FindIf({ rec is it.rec })
			return
		tmpChanges = .changedRecords[i].changes
		if tmpChanges.Member?(member)
			{
			tmpChanges.Delete(member)
			if tmpChanges.Member?('invalidCols')
				tmpChanges.invalidCols.Remove(member)
			}
		if not .recordChanged?(.changedRecords[i])
			.changedRecords.Delete(i)
		}

	ClearAllChanges()
		{
		.changedRecords = Object().Set_default(Object())
		}

	ClearChanges(rec)
		{
		.changedRecords.RemoveIf({ it.rec is rec })
		}

	HasChanges?()
		{
		return .changedRecords.Any?(.recordChanged?)
		}

	AllowLeaving?()
		{
		return .changedRecords.Every?({
			it.changes.Members().Remove('invalidCols').Empty?() })
		}

	RecordChanged?(rec)
		{
		if false is c = .changedRecords.FindOne({ it.rec is rec })
			return false
		return .recordChanged?(c)
		}

	recordChanged?(c)
		{
		if c.rec.New?()
			return true
		return not c.changes.Members().Remove('invalidCols').Empty?()
		}

	RecordDirty?(rec) // public static method
		{
		return rec.vl_list.GetModel().EditModel.RecordChanged?(rec)
		}

	ColumnInvalid?(rec, col)
		{
		if .validRule(rec) isnt ''
			return true
		return .getInvalidCols(rec).Has?(col)
		}

	GetWarningMsg(rec)
		{
		if .warningField is false
			return ''
		return rec[.warningField]
		}

	GetInvalidMsg(rec)
		{
		cols = .getInvalidCols(rec)
		return not cols.Empty?()
			? "Invalid: " $ cols.Map(PromptOrHeading).Join(', ')
			: .validRule(rec)
		}

	GetInvalidInfo(rec)
		{
		validRule = .validRule(rec)
		return Object(:validRule, invalidCols: .getInvalidCols(rec))
		}

	validRule(rec)
		{
		if .RecordChanged?(rec) and .ValidField isnt false and rec.Member?(.ValidField)
			return rec[.ValidField]
		return ''
		}

	getInvalidCols(rec)
		{
		if false is c = .changedRecords.FindOne({ it.rec is rec })
			return #()
		return c.changes.GetDefault('invalidCols', #())
		}

	HasInvalidCols?(rec)
		{
		return not .getInvalidCols(rec).Empty?()
		}

	// use all? option when the value in record is not changed but invalid,
	// like invalid string on number field does not change value from "" but it's invalid
	GetOutstandingChanges(all? = false)
		{
		changes = all? ? .changedRecords : .changedRecords.Filter(.recordChanged?)
		return changes.Map({ it.rec })
		}

	LockRecord(rec)
		{
		if '' is key = .lockKey(rec)
			return true
		handler = .getInitLockHandler()
		if handler.OwnLock?(key)
			return true
		if handler.LockedKeys().Any?({ it isnt key }) or
			false isnt .changedRecords.FindOne({ it.rec isnt rec })
			return 'Please correct the invalid line before editing another record.'
		return .getInitLockHandler().Lock(key)
		}

	LockedKeys()
		{
		return .getInitLockHandler().LockedKeys()
		}

	RecordLocked?(rec)
		{
		if false isnt .changedRecords.FindOne({ it.rec is rec })
			return true
		if '' is key = .lockKey(rec)
			return false
		return .getInitLockHandler().OwnLock?(key)
		}

	HasOtherLockedRecord?(currentRec)
		{
		if false isnt .changedRecords.FindOne({ it.rec isnt currentRec })
			return true
		handler = .getInitLockHandler()
		currentKey = .buildLockKey(currentRec)
		return handler.LockedKeys().Any?({ it isnt currentKey })
		}

	UnlockRecord(rec)
		{
		if '' is key = .lockKey(rec)
			return
		// REFACTOR: remove other places calling ClearChanges when UnlockRecord
		.ClearChanges(rec)
		lock = .getInitLockHandler()
		if lock.OwnLock?(key)
			lock.Unlock(key)
		}

	lockKey(rec)
		{
		return .buildLockKey(rec)
		}

	buildLockKey(rec)
		{
		if .LockKeyField is false or rec is false or not Object?(rec.vl_origin)
			return ''
		return String(rec.vl_origin[.LockKeyField])
		}

	Destroy()
		{
		if .lockHandler isnt false
			.lockHandler.Unlock()
		}
	}
