' Close Photoshop
Set appRef = CreateObject("Photoshop.Application")

Do While appRef.Documents.Count > 0
	appRef.ActiveDocument.Close(wScript.Arguments(0))
Loop

appRef.Quit()