/*
Package cctools provides an SQL-like structure for programming Hyperledger Fabric chaincodes.

Transactions respond with status 200 by default. Errors carry their own status through
errors.ICCError, and a transaction may set a 2xx success status by defining
transactions.Transaction.ResponseRoutine and running it with transactions.RunResponse.
*/
package cctools
