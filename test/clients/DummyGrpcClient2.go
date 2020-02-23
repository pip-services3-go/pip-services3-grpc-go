package test_clients

// let _ = require('lodash');
// let services = require('../../../test/protos/dummies_grpc_pb');
// let messages = require('../../../test/protos/dummies_pb');

// import { FilterParams } from 'pip-services3-commons-node';
// import { PagingParams } from 'pip-services3-commons-node';
// import { DataPage } from 'pip-services3-commons-node';

// import { GrpcClient } from '../../src/clients/GrpcClient';
// import { IDummyClient } from './IDummyClient';
// import { Dummy } from '../Dummy';

// export class DummyGrpcClient2 extends GrpcClient implements IDummyClient {

//     public constructor() {
//         super(services.DummiesClient)
//     }

//     public getDummies(correlationId: string, filter: FilterParams, paging: PagingParams,
//         callback: (err: any, result: DataPage<Dummy>) => void): void {

//         paging = paging || new PagingParams();
//         let pagingParams = new messages.PagingParams();
//         pagingParams.setSkip(paging.skip);
//         pagingParams.setTake(paging.take);
//         pagingParams.setTotal(paging.total);

//         let request = new messages.DummiesPageRequest();
//         request.setPaging(pagingParams);

//         filter = filter || new FilterParams();
//         let filterParams = request.getFilterMap();
//         for (var propName in filter) {
//             if (filter.hasOwnProperty(propName))
//                 filterParams[propName] = filter[propName];
//         }

//         this.call('get_dummies',
//             correlationId,
//             request,
//             (err, result) => {
//                 this.instrument(correlationId, 'dummy.get_page_by_filter');

//                 result = result != null ? result.toObject() : null;
//                 if (result) {
//                     result.data = result.dataList;
//                     delete result.dataList;
//                 }

//                 callback(err, result);
//             }
//         );
//     }

//     public getDummyById(correlationId: string, dummyId: string,
//         callback: (err: any, result: Dummy) => void): void {

//         let request = new messages.DummyIdRequest();
//         request.setDummyId(dummyId);

//         this.call('get_dummy_by_id',
//             correlationId,
//             request,
//             (err, result) => {
//                 this.instrument(correlationId, 'dummy.get_one_by_id');

//                 result = result != null ? result.toObject() : null;
//                 if (result && result.id == "" && result.key == "")
//                     result = null;

//                 callback(err, result);
//             }
//         );
//     }

//     public createDummy(correlationId: string, dummy: any,
//         callback: (err: any, result: Dummy) => void): void {

//         let dummyObj = new messages.Dummy();
//         dummyObj.setId(dummy.id);
//         dummyObj.setKey(dummy.key);
//         dummyObj.setContent(dummy.content);

//         let request = new messages.DummyObjectRequest();
//         request.setDummy(dummyObj);

//         this.call('create_dummy',
//             correlationId,
//             request,
//             (err, result) => {
//                 this.instrument(correlationId, 'dummy.create');

//                 result = result != null ? result.toObject() : null;
//                 if (result && result.id == "" && result.key == "")
//                     result = null;

//                 callback(err, result);
//             }
//         );
//     }

//     public updateDummy(correlationId: string, dummy: any,
//         callback: (err: any, result: Dummy) => void): void {

//         let dummyObj = new messages.Dummy();
//         dummyObj.setId(dummy.id);
//         dummyObj.setKey(dummy.key);
//         dummyObj.setContent(dummy.content);

//         let request = new messages.DummyObjectRequest();
//         request.setDummy(dummyObj);

//         this.call('update_dummy',
//             correlationId,
//             request,
//             (err, result) => {
//                 this.instrument(correlationId, 'dummy.update');

//                 result = result != null ? result.toObject() : null;
//                 if (result && result.id == "" && result.key == "")
//                     result = null;

//                 callback(err, result);
//             }
//         );
//     }

//     public deleteDummy(correlationId: string, dummyId: string,
//         callback: (err: any, result: Dummy) => void): void {

//         let request = new messages.DummyIdRequest();
//         request.setDummyId(dummyId);

//         this.call('delete_dummy_by_id',
//             correlationId,
//             request,
//             (err, result) => {
//                 this.instrument(correlationId, 'dummy.delete_by_id');

//                 result = result != null ? result.toObject() : null;
//                 if (result && result.id == "" && result.key == "")
//                     result = null;

//                 callback(err, result);
//             }
//         );
//     }

// }
