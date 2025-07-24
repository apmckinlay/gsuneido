// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
class
	{
	GetLastRowPos(ob)
		{
		size = ob.Size()
		return size is 0 ? 0 : size-1
		}

	GetNextAvailPos(ob)
		{
		if ob is ""
			return Object(row: 0, pos: 0)

		lastRow = .GetLastRowPos(ob)
		if not ob.Member?(lastRow)
			return Object(row: 0, pos: 0)

		if false is rowPos = .GetNextAvailRowPos(ob[lastRow])
			return Object(row: lastRow+1, pos: 0)

		return Object(row: lastRow, pos: rowPos)
		}

	GetNextAvailRowPos(ob)
		{
		if ob is ""
			return 0

		max = false
		for m in ob.Members()
			{
			pos = m[-1]
			if Number(pos) > max
				max = Number(pos)
			}
		if max is false
			return 0
		return max < 4 ? max+1 : false
		}

	GetValue(ob, row, pos)
		{
		if Number?(pos) or not pos.Prefix?('attachment')
			pos = 'attachment' $ pos
		if ob is "" or ob.Empty?() or not ob.Member?(row) or
			not ob[row].Member?(pos)
			return ""

		return ob[row][pos]
		}

	CompareAttachments(oldAttachmentFieldValue, newAttachmentFieldValue, label = false)
		{
		oldAttachments = .listAttachment(oldAttachmentFieldValue, label)
		newAttachments = .listAttachment(newAttachmentFieldValue, label)
		changed? = oldAttachments.Size() isnt newAttachments.Size()
		if changed? is false
			{
			oldAttachmentsFiles = oldAttachments.Map({ it.fullPath })
			newAttachmentsFiles = newAttachments.Map({ it.fullPath })
			changed? = oldAttachmentsFiles.Difference(newAttachmentsFiles).NotEmpty?() or
				newAttachmentsFiles.Difference(oldAttachmentsFiles).NotEmpty?()
			}
		return Object(:changed?, :newAttachments)
		}

	listAttachment(attachmentFieldValue, label)
		{
		if not Object?(attachmentFieldValue)
			return Object()

		attachments = Object()
		for row in attachmentFieldValue.Members()
			{
			for pos in .. AttachmentsRepeatControl.PerRow
				{
				attachment = OpenImageWithLabelsControl.SplitLabel(
					.GetValue(attachmentFieldValue, row, pos))
				if attachment.file is '' or
					label isnt false and not attachment.labels.Split(', ').Has?(label)
					continue
				attachment.fullPath = Opt(attachment.subfolder, '/') $ attachment.file
				attachment.row = row
				attachment.pos = pos
				attachments.Add(attachment)
				}
			}
		return attachments
		}
	}
