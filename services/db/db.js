const mongoose = require('mongoose');


APIKey = mongoose.Schema({
    username: {
        type: String,
        required: true
    },

    CCCD: {
        type: String,
        required: true
    },
    hashKey: {
        type: String,
        required: true
    },
    information: {
        type: String,
        require: true
    },
    expand: {
        type: String,
        require: true
    }

});

module.exports = mongoose.model('APIKey', APIKey);