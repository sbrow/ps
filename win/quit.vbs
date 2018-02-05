' Close Photoshop
Set appRef = CreateObject("Photoshop.Application")

wScript.echo appRef.Documents.Count
Do While appRef.Documents.Count > 0
	appRef.ActiveDocument.Close(2)
Loop

appRef.Quit()