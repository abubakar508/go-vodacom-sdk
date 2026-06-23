# Endpoint Coverage

All endpoint methods require a valid SessionID, except `GenerateSessionKey` / `GenerateSession`.

| API | Method | SDK method |
|---|---:|---|
| Generate SessionKey | GET | `GenerateSession`, `GenerateSessionKey` |
| C2B Single Stage | POST | `C2BSingleStageWithSession` |
| B2C Single Stage | POST | `B2CSingleStageWithSession` |
| B2B Single Stage | POST | `B2BSingleStageWithSession` |
| Reversal | POST | `ReversalWithSession` |
| Query Transaction Status | GET | `QueryTransactionStatusWithSession` |
| C2B Multi Stage | POST | `C2BMultiStageWithSession` |
| Update Transaction Status | PUT | `UpdateTransactionStatusWithSession` |
| Direct Debit Create | POST | `DirectDebitCreateWithSession` |
| Direct Debit Payment | POST | `DirectDebitPaymentWithSession` |
| Query Beneficiary Name | GET | `QueryBeneficiaryNameWithSession` |
| Query Direct Debit | GET | `QueryDirectDebitWithSession` |
| Cancel Direct Debit | PUT | `CancelDirectDebitWithSession` |

## Endpoint pattern

```text
https://openapi.m-pesa.com/{sandbox|openapi}/ipg/v2/{market}/{endpoint}/
```

Examples:

```text
/sandbox/ipg/v2/vodafoneGHA/c2bPayment/singleStage/
/openapi/ipg/v2/vodacomDRC/directDebitPayment/
```
