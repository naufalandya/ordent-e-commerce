import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { handleTransaction } from './src/services/transaction.service';

const PROTO_PATH = './proto/transaction.proto';

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {});
const transactionProto = grpc.loadPackageDefinition(packageDefinition) as any;

if (!transactionProto.TransactionService || !transactionProto.TransactionService.TransactionService) {
  console.error("TransactionService is not correctly defined in proto file.");
  process.exit(1);
}

const server = new grpc.Server();

const TransactionService = transactionProto.TransactionService.TransactionService;

server.addService(TransactionService.service, {
  HandleTransaction: handleTransaction,
});

const PORT = '50051';
server.bindAsync(`0.0.0.0:${PORT}`, grpc.ServerCredentials.createInsecure(), (error, port) => {
  if (error) {
    console.error(`Server failed to start: ${error}`);
    return;
  }
  console.log(`Server running at http://127.0.0.1:${port}`);
});
