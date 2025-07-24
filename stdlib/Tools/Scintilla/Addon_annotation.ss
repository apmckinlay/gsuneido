// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonIDE
	{
	Setting: ide_show_annotations
	Init()
		{ .SetVisible(.Set) }

	AddAnnotation(line, text)
		{
		.AnnotationSetStyle(line, 0)
		.AnnotationSetText(line, text)
		}

	ClearAllAnnotations()
		{ .AnnotationClearAll() }

	SetVisible(.Set)
		{
		.AnnotationSetVisible(.Set is true
			? SC.ANNOTATION_BOXED
			: SC.ANNOTATION_HIDDEN)
		}

	GetVisible()
		{ return .AnnotationGetVisible() isnt SC.ANNOTATION_HIDDEN }

	ContextMenu()
		{ #("Show/Hide Annotations\tF7") }

	On_ShowHide_Annotations()
		{ .SetVisible(not .GetVisible()) }
	}
