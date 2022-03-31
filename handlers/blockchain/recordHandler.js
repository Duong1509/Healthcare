recordServices = require('../../services/recordServices')
const { json } = require('body-parser');
var xlsx = require('node-xlsx');
const { countDocuments } = require('../../services/db/db');
const APIKey = require('../../services/db/db')
exports.createMB = async function (req, res) {
    try {
        var content = Buffer.from(req.file.buffer).toString("utf-8");
        params = {
            CCCD: req.body.CCCD,

            information: content
        }
        let message =
            await recordServices.Invokecc("CreateMB", params, req.body.CCCD)
        if (message.success) {
            return res.status(200).send(message)
        }
        else {
            return res.status(400).send(message)
        }
    } catch (err) {
        console.log(err)
        return res.status(500).send({
            success: false,
            error: err
        })
    }
}


exports.grantPermission = async function (req, res) {
    try {
        if (req.body.readPermission == "" && req.body.writePermission == "") {
            return res.status(400).send({
                success: false,
                result: "Them nguoi duoc cap quyen"
            })
        }
        if (req.body.readPermission != "") {
            var userdoc = await APIKey.findOne({ hashKey: req.body.readPermission }).exec()
        }
        if(userdoc == null) {
            console.log("khong co user")
            return res.status(400).send({
                success: false,
                result: "Nguoi duoc cap quyen doc khong ton tai"
            })
        }

        if (req.body.writePermission != "") {
            var  userghi = await APIKey.findOne({ hashKey: req.body.writePermission }).exec()
        }

        if(userghi == null) {
            console.log("khong co user")
            return res.status(400).send({
                success: false,
                result: "Nguoi duoc cap quyen ghi khong ton tai"
            })
        }

        params = {
            CCCD: req.body.CCCD,

            readPermission: req.body.readPermission,
            writePermission: req.body.writePermission
        }
        let message =
            await recordServices.Invokecc("GrantPermission", params, req.body.CCCD)
        if (message.success) {
            return res.status(200).send(message)
        }
        else {
            return res.status(400).send(message)
        }
    } catch (err) {

        return res.status(500).send({
            success: false,
            error: err
        })
    }
}

exports.deletePermission = async function (req, res) {
    try {
        if (req.body.readPermission == "" && req.body.writePermission == "") {
            return res.status(400).send({
                success: false,
                result: "Them nguoi bi xoa quyen"
            })
        }
        params = {
            CCCD: req.body.CCCD,

            readPermission: req.body.readPermission,
            writePermission: req.body.writePermission
        }
        let message =
            await recordServices.Invokecc("DeletePermission", params, req.body.CCCD)
        if (message.success) {
            return res.status(200).send(message)
        }
        else {
            return res.status(400).send(message)
        }
    } catch (err) {
        console.log(err)
        return res.status(500).send({
            success: false,
            error: err
        })
    }
}


exports.createRecord = async function (req, res) {
    try {
        var content = Buffer.from(req.file.buffer).toString("utf-8");
        params = {
            CCCD: req.body.CCCD,
            sickType: req.body.sickType,
            data: content,
            medicine: JSON.stringify(JSON.parse(content).Medicine)

        }
        console.log(params.medicine)
        let message =
            await recordServices.Invokecc("CreateRecord", params, req.body.CCCDDoctor)
        if (message.success) {
            return res.status(200).send(message)
        }
        else {
            return res.status(400).send(message)
        }
    } catch (err) {
        console.log(err)
        return res.status(500).send({
            success: false,
            error: err
        })
    }
}

exports.updateRecord = async function (req, res) {
    try {
        var content = Buffer.from(req.file.buffer).toString("utf-8");
        params = {
            CCCD: req.body.CCCD,
            sickType: req.body.sickType,

            data: content,
            time: req.body.time,
            medicine: JSON.stringify(JSON.parse(content).Medicine)

        }
        let message =
            await recordServices.Invokecc("UpdateRecord", params, req.body.CCCDDoctor)
        if (message.success) {
            return res.status(200).send(message)
        }
        else {
            return res.status(400).send(message)
        }
    } catch (err) {
        console.log(err)
        return res.status(500).send({
            success: false,
            error: err
        })
    }
}

exports.deleteRecord = async function (req, res) {
    try {
        params = {
            CCCD: req.body.CCCD,
            sickType: req.body.sickType,
            time: req.body.time

        }
        let message =
            await recordServices.Invokecc("DeleteRecord", params, req.body.CCCDDoctor)
        if (message.success) {
            return res.status(200).send(message)
        }
        else {
            return res.status(400).send(message)
        }
    } catch (err) {
        console.log(err)
        return res.status(500).send({
            success: false,
            error: err
        })
    }
}




exports.readMultiResult = async function (req, res) {
    try {
        params = {
            CCCD: req.body.CCCD,
            sickType: req.body.sickType

        }
        let message =
            await recordServices.Querycc("ReadMultiResult", params, req.body.CCCDDoctor)
        if (message.success) {
            return res.status(200).send(message)
        }
        else {
            return res.status(400).send(message)
        }
    } catch (err) {
        console.log(err)
        return res.status(500).send({
            success: false,
            error: err
        })
    }
}
exports.readRecord = async function (req, res) {
    try {
        params = {
            CCCD: req.body.CCCD,
            sickType: req.body.sickType,
            count: req.body.count

        }
        let message =
            await recordServices.Querycc("ReadRecord", params, req.body.CCCDDoctor)
        if (message.success) {
            return res.status(200).send(message)
        }
        else {
            return res.status(400).send(message)
        }
    } catch (err) {
        console.log(err)
        return res.status(500).send({
            success: false,
            error: err
        })
    }
}
exports.readMedicine = async function (req, res) {
    try {
        params = {
            CCCD: req.body.CCCD

        }
        let message =
            await recordServices.Querycc("ReadMedicine", params, req.body.CCCDDoctor)
        if (message.success) {
            return res.status(200).send(message)
        }
        else {
            return res.status(400).send(message)
        }
    } catch (err) {
        console.log(err)
        return res.status(500).send({
            success: false,
            error: err
        })
    }
}

exports.readAll = async function (req, res) {
    try {
        params = {
            CCCD: req.body.CCCD,
            hashKey: req.body.hashKey
        }
        let message =
            await recordServices.Querycc("ReadAll", params, req.body.CCCD)
        if (message.success) {
            return res.status(200).send(message)
        }
        else {
            return res.status(400).send(message)
        }
    } catch (err) {
        console.log(err)
        return res.status(500).send({
            success: false,
            error: err
        })
    }
}

exports.generateRecord = async function (req, res) {
    try {
        params = {
            CCCD: req.body.CCCD

        }
        let message =
            await recordServices.Invokecc("GenerateRecord", params, req.body.CCCDDoctor)
        if (message.success) {
            return res.status(200).send(message)
        }
        else {
            return res.status(400).send(message)
        }
    } catch (err) {
        console.log(err)
        return res.status(500).send({
            success: false,
            error: err
        })
    }
}

