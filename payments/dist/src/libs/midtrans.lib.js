"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const midtrans_client_1 = __importDefault(require("midtrans-client"));
const dotenv_1 = __importDefault(require("dotenv"));
dotenv_1.default.config();
const midtransServerKey = process.env.MIDTRANS_SERVER_KEY;
const midtransClientKey = process.env.MIDTRANS_CLIENT_KEY;
console.log(midtransClientKey, midtransServerKey);
if (!midtransServerKey || !midtransClientKey) {
    throw new Error('Missing required environment variables for Midtrans.');
}
const snap = new midtrans_client_1.default.Snap({
    isProduction: false,
    serverKey: midtransServerKey,
    clientKey: midtransClientKey
});
exports.default = snap;
