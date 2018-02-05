
Dim app
Set app = CreateObject("Photoshop.Application")
if WScript.Arguments.Count = 0 then
    WScript.Echo "Missing parameters"
else
	path = wScript.Arguments(0)
	folder = wScript.Arguments(1)
	app.DoJavaScriptFile path, Array(Folder)
end if