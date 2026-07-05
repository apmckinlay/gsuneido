// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
#(
	httpsReturn: #(header: 'HTTP/1.1 200 OK
x-amz-id-2: 7C/6bqAJVI3OsvADfYx1BYNDZ/xQLfmlJmiSGZS4WzSBwYpYBHITovjZ9wKrYr3VHuOd/iXbO2w=
x-amz-request-id: 01142F3A6082687B
Date: Fri, 17 Aug 2018 21:03:58 GMT
ETag: "046e80d45d22a34377a87c2506560664"
Content-Length: 0
Server: AmazonS3

')
	emptyListXML: #(content: '<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="https://s3.amazonaws.com/doc/2006-03-01/"><Name>Axon</Name>' $
'<Prefix></Prefix><Marker></Marker><MaxKeys>1000</MaxKeys>' $
'<IsTruncated>false</IsTruncated></ListBucketResult>', header: 'HTTP/1.1 200 OK')

	hasFileXML: #(content: '<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="https://s3.amazonaws.com/doc/2006-03-01/"><Name>Axon</Name>' $
'<Prefix></Prefix><Marker></Marker><MaxKeys>1000</MaxKeys>' $
'<IsTruncated>false</IsTruncated><Contents><Key>test.txt</Key>' $
'<LastModified>2011-11-09T17:16:28.000Z</LastModified>' $
'<ETag>&quot;f7723c8c7130bcd4f0ad7c647d3824bf&quot;</ETag><Size>171</Size>' $
'<Owner><ID>d0c0c41cc87dc6a143e9ec7a06d527e62ea90200b6d8c9a2338924f6323c3b7e</ID>' $
'<DisplayName>axoneta</DisplayName></Owner><StorageClass>STANDARD</StorageClass>' $
'</Contents></ListBucketResult>', header: 'HTTP/1.1 200 OK' )

	hasFileXMLTruncated: #(content: '<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="https://s3.amazonaws.com/doc/2006-03-01/"><Name>Axon</Name>' $
'<Prefix></Prefix><Marker></Marker><MaxKeys>1000</MaxKeys>' $
'<IsTruncated>true</IsTruncated><Contents><Key>test0.txt</Key>' $
'<LastModified>2011-11-09T17:16:28.000Z</LastModified>' $
'<ETag>&quot;f7723c8c7130bcd4f0ad7c647d3824bf&quot;</ETag><Size>171</Size>' $
'<Owner><ID>d0c0c41cc87dc6a143e9ec7a06d527e62ea90200b6d8c9a2338924f6323c3b7e</ID>' $
'<DisplayName>axoneta</DisplayName></Owner><StorageClass>STANDARD</StorageClass>' $
'</Contents></ListBucketResult>', header: 'HTTP/1.1 200 OK' )

	invalidXML: #(content: '<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="https://s3.amazonaws.com/doc/2006-03-01/"><Name>Axon</Name>' $
'<Prefix></Prefix><Marker></Marker><MaxKeys>1000</MaxKeys>' $
'<IsTruncated>false</IsTruncated><Contents><Key>test.txt</Key>' $
'<LastModified>2011-11-09T17:16:28.000Z</LastModified>' $
'&quot;f7723c8c7130bcd4f0ad7c647d3824bf&quot;</ETag><Size>171</Size>' $
'<Owner><ID>d0c0c41cc87dc6a143e9ec7a06d527e62ea90200b6d8c9a2338924f6323c3b7e</ID>' $
'<DisplayName>axoneta</DisplayName></Owner><StorageClass>STANDARD</StorageClass>' $
'</Contents></ListBucketResult>', header: 'HTTP/1.1 200 OK')

	hasFolderFileXML: #(content: '<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="https://s3.amazonaws.com/doc/2006-03-01/"><Name>Axon</Name>' $
'<Prefix>testFolder</Prefix><Marker></Marker><MaxKeys>1000</MaxKeys>' $
'<IsTruncated>false</IsTruncated><Contents><Key>testFolder/</Key>' $
'<LastModified>2011-11-09T17:16:28.000Z</LastModified>' $
'<ETag>&quot;f7723c8c7130bcd4f0ad7c647d3824bf&quot;</ETag><Size>171</Size>' $
'<Owner><ID>d0c0c41cc87dc6a143e9ec7a06d527e62ea90200b6d8c9a2338924f6323c3b7e</ID>' $
'<DisplayName>axoneta</DisplayName></Owner><StorageClass>STANDARD</StorageClass>' $
'</Contents><Contents><Key>testFolder/folderTest.txt</Key>' $
'<LastModified>2011-11-09T17:16:28.000Z</LastModified>' $
'<ETag>&quot;f7723c8c7130bcd4f0ad7c647d3824bf&quot;</ETag><Size>171</Size>' $
'<Owner><ID>d0c0c41cc87dc6a143e9ec7a06d527e62ea90200b6d8c9a2338924f6323c3b7e</ID>' $
'<DisplayName>axoneta</DisplayName></Owner><StorageClass>STANDARD</StorageClass>' $
'</Contents></ListBucketResult>', header: 'HTTP/1.1 200 OK' )

	hasVersionsOfFileXML: #(content: '<?xml version="1.0" encoding="UTF-8"?>
<ListVersionsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Name>Axon</Name>' $
'<Prefix>axon</Prefix><KeyMarker></KeyMarker><VersionIdMarker></VersionIdMarker>' $
'<MaxKeys>1000</MaxKeys><IsTruncated>false</IsTruncated><Version>' $
'<Key>axon.betacad.gpg</Key><VersionId>iMccVr_1xYzUJI40moR0dsGHBlgYDHkU</VersionId>' $
'<IsLatest>true</IsLatest><LastModified>2021-01-11T19:33:41.000Z</LastModified>' $
'<ETag>&quot;18cd086fea3f87f4aadcdace2a520402-3&quot;</ETag><Size>48505180</Size>' $
'<Owner><ID>d0c0c41cc87dc6a143e9ec7a06d527e62ea90200b6d8c9a2338924f6323c3b7e</ID>' $
'<DisplayName>axoneta</DisplayName></Owner><StorageClass>STANDARD</StorageClass>' $
'</Version><Version><Key>axon.betacad.gpg</Key>' $
'<VersionId>70YsAq4iFPYCwHI0pScpa4D_TWG07fod</VersionId><IsLatest>false</IsLatest>' $
'<LastModified>2021-01-08T20:58:25.000Z</LastModified>' $
'<ETag>&quot;0078fb0a375190512d49a136bdba17fa-3&quot;</ETag><Size>48522359</Size>' $
'<Owner><ID>d0c0c41cc87dc6a143e9ec7a06d527e62ea90200b6d8c9a2338924f6323c3b7e</ID>' $
'<DisplayName>axoneta</DisplayName></Owner><StorageClass>STANDARD</StorageClass>' $
'</Version></ListVersionsResult>', header: "HTTP/1.1 200 OK")

	httpsFailure: #(header: "HTTP/1.1 404 Not Found
x-amz-request-id: 6FB389355EFB47FD
x-amz-id-2: +iu4+FltVdCO8QUeQXA18lQ56nspn1/i0O5XayKIGs1DjhyO/g7k8cUCQjWtjz6k8mp2lRamxIw=
Content-Type: application/xml
Transfer-Encoding: chunked
Date: Mon, 20 Aug 2018 17:40:46 GMT
Connection: close
Server: AmazonS3

")

	httpsMetaDataReturn: #(header: 'HTTP/1.1 200 OK
x-amz-id-2: LyJL4ZePqbWDCRqTFC4OTnQ8kfZyelZIHffQ/njOaZt8rRdG6NBQmaV/5vRI3kkHguO9xByB85k=
x-amz-request-id: 3E4018B0A9A9152E
Date: Thu, 21 Jan 2021 14:24:21 GMT
Last-Modified: Tue, 19 Jan 2021 19:26:24 GMT
x-amz-restore: ongoing-request="false", expiry-date="Mon, 25 Jan 2021 00:00:00 GMT"
ETag: "f91841cee42c354f479d078ad19d765e-3"
x-amz-tagging-count: 1
x-amz-storage-class: GLACIER
x-amz-version-id: cx.oBOxbWMpYUGNhTVUfPut8n5jrwSSe
Accept-Ranges: bytes
Content-Type: application/x-www-form-urlencoded; charset=utf-8
Content-Length: 48635752
Server: AmazonS3'
)

	httpRestoreArchivedReturn: #(header: "HTTP/1.1 202 Accepted
x-amz-id-2: i8kYBHcCp5gax+ojoOp5hmcorAdX/rSu6m2vbVJyZlkIgnaOCGwz80wt2uvT6sJAOQp9mksrm4A=
x-amz-request-id: D25FE95E7F0B43FA
Date: Thu, 21 Jan 2021 16:34:21 GMT
x-amz-version-id: cx.oBOxbWMpYUGNhTVUfPut8n5jrwSSe
Content-Length: 0
Server: AmazonS3

", content: "HTTP/1.1 202 Accepted
x-amz-id-2: i8kYBHcCp5gax+ojoOp5hmcorAdX/rSu6m2vbVJyZlkIgnaOCGwz80wt2uvT6sJAOQp9mksrm4A=
x-amz-request-id: D25FE95E7F0B43FA
Date: Thu, 21 Jan 2021 16:34:21 GMT
x-amz-version-id: cx.oBOxbWMpYUGNhTVUfPut8n5jrwSSe
Content-Length: 0
Server: AmazonS3

")
)