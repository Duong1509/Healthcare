const {BlockDecoder} = require('fabric-common');
const fabprotos = require("fabric-protos");

const _validation_codes = {};
for (const key in fabprotos.protos.TxValidationCode) {
	const new_key = fabprotos.protos.TxValidationCode[key];
	_validation_codes[new_key] = key;
}

function convertValidationCode(code) {
	if (typeof code === 'string') {
		return code;
	}
	return _validation_codes[code];
}

function decodeBlock(blockRaw){
    var data= {}

    const block = BlockDecoder.decode(blockRaw);
    data.header = {
        number : block.header.number,
        previous_hash : block.header.previous_hash.toString('base64'),
        data_hash : block.header.data_hash.toString('base64')
    }
    data.payload = []
    for(var index = 0; index < block.data.data.length; index++)
    {
        data.payload[index] = {
            channel_header : block.data.data[index].payload.header.channel_header,
            transaction:{
                is_delete : block.data.data[index].payload.data.actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset[1].rwset.writes[0].is_delete,
                data : block.data.data[index].payload.data.actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset[1].rwset.writes[0].value.toString(),
                creator : block.data.data[index].payload.data.actions[0].header.creator.mspid
            },
            chaincode_spec : block.data.data[index].payload.data.actions[0].payload.chaincode_proposal_payload.input.chaincode_spec,
        }
        data.payload[index].channel_header.extension = block.data.data[index].payload.header.channel_header.extension.toString().replace(/[^\w\s]/gi, '')

        for(var i = 0; i < block.data.data[index].payload.data.actions[0].payload.chaincode_proposal_payload.input.chaincode_spec.input.args.length; i++){
            data.payload[index].chaincode_spec.input.args[i] = block.data.data[index].payload.data.actions[0].payload.chaincode_proposal_payload.input.chaincode_spec.input.args[i].toString()
        }
    }
    console.log(data.payload[0].channel_header.extension)
    return data 
}


function decodeTransaction(blockRaw, TxID){
    var data= {}

    const block = BlockDecoder.decode(blockRaw);
    data = {
        number : block.header.number.toString(),
        txid : TxID,
    }
    for(var index = 0; index < block.data.data.length; index++){
        if(block.data.data[index].payload.header.channel_header.tx_id===TxID){
            data.chaincodename = block.data.data[index].payload.data.actions[0].payload.chaincode_proposal_payload.input.chaincode_spec.chaincode_id.name;
            data.status = block.data.data[index].payload.data.actions[0].payload.action.proposal_response_payload.extension.response.status;
            data.creator_msp_id = block.data.data[index].payload.header.signature_header.creator.mspid;
            data.type = block.data.data[index].payload.header.channel_header.typeString;

            data.proposalhash = block.data.data[index].payload.data.actions[0].payload.action.proposal_response_payload.proposal_hash.toString('hex');

            rwset = block.data.data[index].payload.data.actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset;
            data.read_set = rwset.map(i => ({
                chaincode: i.namespace,
                set: i.rwset.reads
            }));
            data.write_set = rwset.map(i => ({
                chaincode: i.namespace,
                set: i.rwset.writes
            }));
            validation_codes = block.metadata.metadata[
                fabprotos.common.BlockMetadataIndex.TRANSACTIONS_FILTER
            ];
            
            data.validation_code = convertValidationCode(validation_codes[index]);
            break;
        }

    }
    data.write_set[1].set[0].value = data.write_set[1].set[0].value.toString();
    return data 
}

exports.decodeBlock = decodeBlock;
exports.decodeTransaction = decodeTransaction;