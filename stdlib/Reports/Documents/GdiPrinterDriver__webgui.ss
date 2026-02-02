// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HtmlDriver
	{
	New(dc/*unused*/)
		{
		.id = SuRenderBackend().NextId()
		}

	AddPage(dimens)
		{
		super.AddPage(dimens)
		.page = Object(:dimens)
		}

	EndPage()
		{
		SuRenderBackend().RecordAction(false, 'PrintEndPage', [.page, .id])
		}

	Getter_Page()
		{
		return .page
		}

	Finish(status = 0)
		{
		super.Finish(status)
		if status is ReportStatus.SUCCESS
			SuRenderBackend().RecordAction(false, 'PrintEndDoc', [.id])
		else
			{
			SuRenderBackend().CancelAction(false, 'PrintEndPage',
				{ |args| args[1] is .id })
			SuRenderBackend().RecordAction(false, 'PrintCancelDoc', [.id])
			}
		return status
		}
	}