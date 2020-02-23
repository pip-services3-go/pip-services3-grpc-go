package test_services

// let assert = require('chai').assert;
// let grpc = require('grpc');
// var protoLoader = require('@grpc/proto-loader');
// let async = require('async');

// let services = require('../../../src/protos/commandable_grpc_pb');
// let messages = require('../../../src/protos/commandable_pb');

// import { Descriptor } from 'pip-services3-commons-node';
// import { ConfigParams } from 'pip-services3-commons-node';
// import { References } from 'pip-services3-commons-node';

// import { Dummy } from '../Dummy';
// import { DummyController } from '../DummyController';
// import { DummyCommandableGrpcService } from './DummyCommandableGrpcService';

// var grpcConfig = ConfigParams.fromTuples(
//     "connection.protocol", "http",
//     "connection.host", "localhost",
//     "connection.port", 3001
// );

// suite('DummyCommandableGrpcService', ()=> {
//     var _dummy1: Dummy;
//     var _dummy2: Dummy;

//     let service: DummyCommandableGrpcService;

//     let client: any;

//     suiteSetup((done) => {
//         let ctrl = new DummyController();

//         service = new DummyCommandableGrpcService();
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
//             __dirname + "../../../../src/protos/commandable.proto",
//             {
//                 keepCase: true,
//                 // longs: String,
//                 // enums: String,
//                 defaults: true,
//                 oneofs: true
//             }
//         );
//         let clientProto = grpc.loadPackageDefinition(packageDefinition).commandable.Commandable;

//         client = new clientProto('localhost:3001', grpc.credentials.createInsecure());

//         _dummy1 = { id: null, key: "Key 1", content: "Content 1"};
//         _dummy2 = { id: null, key: "Key 2", content: "Content 2"};
//     });

//     test('CRUD Operations', (done) => {
//         var dummy1, dummy2;

//         async.series([
//         // Create one dummy
//             (callback) => {
//                 client.invoke(
//                     {
//                         method: 'dummy.create_dummy',
//                         args_empty: false,
//                         args_json: JSON.stringify({
//                             dummy: _dummy1
//                         })
//                     },
//                     (err, response) => {
//                         assert.isNull(err);

//                         assert.isFalse(response.result_empty);
//                         assert.isString(response.result_json);
//                         let dummy = JSON.parse(response.result_json);

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
//                 client.invoke(
//                     {
//                         method: 'dummy.create_dummy',
//                         args_empty: false,
//                         args_json: JSON.stringify({
//                             dummy: _dummy2
//                         })
//                     },
//                     (err, response) => {
//                         assert.isNull(err);

//                         assert.isFalse(response.result_empty);
//                         assert.isString(response.result_json);
//                         let dummy = JSON.parse(response.result_json);

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
//                 client.invoke(
//                     {
//                         method: 'dummy.get_dummies',
//                         args_empty: false,
//                         args_json: JSON.stringify({})
//                     },
//                     (err, response) => {
//                         assert.isNull(err);

//                         assert.isFalse(response.result_empty);
//                         assert.isString(response.result_json);
//                         let dummies = JSON.parse(response.result_json);

//                         assert.isObject(dummies);
//                         assert.lengthOf(dummies.data, 2);

//                         callback();
//                     }
//                 );
//             },
//         // Update the dummy
//             (callback) => {
//                 dummy1.content = 'Updated Content 1';

//                 client.invoke(
//                     {
//                         method: 'dummy.update_dummy',
//                         args_empty: false,
//                         args_json: JSON.stringify({
//                             dummy: dummy1
//                         })
//                     },
//                     (err, response) => {
//                         assert.isNull(err);

//                         assert.isFalse(response.result_empty);
//                         assert.isString(response.result_json);
//                         let dummy = JSON.parse(response.result_json);

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
//                 client.invoke(
//                     {
//                         method: 'dummy.delete_dummy',
//                         args_empty: false,
//                         args_json: JSON.stringify({
//                             dummy_id: dummy1.id
//                         })
//                     },
//                     (err, response) => {
//                         assert.isNull(err);

//                         assert.isNull(response.error);

//                         callback();
//                     }
//                 );
//             },
//         // Try to get delete dummy
//             (callback) => {
//                 client.invoke(
//                     {
//                         method: 'dummy.get_dummy_by_id',
//                         args_empty: false,
//                         args_json: JSON.stringify({
//                             dummy_id: dummy1.id
//                         })
//                     },
//                     (err, response) => {
//                         assert.isNull(err);

//                         assert.isNull(response.error);
//                         assert.isTrue(response.result_empty);

//                         // assert.isObject(dummy);

//                         callback();
//                     }
//                 );
//             }
//         ], done);
//     });

// });
