// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
// NOTE: Corresponding report requires EmailAttachments: Object() in its Params object
function (params, data, label, attachmentfield)
	{
	if params.ReportDestination isnt 'pdf' or label is false
		return

	for level in data[attachmentfield]
		for field in level.Members().Sort!()
			{
			filePath = OpenImageWithLabelsControl.SplitFullPath(level[field], label)
			if filePath isnt ""
				params.EmailAttachments.AddUnique(filePath)
			}
	}
