# Hotfix: Multi-User Flow & Production Readiness

Read ALL existing Go source files first. Then fix these critical issues.

## 1. Team Member Management (CRITICAL)
Add these APIs:
- `POST /api/v1/teams/:id/members` - Join team (any user can join open teams)
  - Body: `{"user_id": "xxx"}` or no body (join self)
  - Auto-add to team_members table
- `GET /api/v1/teams/:id/members` - List team members
- `DELETE /api/v1/teams/:id/members/:user_id` - Remove member (team admin/creator only)

Need: model, store method, service method, handler, route registration.

## 2. Channel Member Management (CRITICAL)
Add these APIs:
- `POST /api/v1/channels/:id/members` - Join channel
  - Must be team member first
  - Auto-add to channel_members table
- `GET /api/v1/channels/:id/members` - List channel members  
- `DELETE /api/v1/channels/:id/members/:user_id` - Leave/remove

The channel_members table already exists in migration 001. Use it.

## 3. Auto-Join on Create
- When creating a team → auto-add creator as team member with role "admin"
- When creating a channel → auto-add creator as channel member
- This may already be partially implemented. Check and fix.

## 4. Channel Access Fix
- Currently channel operations check membership. This is correct.
- But the join flow doesn't exist. After adding join API, the flow works:
  1. User registers
  2. User joins team
  3. User joins channel
  4. User can send/read messages

## 5. Admin Role
- Team creator should be admin of that team
- Add `PUT /api/v1/users/:id/role` - Set user role (admin only, for system admin)
- Or: first registered user becomes system admin
- Audit log should be accessible to team admins, not just system admins

## 6. Open Team Auto-Join
For better UX, when listing channels of a team, if the team is "open", auto-join the user.
OR: Add a `POST /api/v1/teams/:id/join` endpoint that's simpler.

## 7. Frontend Fixes
In `web/src/`, update:
- Add team join flow in the UI (join button)
- Add channel join flow  
- Show member list
- Handle "not a member" errors gracefully

## 8. Rebuild & Restart
After all fixes:
```bash
go build -o nexus-mm ./cmd/server/
sudo systemctl restart nexus-mm
cd web && npm run build
```

## 9. E2E Verification
Run this test to verify multi-user flow:
```bash
BASE="http://localhost:18065/api/v1"
# Register alice, bob
# Alice creates team
# Bob joins team  
# Alice creates channel
# Bob joins channel
# Both send messages
# Both react
# Thread works cross-user
# Search works
# All should return 200
```

When completely finished:
1. Verify go build passes
2. Verify npm run build passes  
3. Restart service: sudo systemctl restart nexus-mm
4. Run the E2E test above
5. git add, commit, push
6. Run: openclaw system event --text "Done: Nexus-MM hotfix - team/channel member management, multi-user flow fixed, E2E verified" --mode now
