# Chapter 6: Frontend Contribution Areas

## Horizon Dashboard Stack

Horizon's dashboard uses frontend assets in the package, commonly involving:

- Vue components
- JavaScript API calls
- SCSS/CSS
- Blade views
- compiled assets

Exact versions can change, so inspect:

```bash
cat package.json
ls resources/js
ls resources/sass
ls resources/views
```

---

## Good Frontend Contribution Areas

| Area | Difficulty | Examples |
|---|---|---|
| empty states | beginner | better text when no failed jobs |
| loading states | beginner | spinners, disabled controls |
| copy clarity | beginner | clearer labels |
| accessibility | intermediate | keyboard focus, aria labels |
| responsive layout | intermediate | mobile dashboard issues |
| API error handling | intermediate | failed polling request state |
| dashboard performance | advanced | polling and rendering large lists |

---

## Dashboard API Flow

Frontend usually follows:

```text
Vue component
└── HTTP request to Horizon route
    └── controller/API endpoint
        └── repository reads Redis
            └── JSON response
```

If dashboard data is wrong, check both:

- backend JSON response
- frontend rendering logic

---

## Frontend PR Checklist

- reproduce the UI issue
- inspect API response in browser devtools
- identify whether bug is backend or frontend
- make minimal UI change
- build assets if repo expects committed assets
- include screenshot or explanation in PR

---

## Useful UI Improvements

Good areas for first dashboard PR:

- improve empty dashboard states
- add clearer labels for paused supervisors
- make failed-job error display easier to read
- fix small responsive layout problem
- improve keyboard focus on dashboard controls
- improve loading behavior when polling fails

---

## Avoid at First

Avoid starting with:

- total dashboard redesign
- changing build system
- replacing Vue version
- changing routes/API shape
- adding large new feature

These are high-friction PRs.

