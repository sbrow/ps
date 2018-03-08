Set appRef = CreateObject("Photoshop.Application")
dlgMode = 3 'No dialog
set d = CreateObject( "Photoshop.ActionDescriptor" )
Call d.PutEnumerated(appRef.CharIDToTypeID("PGIT"), appRef.CharIDToTypeID("PGIT"), appRef.CharIDToTypeID("PGIN"))
Call d.PutEnumerated(appRef.CharIDToTypeID("PNGf"), appRef.CharIDToTypeID("PNGf"), appRef.CharIDToTypeID("PGAd"))

SET desc = CreateObject( "Photoshop.ActionDescriptor" )
Call desc.PutObject( appRef.CharIDToTypeID("As  "), appRef.CharIDToTypeID("PNGF"), d)
Call desc.PutPath( appRef.CharIDToTypeID("In  "),  wScript.Arguments(0))
Call desc.PutBoolean( appRef.CharIDToTypeID("Cpy "), True )

Call appRef.ExecuteAction(appRef.CharIDToTypeID("save"), desc, dlgMode)