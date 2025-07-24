// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.SpyOn(LockManagerHandler.LockManagerHandler_schedule_lock_renew).Return(0)
		editModel = VirtualListEditModel(#(abc, efg))

		Assert(editModel.HasChanges?() is: false)
		Assert(editModel.ColumnInvalid?([], 'abc') is: false)

		rec = [key1: 'a', key2: 'b', abc: 'invalid']
		editModel.AddChanges(rec, 'abc', 'invalid', #(abc))
		Assert(editModel.HasChanges?())

		rec = [key1: 'a', key2: 'b', abc: 'invalid']
		editModel.ClearChanges(rec)
		Assert(editModel.ColumnInvalid?(rec, 'abc') is: false)
		Assert(editModel.ColumnInvalid?(rec, 'efg') is: false)

		editModel.UnlockRecord(rec)
		}

	Test_lock_unlock()
		{
		.SpyOn(LockManagerHandler.LockManagerHandler_schedule_lock_renew).Return(0)
		editModel = VirtualListEditModel(protectField: 'protect', keys: #(num))
		rec1 = [key1: 'a', num: 'abc', vl_origin: [num: 'abc']]
		rec2 = [key1: 'a', num: 'efg', vl_origin: [num: 'efg']]

		Assert(editModel.HasOtherLockedRecord?(rec1) is: false)
		Assert(editModel.HasOtherLockedRecord?(rec2) is: false)

		.AddTeardown({ .unlock(editModel, rec1) })
		editModel.LockRecord(rec1)
		Assert(editModel.LockedKeys() has: 'abc')

		.AddTeardown({ .unlock(editModel, rec2) })
		Assert(editModel.LockRecord(rec2) is:
			'Please correct the invalid line before editing another record.')
		Assert(editModel.LockedKeys() hasnt: 'efg')

		editModel.UnlockRecord(rec1)
		Assert(editModel.LockedKeys() hasnt: 'abc')

		rec2 = [key1: 'a', num: 'efg', vl_origin: [num: 'efg']]
		.AddTeardown({ .unlock(editModel, rec2) })
		Assert(editModel.LockRecord(rec2))
		Assert(editModel.LockedKeys() has: 'efg')

		Assert(editModel.HasOtherLockedRecord?(rec1))
		Assert(editModel.HasOtherLockedRecord?(rec2) is: false)

		editModel.UnlockRecord(rec2)
		Assert(editModel.LockedKeys() is: #())

		Assert(editModel.HasOtherLockedRecord?(rec1) is: false)
		Assert(editModel.HasOtherLockedRecord?(rec2) is: false)
		}

	unlock(editModel, rec)
		{
		editModel.UnlockRecord(rec)
		}

	Test_buildLockKey()
		{
		mock = Mock()
		mock.LockKeyField = 'name'
		build = VirtualListEditModel.VirtualListEditModel_buildLockKey
		Assert(mock.Eval(build, false) is: '')

		rec =  [name: '', material: 'Boxes', check: '', vl_origin: [name: '']]
		Assert(mock.Eval(build, rec) is: '')

		rec =  [name: 'invalid', material: 'Boxes', vl_origin: [name: 'invalid']]
		Assert(mock.Eval(build, rec) is: 'invalid')

		rec =  [name: 'Ted', material: 'Boxes', check: 'invalid',
			vl_origin: [name: 'Ted']]
		Assert(mock.Eval(build, rec) is: 'Ted')

		mock = Mock()
		mock.LockKeyField = false
		rec = [name: '', material: 'Boxes', check: '']
		mock.Eval(build, rec)

		mock.LockKeyField = 'material'
		rec =  [name: 'Ted', material: 'Boxes', check: 'invalid',
			vl_origin: [material: 'Boxes']]
		Assert(mock.Eval(build, rec) is: 'Boxes')
		}

	Test_Editable?()
		{
		.SpyOn(LockManagerHandler.LockManagerHandler_schedule_lock_renew).Return(0)
		m = VirtualListEditModel()
		Assert(m.Editable?() is: false)

		m = VirtualListEditModel(protectField: 'hello', keys: #(num))
		Assert(m.Editable?())
		}

	Test_AddInvalidCol()
		{
		.SpyOn(LockManagerHandler.LockManagerHandler_schedule_lock_renew).Return(0)
		m = VirtualListEditModel(protectField: 'hello', keys: #(num))
		rec = Record()
		Assert(m.RecordChanged?(rec) is: false)
		Assert(m.HasInvalidCols?(rec) is: false)
		Assert(m.ColumnInvalid?(rec, 'col') is: false)

		m.AddChanges(rec, 'col', 'val')
		Assert(m.RecordChanged?(rec))
		Assert(m.HasInvalidCols?(rec) is: false)
		Assert(m.ColumnInvalid?(rec, 'col') is: false)

		rec = Record()
		m = VirtualListEditModel(protectField: 'hello', keys: #(num))
		m.AddChanges(rec, 'col2', 'val2')
		m.AddInvalidCol(rec, 'col2')
		Assert(m.RecordChanged?(rec))
		Assert(m.HasInvalidCols?(rec))
		Assert(m.ColumnInvalid?(rec, 'col2'))

		m.AddInvalidCol(rec, 'col2')
		Assert(m.RecordChanged?(rec))
		Assert(m.HasInvalidCols?(rec))
		Assert(m.ColumnInvalid?(rec, 'col2'))

		m.RemoveInvalidCol(rec, 'col3')
		Assert(m.RecordChanged?(rec))
		Assert(m.HasInvalidCols?(rec))
		Assert(m.ColumnInvalid?(rec, 'col3') is: false)

		m.RemoveInvalidCol(rec, 'col2')
		Assert(m.RecordChanged?(rec))
		Assert(m.HasInvalidCols?(rec) is: false)
		Assert(m.ColumnInvalid?(rec, 'col2') is: false)
		}

	Test_RecordDirty?()
		{
		.SpyOn(LockManagerHandler.LockManagerHandler_schedule_lock_renew).Return(0)
		_m = VirtualListEditModel(protectField: 'hello', keys: #(num))
		rec = Record()
		rec.vl_list = class
			{
			GetModel()
				{
				return Object(EditModel: _m)
				}
			}
		Assert(VirtualListEditModel.RecordDirty?(rec) is: false)

		_m.AddChanges(rec, 'col', 'val')
		Assert(VirtualListEditModel.RecordDirty?(rec))
		}

	Test_AllLeaving?()
		{
		.SpyOn(LockManagerHandler.LockManagerHandler_schedule_lock_renew).Return(0)
		rec = Record()
		m = VirtualListEditModel()

		Assert(m.AllowLeaving?())

		m.AddInvalidCol(rec, 'col2')
		Assert(m.AllowLeaving?())

		m.AddChanges(rec, 'col2', 'val2')
		Assert(m.AllowLeaving?() is: false)

		m.RemoveInvalidCol(rec, 'col2')
		Assert(m.AllowLeaving?() is: false)
		}
	}
