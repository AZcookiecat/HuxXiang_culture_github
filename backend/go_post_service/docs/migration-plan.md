# Post System Zero-Downtime Migration Plan

## Current Schema

- `users`: `id`, `username`, `email`, `password_hash`, `avatar`, `bio`, `role`, `is_active`, timestamps.
- `community_posts`: `id`, `title`, `content`, `author_id`, `category`, `status`, `view_count`, `like_count`, `comment_count`, timestamps.
- `comments`: `id`, `content`, `author_id`, `post_id`, `parent_id`, timestamps.
- `user_post_likes`: unique pair of `user_id` and `post_id`.

## Compatibility Constraints

- Keep routes under `/api/community`.
- Preserve JSON keys, HTTP methods, and status codes used by Flask today.
- Accept the same JWT bearer token signed by `JWT_SECRET_KEY`.

## Rollout

1. Deploy Gin service in shadow mode against the same database.
2. Mirror read traffic and compare bodies/status codes with golden samples.
3. Route 1%, 10%, 50%, then 100% traffic through the Gin canary.
4. Keep Flask hot and writable until three successful full-link pressure rounds pass.

## Rollback

- Remove the canary route or upstream weight and send all traffic back to Flask.
- Flush post-list caches.
- Leave the shared schema unchanged; this migration is code-first, not schema-first.

## Observability

- Scrape `/metrics` for latency, request count, and breaker state.
- Alert on `5xx`, breaker open rate, DB pool exhaustion, and P99 latency regression.
