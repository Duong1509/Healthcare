const userHandler = require("./handlers/user/userHandler")
const recordHandler = require("./handlers/blockchain/recordHandler")
const upload = require("multer")({
    limits: {
        fileSize: 10 * 1024 * 1024
    }
});
module.exports = function (app) {
    //dang ky
    app.post('/api/user/sign-up', upload.single("information"), userHandler.signUp);
    app.post('/api/user/sign-in', userHandler.signIn);
    // tao benh an
    app.post("/api/blockchain/createMB", upload.single("information"), recordHandler.createMB)
    app.post("/api/blockchain/grantPermission", recordHandler.grantPermission)
    app.post("/api/blockchain/createRecord", upload.single("data"), recordHandler.createRecord)
    app.post("/api/blockchain/generateRecord", recordHandler.generateRecord)
    app.post("/api/blockchain/updateRecord", upload.single("data"), recordHandler.updateRecord)
    app.get("/api/blockchain/readRecord", recordHandler.readRecord)
    app.get("/api/blockchain/readMedicine", recordHandler.readMedicine)
    app.get("/api/blockchain/readAll", recordHandler.readAll)
    app.get("/api/blockchain/readMultiResult", recordHandler.readMultiResult)
    app.delete("/api/blockchain/deleteRecord", recordHandler.deleteRecord)
    app.delete("/api/blockchain/deletePermission", recordHandler.deletePermission)
}