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
    console.log('\x1b[32m%s\x1b[0m %s', 'name:',  nodeName);

    // system_version
    const nodeVersion = await api.rpc.system.version();
    console.log('\x1b[32m%s\x1b[0m %s', 'version:',  nodeVersion);

    // system_chain
    const chain = await api.rpc.system.chain();
    console.log('\x1b[32m%s\x1b[0m %s', 'chain:',  chain);

    // system_chainType
    const chainType = await api.rpc.system.chainType();
    console.log('\x1b[32m%s\x1b[0m %s', 'chainType:',  chainType);
    
    // system_properties
    const properties = await api.rpc.system.properties();
    console.log('\x1b[32m%s\x1b[0m %s', 'properties:',  properties);

    // system_health
    const health = await api.rpc.system.health();
    console.log('\x1b[32m%s\x1b[0m %s', 'health:',  health);

    // system_localPeerId
    const localPeerId = await api.rpc.system.localPeerId();
    console.log('\x1b[32m%s\x1b[0m %s', 'localPeerId:',  localPeerId);

    // system_localListenAddresses
    const localListenAddresses = await api.rpc.system.localListenAddresses();
    console.log('\x1b[32m%s\x1b[0m %s', 'localListenAddresses:',  localListenAddresses);

    // system_peers
    const peers = await api.rpc.system.peers();
    console.log('\x1b[32m%s\x1b[0m %s', 'peers:',  peers);

    // system_networkState
    const networkState = await api.rpc.system.networkState();
    console.log('\x1b[32m%s\x1b[0m %s', 'networkState:',  networkState);

    // system_addReservedPeer
    const reservedPeer = '/ip4/198.51.100.19/tcp/30333/p2p/QmSk5HQbn6LhUwDiNMseVUjuRYhEtYj4aUZ6WfWoGURpdV';
    const addReservedPeer = await api.rpc.system.addReservedPeer(reservedPeer);
    console.log('\x1b[32m%s\x1b[0m %s', 'addReservedPeer:',  addReservedPeer);

    // system_addReservedPeer
    const peerId = 'QmSk5HQbn6LhUwDiNMseVUjuRYhEtYj4aUZ6WfWoGURpdV';
    const removeReservedPeer = await api.rpc.system.removeReservedPeer(peerId);
    console.log('\x1b[32m%s\x1b[0m %s', 'removeReservedPeer:',  removeReservedPeer);

    // system_nodeRoles
    const nodeRoles = await api.rpc.system.nodeRoles();
    console.log('\x1b[32m%s\x1b[0m %s', 'nodeRoles:',  nodeRoles);

    // system_syncState
    const syncState = await api.rpc.system.syncState();
    console.log('\x1b[32m%s\x1b[0m %s', 'syncState:',  syncState);
    
    // system_accountNextIndex
    const ADDR_Alice = '5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY';
    // const aliceNonce = await api.rpc.system.accountNextIndex(ADDR_Alice);
    // console.log(`aliceNonce: ${aliceNonce}`);

    // system_dryRun
    // const extrinsic = '0x010203';
    // const dryRun = await api.rpc.system.dryRun(null);
    // console.log(`dryRun: ${dryRun}`);

    // Babe
    // babe_epochAuthorship
    // const epochAuthorship = api.rpc.babe.epochAuthorship();
    // console.log(`epochAuthorship: ${epochAuthorship}`);

    // Grandpa
    // grandpa_roundState
    const roundState = await api.rpc.grandpa.roundState();
    console.log('\x1b[32m%s\x1b[0m %s', 'roundState:',  roundState);

    // // grandpa_proveFinality
    // const proveBlockNumber = 10;
    // const proveFinality = await api.rpc.grandpa.proveFinality(proveBlockNumber);
    // console.log(`proveFinality: ${proveFinality}`);
   
    // // grandpa_subscribeJustifications
    // const unsubJustifications = await api.rpc.grandpa.subscribeJustifications((justification) => {
    //     console.log(`justification: ${justification}`);
    // });

    // // grandpa_unsubscribeJustifications
    // const justUnsub = await api.rpc.grandpa.unsubscribeJustifications("f00");
    // console.log(`unsubscribe Justifications: ${justUnsub}`);

    // // Author
    // // author_submitExtrinsic
    const testExt = api.tx.system.remark([0, 1, 2]);
    console.log('\x1b[32m%s\x1b[0m %s', 'testExt hash:',  testExt.hash);
    // // const testExt = api.tx.system.setStorage([[0, 1, 2], [3, 4, 5]]);
    // console.log(`testExt: ${testExt}`);
    // const submitExtrinsic = await api.rpc.author.submitExtrinsic(testExt);
    // console.log(`submitExtrinsic: ${submitExtrinsic}`);

    // author_pendingExtrinsics
    const pendingExtrinsics = await api.rpc.author.pendingExtrinsics();
    console.log('\x1b[32m%s\x1b[0m %s', 'pendingExtrinsics:',  pendingExtrinsics);

    // // author_removeExtrinsic
    // const removeExtrinsic = await api.rpc.author.removeExtrinsic(testExt.hash);
    // console.log(`removeExtrinsic: ${removeExtrinsic}`);

    // // author_insertKey
    // // todo: these params cause error:  -32000: could not byteify non 0x prefixed string: 
    // const insertKey = await api.rpc.author.insertKey("dumy", "0x3fb882f70b4ddf5f8923f4a2d3b30a20f79bc0c5de212c1a8977f4972272db8d", "0x5ebf69cfbb4914711f70ff3b9e7455f7d5006b15f220d011387038cf4fb1593e");
    // console.log(`insertKey: ${insertKey}`);

    // author_rotateKeys
    const rotateKeys = await api.rpc.author.rotateKeys();
    console.log('\x1b[32m%s\x1b[0m %s', 'rotateKeys:',  rotateKeys);

    //author_hasSessionKeys
    const sessionKeys = await api.rpc.author.hasSessionKeys(ADDR_Alice)
    console.log('\x1b[32m%s\x1b[0m %s', 'sessionKeys:',  sessionKeys);

    // //author_hasKey
    // const keyring = new Keyring({type: 'sr25519'});
    // const AliceKey = keyring.createFromUri('//Alice');
    // console.log(`AliceKey: ${AliceKey.publicKey}`);
    // const hasKey = await api.rpc.author.hasKey(AliceKey.publicKey, "babe");
    // console.log(`hasKey: ${hasKey}`);

    // // author_submitAndWatchExtrinsic
    // const unsubWatchExtrinsic = await api.rpc.author.submitAndWatchExtrinsic(testExt);
    // console.log(`unsubWatchExtrinsic: ${unsubWatchExtrinsic}`);

    // // author_unwatchExtrinsic
    // const unwatchExtrinsic = await api.rpc.author.unwatchExtrinsic("f00")
    // console.log(`unwatchExtrinsic: ${unwatchExtrinsic}`);

    // Chain
    // chain_getHeader
    const chainGetHeader = await api.rpc.chain.getHeader();
    console.log('\x1b[32m%s\x1b[0m %s', 'chainGetHeader:',  chainGetHeader);

    // chain_getBlock
    const chainGetBlock = await api.rpc.chain.getBlock();
    console.log('\x1b[32m%s\x1b[0m %s', 'chainGetBlock:',  chainGetBlock);

    // chain_getBlockHash
    const chainGetBlockHash = await api.rpc.chain.getBlockHash();
    console.log('\x1b[32m%s\x1b[0m %s', 'chainGetBlockHash:',  chainGetBlockHash);
    
    // chain_getFinalizedHead
    const chainGetFinalizedHead = await api.rpc.chain.getFinalizedHead();
    console.log('\x1b[32m%s\x1b[0m %s', 'chainGetFinalizedHead:',  chainGetFinalizedHead);

    // // chain_subscribeAllHeads
    // let count = 0;
    // const unsubHeads = await api.rpc.chain.subscribeAllHeads((lastHeader) => {
    //     console.log(`last block #${lastHeader.number} has hash ${lastHeader.hash}`);
    //     if (++count === 5) {
    //         unsubHeads();
    //     }
    // });

    // // chain_unsubscribeAllHeads
    // const unsubAllHeads = await api.rpc.chain.unsubcribeAllHeads();
    // console.log(`unsubAllHeads: ${unsubAllHeads}`);

    // chain_subscribeNewHeads/chain_unsubscribeNewHeads
    let count = 0;
    const unsubNewHeads = await api.rpc.chain.subscribeNewHeads((lastHeader) => {
        console.log('\x1b[32m%s\x1b[0m %s', 'last block:', `${lastHeader.number} has hash ${lastHeader.hash}`);
        if (++count === 5) {
            unsubNewHeads();
       }
    });

    // chain_subscribeFinalizedHeads/chain_unsubscribeFinalizedHeads
    let countFinalized = 0;
    const unsubFinalizedHeads = await api.rpc.chain.subscribeFinalizedHeads((lastHeader) => {
        console.log('\x1b[32m%s\x1b[0m %s', 'last finalized block:',  `${lastHeader.number} has hash ${lastHeader.hash}`);
        if (++countFinalized === 5) {
            unsubFinalizedHeads();
       }
    });

    // Offchain
    // offchain_localStorageSet/Get LOCAL
    await api.rpc.offchain.localStorageSet("LOCAL", "0x010203", "0x040506")
    
    const localStorageGet = await api.rpc.offchain.localStorageGet("LOCAL", "0x010203");
    console.log('\x1b[32m%s\x1b[0m %s', 'localStorageGet:',  localStorageGet);

    // offchain_localStorageSet/Get PERSISTENT
    await api.rpc.offchain.localStorageSet("PERSISTENT", "0x010101", "0x040404")
    
    const localStorageGetPersistent = await api.rpc.offchain.localStorageGet("PERSISTENT", "0x010101");
    console.log('\x1b[32m%s\x1b[0m %s', 'localStorageGetPersistent:',  localStorageGetPersistent);

    // State
    // // state_call
    // const stateCall = await api.rpc.state.call("test", "data");
    // console.log(`state_call: ${stateCall}`);

    // state_getPairs
    // prefix returns one result
    const stateGetPairs = await api.rpc.state.getPairs("0x26aa394eea5630e07c48ae0c9558cef7a44704b568d21667356a5a050c118746e333f8c357e331db45010000");
    // prefix returns multiple results
    // const stateGetPairs = await api.rpc.state.getPairs("0x26aa394eea5630e07c48ae0c9558cef7a44704b568d21667356a5a050c118746e3");
    console.log('\x1b[32m%s\x1b[0m %s', 'stateGetPairs:',  stateGetPairs);

    // state_getKeysPaged
    const keysPaged = await api.rpc.state.getKeysPaged(null, 2);
    console.log('\x1b[32m%s\x1b[0m %s', 'keysPaged:',  keysPaged);

    // state_getStorage
    const getStorage = await api.rpc.state.getStorage("0x26aa394eea5630e07c48ae0c9558cef7a44704b568d21667356a5a050c118746e333f8c357e331db45010000");
    console.log('\x1b[32m%s\x1b[0m %s', "getStorage:", getStorage)

    // state_getStorageHash
    const storageHash = await api.rpc.state.getStorageHash("0x26aa394eea5630e07c48ae0c9558cef7a44704b568d21667356a5a050c118746e333f8c357e331db45010000");
    console.log('\x1b[32m%s\x1b[0m %s', "storageHash:", storageHash);

    // state_getStorageSize
    const storageSize = await api.rpc.state.getStorageSize("0x26aa394eea5630e07c48ae0c9558cef7a44704b568d21667356a5a050c118746e333f8c357e331db45010000");
    console.log('\x1b[32m%s\x1b[0m %s', "storageSize:", storageSize);

    // state_getMetadata
    const stateMetadata = await api.rpc.state.getMetadata();
    console.log('\x1b[32m%s\x1b[0m %s', "stateMetadata size:", stateMetadata.encodedLength);

    // state_getRuntimeVersion
    const runtimeVersion = await api.rpc.state.getRuntimeVersion();
    console.log('\x1b[32m%s\x1b[0m %s', "runtimeVersion:", runtimeVersion);
    

    // chain defaults
    const genesisHash = await api.genesisHash;
    console.log(`genesis hash: ${genesisHash}`);

    const runtimeMetadata = await  api.runtimeMetadata;
    // currently not sending runtimeMetadata to console because it's very large, uncomment if you want to see
    console.log(`runtime metadata: ${runtimeMetadata.metadata}`);

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
