set App = CreateObject("Photoshop.Application")
set Doc = App.activeDocument
Doc.Close(CInt(wScript.Arguments(0)))