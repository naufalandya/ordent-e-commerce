"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.handleTransaction = void 0;
const handleTransaction = (call, callback) => {
    const { transactionHistoryId } = call.request;
    console.log(`Received transactionHistoryId: ${transactionHistoryId}`);
    callback(null, { message: `Transaction with ID ${transactionHistoryId} processed successfully.` });
};
exports.handleTransaction = handleTransaction;
