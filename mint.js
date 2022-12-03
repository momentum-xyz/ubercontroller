const out = {
    data: {},
    logs: [],
    error: null
}

let ApiModule

try {
    ApiModule = require('/usr/local/lib/node_modules/@polkadot/api');
} catch (e) {
    out.error = e
    console.log(JSON.stringify(out))
    process.exit(3)
}

const {ApiPromise, WsProvider, Keyring} = ApiModule

async function main() {
    if (process.argv.length !== 4) {
        out.error = `Provide target wallet as first cli argument and admin mnemonic phrase as second`
        console.log(JSON.stringify(out))
        process.exit(1)
    }

    const TARGET_WALLET = process.argv[2]
    log(`target wallet: ${TARGET_WALLET}`)

    const PHRASE = process.argv[3]
    log(`mnemonic phrase: ${PHRASE}`)

    const url = "wss://drive.antst.net:19947"
    // const url = "wss://rpc.polkadot.io"

    log(`Using url: ${url}`)

    const wsProvider = new WsProvider(url);
    const api = await ApiPromise.create({provider: wsProvider});

    const keyring = new Keyring({type: 'sr25519'});
    const newPair = keyring.addFromUri(PHRASE);

    const unsub = await api.tx.balances
        .transfer(TARGET_WALLET, 1_000_000_000_000)
        .signAndSend(newPair, (result) => {
            log(`Current status is ${result.status}`)

            if (result.status.isInBlock) {
                log(`Transaction included at blockHash ${result.status.asInBlock}`)
            } else if (result.status.isFinalized) {
                log(`Transaction finalized at blockHash ${result.status.asFinalized}`);
                unsub();
                api.disconnect()
                console.log(JSON.stringify(out))
            }
        });
}

function log(m) {
    // console.log(m)
    out.logs.push(m)
}

process.on('unhandledRejection', error => {
    log(`unhandledRejection: ${error.message}`);
    out.error = `unhandledRejection: ${error.message}`
    console.log(JSON.stringify(out))
    process.exit(2)
});

main()