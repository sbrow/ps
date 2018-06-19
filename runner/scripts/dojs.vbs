
Dim appRef
Set appRef = CreateObject("Photoshop.Application")
if wScript.Arguments.Count = 0 then
    wScript.Echo "Missing parameters"
else
	' wScript.Echo wScript.Arguments(0)
	' wScript.Echo wScript.Arguments(1)
	path = wScript.Arguments(0)
	args = wScript.Arguments(1)
	error = appRef.DoJavaScriptFile(path, Split(args, ",,"))
	if Not error = "true" and Not error = "[ActionDescriptor]" and Not error = "undefined" Then
		Err.raise 1, "dojs.vbs", error
	end if
end if