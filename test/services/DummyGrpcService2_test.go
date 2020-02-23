package test_services

// let assert = require('chai').assert;
// let grpc = require('grpc');
// var protoLoader = require('@grpc/proto-loader');
// let async = require('async');

// let services = require('../../../test/protos/dummies_grpc_pb');
// let messages = require('../../../test/protos/dummies_pb');

// import { Descriptor } from 'pip-services3-commons-node';
// import { ConfigParams } from 'pip-services3-commons-node';
// import { References } from 'pip-services3-commons-node';

// import { Dummy } from '../Dummy';
// import { DummyController } from '../DummyController';
// import { DummyGrpcService2 } from './DummyGrpcService2';

// var grpcConfig = ConfigParams.fromTuples(
//     "connection.protocol", "http",
//     "connection.host", "localhost",
//     "connection.port", 3000
// );

// suite('DummyGrpcService2', ()=> {
//     var _dummy1: Dummy;
//     var _dummy2: Dummy;

//     let service: DummyGrpcService2;

//     let client: any;

//     suiteSetup((done) => {
//         let ctrl = new DummyController();

//         service = new DummyGrpcService2();
//         service.configure(grpcConfig);

//         let references: References = References.fromTuples(
//             new Descriptor('pip-services-dummies', 'controller', 'default', 'default', '1.0'), ctrl,
//             new Descriptor('pip-services-dummies', 'service', 'grpc', 'default', '1.0'), service
//         );
//         service.setReferences(references);

//         service.open(null, done);
//     });

//     suiteTeardown((done) => {
//         service.close(null, done);
//     });

//     setup(() => {
//         let packageDefinition = protoLoader.loadSync(
//             __dirname + "../../../../test/protos/dummies.proto",
//             {
//                 keepCase: true,
//                 longs: String,
//                 enums: String,
//                 defaults: true,
//                 oneofs: true
//             }
//         );
//         let clientProto = grpc.loadPackageDefinition(packageDefinition).dummies.Dummies;

//         client = new clientProto('localhost:3000', grpc.credentials.createInsecure());

//         _dummy1 = { id: null, key: "Key 1", content: "Content 1"};
//         _dummy2 = { id: null, key: "Key 2", content: "Content 2"};
//     });

//     test('CRUD Operations', (done) => {
//         var dummy1, dummy2;

//         async.series([
//         // Create one dummy
//             (callback) => {
//                 client.create_dummy(
//                     { dummy: _dummy1 },
//                     (err, dummy) => {
//                         assert.isNull(err);

//                         assert.isObject(dummy);
//                         assert.equal(dummy.content, _dummy1.content);
//                         assert.equal(dummy.key, _dummy1.key);

//                         dummy1 = dummy;

//                         callback();
//                     }
//                 );
//             },
//         // Create another dummy
//             (callback) => {
//                 client.create_dummy(
//                     { dummy: _dummy2 },
//                     (err, dummy) => {
//                         assert.isNull(err);

//                         assert.isObject(dummy);
//                         assert.equal(dummy.content, _dummy2.content);
//                         assert.equal(dummy.key, _dummy2.key);

//                         dummy2 = dummy;

//                         callback();
//                     }
//                 );
//             },
//         // Get all dummies
//             (callback) => {
//                 client.get_dummies(
//                     {},
//                     (err, dummies) => {
//                         assert.isNull(err);

//                         assert.isObject(dummies);
//                         assert.lengthOf(dummies.data, 2);

//                         callback();
//                     }
//                 );
//             },
//         // Update the dummy
//             (callback) => {
//                 dummy1.content = 'Updated Content 1';

//                 client.update_dummy(
//                     { dummy: dummy1 },
//                     (err, dummy) => {
//                         assert.isNull(err);

//                         assert.isObject(dummy);
//                         assert.equal(dummy.content, 'Updated Content 1');
//                         assert.equal(dummy.key, _dummy1.key);

//                         dummy1 = dummy;

//                         callback();
//                     }
//                 );
//             },
//         // Delete dummy
//             (callback) => {
//                 client.delete_dummy_by_id(
//                     { dummy_id: dummy1.id },
//                     (err, dummy) => {
//                         assert.isNull(err);

//                         callback();
//                     }
//                 );
//             },
//         // Try to get delete dummy
//             (callback) => {
//                 client.get_dummy_by_id(
//                     { dummy_id: dummy1.id },
//                     (err, dummy) => {
//                         assert.isNull(err);

//                         // assert.isObject(dummy);

//                         callback();
//                     }
//                 );
//             }
//         ], done);
//     });

// });
