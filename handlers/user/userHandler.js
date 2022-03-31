const users = require('../../services/userServices')

exports.signIn = async function(req, res){

}

exports.signUp = async function(req, res){
    try {
        var message = await users.RegisterUsers(req)
        console.log(message)
        if (!message.success){
            console.log(message)
            return res.status(400).json({
                success: false,
                result: message.message
            })
        }
        else{
            return res.status(200).json({
                success: true,
                result: message.message
            })
        }
       

    }
    catch(err){
        console.log(err)
        return res.status(500).send({
            success: false,
            errors: err
        })

    }
}