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

const collectionId = 3

async function main() {

    if (process.argv.length !== 6) {
        const m = `Provide target wallet as first cli argument, admin mnemonic phrase as second, meta as third, userID as fourth`
        exitWithError(m)
    }

    const WALLET = process.argv[2]
    log(`wallet: ${WALLET}`)

    const PHRASE = process.argv[3]
    log(`mnemonic phrase: ${PHRASE.substring(0, 3)} ... ${PHRASE.substring(PHRASE.length - 3)}; length: ${PHRASE.length}`)

    const META = process.argv[4]
    log(`NFT metadata: ${META}`)

    const USER_ID = process.argv[5]
    log(`UserID: ${USER_ID}`)

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


    const url = "wss://drive.antst.net:19947"
    // const url = "wss://rpc.polkadot.io"

    log(`Using url: ${url}`)

    const wsProvider = new WsProvider(url);
    const api = await ApiPromise.create({provider: wsProvider});

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

    const keyring = new Keyring({type: 'sr25519'});
    const newPair = keyring.addFromUri(PHRASE);

    log(newPair.address)

    await setMeta(api, itemId, USER_ID, name, image, newPair)
    out.data = {userID: USER_ID, name, image}
    log("disconnect")
    console.log(JSON.stringify(out))
    await api.disconnect()
}

async function setMeta(api, item_id, user_id, name, image, admin_pair) {

    const isFrozen = false
    const meta = [user_id, name, 0, image]

    return new Promise((resolve, reject) => {
        api.tx.uniques
            .setMetadata(collectionId, item_id, JSON.stringify(meta), isFrozen)
            .signAndSend(admin_pair, (result) => {
                log(`Current status is ${result.status}`);

                if (result.status.isInBlock) {
                    log(`Transaction included at blockHash ${result.status.asInBlock}`);
                } else if (result.status.isFinalized) {
                    log(`Transaction finalized at blockHash ${result.status.asFinalized}`);
                    // unsub();
                    resolve(user_id)
                }
            });
    })
}

function exitWithError(message) {
    out.error = message
    console.log(JSON.stringify(out))
    process.exit(0)
}

function log(m) {
    // console.log(m)
    out.logs.push(m)
}

main()