package test_clients

// import { Descriptor } from 'pip-services3-commons-node';
// import { ConfigParams } from 'pip-services3-commons-node';
// import { References } from 'pip-services3-commons-node';

// import { DummyController } from '../DummyController';
// import { DummyGrpcService } from '../services/DummyGrpcService';
// import { DummyGrpcClient } from './DummyGrpcClient';
// import { DummyClientFixture } from './DummyClientFixture';

// var grpcConfig = ConfigParams.fromTuples(
//     "connection.protocol", "http",
//     "connection.host", "localhost",
//     "connection.port", 3000
// );

// suite('DummyGrpcClient', ()=> {
//     let service: DummyGrpcService;
//     let client: DummyGrpcClient;
//     let fixture: DummyClientFixture;

//     suiteSetup((done) => {
//         let ctrl = new DummyController();

//         service = new DummyGrpcService();
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

//     setup((done) => {
//         client = new DummyGrpcClient();
//         fixture = new DummyClientFixture(client);

//         client.configure(grpcConfig);
//         client.setReferences(new References());
//         client.open(null, done);
//     });

//     test('CRUD Operations', (done) => {
//         fixture.testCrudOperations(done);
//     });

// });
