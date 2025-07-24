// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function (page = 'StartPage')
	{
	WikiEnsure()
	if '' isnt err = WikiTitleValid?(page)
		return err
	lastViewed = Date().Format('yyyyMM')
	if false is x = Query1('wiki', name: page)
		{
		x = Record(text: "This page does not exist, use the Edit button to create it.")
		lastViewed = ''
		}

	vars = Object().Set_default('')
	vars.page = page
	vars.title = WikiFormatTitle(page)
	vars.body = Xml('div', WikiFormat(x.text), id: 'body')
	vars.action =
		'<div class="noPrint">
			[<a href="Wiki?append=' $ page $ '">Add a comment</a>]
			[<a href="Wiki?edit=' $ page $ '">Edit this page</a>]</div>'
	vars.smallaction =
		'<span class="noPrint" style="font-size: 80%;">
			[<a href="Wiki?edit=' $ page $ '">Edit</a>]&nbsp;
			[<a href="Wiki?append=' $ page $ '">Comment</a>]&nbsp;
			</span>'
	vars.orphaned = x.orphaned is true
		? '<span style="margin-left: .5em;">' $ WikiOrphans.OrphanedMessage $ '</span>'
		: ''

	if Date?(x.created)
		vars.created = "Created " $ x.created.Format("yyyy-MM-dd")
	if Date?(x.edited)
		vars.lastedit = "Last edited " $ x.edited.Format("yyyy-MM-dd")

	// update last viewed
	if x.lastviewed isnt lastViewed
		QueryApply('wiki where name is ' $ Display(page), update:)
			{
			it.lastviewed = lastViewed
			it.Update()
			}

	return WikiTemplate().Replace('\$[a-z]+', { |var| vars[var[1 ..]] } )
	}
