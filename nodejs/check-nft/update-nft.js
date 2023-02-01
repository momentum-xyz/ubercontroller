const out = {
    data: {},
    logs: [],
    error: null
}

let ApiModule
let UtilCryptoModule
let UtilModule


try {
    ApiModule = require('@polkadot/api');
    UtilCryptoModule = require('@polkadot/util-crypto');
    UtilModule = require('@polkadot/util');
} catch (e) {
    exitWithError(e.toString())
}

const {ApiPromise, WsProvider, Keyring} = ApiModule
const {decodeAddress, encodeAddress} = UtilCryptoModule

async function main() {


    if (process.argv.length !== 5) {
        const m = `Provide target wallet as first cli argument, admin mnemonic phrase as second, meta as third,`
        exitWithError(m)
    }

    const TARGET_WALLET = process.argv[2]
    log(`target wallet: ${TARGET_WALLET}`)

    const PHRASE = process.argv[3]
    log(`mnemonic phrase: ${PHRASE.substring(0, 3)} ... ${PHRASE.substring(PHRASE.length - 3)}; length: ${PHRASE.length}`)


    const META = process.argv[4]
    log(`NFT metadata: ${META}`)

    let m
    try {
        m = JSON.parse(META)
    } catch (e) {
        exitWithError(e.message)
    }

    const name = m.name
    const image = m.image

    if (!name || !image) {
        exitWithError(`META must contain name and image`)
    }

    log("disconnect")
    out.data = {userID: userID, name, image}
    console.log(JSON.stringify(out))
    //await api.disconnect()
}

function exitWithError(message) {
    out.error = message
    console.log(JSON.stringify(out))
    process.exit(0)
}

main()