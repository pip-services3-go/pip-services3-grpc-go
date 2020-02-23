package test_services

// let _ = require('lodash');
// let services = require('../../../test/protos/dummies_grpc_pb');
// let messages = require('../../../test/protos/dummies_pb');

// import { IReferences } from 'pip-services3-commons-node';
// import { Descriptor } from 'pip-services3-commons-node';
// import { DataPage } from 'pip-services3-commons-node';
// import { FilterParams } from 'pip-services3-commons-node';
// import { PagingParams } from 'pip-services3-commons-node';
// import { ObjectSchema } from 'pip-services3-commons-node';
// import { TypeCode } from 'pip-services3-commons-node';
// import { FilterParamsSchema } from 'pip-services3-commons-node';
// import { PagingParamsSchema } from 'pip-services3-commons-node';

// import { Dummy } from '../Dummy';
// import { DummySchema } from '../DummySchema';
// import { GrpcService } from '../../src/services/GrpcService';
// import { IDummyController } from '../IDummyController';

// export class DummyGrpcService2 extends GrpcService {
//     private _controller: IDummyController;
//     private _numberOfCalls: number = 0;

//     public constructor() {
//         super(services.DummiesService);
//         this._dependencyResolver.put('controller', new Descriptor("pip-services-dummies", "controller", "default", "*", "*"));
//     }

// 	public setReferences(references: IReferences): void {
// 		super.setReferences(references);
//         this._controller = this._dependencyResolver.getOneRequired<IDummyController>('controller');
//     }

//     public getNumberOfCalls(): number {
//         return this._numberOfCalls;
//     }

//     private incrementNumberOfCalls(
//         call: any, callback: (err: any, result: any) => void, next: () => void) {
//         this._numberOfCalls++;
//         next();
//     }

//     private dummyToObject(dummy: Dummy): any {
//         let obj = new messages.Dummy();

//         if (dummy) {
//             obj.setId(dummy.id);
//             obj.setKey(dummy.key);
//             obj.setContent(dummy.content);
//         }

//         return obj;
//     }

//     private dummyPageToObject(page: DataPage<Dummy>): any {
//         let obj = new messages.DummiesPage();

//         if (page) {
//             obj.setTotal(page.total);
//             let data = _.map(page.data, this.dummyToObject);
//             obj.setDataList(data);
//         }

//         return obj;
//     }

//     private getPageByFilter(call: any, callback: any) {
//         let request = call.request.toObject();
//         let filter = FilterParams.fromValue(request.filterMap);
//         let paging = PagingParams.fromValue(call.request.paging);

//         this._controller.getPageByFilter(
//             call.request.correlation_id,
//             filter,
//             paging,
//             (err, page) => {
//                 callback(err, this.dummyPageToObject(page));
//             }
//         );
//     }

//     private getOneById(call: any, callback: any) {
//         let request = call.request.toObject();

//         this._controller.getOneById(
//             request.correlation_id,
//             request.dummy_id,
//             (err, result) => {
//                 callback(err, this.dummyToObject(result));
//             }
//         );
//     }

//     private create(call: any, callback: any) {
//         let request = call.request.toObject();

//         this._controller.create(
//             request.correlation_id,
//             request.dummy,
//             (err, result) => {
//                 callback(err, this.dummyToObject(result));
//             }
//         );
//     }

//     private update(call: any, callback: any) {
//         let request = call.request.toObject();

//         this._controller.update(
//             request.correlation_id,
//             request.dummy,
//             (err, result) => {
//                 callback(err, this.dummyToObject(result));
//             }
//         );
//     }

//     private deleteById(call: any, callback: any) {
//         let request = call.request.toObject();

//         this._controller.deleteById(
//             request.correlation_id,
//             request.dummy_id,
//             (err, result) => {
//                 callback(err, this.dummyToObject(result));
//             }
//         );
//     }

//     public register() {
//         this.registerInterceptor(this.incrementNumberOfCalls);

//         this.registerMethod(
//             'get_dummies',
//             null,
//             // new ObjectSchema(true)
//             //     .withOptionalProperty("paging", new PagingParamsSchema())
//             //     .withOptionalProperty("filter", new FilterParamsSchema()),
//             this.getPageByFilter
//         );

//         this.registerMethod(
//             'get_dummy_by_id',
//             null,
//             // new ObjectSchema(true)
//             //     .withRequiredProperty("dummy_id", TypeCode.String),
//             this.getOneById
//         );

//         this.registerMethod(
//             'create_dummy',
//             null,
//             // new ObjectSchema(true)
//             //     .withRequiredProperty("dummy", new DummySchema()),
//             this.create
//         );

//         this.registerMethod(
//             'update_dummy',
//             null,
//             // new ObjectSchema(true)
//             //     .withRequiredProperty("dummy", new DummySchema()),
//             this.update
//         );

//         this.registerMethod(
//             'delete_dummy_by_id',
//             null,
//             // new ObjectSchema(true)
//             //     .withRequiredProperty("dummy_id", TypeCode.String),
//             this.deleteById
//         );
//     }
// }
