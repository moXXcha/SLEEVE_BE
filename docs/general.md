# AI エージェント向け開発ガイドライン

## Ⅰ. 基本的な役割

あなたは、堅牢で強固なシステムを開発する AI エージェントです。以下の手順と規約を厳守し、品質の高いコードを効率的に生成してください。

---

## 🎯 タスク実行の基本原則

タスクを遂行する際は、以下の7つの基本原則を必ず守ってください：

1. **ドキュメント優先**: 全てのドキュメントを読み、理解してから開始
2. **承認必須**: 作業計画書とDB変更計画書は必ず承認を得る
3. **テスト駆動**: テストを先に書き、実装はテストを通すために行う
4. **段階的実装**: 小さな単位で実装し、こまめにコミット
5. **品質維持**: Linter、テスト、レビューを経て品質を保つ
6. **ドキュメント更新**: DB変更時は必ず `docs/db_schema.md` を更新
7. **レビュー待機**: PR作成後は他の作業に進まず、レビューを待つ

これらの原則は、高品質なコードを効率的に開発し、チーム全体の生産性を高めるために不可欠です。

---

## Ⅱ. 作業手順の構造

### 作業の2階層構造

SLEEVE プロジェクトの開発作業は、**2つのレベル**で構成されています：

```
┌─────────────────────────────────────────────────────────────┐
│ Level 1: タスク全体の作業手順（親タスク: Story/Task）          │
│  - Jira Story/Task に対応                                      │
│  - 1つのTopic Branch（または単一のFeature Branch）              │
│  - 複数のSubtaskをまとめる親レベルの作業                         │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│ Level 2: Subtask内の作業手順（各Subtask: Feature Branch）      │
│  - Jira Subtask に対応                                        │
│  - 1つのFeature Branch                                        │
│  - 実際のコーディング・テスト・PRを行う                          │
└─────────────────────────────────────────────────────────────┘
```

### 全体フロー図

```
[タスク開始]
     ↓
[Level 1: ドキュメント読解・作業計画]
     ↓
     ├─ 単一タスクの場合
     │   ↓
     │  [Feature Branch作成] (Level 1)
     │   ↓
     │  [Subtaskとして実装] (Level 2: 1つだけ)
     │   ↓
     │  [PR → main] (Level 1)
     │   ↓
     │  [完了]
     │
     └─ 複数Subtaskに分割する場合
         ↓
        [Topic Branch作成] (Level 1)
         ↓
        [全体テスト定義] (Level 1)
         ↓
        ┌────────────────────────┐
        │ Subtask 1 (Level 2)     │→ [PR → Topic]
        │ Subtask 2 (Level 2)     │→ [PR → Topic]
        │ Subtask 3 (Level 2)     │→ [PR → 依存先 → Topic]
        └────────────────────────┘
         ↓
        [全Subtask完了確認] (Level 1)
         ↓
        [PR: Topic → main] (Level 1)
         ↓
        [実装詳細書記載] (Level 1)
         ↓
        [完了]
```

---

## Ⅲ. Level 1: タスク全体の作業手順

このレベルは、**Jira の Story/Task に対応**します。複数のSubtaskを統括し、機能全体のゴールを管理します。

### 0. ドキュメントの完全読解

- **必須作業:** タスクを開始する前に、`docs/`ディレクトリ配下の全てのドキュメントを読み込み、プロジェクトの規約、アーキテクチャ、既存の決定事項を完全に理解してください。
- また、現ディレクトリの install されているモジュールなどを確認して重複して install が行われないようにしてください

### 1. 作業計画の策定と承認

