# Authority API (authority/v1/authority.proto)

```protobuf
syntax = "proto3";

package authority.v1;

option go_package = "github.com/yourorg/youtube-analytics/proto/authority/v1;authorityv1";

message Account {
  string id = 1;             // uuid v7
  string email = 2;
  bool   email_verified = 3;
  string display_name = 4;
  string photo_url = 5;
  bool   is_active = 6;
  string last_login_at = 7;  // RFC3339
  repeated string roles = 8; // ["user"], extendable
}

// GetAccount
message GetAccountRequest {}
message GetAccountResponse { Account account = 1; }

// SignUp
message SignUpRequest { string email = 1; string password = 2; string display_name = 3; }
message SignUpResponse { Account account = 1; string id_token = 2; string refresh_token = 3; }

// SignIn
message SignInRequest { string email = 1; string password = 2; }
message SignInResponse { string id_token = 1; string refresh_token = 2; int64 expires_in = 3; }

// SignOut
message SignOutRequest { string refresh_token = 1; }
message SignOutResponse { bool success = 1; }

// ResetPassword
message ResetPasswordRequest { string email = 1; }
message ResetPasswordResponse { bool email_sent = 1; }

service AuthorityService {
  rpc GetAccount (GetAccountRequest) returns (GetAccountResponse);
  rpc SignUp (SignUpRequest) returns (SignUpResponse);
  rpc SignIn (SignInRequest) returns (SignInResponse);
  rpc SignOut (SignOutRequest) returns (SignOutResponse);
  rpc ResetPassword (ResetPasswordRequest) returns (ResetPasswordResponse);
}
```

## Notes

- Roles: `Account.roles` returns assigned role names; MVP defaults to ["user"].
- Identity management endpoints (LinkIdentity, AssignRole) are deferred for MVP.
- Use Identity Platform (e.g., Firebase) for authentication flows; this service acts as a facade for account retrieval and basic ID flows.
