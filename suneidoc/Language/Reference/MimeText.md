### MimeText

``` suneido
(text = "", subtype = "plain", charset = "us-ascii")
```

A sub-class of [MimeBase](<MimeBase.md>) for a main type of "text". Sets Content-Transfer-Encoding to 7bit

For example:

``` suneido
MimeText("hello world").
    To("joe@hotmail.com").
    From('sue@mail.com').
    Subject('test').
    Date().
    Message_ID().
    ToString()
```

would produce:

``` suneido
Content-Type: text/plain; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit
To: joe@hotmail.com
From: sue@mail.com
Date: Fri, 19 Oct 2007 14:41:02 -0600
Subject: test
Message-ID: <2d91ef3e.5685.4d7c.b664.cdac914677f0@suneido.com>

hello world
```