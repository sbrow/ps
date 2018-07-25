' Open photoshop.
Set app = CreateObject("Photoshop.Application")
if WScript.Arguments.Count = 0 then
    WScript.Echo "Missing parameters"
else
path = wScript.Arguments(0)
Set doc = app.Open(path)
wScript.echo doc.Name
end if