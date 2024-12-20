"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.handleTransaction = void 0;
const midtrans_lib_1 = __importDefault(require("../libs/midtrans.lib"));
const prisma_lib_1 = __importDefault(require("../libs/prisma.lib"));
const handleTransaction = (call, callback) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { transactionHistoryId, total } = call.request;
        const orderDetailsMidtrans = {
            transaction_details: {
                order_id: "order-id-node-" + Math.round(new Date().getTime() / 1000),
                gross_amount: Number(total),
            },
            credit_card: {
                secure: true,
            },
        };
        const transactionMidtrans = yield midtrans_lib_1.default.createTransaction(orderDetailsMidtrans);
        yield prisma_lib_1.default.transactionHistory.update({
            where: {
                id: transactionHistoryId
            },
            data: {
                midtrans_order_id: orderDetailsMidtrans.transaction_details.order_id,
                payment: transactionMidtrans.redirect_url
            }
        });
        console.log(`Received transactionHistoryId: ${transactionHistoryId}`);
        callback(null, { message: `Transaction with ID ${transactionHistoryId} processed successfully.` });
    }
    catch (err) {
        throw err;
    }
});
exports.handleTransaction = handleTransaction;
