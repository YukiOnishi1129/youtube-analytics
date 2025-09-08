# Backend Architecture — Final Plan (Clean Architecture + DDD)

## 1) 全体方針
- Clean Architecture（依存は内側へ）
  - domain（Entity/VO/Domain Service）
  - port/input（UseCase 公開 I/F）
  - port/output/gateway（DB/外部 API 抽象）
  - port/output/presenter（出力境界）
  - usecase（Interactor: Input Port 実装）
  - adapter（controller/presenter/gateway の具象）
  - driver（起動・設定・接続・セキュリティ・観測）
  - cmd（Composition Root）
- DDD：集約＝keyword / video / metric / channel / account
- 1集約 = 1パッケージ、ファイルは“概念ごと”（例：keyword.go, keyword_service.go）

## 2) 認証・認可（Identity Platform + Cloud Run）
- ユーザー API：各サービスが ID トークンを自前検証（go-oidc, JWKS キャッシュ）
- 内部 API（Cloud Tasks/Scheduler）：Cloud Run IAM のみ許可 + OIDC 二重検証
- メソッドポリシー：PUBLIC / USER_ID_TOKEN / SERVICE_OIDC（Interceptor で制御）
- 共通化：`services/pkg/identityauth` に JWKS, Verify, gRPC Interceptor を実装

## 3) ディレクトリ構成（モノレポ / go.work）

```
youtube-analytics/
├─ proto/                              # .proto（buf 管理）
├─ services/
│   ├─ go.work
│   ├─ pkg/
│   │   ├─ identityauth/
│   │   └─ pb/
│   ├─ ingestion-service/
│   ├─ analytics-service/
│   └─ authority-service/
└─ web/
```

同名 package（keyword 等）が層ごとに存在しても OK。import 時は alias で区別（例：domKeyword, ucKeyword, inKeyword）。

## 4) proto & 生成物
- buf generate
- Go → `services/pkg/pb`
- TS → `web/client/src/external/client/grpc`

## 5) Docker / CI
- サービスごとに Dockerfile（prod/dev）
- Cloud Run：minInstances=0（必要に応じ 1）、未認証呼び出し不可、Secrets は SM→env
- GitHub Actions：サービス単位で build/deploy（行列化可）

## 6) テスト戦略
- Unit（domain） / UseCase（port 差替） / Component（handler→usecase→repo）
- Contract（buf breaking） / Service E2E（少数） / System E2E（staging のスモーク）
- 冪等性：同一入力の複数回実行で 1 回分の効果のみ

## 7) バッチ & 冪等性
- Cloud Tasks → `/snapshot`（ingestion）
- TaskID = `snap:{videoId}:{cp}` / DB は UNIQUE + upsert

## 8) ドメイン（抜粋）
- Keyword：name, filterType, pattern, enabled / PatternBuilder.Build()
- Channel + ChannelSnapshot（登録者推移）
- Video + VideoSnapshot（0/3/6/12/24/48/72/168h）
- Metric（CP 別の読み取りモデル：growth/ratio/Wilson/LPS/Exclude）
- History（TopN の凍結）
- Account / Identity / Role（email 一意, identity 重複禁止, 非アクティブはログイン不可）

## 9) 実装ルール（読みやすさ優先）
- 1 ドメイン=1 パッケージ、ファイルは概念で分割
- 不変条件はコンストラクタ、状態遷移はメソッド
- 集約間参照は ID で繋ぐ
- Repository 抽象は `port/output/gateway`（UseCase 所有）
- Import alias で同名パッケージの混同を避ける

## 10) 設定値（最低限）
- OIDC：`IDP_ISSUER`, `IDP_AUDIENCE`, `JWKS_CACHE_TTL`
- 内部 API：`TASKS_AUDIENCE`, `INTERNAL_SECRET`
- 指標：`LIKES_PER_SUBSCRIPTION_SCALE=1000`, `LIKES_PER_SUBSCRIPTION_OFFSET=500`

