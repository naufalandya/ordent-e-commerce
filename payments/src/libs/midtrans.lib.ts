import midtransClient from 'midtrans-client';
import dotenv from 'dotenv';

dotenv.config();

const midtransServerKey: string | undefined = process.env.MIDTRANS_SERVER_KEY;
const midtransClientKey: string | undefined = process.env.MIDTRANS_CLIENT_KEY;

console.log(midtransClientKey, midtransServerKey)

if (!midtransServerKey || !midtransClientKey) {
  throw new Error('Missing required environment variables for Midtrans.');
}

const snap = new midtransClient.Snap({
  isProduction: false,
  serverKey: midtransServerKey,
  clientKey: midtransClientKey
});

export default snap;
