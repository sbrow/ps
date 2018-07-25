DIM objApp
SET objApp = CreateObject("Photoshop.Application")
REM Use dialog mode 3 for show no dialogs
DIM dialogMode
dialogMode = 3
DIM idsetd
idsetd = objApp.CharIDToTypeID("setd")
    DIM desc2
    SET desc2 = CreateObject("Photoshop.ActionDescriptor")
    DIM idnull
    idnull = objApp.CharIDToTypeID("null")
        DIM ref2
        SET ref2 = CreateObject("Photoshop.ActionReference")
        DIM idPrpr
        idPrpr = objApp.CharIDToTypeID("Prpr")
        DIM idLefx
        idLefx = objApp.CharIDToTypeID("Lefx")
        Call ref2.PutProperty(idPrpr, idLefx)
        DIM idLyr
        idLyr = objApp.CharIDToTypeID("Lyr ")
        DIM idOrdn
        idOrdn = objApp.CharIDToTypeID("Ordn")
        DIM idTrgt
        idTrgt = objApp.CharIDToTypeID("Trgt")
        Call ref2.PutEnumerated(idLyr, idOrdn, idTrgt)
    Call desc2.PutReference(idnull, ref2)
    DIM idT
    idT = objApp.CharIDToTypeID("T   ")
        DIM desc3
        SET desc3 = CreateObject("Photoshop.ActionDescriptor")
        DIM idScl
        idScl = objApp.CharIDToTypeID("Scl ")
        DIM idPrc
        idPrc = objApp.CharIDToTypeID("#Prc")
        Call desc3.PutUnitDouble(idScl, idPrc, 416.666667)
        DIM idSoFi
        idSoFi = objApp.CharIDToTypeID("SoFi")
            DIM desc4
            SET desc4 = CreateObject("Photoshop.ActionDescriptor")
            DIM idenab
            idenab = objApp.CharIDToTypeID("enab")
            Call desc4.PutBoolean(idenab, True)
            DIM idMd
            idMd = objApp.CharIDToTypeID("Md  ")
            DIM idBlnM
            idBlnM = objApp.CharIDToTypeID("BlnM")
            DIM idNrml
            idNrml = objApp.CharIDToTypeID("Nrml")
            Call desc4.PutEnumerated(idMd, idBlnM, idNrml)
            DIM idOpct
            idOpct = objApp.CharIDToTypeID("Opct")
            idPrc = objApp.CharIDToTypeID("#Prc")
            Call desc4.PutUnitDouble(idOpct, idPrc, 100.000000)
            DIM idClr
            idClr = objApp.CharIDToTypeID("Clr ")
                DIM desc5
                SET desc5 = CreateObject("Photoshop.ActionDescriptor")
                DIM idRd
                idRd = objApp.CharIDToTypeID("Rd  ")
                Call desc5.PutDouble(idRd, CInt(wScript.Arguments(0)))
                DIM idGrn
                idGrn = objApp.CharIDToTypeID("Grn ")
                Call desc5.PutDouble(idGrn,CInt(wScript.Arguments(1)))
                DIM idBl
                idBl = objApp.CharIDToTypeID("Bl  ")
                Call desc5.PutDouble(idBl, CInt(wScript.Arguments(2)))
            DIM idRGBC
            idRGBC = objApp.CharIDToTypeID("RGBC")
            Call desc4.PutObject(idClr, idRGBC, desc5)
        idSoFi = objApp.CharIDToTypeID("SoFi")
        Call desc3.PutObject(idSoFi, idSoFi, desc4)
        DIM idFrFX
        idFrFX = objApp.CharIDToTypeID("FrFX")
            DIM desc6
            SET desc6 = CreateObject("Photoshop.ActionDescriptor")
            idenab = objApp.CharIDToTypeID("enab")
            Call desc6.PutBoolean(idenab, True)
            DIM idStyl
            idStyl = objApp.CharIDToTypeID("Styl")
            DIM idFStl
            idFStl = objApp.CharIDToTypeID("FStl")
            DIM idOutF
            idOutF = objApp.CharIDToTypeID("OutF")
            Call desc6.PutEnumerated(idStyl, idFStl, idOutF)
            DIM idPntT
            idPntT = objApp.CharIDToTypeID("PntT")
            DIM idFrFl
            idFrFl = objApp.CharIDToTypeID("FrFl")
            DIM idSClr
            idSClr = objApp.CharIDToTypeID("SClr")
            Call desc6.PutEnumerated(idPntT, idFrFl, idSClr)
            idMd = objApp.CharIDToTypeID("Md  ")
            idBlnM = objApp.CharIDToTypeID("BlnM")
            idNrml = objApp.CharIDToTypeID("Nrml")
            Call desc6.PutEnumerated(idMd, idBlnM, idNrml)
            idOpct = objApp.CharIDToTypeID("Opct")
            idPrc = objApp.CharIDToTypeID("#Prc")
            Call desc6.PutUnitDouble(idOpct, idPrc, 100.000000)
            DIM idSz
            idSz = objApp.CharIDToTypeID("Sz  ")
            DIM idPxl
            idPxl = objApp.CharIDToTypeID("#Pxl")
            Call desc6.PutUnitDouble(idSz, idPxl, CLng(wScript.Arguments(3)))
            idClr = objApp.CharIDToTypeID("Clr ")
                DIM desc7
                SET desc7 = CreateObject("Photoshop.ActionDescriptor")
                idRd = objApp.CharIDToTypeID("Rd  ")
                Call desc7.PutDouble(idRd, CInt(wScript.Arguments(4)))
                idGrn = objApp.CharIDToTypeID("Grn ")
                Call desc7.PutDouble(idGrn, CInt(wScript.Arguments(5)))
                idBl = objApp.CharIDToTypeID("Bl  ")
                Call desc7.PutDouble(idBl, CInt(wScript.Arguments(6)))
            idRGBC = objApp.CharIDToTypeID("RGBC")
            Call desc6.PutObject(idClr, idRGBC, desc7)
        idFrFX = objApp.CharIDToTypeID("FrFX")
        Call desc3.PutObject(idFrFX, idFrFX, desc6)
    idLefx = objApp.CharIDToTypeID("Lefx")
    Call desc2.PutObject(idT, idLefx, desc3)
Call objApp.ExecuteAction(idsetd, desc2, dialogMode)

