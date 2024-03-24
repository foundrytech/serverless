# Email Verification Cloud Function

This repository hosts the code for a serverless Google Cloud Function designed to handle email verification for new user accounts. The function is triggered by a Pub/Sub event whenever a new user account is created.

## Functionality:

1. **Email Verification Link:**
   - Generates and emails a verification link to the user's provided email address.
   - The verification link is valid for 2 minutes to expedite the demo process.
   - Expired links are rendered invalid and cannot be used for verification.

2. **Email Tracking:**
   - Tracks all emails sent within a CloudSQL instance.
   - Utilizes the same CloudSQL instance and database used by the web application for seamless integration.
