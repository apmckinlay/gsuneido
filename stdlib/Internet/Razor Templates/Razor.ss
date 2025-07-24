// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// see:
//	http://weblogs.asp.net/scottgu/archive/2010/07/02/introducing-razor.aspx
//	http://vibrantcode.com/blog/2010/7/3/introducing-razor-a-new-view-engine-for-aspnet.html
// 	http://vibrantcode.com/blog/2010/7/5/inside-razor-part-1-recursive-ping-pong.html
//	http://vibrantcode.com/blog/2010/7/12/inside-razor-part-2-expressions.html
//	http://www.ironshay.com/post/The-Razor-View-Engine-Basics.aspx
class
	{
	// TODO cache compiled templates
	CallClass(template, context = #())
		{
		src = .Translate(template)
		Dbg("SOURCE:\n" $ src.RightTrim() $ "\n-------")
		fn = ("function () {\ns = ''\n" $ src $ "\n }").Compile()
		return context.Eval(fn)
		}
	Translate(template)
		{
		template = template.Trim()
		if template is ""
			return ""
		dest = new .dest
		RazorHtml(dest, template)
		return dest.Output
		}
	dest: class
		{
		Output: ""
		Html(s)
			{
			Dbg('>>>html': Display(s))
			.Output $= 's $= ' $ Display(s) $ '\n'
			}
		Expr(s) // html encoded
			{
			Dbg('>>>expr': Display(s))
			.Output $= 's $= H(' $ s $ ')\n'
			}
		Code(s)
			{
			Dbg('>>>code': Display(s))
			if s isnt ""
				.Output $= s $ '\n'
			}
		CodeFragment(s)
			{
			Dbg('>>>frag': Display(s))
			.Output $= s
			}
		}
	}