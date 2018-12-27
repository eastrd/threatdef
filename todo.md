# TODO List

Backend: Receiver (Golang)
- Add "Create Table if Not Exist" function into each handler (DONE)
- Add Command handler (DONE)
- Make into prod ready (DONE)
- Add Payload handler
- Add IP location table
- Collect & Refractor MySQL functions
    - SQLX library?
    - Using Struct?

Backend: RestAPI (Golang)
- Write endpoints for:
    - Return latest X tunnel records based on decreasing ID in JSON format (DONE)
    - Return latest X cmd records based on decreasing ID in JSON format
    - Return top X login combinations based on decreasing num_attempts in JSON format
    - Design a way to show live attacks
    - Design a way to show heatmap of attacks based on regions

Sender (Python3)
- Check if scheduling function works (DONE)
- Make interaction with backend non-blocking
- Make sender able to send past historical logs to backend receiver


Frontend (Vue.js)
- Redesign UX & UI