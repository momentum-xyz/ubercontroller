const out = {
    data: {},
    logs: [],
    error: null
}

let ApiModule
let UUIDModule
let UtilCryptoModule

try {
    ApiModule = require('@polkadot/api');
    UtilCryptoModule = require('@polkadot/util-crypto');
    UtilModule = require('@polkadot/util');
} catch (e) {
    exitWithError(33, e.toString())
}

const {ApiPromise, WsProvider, Keyring} = ApiModule
// const {decodeAddress} = UtilModule
const {decodeAddress, encodeAddress} = UtilCryptoModule


async function main() {

    if (process.argv.length !== 3) {
        exitWithError(11, `Provide target wallet as first cli argument`)
    }

    const TARGET_HEX_WALLET = process.argv[2]
    log(`target wallet: ${TARGET_HEX_WALLET}`)

    let WALLET
    try {
        WALLET = encodeAddress(TARGET_HEX_WALLET)
    } catch (e) {
        exitWithError(10, `Can not encode wallet (${TARGET_HEX_WALLET}) to ss58 format`)
    }
    log(`Encoded wallet=${WALLET}`)


    const url = "wss://drive.antst.net:19947"
    // const url = "wss://rpc.polkadot.io"

    log(`Using url: ${url}`)

    const wsProvider = new WsProvider(url);
    const api = await ApiPromise.create({provider: wsProvider});

    // To get details about collection by collection ID
    // const r = await api.query.uniques.class(0)

    // const r = await api.query.uniques.account(BOB, 0, null)
    // const r = await api.query.collatorSelection.lastAuthoredBlock(BOB)
    // const r = await api.query.collatorSelection.lastAuthoredBlock()

    // const BOB = "5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty"
    // const FERDIE = "5CiPPseXPECbkjWCa6MnjNokrgYjMqmKndv2rSnekmSK2DjL"
    const collectionId = 0
    // https://github.com/Phala-Network/khala-parachain/blob/bdd2daada48536e74800afd9d8bcbfb44b240bcb/scripts/js/pw/util/fetch.js#L54
    // const r = await api.query.uniques.account.entries(FERDIE, collectionId);
    const r = await api.query.uniques.account.entries(WALLET, collectionId);

    let itemId = null
    const i = r.map(([key, _value]) =>
        [key.args[0].toString(), key.args[1].toNumber(), key.args[2].toNumber()]
    )

    if (i.length === 0) {
        exitWithError(1, `UserID not found for wallet ${WALLET}`)
    }

    try {
        itemId = i[0][2]
    } catch (e) {
        exitWithError(1, `Can not get itemId from response`)
    }

    log(`ItemID=${itemId} collectionID=${collectionId}`)

    const r2 = await api.query.uniques.instanceMetadataOf(collectionId, itemId);
    // const r2 = await api.query.uniques.instanceMetadataOf(collectionId, 77227733);

    const meta = r2.toHuman()

    if (meta === null) {
        exitWithError(2, `No metadata for itemID=${itemId} collectionID=${collectionId} wallet=${WALLET}`)
    }

    let data

    try {
        data = JSON.parse(meta.data)
    } catch (e) {
        exitWithError(3, `Can not parse to JSON metadata for itemID=${itemId} collectionID=${collectionId} wallet=${WALLET}`)
    }

    out.data = data
    console.log(JSON.stringify(out))
    process.exit(0)
    await api.disconnect()

    // if (!Array.isArray(data)) {
    //     exitWithError(4, `Metadata is not array for itemID=${itemId} collectionID=${collectionId} wallet=${WALLET}`)
    // }
    //
    // if (!data[0]) {
    //     exitWithError(4, `Metadata array is empty for itemID=${itemId} collectionID=${collectionId} wallet=${WALLET}`)
    // }
    //
    // out.data.userUUID = data[0]
    // console.log(JSON.stringify(out))
    // process.exit(0)
    // await api.disconnect()
}

main()

function exitWithError(code, message) {
    out.error = message
    console.log(JSON.stringify(out))
    process.exit(code)
}

function log(m) {
    // console.log(m)
    out.logs.push(m)
}