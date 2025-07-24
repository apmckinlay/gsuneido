// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// this is tests for built-in methods
// Records_Test is for methods defined in Records
Test
	{
	Test_preset()
		{
		observed = false
		x = Record()
		x.Observer({ |member| observed = member })
		x.PreSet('age', 23)
		Assert(x hasMember: 'age')
		Assert(x.age is: 23)
		Assert(observed is: false)
		x.age = 99
		Assert(observed is: 'age')
		}
	Test_copydeps()
		{
		x = Record()
		x.name = "Fred"
		x.test_amount = 100
		Assert(x.test_simple_pull is: "Fred and 100")
		y = x.Copy()
		y.test_amount = 200
		Assert(y.test_simple_pull is: "Fred and 200")
		}
	Test_removeObservers()
		{
		// tests removing existing and non existing observers
		x = Record()
		f = function () { throw "can't happen" }
		x.RemoveObserver(f)
		x.Observer(f)
		x.RemoveObserver(f)
		x.a = "test"
		x.RemoveObserver(f)
		}
	Test_Delete_observer()
		{
		r = Record()
		flag = false
		r.Observer({|member| flag = member })

		r.x = 123
		Assert(flag is: #x)

		flag = false
		r.Delete(#x)
		Assert(flag is: #x)

		r.y = 123
		flag = false
		r.Erase(#y)
		Assert(flag is: #y)
		}
	Test_Delete_invalidate()
		{
		r = Record(d: 123)
		r.AttachRule(#rule, function ()
			{
			.d // create dependency
			Timestamp()
			})
		t = r.rule
		Assert(r.GetDeps(#rule) is: 'd')
		Assert(r.rule is: t) // rule not triggered again
		r.Delete(#d) // should invalidate rule
		Assert(r.rule greaterThan: t) // rule triggered
		}
	Test_assign_equal_bug()
		{
		r = Record()
		log = ""
		r.Observer({|member| log $= member })
		x = Object()
		y = Object()
		r.a = x
		Assert(log is: "a")
		r.a = y // previously ignored because x == y
		x[0] = 'x'
		y[0] = 'y'
		Assert(r.a is: #(y)) // previously gave #(x)
		}
	Test_GetDefault()
		{
		r = Record()
		r.AttachRule(#foo, function () { 123 })
		Assert(r.GetDefault(#foo, false) is: 123)
		Assert(r.Member?(#foo), msg: 'member? foo')
		r = Record(foo: 123)
		r.AttachRule(#foo, function () { 456 })
		Assert(r.GetDefault(#foo, false) is: 123)
		r.Invalidate(#foo)
		Assert(r.GetDefault(#foo, false) is: 456)
		}
	Test_GetDefault_Deps()
		{
		r = Record()
		r.AttachRule(#foo, function() { .GetDefault(#bar, 999) })
		Assert(r.foo is: 999)
		Assert(r.GetDeps(#foo) is: "bar")
		r = Record(bar: 123)
		r.AttachRule(#foo, function() { .GetDefault(#bar, 999) })
		Assert(r.foo is: 123)
		Assert(r.GetDeps(#foo) is: "bar")
		}
	Test_MemberQ_Deps()
		{
		// apm - I'm not sure if this is "correct" but it's the current behavior
		r = Record()
		r.AttachRule(#foo, function() { .Member?(#bar) })
		Assert(r.foo is: false)
		Assert(r.GetDeps(#foo) is: "")
		r = Record(bar: 123)
		r.AttachRule(#foo, function() { .Member?(#bar) })
		Assert(r.foo)
		Assert(r.GetDeps(#foo) is: "")
		}
	Test_DeleteAll_vs_Clear()
		{
		// Clear() makes the record the same as Record()
		x = QueryFirst('columns sort column')
		Assert(x.New?() is: false, msg: 'clear false')
		x.Clear()
		Assert(x.New?(), msg: 'clear true')
		Assert(x isSize: 0)

		// Delete(all:) is the equivalent to Delete of each member
		x = QueryFirst('columns sort column')
		Assert(x.New?() is: false, msg: 'new? false')
		x.Delete(all:)
		Assert(x.New?() is: false, msg: 'delete all false')
		Assert(x isSize: 0)
		}
	Test_New?()
		{
		Assert(Record().New?(), msg: 'record new?')

		Assert(QueryFirst('tables sort table').New?() is: false, msg: 'tables new?')
		Assert(QueryFirst('stdlib sort name').New?() is: false, msg: 'stdlib new?')
		}
	Test_Invalidate()
		{
		Suneido.tmp = 0
		r = Record(foo: 123)
		r.AttachRule(#rul, function(){ ++Suneido.tmp; .foo })
		r.rul
		r.Invalidate(#foo)
		r.rul
		r.Invalidate(#foo)
		r.rul
		Assert(Suneido.tmp is: 3)
		}
	}
