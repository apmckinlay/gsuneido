// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(sourceTran, attachmentField, labelFilter = false)
		{
		attachments = Object()
		.loop(sourceTran, attachmentField)
			{
			if "" isnt path = OpenImageWithLabelsControl.SplitFullPath(it, labelFilter)
				attachments.Add(path)
			}
		return attachments
		}

	ListWithLabels(sourceTran, attachmentField)
		{
		attachments = Object()
		.loop(sourceTran, attachmentField)
			{
			pathOb = OpenImageWithLabelsControl.SplitLabel(it)
			path = OpenImageWithLabelsControl.FullPath(pathOb.file, pathOb.subfolder)
			if path isnt ""
				attachments.Add(Object(:path, labels: pathOb.labels))
			}
		return attachments
		}

	loop(sourceTran, attachmentField, block)
		{
		for level in sourceTran[attachmentField]
			for m in level.Members().Sort!()
				block(level[m])
		}
	}