- Markdown 形式で詳細な作業計画書を作成します。
- ファイル名は `docs/tasks/[JiraのID]_[タスク名]/task_plan.md` としてください。（例: `docs/tasks/PROJ-123_implement_user_auth/task_plan.md`）
- **テンプレート**: `docs/templates/task_plan.md` を参照してください
- 特に、Jira チケットの受け入れ条件を基にした**テストの通過基準**を明確かつ詳細に記述してください。これが実装のゴールとなります。
- **ブランチ戦略の考案**: 作業計画書内で、このタスクのブランチ戦略を定義してください。単一タスクで完結するか、複数のサブタスクに分割するかを判断し、分割する場合は各サブタスクの依存関係、ブランチ構造、マージ順序を明確に記述してください。詳細は`docs/github.md`を参照してください。
- **Jira Subtask内容の作成**（複数サブタスクに分割する場合）: `docs/tasks/[JiraのID]_[タスク名]/subtask_list.md`にSubtask内容（タイトル、説明、依存関係、受け入れ条件）を作成してください。テンプレートは `docs/templates/subtask_list.md` を参照してください。作業者がこれを基にJiraでSubtaskを手動作成します。作成後、Subtask IDを受け取り、作業計画書に反映してください。
- **メッセージ定義の考案**: 作業計画書内で、この機能で使用するエラーメッセージとログメッセージを考案し、`docs/messages/[feature]/errors.md`と`docs/messages/[feature]/logs.md`に記述してください。各メッセージの出力タイミング、関連関数、HTTPステータス（該当する場合）を明記してください。
- **DB変更の確認**: テーブルの作成・編集・削除が必要な場合は、以下を実施してください：
  1. `docs/db_schema.md` で現在のスキーマを確認
  2. `docs/DB_manual.md` の手順に従って変更計画書を作成
  3. 作業者の承認を得る
  4. **実装後は必ず `docs/db_schema.md` を更新**（必須）
- **作業者への確認:** 計画書を作成したら、作業者に内容の確認を促してください。承認が得られたら、次のステップに進みます。
- テスト通過基準は結合テストの通過基準と考えてください

### 2. ブランチ戦略の実行

作業計画書で定義したブランチ戦略に従い、親ブランチを作成します。

#### 2-1. 単一タスクの場合

単一タスク（Subtaskに分割しない）場合は、Feature Branchを作成し、そのまま**Level 2のSubtask作業手順**に進みます。

```bash
git checkout main
git pull origin main
git checkout -b feature/[JiraのID]_[機能名]
git push -u origin feature/[JiraのID]_[機能名]
```

**次のステップ**: **Level 2（Subtask内の作業手順）**に進んでください。

#### 2-2. 複数サブタスクに分割する場合

複数のSubtaskに分割する場合は、Topic Branchを作成します。

```bash
# Topic Branchを作成
git checkout main
git pull origin main
git checkout -b topic/[JiraのID]_[機能名]
git push -u origin topic/[JiraのID]_[機能名]
```

各SubtaskのFeature Branchは、**Level 2の実装時**に作成します。

### 3. 全体テストの定義（複数Subtaskの場合のみ）

複数Subtaskに分割した場合、作業計画書で定義した「テストの通過基準」と Jira チケットの「受け入れ条件」を基に、**結合テスト**のゴールを定義します。

#### 結合テストのゴール定義（Level 1で実施）

**このステップで行うこと:**
- Goの標準テスト機能を使用し、実装したい機能全体の動作テストの**ゴール（何をテストするか）**を定義します
- 作業計画書に明記してあるテスト通過基準を検証できるように、結合テストの**テスト項目**を決定します
- 結合テストは`tests/integration/[feature]/`ディレクトリ内に配置する予定であることを確認します

**このステップで行わないこと:**
- 結合テストの**詳細実装**（Goコードの作成）はまだ行いません
- 各Subtaskで単体テストと実装を進めながら、必要に応じて結合テストを追加していきます

**結合テストのゴール定義例:**
```
例: ユーザー認証機能の結合テスト
- ユーザー登録から認証トークン取得までの一連のフローが正常に動作すること
- 不正なメールアドレスでの登録が適切にエラーとなること
- 存在しないユーザーでのログインが適切にエラーとなること
- 認証トークンを使ったAPI呼び出しが正常に動作すること
```

### 4. 各Subtaskの実行

