## Table of Contents

- [Features Completed](#features-completed)
- [How To Run](#how-to-run)
- [BE API](#be-api)


## Features Completed
**FE**

Employee Interface (Svelte)
- View assigned shifts []
- View available (unassigned) shifts []
- Request to work on available shifts (creates a pending request) []
- See status of requests (pending, approved, rejected) []
- (Optional) Set availability calendar []

Admin Interface (Svelte or API tool)
- Create/manage shift schedules: date, start_time, end_time, role, location [x]
- See all pending requests from workers [x]
- Approve/reject shift requests [x]
- View finalized rosters (approved shift assignments) []
- Edit/delete shifts and assignments  (50%)

**BE**

Backend (Go)
- RESTful Endpoints:
- CRUD for shifts [x]
- Worker shift signup [x]
- Admin approval/rejection [x]
- View assignments by worker or by day (50%)
- Use SQLite (recommended for quick setup)
- Use only raw SQL if possible
- Enforce basic business requirements (see below)

Business Requirements
These are basic business requirements:
- Workers cannot request shifts already assigned to someone else [x]
- No overlapping shift requests allowed per worker [x]
- Max 1 shift per day, max 5 shifts per week per worker [x]
- Admin can override or reassign approved shifts []
- Conflict checking must occur on both worker request and admin approval (50%)
-  Shift times are stored and compared in UTC [x]


## How To Run

To clone and run this application, you'll need [Git](https://git-scm.com) and [Android Studio](https://developer.android.com/studio) installed on your computer. From your command line:

```bash
# Clone this repository
git clone https://github.com/andibalo/payd-test
```

Then open the project and run 
```bash
docker-compose up
```

This however will only run the backend as I had trouble dockerizing the sveltekit app and ran out of time. In order to run the frontend, we will need to run it locally. Ensure you have pnpm and node >= v22.15.1 installed.

```bash
pnpm install
pnpm run dev
```

## BE API

The postman collection json for this service can be found at /backend/postman. You can download it and import it on your local postman application and use the authroization static token RMS_ADMIN test the APIs.


## License

MIT