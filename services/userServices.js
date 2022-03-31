const FabricCaServices = require('fabric-ca-client');
const { Wallets, X509WalletMixin } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const APIKey = require('./db/db')
const crypto = require('crypto');
const mongoose = require('mongoose')
const xlsx = require("node-xlsx")



exports.RegisterUsers = async function RegisterUsers(req) {
    

    try {
        const ccpPath = path.resolve(__dirname, '..', 'blockchain', 'organizations', 'peerOrganizations', 'issuer.com', 'connection-issuer.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf-8'));
        const caUrl = ccp.certificateAuthorities["ca.issuer.com"].url;
        const ca = new FabricCaServices(caUrl);
        const walletPath = path.join(process.cwd(),'blockchain','wallets');
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        var res ={
            success : false,
            message: ""
        }

        const userExist = await wallet.get(req.body.CCCD);
        if (userExist) {
            
            console.log('An identity for the "user" already exists in the wallets')
            res.message = 'Nguoi dung da ton tai'
            return res
            
        }
        const adminIdentity = await wallet.get('admin');
        if (!adminIdentity) {
            console.log('An identity for the admin user "admin" not already exists in the wallets')
            res.message = 'An identity for the admin user "admin" not already exists in the wallets'
            return res
        }

        const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);

        const adminUser = await provider.getUserContext(adminIdentity, 'admin');
        var hashkey = crypto.createHash('md5').update(req.body.CCCD + Date.now()).digest('hex')
        const secret = await ca.register({
            affiliation: 'issuer.department',
            enrollmentID: req.body.CCCD,
            role: 'client',
            attrs: [{"name": "permission",
                    "value": hashkey, 
                    "ecert": true}],
        }, adminUser);

        const enrollment = await ca.enroll({
            enrollmentID: req.body.CCCD,
            enrollmentSecret: secret,
            attr_reqs: [{ name: "permission", optional: false }]
        });

        const X509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: 'issuer',
            type: 'X.509',
        };
        // cho vao vi
        await wallet.put(req.body.CCCD, X509Identity);

        
        // xong 1 buoc
        //cho vao mongodb
        var content = Buffer.from(req.file.buffer).toString("utf-8");
        
                var apiKey = new (APIKey)
                apiKey.username = req.body.username
                apiKey.CCCD = req.body.CCCD
                apiKey.hashKey = hashkey
                apiKey.information = content
                apiKey.expand = req.body.expand
                apiKey.save(function(err){
                    if(err) {
                        throw   err
                    }
                })
                res.success = true
                res.message = hashkey
                return res;
    } catch (err) {
        console.log(err)
        res.message = err.toString()
        return res
    }
}
