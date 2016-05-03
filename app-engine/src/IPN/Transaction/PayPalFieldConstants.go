package TransActionDao

// https://developer.paypal.com/docs/classic/ipn/integration-guide/IPNandPDTVariables/

// Transaction and notification-related variables
const (
	FIELD_BUSINESS = "business" // Email address or account ID of the payment recipient (that is, the merchant). Equivalent to the values of receiver_email (if payment is sent to primary account) and business set in the Website Payment HTML. Note: The value of this variable is normalized to lowercase characters. Length: 127 characters
	FIELD_CHARSET = "charset" // Character set
	FIELD_CUSTOM = "custom" // Custom value as passed by you, the merchant. These are pass-through variables that are never presented to your customer  Length: 255 characters
	FIELD_IPN_TRACK_ID = "ipn_track_id" // Internal; only for use by MTS and DTS
	FIELD_NOTIFY_VERSION = "notify_version" // Message's version number
	FIELD_PARENT_TXN_ID = "parent_txn_id" // In the case of a refund, reversal, or canceled reversal, this variable contains the txn_id of the original transaction, while txn_id contains a new ID for the new transaction.  Length: 19 characters
	FIELD_RECEIPT_ID = "receipt_id" // Unique ID generated during guest checkout (payment by credit card without logging in).
	FIELD_RECEIVER_EMAIL = "receiver_email" // Primary email address of the payment recipient (that is, the merchant). If the payment is sent to a non-primary email address on your PayPal account, the receiver_email is still your primary email. Note: The value of this variable is normalized to lowercase characters. Length: 127 characters
	FIELD_RECEIVER_ID = "receiver_id" // Unique account ID of the payment recipient (i.e., the merchant). This is the same as the recipient's referral ID.  Length: 13 characters
	FIELD_RESEND = "resend" // Whether this IPN message was resent (equals true); otherwise, this is the original message.
	FIELD_RESIDENCE_COUNTRY = "residence_country" // ISO 3166 country code associated with the country of residence  Length: 2 characters
	FIELD_TEST_IPN = "test_ipn" // Whether the message is a test message. It is one of the following values: 1 — the message is directed to the Sandbox
	FIELD_TXN_ID = "txn_id" // The merchant's original transaction identification number for the payment from the buyer, against which the case was registered.
	FIELD_TXN_TYPE = "txn_type" // The kind of transaction for which the IPN message was sent.
	FIELD_VERIFY_SIGN = "verify_sign" // Encrypted string used to validate the authenticity of the transaction
)

//Buyer Information
const (
	FIELD_ADDRESS_COUNTRY = "address_country" // Country of customer's address  Length: 64 characters
	FIELD_ADDRESS_CITY = "address_city" // City of customer's address  Length: 40 characters
	FIELD_ADDRESS_COUNTRY_CODE = "address_country_code" // ISO 3166 country code associated with customer's address  Length: 2 characters
	FIELD_ADDRESS_NAME = "address_name" // Name used with address (included when the customer provides a Gift Address)  Length: 128 characters
	FIELD_ADDRESS_STATE = "address_state" // State of customer's address  Length: 40 characters
	FIELD_ADDRESS_STATUS = "address_status" // Whether the customer provided a confirmed address. It is one of the following values: confirmed — Customer provided a confirmed address. unconfirmed — Customer provided an unconfirmed address.
	FIELD_ADDRESS_STREET = "address_street" // Customer's street address.  Length: 200 characters
	FIELD_ADDRESS_ZIP = "address_zip" // Zip code of customer's address.  Length: 20 characters
	FIELD_CONTACT_PHONE = "contact_phone" // Customer's telephone number.  Length: 20 characters
	FIELD_FIRST_NAME = "first_name" // Customer's first name  Length: 64 characters
	FIELD_LAST_NAME = "last_name" // Customer's last name  Length: 64 characters
	FIELD_PAYER_BUSINESS_NAME = "payer_business_name" // Customer's company name, if customer is a business  Length: 127 characters
	FIELD_PAYER_EMAIL = "payer_email" // Customer's primary email address. Use this email to provide any credits.  Length: 127 characters
	FIELD_PAYER_ID = "payer_id" // Unique customer ID.  Length: 13 characters
)

