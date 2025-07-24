// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_stop()
		{
		table = 'table'
		list = new .local_list(false)
		sgm = SvcGetMaster
			{
			SvcGetMaster_runningTestRunner()
				{
				return false
				}
			SvcGetMaster_needNew()
				{
				return true
				}
			}

		Assert(sgm.SvcGetMaster_stop(list, table) is: false)
		table = ''
		Assert(sgm.SvcGetMaster_stop(list, table))

		table = 'table'
		empty_list = new .local_list(true)
		Assert(sgm.SvcGetMaster_stop(empty_list, table))
		}

	Test_forceStop()
		{
		sgm = SvcGetMaster
			{
			SvcGetMaster_startNewCheck()
				{
				}
			SvcGetMaster_needNew()
				{
				return true
				}
			}
		forceStop = sgm.SvcGetMaster_forceStop
		Assert(forceStop())

		sgm = SvcGetMaster
			{
			SvcGetMaster_startNewCheck()
				{
				}
			SvcGetMaster_needNew()
				{
				return false
				}
			}
		forceStop = sgm.SvcGetMaster_forceStop
		Assert(forceStop() is: false)
		}

	local_list: class
		{
		New(.listEmpty)
			{
			}
		Get()
			{
			return .listEmpty
				? #()
				: #((svc_name: 'name', svc_lib: 'lib', svc_type: ' '))
			}
		}
	}
