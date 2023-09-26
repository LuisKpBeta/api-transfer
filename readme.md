# Transfer API

This API was created based on the (DevGym)[https://app.devgym.com.br/challenges/9af13172-e1fe-4c2e-ac10-cb6b0bcf2efc] challenge  

## Requirements

- [x] Create an endpoint that receives a user ID and returns their balance.
- [x] Validate if the source user has sufficient balance before performing a transfer.
- [x] Create an endpoint that receives two user IDs and a monetary value representing the transfer between them.
- [x] Consider the possibility of concurrent transfers where two people transfer money to a third person at the same time.
- [x] If a transfer fails, the balance of the source user should be restored.

Note: There is no need to create endpoints for creating users; populate the database in a way that both users exist, and transfers can be made between them.

## Main Challenges

- Ensure that two transfers executed at the same time will be processed correctly.
- Forcing a failure operation in PostgreSQL to simulate a failed operation, ensuring that if failures occur during the transaction, the entire operation will be rolled back.