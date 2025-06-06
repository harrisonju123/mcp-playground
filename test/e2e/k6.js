import grpc from 'k6/net/grpc';
import { check, fail } from 'k6';

const client = new grpc.Client();
client.load([], 'aggregator.proto');
let connectedPerVU = {};

export default () => {
    if (!connectedPerVU[__VU]) {
        client.connect('localhost:50051', { plaintext: true });
        connectedPerVU[__VU] = true;
    }

    const res = client.invoke('aggregator.Aggregator/CallTool', {
        name: 'echo', argsJson: { msg: 'after-failover' },
    });

    if (!check(res, { 'status is OK': (r) => r && r.status === grpc.StatusOK })) {
        console.error(`✗ RPC response after failover: ${JSON.stringify(res)}`);
        fail('Post-failover check did not return OK');
    }
};