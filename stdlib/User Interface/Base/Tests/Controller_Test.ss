// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	ctrlr: Controller
		{
		New2() { }
		Two() { 2 }
		}
	Test_redir()
		{
		c = new .ctrlr
		Assert(c.GetRedir(#a) is: false)
		c.Redir(#a, false)
		Assert(c.GetRedir(#a) is: false)
		ctrl = Mock()
		c.Redir(#a, ctrl)
		ctrl2 = Mock()
		c.Redir(#b, ctrl2)
		Assert(c.GetRedir(#a) is: ctrl)
		Assert(c.GetRedir(#b) is: ctrl2)
		c.Redir(#b, ctrl)
		Assert(c.GetRedir(#b) is: ctrl)
		c.DeleteRedir(#b)
		Assert(c.GetRedir(#b) is: false)
		c.RemoveRedir(ctrl)
		Assert(c.GetRedir(#a) is: false)
		}
	ctrl: Control
		{
		New2() { }
		}
	Test_send()
		{
		ctrlr = new .ctrlr
		ctrl = new .ctrl
		ctrl.Controller = ctrlr
		Assert(ctrl.Send(#x) is: 0)
		Assert(ctrl.Send(#Two) is: 2)
		ctrl2 = class { Two() { "two" } }
		ctrlr.Redir(#Two, ctrl2)
		Assert(ctrl.Send(#Two) is: "two")
		}
	}
