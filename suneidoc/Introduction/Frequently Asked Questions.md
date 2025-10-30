## Frequently Asked Questions

[What is Suneido&trade;?](<#What is Suneido>)   
[Why integrated?](<#Why integrated>)   
[What does it run on?](<#What does it run on>)   
[What is it for?](<#What is it for>)   
[How much does it cost?](<#How much does it cost>)   
[Why another language?](<#Why another language>)   
[Why another database?](<#Why another database>)   
[Why not SQL?](<#Why not SQL>)   
[Can I use a different database with Suneido?](<#Can I use a different database with Suneido>)   
[Can I access Suneido databases from other languages?](<#Can I access Suneido databases from other languages>)   
[What does the name "Suneido" mean?](<#What does the name Suneido mean>)   
[Can I download a copy of the on-line Manual?](<#Can I download a copy of the on-line Manual>)   
[How do I get started with Suneido?](<#How do I get started with Suneido>)   
[Can I sell applications written with Suneido?](<#Can I sell applications written with Suneido>)   
[Can Suneido produce stand-alone executables?](<#Can Suneido produce stand-alone executables>)   
<span id="What is Suneido">What is Suneido?</span>
: Suneido&trade; is a complete, integrated application platform that includes an object-oriented language, client-server database, and user interface and reporting frameworks to enable you to create, deploy, and maintain applications easier and quicker, with more features and less code. Suneido is Open Source Software - it is provided free, with complete source code, to give programmers full control over their development projects.

<span id="Why integrated">Why integrated?</span>
: "*The real frustrating thing with modern day programming - where you try to assemble a bunch of components, which should fit together - is making the tools work together as gracefully as possible. Say you have a rough outline for assembling a Web front end or a Visual Basic front end with Java glue, to a database back end ... two to three to four things that combine together. It's bad now and going to get worse.*" - Jon Bentley, member of the technical staff at Bell Labs and author of the classic "Programming Pearls"

<span id="What does it run on">What does it run on?</span>
: Currently (2024) Suneido runs on Windows, Linux, and MacOS. However the IDE only runs on Windows. Application GUI's can run with a web browser client using suneido.js

<span id="What is it for">What is it for?</span>
: Suneido is designed to build, deploy, and maintain information systems such as management information systems, accounting systems, or vertical applications. It is suitable for most applications that require a mix of user interface, database, and reporting.

<span id="How much does it cost">How much does it cost?</span>
: Nothing. Suneido is **free**. This is the complete product including all the source code. It's not a trial version or shareware.

<span id="Why another language">Why another language?</span>
: To be most effective for a given type of work a language has to be designed for that work. It's a question of having the right tool for the job. We needed a simple, easy to learn, powerful language that was dynamic, like Lisp or Smalltalk, but had a familiar syntax like C++ or Java. We also needed tight integration with the database, and the right facilities to build application frameworks. And we needed features like garbage collection and exception handling.

<span id="Why another database">Why another database?</span>
: We didn't feel we could attain the tight integration we wanted if we used a third party database. We also wanted to keep Suneido easy to learn and deploy.

<span id="Why not SQL">Why not SQL?</span>
: Suneido uses a simple yet powerful relational algebra query language rather than SQL.   
"*... it cannot be denied that SQL in its present form leaves rather a lot to be desired ... the language is filled with numerous restrictions, ad hoc constructs, and annoying special rules. These factors in turn make the language hard to define, describe, teach, learn, remember, apply, and implement.*" C.J.Date, author of the classic "Introduction to Database Systems".

<span id="Can I use a different database with Suneido">Can I use a different database with Suneido?</span>
: Certainly. We haven't written any interfaces to other databases, but you are welcome to, and we look forward to contributions in this area. However, unless you
have specific requirements we encourage you to consider using Suneido's database. Its tight integration with the language and the libraries make it an obvious choice.

<span id="Can I access Suneido databases from other languages">Can I access Suneido databases from other languages?</span>
: Suneido's client and server communicate using a straightforward protocol across TCP/IP. It would be relatively easy to write an interface for another language using Suneido's client code as a base. Or a simple Suneido server could be written to communicate using a different protocol, similar to the Suneido SMTP, POP3, and HTTP servers.

<span id="What does the name Suneido mean">What does the name "Suneido" mean?</span>
: "Suneido" comes from the ancient Greek meaning roughly "absolute knowledge" - knowledge that you know that you know. We picked it because it sounded good and no one else was using it. Its oriental "sound" led to the Karate font and the dragons.

<span id="Can I download a copy of the on-line Manual">Can I download a copy of the on-line Manual?</span>
: The Suneido download includes the same documentation as the on-line manual. You can access it from the IDE windows via the Help menu or by pressing F1. If you want a standalone copy of the documentation choose Edit a Book from the IDE menu, enter "suneidoc" as the name of the book, choose Export Html from the File menu of Book Edit, and enter the name of the directory to be created containing the manual as individual html files. You can then access these files using a browser such as Amaya, Netscape, or Internet Explorer. If you can't run Suneido for some reason, contact us and we'll be happy to supply you with a copy of the html files.

<span id="How do I get started with Suneido">How do I get started with Suneido?</span>
: You probably want to start by going through the 
[Getting Started](<../Getting Started.md>) section of the manual. You might want to also have a look at the 
[Tools](<../Tools.md>) section of the manual. Once you've covered the basics you may want to have a look at some of the sample applications in the 
[Download Library](<http://www.suneido.com/index.php?option=com_content&task=view&id=61&Itemid=45>). Don't overlook the standard library (stdlib) as a source of examples - it includes the source code for the entire development environment. If you've got questions (probably lots!) try posting them on the 
[Forum](<http://www.suneido.com/index.php?option=com_joomlaboard&Itemid=52>) Good luck!

<span id="Can I sell applications written with Suneido">Can I sell applications written with Suneido?</span>
:  Yes. There are very few restrictions on what you can do with Suneido. There are no fees or additional licenses required to use or distribute Suneido. The C++ source code for the Suneido executable itself is under the Gnu Public License (GPL). This means that if you include Suneido in your package, you must either include the source code or make it available (i.e. give them a link to the Suneido website). And if you modify or add to the C++ code, you must make your modifications available, also under the GPL. Any code you write in Suneido however, is yours - you can do whatever you like with it and it does not have to be under the GPL. For specifics, read the 
[GPL License](<GPL License.md>).

<span id="Can Suneido produce stand-alone executables">Can Suneido produce stand-alone executables?</span>
: The short answer is "No". However, it isn't really necessary since a Suneido application only requires two files - suneido.exe and suneido.db - the database, which includes the application code. Even if you combined suneido.exe and your application code into one executable, you'd still need a database file, and so you'd still have two files. In a multi-user situation, Suneido runs client-server, the clients just need suneido.exe (usually from a shared network location). The server instance of suneido.exe is the only one that accesses the database (suneido.db). The clients use TCP/IP to talk to the server (to access the application code and the data).