各Subtaskを順次実行します。各Subtaskでは、**Level 2（Subtask内の作業手順）**に従って作業を進めてください。

- 依存関係のないSubtaskは並行して作業可能
- 依存関係のあるSubtaskは、依存先の完了を待ってから開始

**Subtask実行中の結合テスト作成（Level 2で実施）:**

各Subtaskの実装中に、以下のタイミングで結合テストを**追加実装**してください：

1. **Subtaskの実装が完了し、単体テストが通過した時点**
   - そのSubtaskに関連する結合テストケースが明確になった場合
   - `tests/integration/[feature]/`ディレクトリに結合テストを追加実装
   - 実装した結合テストが通過することを確認

2. **複数のSubtaskの統合が完了した時点**
   - 複数の機能が連携して動作する結合テストケースを追加実装
   - 例: UserモデルとログインAPIの統合テスト

3. **全Subtask完了後（Level 1に戻った時点）**
   - 最終的な結合テストを実行し、全ての結合テストが通過することを確認
   - 詳細は「5. 全Subtask完了後の処理」を参照

**Subtask実行中の注意点:**
- 各SubtaskのPRは、作業計画書で定義したマージ先（Topic Branchまたは依存先Feature Branch）に作成してください
- 結合テストを追加した場合は、PR作成前に結合テストが通過することを確認してください

### 5. 全Subtask完了後の処理

全てのSubtaskが完了したら、以下の処理を行います。

#### 5-1. Topic Branch → main のPR作成（複数Subtaskの場合のみ）

全てのSubtaskがTopic Branchにマージされたら、Topic BranchからmainへのPRを作成します。

**GitHub CLI（`gh`）を使用してPR作成**:

```bash
gh pr create --base main --title "[JiraのID]: [タスク名]" --body-file docs/tasks/[JiraのID]_[タスク名]/topic_pr_description.md
```

**PR作成が失敗する場合**:

PR内容を`docs/tasks/[JiraのID]_[タスク名]/topic_pr_description.md`に作成し、作業者に手動作成を依頼してください。テンプレートは `docs/templates/topic_pr_description.md` を参照してください。

#### 5-2. 最終結合テストの実行

Topic Branch（または単一のFeature Branch）で、全体の結合テストを実行します。

```bash
# app/ディレクトリで全テストを実行
cd app
go test ./...

# tests/ディレクトリで結合テストを実行
cd ../tests
go test ./...
```

全てのテストが通過することを確認してください。

### 6. 実装詳細書の記載

全ての実装が完了したら、実装詳細書を記載してください。

実装詳細書は `docs/tasks/[JiraのID]_[タスク名]/implementation_details.md` として作成します。テンプレートは `docs/templates/implementation_details.md` を参照してください。

---

## Ⅳ. Level 2: Subtask内の作業手順

このレベルは、**Jira の Subtask に対応**します。1つのFeature Branchで、実際のコーディング・テスト・PRを行います。

**重要**: 単一タスクの場合も、このLevel 2の手順に従って作業を進めてください（Subtaskが1つだけと考えます）。

### 1. Feature Branchの作成

各Subtaskの実装を開始する際に、Feature Branchを作成します。

#### 1-1. 依存関係がない場合

Topic Branch（複数Subtaskの場合）またはmain（単一タスクの場合は既に作成済み）から分岐します。

```bash
# 複数Subtaskの場合
git checkout topic/[親JiraのID]_[機能名]
git pull origin topic/[親JiraのID]_[機能名]
git checkout -b feature/[SubtaskID]_[機能名]
git push -u origin feature/[SubtaskID]_[機能名]
```

#### 1-2. 依存関係がある場合

依存先のFeature Branchから分岐します。

```bash
# 依存先のFeature Branchから分岐
git checkout feature/[依存先SubtaskID]_[機能名]
git pull origin feature/[依存先SubtaskID]_[機能名]
git checkout -b feature/[SubtaskID]_[機能名]
git push -u origin feature/[SubtaskID]_[機能名]
```

### 2. Subtask用テストの作成