// Payment information variables
const (
	FIELD_AUTH_AMOUNT = "auth_amount" // Authorization amount
	FIELD_AUTH_EXP = "auth_exp" // Authorization expiration date and time, in the following format: HH:MM:SS DD Mmm YY, YYYY PST  Length: 28 characters
	FIELD_AUTH_ID = "auth_id" // Authorization identification number  Length: 19 characters
	FIELD_AUTH_STATUS = "auth_status" // Status of authorization
	FIELD_ECHECK_TIME_PROCESSED = "echeck_time_processed" // The time an eCheck was processed; for example, when the status changes to Success or Completed. The format is as follows: hh:mm:ss MM DD, YYYY ZONE, e.g. 04:55:30 May 26, 2011 PDT.
	FIELD_EXCHANGE_RATE = "exchange_rate" // Exchange rate used if a currency conversion occurred.
	FIELD_FRAUD_MANAGEMENT_PENDING_FILTERS_X = "fraud_management_pending_filters_x" // One or more filters that identify a triggering action associated with one of the following payment_status values: Pending, Completed, Denied, where x is a number starting with 1 that makes the IPN variable name unique; x is not the filter's ID number. The filters and their ID numbers are as follows: 1 - AVS No Match 2 - AVS Partial Match 3 - AVS Unavailable/Unsupported 4 - Card Security Code (CSC) Mismatch 5 - Maximum Transaction Amount 6 - Unconfirmed Address 7 - Country Monitor 8 - Large Order Number 9 - Billing/Shipping Address Mismatch 10 - Risky ZIP Code 11 - Suspected Freight Forwarder Check 12 - Total Purchase Price Minimum 13 - IP Address Velocity 14 - Risky Email Address Domain Check 15 - Risky Bank Identification Number (BIN) Check 16 - Risky IP Address Range 17 - PayPal Fraud Model
	FIELD_INVOICE = "invoice" // Pass-through variable you can use to identify your Invoice Number for this purchase. If omitted, no variable is passed back.  Length: 127 characters
	FIELD_ITEM_NAMEX = "item_namex" // Item name as passed by you, the merchant. Or, if not passed by you, as entered by your customer. If this is a shopping cart transaction, PayPal will append the number of the item (e.g., item_name1, item_name2, and so forth).  Length: 127 characters
	FIELD_ITEM_NUMBERX = "item_numberx" // Pass-through variable for you to track purchases. It will get passed back to you at the completion of the payment. If omitted, no variable will be passed back to you. If this is a shopping cart transaction, PayPal will append the number of the item (e.g., item_number1, item_number2, and so forth)  Length: 127 characters
	FIELD_MC_CURRENCY = "mc_currency" // For payment IPN notifications, this is the currency of the payment. For non-payment subscription IPN notifications (i.e., txn_type= signup, cancel, failed, eot, or modify), this is the currency of the subscription. For payment subscription IPN notifications, it is the currency of the payment (i.e., txn_type = subscr_payment)
	FIELD_MC_FEE = "mc_fee" // Transaction fee associated with the payment. mc_gross minus mc_fee equals the amount deposited into the receiver_email account. Equivalent to payment_fee for USD payments. If this amount is negative, it signifies a refund or reversal, and either of those payment statuses can be for the full or partial amount of the original transaction fee.
	FIELD_MC_GROSS = "mc_gross" // Full amount of the customer's payment, before transaction fee is subtracted. Equivalent to payment_gross for USD payments. If this amount is negative, it signifies a refund or reversal, and either of those payment statuses can be for the full or partial amount of the original transaction.
	FIELD_MC_GROSS_X = "mc_gross_x" // The amount is in the currency of mc_currency, where x is the shopping cart detail item number. The sum of mc_gross_x should total mc_gross.
	FIELD_MC_HANDLING = "mc_handling" // Total handling amount associated with the transaction.
	FIELD_MC_SHIPPING = "mc_shipping" // Total shipping amount associated with the transaction.
	FIELD_MC_SHIPPINGX = "mc_shippingx" // This is the combined total of shipping1 and shipping2 Website Payments Standard variables, where x is the shopping cart detail item number. The shippingx variable is only shown when the merchant applies a shipping amount for a specific item. Because profile shipping might apply, the sum of shippingx might not be equal to shipping.
	FIELD_MEMO = "memo" // Memo as entered by your customer in PayPal Website Payments note field.  Length: 255 characters
	FIELD_NUM_CART_ITEMS = "num_cart_items" // If this is a PayPal Shopping Cart transaction, number of items in cart.
	FIELD_OPTION_NAME1 = "option_name1" // Option 1 name as requested by you. PayPal appends the number of the item where x represents the number of the shopping cart detail item (e.g., option_name1, option_name2).  Length: 64 characters
	FIELD_OPTION_NAME2 = "option_name2" // Option 2 name as requested by you. PayPal appends the number of the item where x represents the number of the shopping cart detail item (e.g., option_name2, option_name2).  Length: 64 characters
	FIELD_OPTION_SELECTION1 = "option_selection1" // Option 1 choice as entered by your customer.  PayPal appends the number of the item where x represents the number of the shopping cart detail item (e.g., option_selection1, option_selection2).  Length: 200 characters
	FIELD_OPTION_SELECTION2 = "option_selection2" // Option 2 choice as entered by your customer.  PayPal appends the number of the item where x represents the number of the shopping cart detail item (e.g., option_selection1, option_selection2).  Length: 200 characters
	FIELD_PAYER_STATUS = "payer_status" // Whether the customer has a verified PayPal account. verified — Customer has a verified PayPal account. unverified — Customer has an unverified PayPal account.
	FIELD_PAYMENT_DATE = "payment_date" // Time/Date stamp generated by PayPal, in the following format: HH:MM:SS Mmm DD, YYYY PDT  Length: 28 characters
	FIELD_PAYMENT_FEE = "payment_fee" // USD transaction fee associated with the payment. payment_gross minus payment_fee equals the amount deposited into the receiver email account. Is empty for non-USD payments. If this amount is negative, it signifies a refund or reversal, and either of those payment statuses can be for the full or partial amount of the original transaction fee. Note: This is a deprecated field. Use mc_fee instead.
	FIELD_PAYMENT_FEE_X = "payment_fee_x" // If the payment is USD, then the value is the same as that for mc_fee_x, where x is the record number; if the currency is not USD, then this is an empty string. Note: This is a deprecated field. Use mc_fee_x instead.
	FIELD_PAYMENT_GROSS = "payment_gross" // Full USD amount of the customer's payment, before transaction fee is subtracted. Will be empty for non-USD payments. This is a legacy field replaced by mc_gross. If this amount is negative, it signifies a refund or reversal, and either of those payment statuses can be for the full or partial amount of the original transaction.
	FIELD_PAYMENT_GROSS_X = "payment_gross_x" // If the payment is USD, then the value for this is the same as that for the mc_gross_x, where x is the record number the mass pay item. If the currency is not USD, this is an empty string. Note: This is a deprecated field. Use mc_gross_x instead.
	FIELD_PAYMENT_STATUS = "payment_status" // The status of the payment: Canceled_Reversal: A reversal has been canceled. For example, you won a dispute with the customer, and the funds for the transaction that was reversed have been returned to you.  Completed: The payment has been completed, and the funds have been added successfully to your account balance.  Created: A German ELV payment is made using Express Checkout.  Denied: The payment was denied. This happens only if the payment was previously pending because of one of the reasons listed for the pending_reason variable or the Fraud_Management_Filters_x variable.  Expired: This authorization has expired and cannot be captured.  Failed: The payment has failed. This happens only if the payment was made from your customer's bank account.  Pending: The payment is pending. See pending_reason for more information.  Refunded: You refunded the payment.  Reversed: A payment was reversed due to a chargeback or other type of reversal. The funds have been removed from your account balance and returned to the buyer. The reason for the reversal is specified in the ReasonCode element.  Processed: A payment has been accepted.  Voided: This authorization has been voided.
	FIELD_PAYMENT_TYPE = "payment_type" // echeck: This payment was funded with an eCheck.  instant: This payment was funded with PayPal balance, credit card, or Instant Transfer.
	FIELD_PENDING_REASON = "pending_reason" // This variable is set only if payment_status is Pending.  address: The payment is pending because your customer did not include a confirmed shipping address and your Payment Receiving Preferences is set yo allow you to manually accept or deny each of these payments. To change your preference, go to the Preferences section of your Profile.  authorization: You set the payment action to Authorization and have not yet captured funds.  echeck: The payment is pending because it was made by an eCheck that has not yet cleared.  intl: The payment is pending because you hold a non-U.S. account and do not have a withdrawal mechanism. You must manually accept or deny this payment from your Account Overview.  multi_currency: You do not have a balance in the currency sent, and you do not have your profiles's Payment Receiving Preferences option set to automatically convert and accept this payment. As a result, you must manually accept or deny this payment.  order: You set the payment action to Order and have not yet captured funds.  paymentreview: The payment is pending while it is reviewed by PayPal for risk.  regulatory_review: The payment is pending because PayPal is reviewing it for compliance with government regulations. PayPal will complete this review within 72 hours. When the review is complete, you will receive a second IPN message whose payment_status/reason code variables indicate the result.  unilateral: The payment is pending because it was made to an email address that is not yet registered or confirmed.  upgrade: The payment is pending because it was made via credit card and you must upgrade your account to Business or Premier status before you can receive the funds. upgrade can also mean that you have reached the monthly limit for transactions on your account.  verify: The payment is pending because you are not yet verified. You must verify your account before you can accept this payment.  other: The payment is pending for a reason other than those listed above. For more information, contact PayPal Customer Service.
	FIELD_PROTECTION_ELIGIBILITY = "protection_eligibility" // ExpandedSellerProtection: Seller is protected by Expanded seller protection SellerProtection: Seller is protected by PayPal's Seller Protection Policy  None: Seller is not protected under Expanded seller protection nor the Seller Protection Policy
	FIELD_QUANTITY = "quantity" // Quantity as entered by your customer or as passed by you, the merchant. If this is a shopping cart transaction, PayPal appends the number of the item (e.g. quantity1, quantity2).
	FIELD_REASON_CODE = "reason_code" // This variable is set if payment_status is Reversed, Refunded, Canceled_Reversal, or Denied.  adjustment_reversal: Reversal of an adjustment.  admin_fraud_reversal: The transaction has been reversed due to fraud detected by PayPal administrators.  admin_reversal: The transaction has been reversed by PayPal administrators.  buyer-complaint: The transaction has been reversed due to a complaint from your customer.  chargeback: The transaction has been reversed due to a chargeback by your customer.  chargeback_reimbursement: Reimbursement for a chargeback.  chargeback_settlement: Settlement of a chargeback.  guarantee: The transaction has been reversed because your customer exercised a money-back guarantee.  other: Unspecified reason.  refund: The transaction has been reversed because you gave the customer a refund.  regulatory_block: PayPal blocked the transaction due to a violation of a government regulation. In this case, payment_status is Denied.  regulatory_reject: PayPal rejected the transaction due to a violation of a government regulation and returned the funds to the buyer. In this case, payment_status is Denied.  regulatory_review_exceeding_sla: PayPal did not complete the review for compliance with government regulations within 72 hours, as required. Consequently, PayPal auto-reversed the transaction and returned the funds to the buyer. In this case, payment_status is Denied. Note that "sla" stand for "service level agreement".  unauthorized_claim: The transaction has been reversed because it was not authorized by the buyer.  unauthorized_spoof: The transaction has been reversed due to a customer dispute in which an unauthorized spoof is suspected.  Note: Additional codes may be returned.
	FIELD_REMAINING_SETTLE = "remaining_settle" // Remaining amount that can be captured with Authorization and Capture
	FIELD_SETTLE_AMOUNT = "settle_amount" // Amount that is deposited into the account's primary balance after a currency conversion from automatic conversion (through your Payment Receiving Preferences) or manual conversion (through manually accepting a payment).
	FIELD_SETTLE_CURRENCY = "settle_currency" // Currency of settle_amount.
	FIELD_SHIPPING = "shipping" // Shipping charges associated with this transaction.  Format: unsigned, no currency symbol, two decimal places.
	FIELD_SHIPPING_METHOD = "shipping_method" // The name of a shipping method from the Shipping Calculations section of the merchant's account profile. The buyer selected the named shipping method for this transaction.
	FIELD_TAX = "tax" // Amount of tax charged on payment. PayPal appends the number of the item (e.g., item_name1, item_name2). The taxx variable is included only if there was a specific tax amount applied to a particular shopping cart item. Because total tax may apply to other items in the cart, the sum of taxx might not total to tax.
	FIELD_TRANSACTION_ENTITY = "transaction_entity" // Authorization and Capture transaction entity
)

