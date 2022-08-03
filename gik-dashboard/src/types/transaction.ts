export interface Transaction {
    ID: number;
    timestamp: number;
    clientId: number;
    type: boolean;
    signerId: number;
    totalQuantity: number;
}

export interface TransactionItem {
    ID: number;
    transactionId: number;
    quantity: number;
    productId: number;
}