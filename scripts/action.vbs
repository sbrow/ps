set appRef = CreateObject("Photoshop.Application")
' No dialogs'
dlgMode = 3

set desc = CreateObject( "Photoshop.ActionDescriptor" )
set ref = CreateObject( "Photoshop.ActionReference" )
Call ref.PutName(appRef.CharIDToTypeID("Actn"), wScript.Arguments(1))
Call ref.PutName(appRef.CharIDToTypeID("ASet"), wScript.Arguments(0))
Call desc.PutReference(appRef.CharIDToTypeID("null"), ref)
Call appRef.ExecuteAction(appRef.CharIDToTypeID("Ply "), desc, dlgMode)