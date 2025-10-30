## Working With Images

**Category:** User Interface

**Problem**

How to store images in a table and retrieve and display them.

**Ingredients**

[Image](<../Language/Reference/Image.md>),
[ImageControl](<../User Interfaces/Reference/ImageControl.md>), 
[OpenImageControl](<../User Interfaces/Reference/OpenImageControl.md>), 
[AccessControl](<../User Interfaces/Reference/AccessControl.md>),
[ImageFormat](<../Reports/Reference/ImageFormat.md>),
[Query1](<../Database/Reference/Query1.md>),
[OpenFileControl](<../User Interfaces/Reference/OpenFileControl.md>)

**Recipe**

The easiest way to handle images with Suneido is to use [Image](<../Language/Reference/Image.md>). Image uses the Windows OleLoadPicture, which handles bmp, gif, jpg, ico, emf, and wmf images (but not png). However, normally you won't want to use Image directly; you'll want to use the classes that are built on top of Image.
[ImageControl](<../User Interfaces/Reference/ImageControl.md>)
: Displays an image on the screen.

[OpenImageControl](<../User Interfaces/Reference/OpenImageControl.md>)
: Displays an image on the screen and allows the user to choose an image from a file.

[ImageFormat](<../Reports/Reference/ImageFormat.md>)
: Prints an image on a report.

For example, to display an image file using ImageControl:

``` suneido
Window(#(Image 'c:/windows/clouds.bmp'))
```

Note: Depending on which version of Windows you are running, you may need to use different file names to make these examples work.

Or to allow choosing the image with OpenImageControl:

``` suneido
Window(#(OpenImage xmin: 200, ymin: 200))
```

You can use an Access Control to create and view images in a table. First, create a table. Enter and run this from QueryView:

``` suneido
create imagetable (image_name, image_field) key(image_name)
```

Then enter this in a library (e.g. mylib) as My_ImageAccess:

``` suneido
#(Access
    'imagetable'
    (Form
        (image_name)
        (image_field))
    )
```

You will also need Field_ definitions:
<pre>
<b>Field_image_name</b>
Field_string
    {
    Prompt: 'Image Name'
    }

<b>Field_image_field</b>
Field_string
    {
    Prompt: ''
    Heading: 'Image'
    Control: (OpenImage xmin: 200 ymin: 200)
    }
</pre>

Now you should be able to run it with:

``` suneido
Window(My_ImageAccess)
```

Add an image to the table with a name of "sample". You can then retrieve and display an image from the table using Query1 and ImageControl. For example:

``` suneido
imagename = "sample"
x = Query1("imagetable where image_name is " $ Display(imagename))
Window(Object("Image", x.image_field))
```

**Discussion**

Images can be stored in tables in two ways.

The first way is to store the actual image in the table.  A string containing image information will be stored in the table.  If you look at the "imagebook" table, for example, in Query View, you will see that it contains several images.  The name of the image is stored in the name field.  The actual image is in the text field.

The second way is to store the path to the image.  If the image to be displayed was always going to be in the same place, the path can be stored instead of the actual image (e.g. "C:\WINDOWS\Clouds.bmp").

ImageControl can be used to display images that are stored as a path or as the actual image.  OpenImageControl uses a combination of ImageControl and OpenFileControl.  By double-clicking on the image, the user is allowed to browse to an image file that they would like to choose.  They can also select whether the path to the image or the actual image should be stored.