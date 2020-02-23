package test_clients

import { FilterParams } from 'pip-services3-commons-node';
import { PagingParams } from 'pip-services3-commons-node';
import { DataPage } from 'pip-services3-commons-node';

import { GrpcClient } from '../../src/clients/GrpcClient';
import { IDummyClient } from './IDummyClient';
import { Dummy } from '../Dummy';

// export class DummyGrpcClient extends GrpcClient implements IDummyClient {
        
//     public constructor() {
//         super(__dirname + "../../../../test/protos/dummies.proto", "dummies.Dummies")
//     }

//     public getDummies(correlationId: string, filter: FilterParams, paging: PagingParams,
//         callback: (err: any, result: DataPage<Dummy>) => void): void {
//         this.call('get_dummies',
//             correlationId, 
//             { 
//                 filter: filter,
//                 paging: paging
//             },
//             (err, result) => {
//                 this.instrument(correlationId, 'dummy.get_page_by_filter');
//                 callback(err, result);
//             }
//         );
//     }

//     public getDummyById(correlationId: string, dummyId: string,
//         callback: (err: any, result: Dummy) => void): void {
//         this.call('get_dummy_by_id',
//             correlationId,
//             {
//                 dummy_id: dummyId
//             }, 
//             (err, result) => {
//                 this.instrument(correlationId, 'dummy.get_one_by_id');

//                 if (result && result.id == "" && result.key == "")
//                     result = null;

//                 callback(err, result);
//             }
//         );        
//     }

//     public createDummy(correlationId: string, dummy: any,
//         callback: (err: any, result: Dummy) => void): void {
//         this.call('create_dummy',
//             correlationId,
//             {
//                 dummy: dummy
//             }, 
//             (err, result) => {
//                 this.instrument(correlationId, 'dummy.create');

//                 if (result && result.id == "" && result.key == "")
//                     result = null;

//                 callback(err, result);
//             }
//         );
//     }

//     public updateDummy(correlationId: string, dummy: any,
//         callback: (err: any, result: Dummy) => void): void {
//         this.call('update_dummy',
//             correlationId, 
//             {
//                 dummy: dummy
//             }, 
//             (err, result) => {
//                 this.instrument(correlationId, 'dummy.update');

//                 if (result && result.id == "" && result.key == "")
//                     result = null;

//                 callback(err, result);
//             }
//         );
//     }

//     public deleteDummy(correlationId: string, dummyId: string,
//         callback: (err: any, result: Dummy) => void): void {
//         this.call('delete_dummy_by_id',
//             correlationId, 
//             {
//                 dummy_id: dummyId
//             }, 
//             (err, result) => {
//                 this.instrument(correlationId, 'dummy.delete_by_id');

//                 if (result && result.id == "" && result.key == "")
//                     result = null;

//                 callback(err, result);
//             }
//         );
//     }
  
// }