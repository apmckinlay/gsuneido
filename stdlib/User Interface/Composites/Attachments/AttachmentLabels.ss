// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// can possibly refactor to share a class with the EmailAddresses class
class
	{
	Ensure()
		{
		extra = GetContributions('ExtraAttachmentLabelFields')
		fields = extra.Map({ it.field }).Join(', ')
		indexes = extra.Map({ it.index }).Join('\r\n')
		Database('ensure attachment_labels (attachlbl_label' $ Opt(', ', fields) $ ')
			key(attachlbl_label)' $ Opt('\r\n', indexes))
		}

	GetLabels(prefix, limit)
		{
		list = QueryAll('attachment_labels
			where attachlbl_label >= ' $ Display(prefix) $
			' and attachlbl_label < ' $ Display(prefix $ '~'),
			:limit)
		return list.Map!({ it.attachlbl_label }).UniqueValues().Sort!()
		}

	DeleteLabels(label, t)
		{
		DoWithTran(t, update:)
			{|t|
			t.QueryDo('delete attachment_labels
				where attachlbl_label is ' $ Display(label))
			}
		}

	OutputLabels(label, t = false)
		{
		extra = GetContributions('ExtraAttachmentLabelFields')
		rec = [attachlbl_label: label.Lower()]
		extra.Each({ (it.beforeOutput)(rec) })
		DoWithTran(t, update:)
			{|t|
			try
				t.QueryOutput('attachment_labels', rec)
			catch (unused, 'duplicate key')
				; // ignore
			}
		}
	}
