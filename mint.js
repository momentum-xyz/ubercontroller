const out = {
    data: {},
    logs: [],
    error: null
}

let ApiModule
let UUIDModule

try {
    ApiModule = require('/usr/local/lib/node_modules/@polkadot/api');
    UUIDModule = require('/usr/local/lib/node_modules/uuid');
} catch (e) {
    exitWithError(e.message)
}

const {ApiPromise, WsProvider, Keyring} = ApiModule
const {v4} = UUIDModule

const collectionId = 3

async function hasPayment(block_hash, from_wallet, amount, api, admin_pair) {
    const signedBlock = await api.rpc.chain.getBlock(block_hash);
    // console.log(signedBlock)
    // console.log(JSON.stringify(signedBlock))

    // the information for each of the contained extrinsics

    for (const ex of signedBlock.block.extrinsics) {
        const info = ex.toHuman()

        if (info.signer &&
            info.signer.Id &&
            info.signer.Id.toString &&
            info.signer.Id.toString() === from_wallet) {

            if (!info.method) {
                continue
            }

            const {args, method, section} = info.method

            if (section && section === "balances") {
                if (method && method.includes("transfer")) {
                    if (args && args.dest) {
                        if (args.dest.Id === admin_pair.address) {
                            return true
                        }
                    }
                }
            }
        }
    }

    return false
}

async function mint(api, owner_wallet, admin_pair) {

    const item = getRandom(1, 1_000_000_000)
    log("itemID=" + item)


    return new Promise((resolve, reject) => {
        api.tx.uniques
            .mint(collectionId, item, owner_wallet)
            .signAndSend(admin_pair, (result) => {
                log(`Current status is ${result.status}`);

                if (result.status.isInBlock) {
                    log(`Transaction included at blockHash ${result.status.asInBlock}`);
                } else if (result.status.isFinalized) {
                    log(`Transaction finalized at blockHash ${result.status.asFinalized}`);
                    // unsub();
                    resolve(item)
                }
            });
    })

}

async function setMeta(api, item_id, name, image, admin_pair) {

    const isFrozen = false
    const uuid = v4()

    // ["d83670c7-a120-47a4-892d-f9ec75604f74","Mitia",0,"https://picsum.photos/102"]
    const meta = [uuid, name, 0, image]

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
                    resolve(uuid)
                }
            });
    })
}

async function main() {

    if (process.argv.length !== 6) {
        const m = `Provide target wallet as first cli argument, admin mnemonic phrase as second, meta as third, block_hash as forth`
        exitWithError(m)
    }

    const TARGET_WALLET = process.argv[2]
    log(`target wallet: ${TARGET_WALLET}`)

    const PHRASE = process.argv[3]
    log(`mnemonic phrase: ${PHRASE.substring(0, 3)} ... ${PHRASE.substring(PHRASE.length - 3)}; length: ${PHRASE.length}`)


    const META = process.argv[4]
    log(`NFT metadata: ${META}`)

    const BLOCK_HASH = process.argv[5]
    log(`block hash: ${BLOCK_HASH}`)

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

    const keyring = new Keyring({type: 'sr25519'});
    const newPair = keyring.addFromUri(PHRASE);

    log(newPair.address)

    const amount = 1
    const flag = await hasPayment(BLOCK_HASH, TARGET_WALLET, amount, api, newPair)
    if (!flag) {
        exitWithError(`Payment to admin wallet from wallet=${TARGET_WALLET} not found in given block blockHash=${BLOCK_HASH}`)
    }

    const itemID = await mint(api, TARGET_WALLET, newPair)
    log(` Item minted: itemID=${itemID}`)

    const userID = await setMeta(api, itemID, name, image, newPair)

    log("disconnect")
    out.data = {userID: userID, name, image}
    console.log(JSON.stringify(out))
    await api.disconnect()

}

function log(m) {
    // console.log(m)
    out.logs.push(m)
}

function getRandom(min, max) {
    return Math.floor(Math.random() * (max - min) + min);
}

process.on('unhandledRejection', error => {
    exitWithError(`unhandledRejection: ${error.message}`)
});

function exitWithError(message) {
    out.error = message
    console.log(JSON.stringify(out))
    process.exit(0)
}

main()