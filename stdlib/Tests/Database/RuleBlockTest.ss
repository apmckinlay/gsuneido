// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		// need to wrap in function so blocks don't get "this" from test class
		function ()
			{
			test = function (rule)
				{
				rec = Record(sub_rec: [b: 'bbb'])
				rec.AttachRule('a', rule)
				Assert(rec.a is: 'bbb')
				Assert(rec.GetDeps('a') is: "sub_rec")
				}
			test(function () { .sub_rec.b })	// function
			test({ .sub_rec.b })				// "static" block
			fld = 'b'
			test({ .sub_rec[fld] })				// block with context
			ruleFactory = function (f) { return { .sub_rec[f] } }
			test(ruleFactory('b'))				// persisted block
			}()
		}
	}