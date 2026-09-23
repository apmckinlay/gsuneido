// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
#((menuoption: "Edit Comment",
		menufunc: function(sel, source)
			{
			svc = SvcClient()
			if "" is msg = svc.AllowEditCommit(sel.table, sel.name, sel.When)
				{
				comment = ZoomControl(source.WindowHwnd(), sel.Comment, false).Trim()
				if comment isnt sel.Comment
					if "" is
						msg = svc.UpdateComment(sel.table, sel.name, sel.When, comment)
						sel.Comment = comment
				}
			if msg isnt ""
				source.AlertError(source.Title, msg)
			return Sys.SuneidoJs?()
			}))
