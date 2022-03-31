const express = require('express')
const mongoose = require('mongoose')
const bodyParser = require('body-parser')
const helmet = require('helmet')
const morgan = require('morgan')
const cors = require('cors')
require('dotenv').config()
var app = express()
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({
    extended: false
}))
app.options('*', cors({
    credentials: true
}));
app.use(helmet.hidePoweredBy());
mongoose.connect(process.env.DB_ADDR, {
    useNewUrlParser: true,
    useUnifiedTopology: true
})

app.use(cors())

app.use(morgan('combined'))
require("./route")(app)
app.get("/", function (req, res) {
    console.log(req.body.username)
    return res.status(200).send({
        success: true,
        message: "ok"
    })
})
app.listen(10001, () => {
    console.log("server is listening on port 10001")
})

