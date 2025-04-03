# User Onboarding Status Changes: Stage 1 - SIM_VERIFICATION

This document outlines the changes to the user onboarding process for **Stage 1: SIM_VERIFICATION**, including the steps, API calls, and expected `useronboarding status` API responses. These updates are intended for frontend implementation.

---

## Overview of Stage 1: SIM_VERIFICATION

Stage 1 consists of two sequential steps:
1. **AUTHORIZATION**: Handles user authentication and SIM verification.
2. **PERSONAL_DETAILS**: Collects and verifies user personal information.

Below are the details for each step, including API calls and status updates.

---

## 1. AUTHORIZATION

### Description
This is the first step of Stage 1, where the user authenticates and verifies their SIM. It involves three API calls:
- `/authorization`
- `/initiate-sim-verification`
- `/sim-verification-status`

### Failure Scenario
If any of the above API calls fail, the `useronboarding status` API response will reflect the incomplete state:
```json
{
  "user_id": "d6279cf6-601c-4",
  "current_stage_name": "SIM_VERIFICATION",
  "current_step_name": "AUTHORIZATION",
  "is_sim_verification_complete": false
  // Rest of the data remains unchanged
}

### Success Scenario
If the SIM verification is successful, the `useronboarding status` API response will reflect the completed state:
```json
{
  "user_id": "d6279cf6-601c-4",
  "current_stage_name": "SIM_VERIFICATION",
  "current_step_name": "PERSONAL_DETAILS",
  "is_sim_verification_complete": true
  // Rest of the data remains unchanged
}

## 2. PERSONAL_DETAILS

### Description
This is the second step of Stage 1, where the user provides personal information. It involves three API calls:
- `/email/sendverification`
- `/email/verification-status`
- `/onboarding/personal-information`

### Failure Scenario 
If any of the above API calls fail, the `useronboarding status` API response will reflect the incomplete state:
```json
{
  "user_id": "d6279cf6-601c-4",
  "current_stage_name": "SIM_VERIFICATION",
  "current_step_name": "PERSONAL_DETAILS",
  "is_sim_verification_complete": true
  // Rest of the data remains unchanged
}  

### Success Scenario
If the personal information is successfully updated, the `useronboarding status` API response will reflect the completed state:
```json
{
  "user_id": "d6279cf6-601c-4",
  "current_stage_name": "SIM_VERIFICATION",
  "current_step_name": "",
  "is_sim_verification_complete": true
  // Rest of the data remains unchanged
}
```