このSubtaskで実装する機能に対する**単体テスト**を作成します。

#### 単体テストの作成

- Go の標準テスト機能を使用し、UseCase 層や Repository 層のロジックを個別に検証するテストを作成します
- 単体テストはテストする関数を含むファイルの隣に`_test.go`として配置してください

**TDD（テスト駆動開発）の推奨**:
1. テストを先に書く（Red）
2. 最小限の実装でテストを通す（Green）
3. リファクタリング（Refactor）

詳細な手順は `docs/task_procedure.md` の「2.2.3 段階的実装手順」を参照してください。

### 3. 実装

承認された作業計画に基づき、テストが全て通過するように実装を行います。

#### 3-1. メッセージの実装

- 作業計画書の`docs/messages/[feature]/`で考案したメッセージを基に、`app/messages/[feature]/errors.go`と`app/messages/[feature]/logs.go`にGoの定数として実装してください
- 機能固有のメッセージは `app/messages/[feature]/` に配置し、複数機能で共通で使用するメッセージは `app/messages/common/` に配置します
- コードからはこの定義を参照するように実装してください

#### 3-2. コーディング

- `docs/BE_coding_roles.md` のコーディング規約を厳守してください
- `docs/coding_guide.md` のコーディングスタイルを厳守してください
- 実装の詳細な手順は `docs/task_procedure.md` の「2. 実装実行」を参照してください

#### 3-3. Linterの実行

作成が終わったら、`app/`内で`golangci-lint run --fix`を実行して、linter の実行と format を行ってください。

```bash
cd app
golangci-lint run --fix
```

### 4. テストとデバッグ

作成した単体テストを実行し、全てのテストが通過することを確認します。

```bash
# 単体テストの実行
go test ./...
```

#### テスト失敗時（実装中のデバッグ）

1. **調査:** エラーの原因を調査します。
2. **デバッグ調査書の作成:**
   - `docs/tasks/[JiraのID]_[タスク名]/debug.md` というファイル名で、調査内容とデバッグ方針を記録します。
   - （例: `docs/tasks/PROJ-123_implement_user_auth/debug.md`）
   - テンプレートは `docs/templates/debug.md` を参照してください。
3. **作業者への確認:** デバッグ調査書を作成後、作業者に内容の確認を促し、承認を得てからデバッグ作業を行ってください。
4. **繰り返し:** 全てのテストが成功するまで、このプロセスを繰り返します。

詳細な修正手順は `docs/task_procedure.md` の「2.2.7 テスト失敗・実装不具合時の修正手順」を参照してください。

### 5. PR作成前の準備

**重要**: Subtaskの実装とテストが完了したら、PR作成前に以下の手順を必ず実行してください。

#### 5.1 最終確認

```bash
# 1. 変更されたファイルを確認（commit漏れがないかチェック）
git status

# 2. 全てのテストが通過することを確認
go test ./...

# 3. Linterエラーがないことを確認
cd app
golangci-lint run --fix
cd ..
```

#### 5.2 コミットとプッシュ

```bash
# 1. 未コミットの変更がある場合はコミット
git add [ファイルパス]
git commit -m "[適切なコミットメッセージ]"

# 2. リモートにプッシュ
git push origin feature/[SubtaskID]_[機能名]

# 3. プッシュ後、commit漏れがないことを最終確認
git status  # "nothing to commit, working tree clean" と表示されればOK
```

**チェックリスト**:
- [ ] `git status` で未コミットの変更がない
- [ ] 全ての単体テストが通過
- [ ] `golangci-lint run --fix` でエラーがない
- [ ] 全ての変更がリモートにプッシュされている
- [ ] コミットメッセージが規約に従っている
- [ ] DB変更を行った場合は `docs/db_schema.md` を更新している

### 6. PR作成

PR作成前の準備が完了したら、その時点でPRを作成します。

**GitHub CLI（`gh`）を使用してPR作成**:

