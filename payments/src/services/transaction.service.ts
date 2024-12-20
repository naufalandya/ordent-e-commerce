import * as grpc from '@grpc/grpc-js';
import snap from '../libs/midtrans.lib';
import prisma from '../libs/prisma.lib';

export const handleTransaction = async (call: grpc.ServerUnaryCall<any, any>, callback: grpc.sendUnaryData<any>) => {

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

      const transactionMidtrans = await snap.createTransaction(
        orderDetailsMidtrans
      );

      await prisma.transactionHistory.update({
        where : { 
          id : transactionHistoryId
        },
        data : {
          midtrans_order_id : orderDetailsMidtrans.transaction_details.order_id,
          payment : transactionMidtrans.redirect_url
        }
      })
    
      console.log(`Received transactionHistoryId: ${transactionHistoryId}`);
      callback(null, { message: `Transaction with ID ${transactionHistoryId} processed successfully.` });
  } catch (err) {
    throw err
  }
};