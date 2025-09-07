# Authority API (authority/v1/authority.proto)

```protobuf
syntax = "proto3";

package authority.v1;

option go_package = "github.com/yourorg/youtube-analytics/proto/authority/v1;authorityv1";

message Identity {
  string provider = 1;       // 'google' | 'password' | 'github' ...
  string provider_uid = 2;   // arbitrary
}

message Account {
  string id = 1;             // uuid v7
  string email = 2;
  bool   email_verified = 3;
  string display_name = 4;
  string photo_url = 5;
  repeated Identity identities = 6;
}

message GetMeRequest {}
message GetMeResponse { 
  Account account = 1; 
}

service AuthorityService {
  rpc GetMe (GetMeRequest) returns (GetMeResponse);
}
```

## Notes

- **Date format**: Uses RFC3339 string throughout (easy to handle between frontend/backend)
- **pattern field**: Intended to be shown only in advanced mode (hidden in normal UI)
- **hide_low_sample**: Server-side default set to true for stable UX
- **Code generation**: Generate code with buf/protoc to implement handlers for each service
