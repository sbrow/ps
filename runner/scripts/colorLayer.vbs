DIM objApp
SET objApp = CreateObject("Photoshop.Application")
DIM dialogMode
dialogMode = 3
DIM idsetd
idsetd = objApp.CharIDToTypeID("setd")
    DIM desc134
    SET desc134 = CreateObject("Photoshop.ActionDescriptor")
    DIM idnull
    idnull = objApp.CharIDToTypeID("null")
        DIM ref44
        SET ref44 = CreateObject("Photoshop.ActionReference")
        DIM idPrpr
        idPrpr = objApp.CharIDToTypeID("Prpr")
        DIM idLefx
        idLefx = objApp.CharIDToTypeID("Lefx")
        Call ref44.PutProperty(idPrpr, idLefx)
        DIM idLyr
        idLyr = objApp.CharIDToTypeID("Lyr ")
        DIM idOrdn
        idOrdn = objApp.CharIDToTypeID("Ordn")
        DIM idTrgt
        idTrgt = objApp.CharIDToTypeID("Trgt")
        Call ref44.PutEnumerated(idLyr, idOrdn, idTrgt)
    Call desc134.PutReference(idnull, ref44)
    DIM idT
    idT = objApp.CharIDToTypeID("T   ")
        DIM desc135
        SET desc135 = CreateObject("Photoshop.ActionDescriptor")
        DIM idScl
        idScl = objApp.CharIDToTypeID("Scl ")
        DIM idPrc
        idPrc = objApp.CharIDToTypeID("#Prc")
        Call desc135.PutUnitDouble(idScl, idPrc, 416.666667)
        DIM idSoFi
        idSoFi = objApp.CharIDToTypeID("SoFi")
            DIM desc136
            SET desc136 = CreateObject("Photoshop.ActionDescriptor")
            DIM idenab
            idenab = objApp.CharIDToTypeID("enab")
            Call desc136.PutBoolean(idenab, True)
            DIM idMd
            idMd = objApp.CharIDToTypeID("Md  ")
            DIM idBlnM
            idBlnM = objApp.CharIDToTypeID("BlnM")
            DIM idNrml
            idNrml = objApp.CharIDToTypeID("Nrml")
            Call desc136.PutEnumerated(idMd, idBlnM, idNrml)
            DIM idOpct
            idOpct = objApp.CharIDToTypeID("Opct")
            idPrc = objApp.CharIDToTypeID("#Prc")
            Call desc136.PutUnitDouble(idOpct, idPrc, 100.000000)
            DIM idClr
            idClr = objApp.CharIDToTypeID("Clr ")
                DIM desc137
                SET desc137 = CreateObject("Photoshop.ActionDescriptor")
                DIM idRd
                idRd = objApp.CharIDToTypeID("Rd  ")
                Call desc137.PutDouble(idRd, CInt(wScript.Arguments(0)))
                ' Call desc137.PutDouble(idRd, 255)
                DIM idGrn
                idGrn = objApp.CharIDToTypeID("Grn ")
                Call desc137.PutDouble(idGrn, Cint(wScript.Arguments(1)))
                ' Call desc137.PutDouble(idGrn, 255)
                DIM idBl
                idBl = objApp.CharIDToTypeID("Bl  ")
                Call desc137.PutDouble(idBl, CInt(wScript.Arguments(2)))
                ' Call desc137.PutDouble(idBl, 255)
            DIM idRGBC
            idRGBC = objApp.CharIDToTypeID("RGBC")
            Call desc136.PutObject(idClr, idRGBC, desc137)
        idSoFi = objApp.CharIDToTypeID("SoFi")
        Call desc135.PutObject(idSoFi, idSoFi, desc136)
    idLefx = objApp.CharIDToTypeID("Lefx")
    Call desc134.PutObject(idT, idLefx, desc135)
Call objApp.ExecuteAction(idsetd, desc134, dialogMode)