```bash
# 依存関係がない場合（topic branchへのPR）
gh pr create --base topic/[親ID]_[機能名] --title "[SubtaskID]: [タスク名]" --body-file docs/tasks/[JiraのID]_[タスク名]/pr_description_[SubtaskID].md

# 依存関係がある場合（依存先feature branchへのPR）
gh pr create --base feature/[依存先SubtaskID]_[機能名] --title "[SubtaskID]: [タスク名]" --body-file docs/tasks/[JiraのID]_[タスク名]/pr_description_[SubtaskID].md

# 単一タスクの場合（mainへのPR）
gh pr create --base main --title "[JiraのID]: [タスク名]" --body-file docs/tasks/[JiraのID]_[タスク名]/pr_description.md
```

**PR作成が失敗する場合**:

PR内容を`docs/tasks/[JiraのID]_[タスク名]/pr_description_[SubtaskID].md`に作成し、作業者に手動作成を依頼してください。テンプレートは `.github/pull_request_template.md` を参照してください。

#### 5.1 レビュー待機フロー（現状の運用）

**⚠️ 注意：以下は現状の運用です。メンバーが増えた際は変更される可能性があります。**

PR作成後は、以下のフローに従ってください：

1. **PR作成後**
   - オーナーに通知
   - レビュー待ちの状態で待機

2. **❌ レビュー待ちの間、他のSubtaskに進まない**
   - PR作成後は、オーナーのレビュー・承認・マージを完全に待つ
   - レビュー中に他の作業を進めない

3. **修正依頼があった場合**
   - 指摘された箇所を修正
   - 追加のコミットを作成
   - プッシュして再レビュー依頼

4. **承認された場合**
   - オーナーがマージを実行
   - マージ完了を確認

**理由:**
- 現状は少人数開発のため、並行作業によるコンフリクトを避ける
- レビューでの修正指示を即座に反映できるようにする
- ブランチ管理をシンプルに保つ

**将来の変更予定:**
- メンバーが増えた際は、依存関係のないSubtaskを並行開発できるように変更する可能性があります

詳細は `docs/github.md` の「レビュー待機フロー（現状の運用）」を参照してください。

### 7. マージ後の処理

PRがレビューされ、承認され、オーナーがマージを実行したら、以下の処理を行います。

#### 7.1 Feature Branchの削除

```bash
# ローカルブランチを削除
git branch -d feature/[SubtaskID]_[機能名]

# リモートブランチを削除（GitHubのPRマージ時に自動削除される場合もある）
git push origin --delete feature/[SubtaskID]_[機能名]
```

#### 7.2 Topic/Main Branchの最新化

**重要:** マージ後、必ず親ブランチ（TopicまたはMain）を最新化してください。

```bash
# Topic Branchを最新化（複数Subtaskの場合）
git checkout topic/[親ID]_[機能名]
git pull origin topic/[親ID]_[機能名]

# Main Branchを最新化（単一タスクの場合）
git checkout main
git pull origin main
```

#### 7.3 次のSubtaskへの準備

次のSubtaskに進む前に、以下を確認してください：

- 前のSubtaskのFeature Branchが削除されている
- Topic/Main Branchが最新化されている
- 作業ディレクトリがクリーンな状態になっている

#### 7.4 次のステップ

- **次のSubtaskがある場合**: Level 2のステップ1（Feature Branchの作成）に戻る
- **全Subtask完了の場合**: **Level 1のステップ5**に戻る

---

## Ⅴ. 関連ドキュメント

詳細な実装手順とテンプレートについては、以下のドキュメントを参照してください：

- **実装手順**: `docs/task_procedure.md` - 詳細な実装手順とテンプレート
- **コーディング規約**: `docs/BE_coding_roles.md` - アーキテクチャとコーディング規約
- **コーディングスタイル**: `docs/coding_guide.md` - コーディングスタイルガイド
- **Git 運用**: `docs/github.md` - Git Flow 運用とコミット規約
- **DB 操作**: `docs/DB_manual.md` - データベース操作手順
- **Linter/Formatter**: `docs/lint_format_manual.md` - コード品質管理
