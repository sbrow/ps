
Dim appRef
Set appRef = CreateObject("Photoshop.Application")
if wScript.Arguments.Count = 0 then
    wScript.Echo "Missing parameters"
else
	path = wScript.Arguments(0)
	args = wScript.Arguments(1)
	appRef.DoJavaScriptFile path, Split(args, ",")
end if