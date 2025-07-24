// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// GetDeviceCaps defines
#(
DRIVERVERSION: 0     /* Device driver version                    */
TECHNOLOGY:    2     /* Device classification                    */
HORZSIZE:      4     /* Horizontal size in millimeters           */
VERTSIZE:      6     /* Vertical size in millimeters             */
HORZRES:       8     /* Horizontal width in pixels               */
VERTRES:       10    /* Vertical height in pixels                */
BITSPIXEL:     12    /* Number of bits per pixel                 */
PLANES:        14    /* Number of planes                         */
NUMBRUSHES:    16    /* Number of brushes the device has         */
NUMPENS:       18    /* Number of pens the device has            */
NUMMARKERS:    20    /* Number of markers the device has         */
NUMFONTS:      22    /* Number of fonts the device has           */
NUMCOLORS:     24    /* Number of colors the device supports     */
PDEVICESIZE:   26    /* Size required for device descriptor      */
CURVECAPS:     28    /* Curve capabilities                       */
LINECAPS:      30    /* Line capabilities                        */
POLYGONALCAPS: 32    /* Polygonal capabilities                   */
TEXTCAPS:      34    /* Text capabilities                        */
CLIPCAPS:      36    /* Clipping capabilities                    */
RASTERCAPS:    38    /* Bitblt capabilities                      */
ASPECTX:       40    /* Length of the X leg                      */
ASPECTY:       42    /* Length of the Y leg                      */
ASPECTXY:      44    /* Length of the hypotenuse                 */

LOGPIXELSX:    88    /* Logical pixels/inch in X                 */
LOGPIXELSY:    90    /* Logical pixels/inch in Y                 */

SIZEPALETTE:  104    /* Number of entries in physical palette    */
NUMRESERVED:  106    /* Number of reserved entries in palette    */
COLORRES:     108    /* Actual color resolution                  */

PHYSICALWIDTH:   110 /* Physical Width in device units           */
PHYSICALHEIGHT:  111 /* Physical Height in device units          */
PHYSICALOFFSETX: 112 /* Physical Printable Area x margin         */
PHYSICALOFFSETY: 113 /* Physical Printable Area y margin         */
SCALINGFACTORX:  114 /* Scaling factor x                         */
SCALINGFACTORY:  115 /* Scaling factor y                         */
)