const (
	// IPN transaction types (FIELD_TXN_TYPE)
	// Credit card chargeback if the case_type variable contains chargeback
	TXN_TYPE_ADJUSTMENT = "adjustment" // A dispute has been resolved and closed
	TXN_TYPE_CART = "cart" // Payment received for multiple items; source is Express Checkout or the PayPal Shopping Cart.
	TXN_TYPE_EXPRESS_CHECKOUT = "express_checkout" // Payment received for a single item; source is Express Checkout
	TXN_TYPE_MASSPAY = "masspay" // Payment sent using Mass Pay
	TXN_TYPE_MERCH_PMT = "merch_pmt" // Monthly subscription paid for Website Payments Pro, Reference transactions, or Billing Agreement payments
	TXN_TYPE_MP_CANCEL = "mp_cancel" // Billing agreement cancelled
	TXN_TYPE_MP_SIGNUP = "mp_signup" // Created a billing agreement
	TXN_TYPE_NEW_CASE = "new_case" // A new dispute was filed
	TXN_TYPE_PAYOUT = "payout" // A payout related to a global shipping transaction was completed.
	TXN_TYPE_PRO_HOSTED = "pro_hosted" // Payment received; source is Website Payments Pro Hosted Solution.
	TXN_TYPE_RECURRING_PAYMENT = "recurring_payment" // Recurring payment received
	TXN_TYPE_RECURRING_PAYMENT_EXPIRED = "recurring_payment_expired" // Recurring payment expired
	TXN_TYPE_RECURRING_PAYMENT_FAILED = "recurring_payment_failed" // Recurring payment failed  This transaction type is sent if:  The attempt to collect a recurring payment fails The "max failed payments" setting in the customer's recurring payment profile is 0  In this case, PayPal tries to collect the recurring payment an unlimited number of times without ever suspending the customer's recurring payments profile.
	TXN_TYPE_RECURRING_PAYMENT_PROFILE_CANCEL = "recurring_payment_profile_cancel" // Recurring payment profile canceled
	TXN_TYPE_RECURRING_PAYMENT_PROFILE_CREATED = "recurring_payment_profile_created" // Recurring payment profile created
	TXN_TYPE_RECURRING_PAYMENT_SKIPPED = "recurring_payment_skipped" // Recurring payment skipped; it will be retried up to 3 times, 5 days apart
	TXN_TYPE_RECURRING_PAYMENT_SUSPENDED = "recurring_payment_suspended" // Recurring payment suspended  This transaction type is sent if PayPal tried to collect a recurring payment, but the related recurring payments profile has been suspended.
	TXN_TYPE_RECURRING_PAYMENT_SUSPENDED_DUE_TO_MAX_FAILED_PAYMENT = "recurring_payment_suspended_due_to_max_failed_payment" // Recurring payment failed and the related recurring payment profile has been suspended  This transaction type is sent if: PayPal's attempt to collect a recurring payment failed The "max failed payments" setting in the customer's recurring payment profile is 1 or greater the number of attempts to collect payment has exceeded the value specified for "max failed payments"  In this case, PayPal suspends the customer's recurring payment profile.
	TXN_TYPE_SEND_MONEY = "send_money" // Payment received; source is the Send Money tab on the PayPal website
	TXN_TYPE_SUBSCR_CANCEL = "subscr_cancel" // Subscription canceled
	TXN_TYPE_SUBSCR_EOT = "subscr_eot" // Subscription expired
	TXN_TYPE_SUBSCR_FAILED = "subscr_failed" // Subscription payment failed
	TXN_TYPE_SUBSCR_MODIFY = "subscr_modify" // Subscription modified
	TXN_TYPE_SUBSCR_PAYMENT = "subscr_payment" // Subscription payment received
	TXN_TYPE_SUBSCR_SIGNUP = "subscr_signup" // Subscription started
	TXN_TYPE_VIRTUAL_TERMINAL = "virtual_terminal" // Payment received; source is Virtual Terminal
	TXN_TYPE_WEB_ACCEPT = "web_accept" // Payment received; source is any of the following: A Direct Credit Card (Pro) transaction A Buy Now, Donation or Smart Logo for eBay auctions button
)

const (
	STATUS_COMPLEATED = "Completed"
)