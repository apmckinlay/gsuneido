// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	render(doc)
		{
		return AstFmtRender().Render(doc)
		}

	cursor()
		{
		return [i: 0, done: 0, blank: false]
		}

	Test_commentIn()
		{
		cm = AstFmtComments("x\n// one\n// two\ny = 1")
		Assert(cm.CommentIn?(0, 2), is: false)
		Assert(cm.CommentIn?(0, 3))
		Assert(cm.CommentIn?(15, 99), is: false)
		}

	Test_leading()
		{
		cm = AstFmtComments("x\n// one\n// two\ny = 1")
		cur = .cursor()
		Assert(.render(cm.Leading(cur, 16)), is: "// one\n// two\n")
		Assert(cur.done, is: 16)
		}

	Test_trailing()
		{
		cm = AstFmtComments("f() // note\ng()")
		cur = .cursor()
		cm.SkipTo(cur, 3)
		Assert(cur.done, is: 3)
		Assert(.render(cm.Trailing(cur, 3)), is: " // note")
		}

	Test_blankDetection()
		{
		cm = AstFmtComments("x\n\n\ny")
		cur = .cursor()
		cm.SkipTo(cur, 1)
		cm.Trailing(cur, 1, allowBlank:)
		Assert(cur.blank)
		cm2 = AstFmtComments("x\ny")
		cur2 = .cursor()
		cm2.SkipTo(cur2, 1)
		cm2.Trailing(cur2, 1, allowBlank:)
		Assert(cur2.blank, is: false)
		}

	Test_unusedSuppressed()
		{
		src = "x /*unused*/ y"
		cur = .cursor()
		Assert(.render(AstFmtComments(src).Leading(cur, 13, unusedParam:)), is: "")
		cur = .cursor()
		Assert(.render(AstFmtComments(src).Leading(cur, 13)), is: "/*unused*/ ")
		}
	}
