Function IsProcessRunning( strComputer, strProcess )
    Dim Process, strObject
    IsProcessRunning = False
    strObject   = "winmgmts://" & strComputer
    For Each Process in GetObject( strObject ).InstancesOf( "win32_process" )
    If UCase( Process.name ) = UCase( strProcess ) Then
        IsProcessRunning = True
        Exit Function
    End If
    Next
End Function
wScript.Echo IsProcessRunning(".", "Photoshop.exe")
