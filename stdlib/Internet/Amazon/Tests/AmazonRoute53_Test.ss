// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_extractMapping()
		{
		extract = AmazonRoute53.AmazonRoute53_extractMapping
		mapping = Object()
		next = extract(.response $ .footFirst $ .end, mapping)
		Assert(mapping is: #(
			'10 inbound-smtp.us-east-1.amazonaws.com': 'example.net'
			'99.99.99.99': 'sub.example.net'))
		Assert(next is: #(type: 'A', name: 'sub.example.net.'))

		mapping = Object()
		next = extract(.response $ .footLast $ .end, mapping)
		Assert(mapping is: #(
			'10 inbound-smtp.us-east-1.amazonaws.com': 'example.net'
			'99.99.99.99': 'sub.example.net'))
		Assert(next is: false)
		}

	response: `<listresourcerecordsetsresponse
		xmlns="https://route53.amazonaws.com/doc/2013-04-01/">
	<resourcerecordsets>
		<resourcerecordset>
			<name>
				example.net.
			</name>
			<type>
				MX
			</type>
			<ttl>
				1800
			</ttl>
			<resourcerecords>
				<resourcerecord>
					<value>
						10 inbound-smtp.us-east-1.amazonaws.com
					</value>
				</resourcerecord>
			</resourcerecords>
		</resourcerecordset>
		<resourcerecordset>
			<name>
				sub.example.net.
			</name>
			<type>
				A
			</type>
			<ttl>
				300
			</ttl>
			<resourcerecords>
				<resourcerecord>
					<value>
						99.99.99.99
					</value>
				</resourcerecord>
			</resourcerecords>
		</resourcerecordset>
	</resourcerecordsets>`

	footFirst: '<istruncated>
		true
	</istruncated>
	<nextrecordname>
		sub.example.net.
	</nextrecordname>
	<nextrecordtype>
		A
	</nextrecordtype>'

	footLast: '<IsTruncated>false</IsTruncated>'

	end: '	<maxitems>
		100
		</maxitems>
	</listresourcerecordsetsresponse>'
	}