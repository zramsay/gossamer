// Import
const { ApiPromise, WsProvider } = require('@polkadot/api');
const {Keyring } = require('@polkadot/keyring');

async function main() {
    // Construct

    const wsProvider = new WsProvider('ws://127.0.0.1:8546');
    // const wsProvider = new WsProvider('ws://127.0.0.1:9944');
    const api = await ApiPromise.create({ provider: wsProvider });

    // system calls
    // system_name
    const nodeName = await api.rpc.system.name();
    console.log(`name: ${nodeName}`);

    // system_version
    const nodeVersion = await api.rpc.system.version();
    console.log(`version: ${nodeVersion}`);

    // system_chain
    const chain = await api.rpc.system.chain();
    console.log(`chain: ${chain}`);

    // system_chainType
    const chainType = await api.rpc.system.chainType();
    console.log(`chainType: ${chainType}`);
    
    // system_properties
    const properties = await api.rpc.system.properties();
    console.log(`properties: ${properties}`);

    // system_health
    const health = await api.rpc.system.health();
    console.log(`health: ${health}`);

    // system_localPeerId
    const localPeerId = await api.rpc.system.localPeerId();
    console.log(`localPeerId: ${localPeerId}`);

    // system_localListenAddresses
    const localListenAddresses = await api.rpc.system.localListenAddresses();
    console.log(`localListenAddresses: ${localListenAddresses}`);

    // system_peers
    const peers = await api.rpc.system.peers();
    console.log(`peers: ${peers}`);

    // system_networkState
    const networkState = await api.rpc.system.networkState();
    console.log(`networkState: ${networkState}`);

    // system_addReservedPeer
    const reservedPeer = '/ip4/198.51.100.19/tcp/30333/p2p/QmSk5HQbn6LhUwDiNMseVUjuRYhEtYj4aUZ6WfWoGURpdV';
    const addReservedPeer = await api.rpc.system.addReservedPeer(reservedPeer);
    console.log(`addReservedPeer: ${addReservedPeer}`);

    // system_addReservedPeer
    const peerId = 'QmSk5HQbn6LhUwDiNMseVUjuRYhEtYj4aUZ6WfWoGURpdV';
    const removeReservedPeer = await api.rpc.system.removeReservedPeer(peerId);
    console.log(`removeReservedPeer: ${removeReservedPeer}`);

    // system_nodeRoles
    const nodeRoles = await api.rpc.system.nodeRoles();
    console.log(`nodeRoles: ${nodeRoles}`);

    // system_syncState
    const syncState = await api.rpc.system.syncState();
    console.log(`syncState: ${syncState}`);
    
    // system_accountNextIndex
    const ADDR_Alice = '5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY';
    // const aliceNonce = await api.rpc.system.accountNextIndex(ADDR_Alice);
    // console.log(`aliceNonce: ${aliceNonce}`);

    // system_dryRun
    const extrinsic = '0x010203';
    const dryRun = await api.rpc.system.dryRun(null);
    console.log(`dryRun: ${dryRun}`);

    // chain defaults
    const genesisHash = await api.genesisHash;
    console.log(`genesis hash: ${genesisHash}`);

    const runtimeMetadata = await  api.runtimeMetadata;
    // currently not sending runtimeMetadata to console because it's very large, uncomment if you want to see
    console.log(`runtime metadata: ${runtimeMetadata.metadata}`);

    const runtimeVersion = await api.runtimeVersion;
    console.log(`runtime version: ${runtimeVersion}`);

    const libraryInfo = await api.libraryInfo;
    console.log(`library info: ${libraryInfo}`);

    // //Basic queries
    // const now = await api.query.timestamp.now();
    // console.log(`timestamp now: ${now}`);

    // // Retrieve the account balance & nonce via the system module
    // const ADDR_Alice = '5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY';
    // const { nonce, data: balance } = await api.query.system.account(ADDR_Alice);
    // console.log(`Alice: balance of ${balance.free} and a nonce of ${nonce}`)

    // // RPC queries
    // const chain = await api.rpc.system.chain();
    // console.log(`system chain: ${chain}`);

    // const sysProperties = await api.rpc.system.properties();
    // console.log(`system properties: ${sysProperties}`);

    // const chainType = await api.rpc.system.chainType();
    // console.log(`system chainType: ${chainType}`);

    // const header = await api.rpc.chain.getHeader();
    // console.log(`header ${header}`);

    // // Subscribe to the new headers
    // // TODO: Issue: chain.subscribeNewHeads is returning values twice for each result new head.
    // let count = 0;
    // const unsubHeads = await api.rpc.chain.subscribeNewHeads((lastHeader) => {
    //     console.log(`${chain}: last block #${lastHeader.number} has hash ${lastHeader.hash}`);
    //     if (++count === 5) {
    //         unsubHeads();
    //     }
    // });

    // const blockHash = await api.rpc.chain.getBlockHash();
    // console.log(`current blockhash ${blockHash}`);

    // const block = await api.rpc.chain.getBlock(blockHash);
    // console.log(`current block: ${block}`);

    // // Simple transaction
    // // TODO Issue:  This currently fails with error: RPC-CORE: submitExtrinsic(extrinsic: Extrinsic): Hash:: -32000: validator: (nil *modules.Extrinsic): null
    // const keyring = new Keyring({type: 'sr25519' });
    // const aliceKey = keyring.addFromUri('//Alice',  { name: 'Alice default' });
    // console.log(`${aliceKey.meta.name}: has address ${aliceKey.address} with publicKey [${aliceKey.publicKey}]`);

    // const ADDR_Bob = '0x90b5ab205c6974c9ea841be688864633dc9ca8a357843eeacf2314649965fe22';

    // const transfer = await api.tx.balances.transfer(ADDR_Bob, 12345)
    //     .signAndSend(aliceKey);

    // console.log(`hxHash ${transfer}`);

}


main().catch(console.error);
