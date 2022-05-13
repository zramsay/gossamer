// Import
const { ApiPromise, WsProvider } = require('@polkadot/api');
const {Keyring } = require('@polkadot/keyring');

async function main() {
    const wsProvider = new WsProvider('ws://127.0.0.1:8546');
    const api = await ApiPromise.create({ provider: wsProvider });

    // system_version
    const nodeVersion = await api.rpc.system.version();
    console.log('\x1b[32m%s\x1b[0m %s', 'version:',  nodeVersion);

    // state_getRuntimeVersion
    const runtimeVersion = await api.rpc.state.getRuntimeVersion();
    console.log('\x1b[32m%s\x1b[0m %s', "runtimeVersion:", runtimeVersion);
    
    // wait for block 1
    const unsub = await api.rpc.chain.subscribeNewHeads(async (lastHeader) => {
        console.log(`latest block: #${lastHeader.number} `);
        if (lastHeader.number > 1) {
            unsub()
        }
    });

    // Simple transaction
    const keyring = new Keyring({type: 'sr25519' });
    const aliceKey = keyring.addFromUri('//Alice',  { name: 'Alice default' });
    console.log(`${aliceKey.meta.name}: has address ${aliceKey.address} with publicKey [${aliceKey.publicKey}]`);

    const bobKey = keyring.addFromUri('//Bob', {name: 'Bob default'});
    console.log(`${bobKey.meta.name}: has address ${bobKey.address} with publicKey [${bobKey.publicKey}], ${toHexString(bobKey.publicKey)}`);

    const ADDR_Bob = '0x90b5ab205c6974c9ea841be688864633dc9ca8a357843eeacf2314649965fe22';
    // bob 5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty

    // const transfer = await api.tx.balances.transfer(bobKey.address, 12345).signAndSend(aliceKey);
// NOTE, this works with runtime spec version 264
    const transfer = await api.tx.balances.transfer(bobKey.address, 12345).paymentInfo(aliceKey);
    console.log(`transaction hash: ${transfer}`);
}

main().catch(console.error);

function toHexString(byteArray) {
    return Array.from(byteArray, function(byte) {
        return ('0' + (byte & 0xFF).toString(16)).slice(-2);
    }).join('')
}
