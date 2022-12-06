const out = {
    data: {},
    logs: [],
    error: null
}

let ApiModule
let UtilCryptoModule

try {
    ApiModule = require('@polkadot/api');
    UtilCryptoModule = require('@polkadot/util-crypto');
} catch (e) {
    exitWithError(e.toString())
}

const {ApiPromise, WsProvider, Keyring} = ApiModule
const {decodeAddress, encodeAddress} = UtilCryptoModule


async function main() {

    if (process.argv.length !== 3) {
        exitWithError(`Provide target wallet as first cli argument`)
    }

    const TARGET_HEX_WALLET = process.argv[2]
    log(`target wallet: ${TARGET_HEX_WALLET}`)

    let WALLET
    try {
        WALLET = encodeAddress(TARGET_HEX_WALLET)
    } catch (e) {
        exitWithError(`Can not encode wallet (${TARGET_HEX_WALLET}) to ss58 format`)
    }
    log(`Encoded wallet=${WALLET}`)


    const url = "wss://drive.antst.net:19947"
    // const url = "wss://rpc.polkadot.io"

    log(`Using url: ${url}`)

    const wsProvider = new WsProvider(url);
    const api = await ApiPromise.create({provider: wsProvider});

    const collectionId = 0
    // https://github.com/Phala-Network/khala-parachain/blob/bdd2daada48536e74800afd9d8bcbfb44b240bcb/scripts/js/pw/util/fetch.js#L54
    const r = await api.query.uniques.account.entries(WALLET, collectionId);

    let itemId = null
    const i = r.map(([key, _value]) =>
        [key.args[0].toString(), key.args[1].toNumber(), key.args[2].toNumber()]
    )

    if (i.length === 0) {
        exitWithError(`UserID not found for wallet ${WALLET}`)
    }

    try {
        itemId = i[0][2]
    } catch (e) {
        exitWithError(`Can not get itemId from response`)
    }

    log(`ItemID=${itemId} collectionID=${collectionId}`)

    const r2 = await api.query.uniques.instanceMetadataOf(collectionId, itemId);

    const meta = r2.toHuman()

    if (meta === null) {
        exitWithError(`No metadata for itemID=${itemId} collectionID=${collectionId} wallet=${WALLET}`)
    }

    let data

    try {
        data = JSON.parse(meta.data)
    } catch (e) {
        exitWithError(`Can not parse to JSON metadata for itemID=${itemId} collectionID=${collectionId} wallet=${WALLET}`)
    }

    out.data = data
    console.log(JSON.stringify(out))
    process.exit(0)
    await api.disconnect()

    // if (!Array.isArray(data)) {
    //     exitWithError(`Metadata is not array for itemID=${itemId} collectionID=${collectionId} wallet=${WALLET}`)
    // }
    //
    // if (!data[0]) {
    //     exitWithError(`Metadata array is empty for itemID=${itemId} collectionID=${collectionId} wallet=${WALLET}`)
    // }
    //
    // out.data.userUUID = data[0]
    // console.log(JSON.stringify(out))
    // process.exit(0)
    // await api.disconnect()
}

main()

function exitWithError(message) {
    out.error = message
    console.log(JSON.stringify(out))
    process.exit(0)
}

function log(m) {
    // console.log(m)
    out.logs.push(m)
}