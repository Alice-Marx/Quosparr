[English](./README.md) · [简体中文](./README_ZH.md) · [繁體中文](./README_ZH_HANT.md) · [한국어](./README_KO.md) · 日本語 · [Русский](./README_RU.md) · [Français](./README_FR.md)

<div align="center">

<img src="Resources/AppIcon.png" alt="Quosparr logo" width="128" height="128">

# Quosparr

<sub><i>名前は <b>Quo</b>ta（クォータ）+ <b>Sparr</b>ow（スズメ）から</i></sub>

メニューバーにとまる小さなスズメが、Codex、Claude、Grok のクォータを見守ります。

Quosparr は macOS、Windows、Linux 向けのネイティブなメニューバー / システムトレイアプリです。Codex、Claude Code、Grok のサブスクリプションクォータ、reset credits、API 換算コスト、ローカルセッションを表示し、プロバイダーごとに独立したアカウント管理を備えています。3 社のグローバル Skill、MCP サーバー、ローカルセッションもまとめて管理できます。

</div>

> このリポジトリは [Licoy/Quosparr](https://github.com/Licoy/Quosparr) の [Alice-Marx によるフォーク](https://github.com/Alice-Marx/Quosparr) で、GNU Affero General Public License v3.0（[LICENSE](LICENSE)）に従います。Alice が加えた変更は、独立した **Grok Bot アシスタント**の週次クォータパネルと、独立した **Cursor** 使用量パネルです。どちらにも**公開されたサードパーティ課金 API はありません**。各パネルは自分のローカル JSON スナップショットだけを読みます。Cookie、Grok Bot の設定ページ、Cursor ダッシュボードの HTML は取得せず、このフォークにはどちらのスナップショットを書くルーチンも含まれません。

## ハイライト

**Codex**

- 5 時間、週間、追加のクォータウィンドウの残りに加え、現在のアカウント、サブスクリプション階層、有効期限を表示します。
- reset credits の数、状態、有効期限を表示します。
- 公式の週次使用率と日次 Credits から今週の枠を推定し、過去の週と比較して枠の引き下げを見つけやすくします。
- Codex リセット更新: 今日 Codex のレート制限がリセットされたかを確認でき、専用の更新間隔を設定できます（データは [Did Codex Reset](https://didcodexreset.com) から取得）。
- プロバイダーアカウント: ブラウザでのサインイン、ローカル `auth.json` の取り込み、token / JSON の貼り付け（`/auth/session` を含む）、ファイルの取り込み、API キーの追加に対応します。確認後に `~/.codex/auth.json` を原子的に書き込んで切り替えます。
- API アカウント: Base URL + API キー（Base URL が空なら公式 OpenAI API）で、サイトからモデル一覧を取得して選べます。切り替え時に `auth.json` と `config.toml` の `quosparr` プロバイダーを書き込み（事前にバックアップ）、new-api 系や OpenAI 互換の中継サイトの残高をアカウント一覧、パネル、ステータスバー、フローティングウィンドウ、ウィジェットに表示します。 `config.toml` で既に中継サイトを選んでいるキーログインは、そのサイトとモデルで自動的に取り込みます。 使用中のアカウントを編集するとキーを差し替えられ（空欄で維持）、設定ファイルをすぐ書き換えるか確認します。書き換えた後は Codex を再起動してください。
- 最近のセッション、プロジェクト、状態を表示し、ローカルのセッションインデックスを修復できます。

**Claude**

- Claude Code のサブスクリプション使用量として、5 時間枠、7 日枠、モデル別の週次枠、追加使用量を表示します。この値は公開された課金 API ではありません。
- Claude Code のローカルセッションログから token 使用量、API 換算コスト（`anthropic-builtin-2026-10-14`）、最近のセッション、出力速度を集計します。同じ応答は一度だけ数えます。未知のモデルに正確な費用は出しません。ログに先頭 token の時刻がないため、その列は「—」で、速度は端から端までの時間で計算します。
- サブスクリプションアカウントは CLI を先に実行せずブラウザでログインでき、現在のログイン、資格情報の貼り付け、ファイルからの取り込みも使えます。ブラウザログインはキャンセルとコールバックの貼り付けに対応し、アカウントライブラリにだけ追加します。Claude Code で使うには明示的に切り替えてください。Quosparr は現在のアカウントを更新しません。期限切れなら `claude` を一度実行します。他の管理アカウントはライブラリのコピーだけを更新します。
- macOS では公式資格情報はキーチェーンにあり、`/usr/bin/security` で読み書きします。Windows と Linux では `<config>/.credentials.json` です。切り替えは `claudeAiOauth` だけを置き換え、MCP トークンは残し、`.claude.json` の `oauthAccount` だけを更新します。実行中の Claude Code は終了せず、以前のログインを書き戻すことがあるため、切り替え前に警告します。
- API アカウント: Base URL + API キーで、サイトからモデル一覧を取得して Opus / Sonnet / Haiku / Fable のモデルを個別に割り当てられます。切り替え時は `ANTHROPIC_BASE_URL`、`ANTHROPIC_AUTH_TOKEN` と割り当てだけを `settings.json` に書き込み（事前にバックアップ）、サブスクリプションのログインはそのままです。サブスクリプションアカウントに戻すとこれらのキーを削除します。中継サイトの残高は Claude の使用量を表示するすべての場所に出ます。 `settings.json` に既に設定された中継構成は自動で取り込みます。 使用中のアカウントを編集するとキーを差し替えられ（空欄で維持）、`settings.json` をすぐ書き換えるか確認します。書き換えた後は Claude Code を再起動してください。

- **Claude Desktop のサブスクリプション**：両方のアカウント一覧から紐付けできます、初回は表示されたアカウントと組織に手動でログインして Quosparr で明示的に確認します、毎回作業を保存し、Desktop を終了してログイン・設定・ローカル状態を含む構成全体を保存・復元することへの同意が必要です、完了またはキャンセルで以前の構成に戻ります、連動切り替えは既定でオフで、現在の CLI アカウントも紐付け済み Desktop を単独で開けます、CLI と Desktop の結果は別々に表示し、起動成功は本人確認を意味しません。
- Desktop のネイティブ暗号化ログインを含むプロファイルは `~/.quosparr/claude-desktop` に非公開保存します、Windows は現在のユーザーだけの protected ACL を使用します、公式データパスを保持し、Quosparr は復号せずアカウントのエクスポートにも含めません、元の構成はアカウントを取り込まず別に保持し、ブラウザや他アプリのデータは読みません、プロファイルと Desktop データは同一ボリュームが必要で、中断後は明示的に復元できます、API アカウントは対象外で、macOS・Windows・Linux の実機および Desktop バージョン互換性の検証は未完了です。

**Grok**

- 公式 CLI の chat-proxy 課金 API から、含まれるクォータ、Build / Imagine / Chat の内訳、週 / 月の周期、プリペイド残高、オンデマンド使用量を表示します。
- 複数の OAuth / SuperGrok アカウントに対応し、隔離されたサインイン、現在のログインの取り込み、資格情報の貼り付け、更新、エイリアス、並べ替え、明示的な切り替えができます。
- ローカルの Grok CLI ログから最近のセッションを表示します。

**Grok Bot アシスタント**

- パネルのタブは **Grok Bot**（简体中文: **Grok Bot 助手**）で、常に表示され、**Grok** CLI のタブとは別です。Grok CLI のアカウント、セッション、課金画面ではありません。
- 使用率、リセット時刻、オンデマンドの要約を表示します。スナップショットに内訳があれば、割合の高いボットとアクティビティを表示します。
- 更新はスナップショットファイルの更新時刻を監視し、ヘッダーの更新もそのファイルだけを読み直します。スナップショットが無いときは、空の状態が取り込みパスを説明します。
- **JSON を取り込む**は、選んだファイルを検証してからスナップショットのパスへコピーします。トークンを含む未知のフィールドは捨てられます。組み込みの書き込み処理はありません。[`examples/grokbot-usage.json`](examples/grokbot-usage.json) をコピーするか、自分でファイルを置き換えてください。後から加えるルーチンは、同じディレクトリに一時ファイルを書いてから名前を変更して置き換えてください。

**Cursor 使用量**

- パネルのタブは **Cursor**（简体中文: **Cursor 额度**）で、常に表示され、**Grok Bot** と **Grok** CLI のタブとは別です。Cursor IDE / Cloud Agents / Cursor 製品の使用量であり、Grok Bot アシスタントの週次クォータでも xAI Grok CLI でもありません。
- 含まれる使用率、リセット時刻、オンデマンドの状態、任意のオンデマンド金額を表示します。`includedPools` があるとき、メインの表示は Cursor モデルとその他のモデルの 2 本の進捗です。そのオブジェクトが無く内訳があるときは、割合の高い製品を表示します。任意の `grokBotWeeklyRef` は Grok Bot アシスタントのパネルを指し、週次の数値はコピーしません。任意の `periodStart` と `history` は消費速度と使用量の見積もりに使います。
- 更新はスナップショットファイルの更新時刻を監視し、ヘッダーの更新もそのファイルだけを読み直します。スナップショットが無いときは、空の状態が取り込みパスと、公式 API がないことを説明します。
- **JSON を取り込む**は、選んだファイルを検証してからスナップショットのパスへコピーします。トークンを含む未知のフィールドは捨てられます。組み込みの書き込み処理はありません。[`examples/cursor-usage.json`](examples/cursor-usage.json) をコピーするか、自分でファイルを置き換えてください。後から加えるルーチンは、同じディレクトリに一時ファイルを書いてから名前を変更して置き換えてください。

**3 つのプロバイダー共通**

- 設定の「モデル価格」は「詳細設定」の直前で同じグループに配置；models.dev の公開カタログを既定で1時間ごとに自動同期し、無効化や15分〜7日のプリセット・カスタム間隔に対応、変更は即時保存し、自動同期が無効でも手動同期できます
- モデル価格は100万トークンあたりUSDの参考料金で、画像・動画など全モダリティの料金は対象外；検証済みのコンテキスト段階、キャッシュ書込TTL、Fast、地域料金を保持するため、カタログ参考値と適用値が異なる場合があり、詳細に出典とバージョンを表示、以下の日付付き内蔵表はオフラインの基準として保持します；OpenAI、Anthropic、xAI の元のプロバイダーIDに属する項目と内蔵料金だけを固定の Codex / Claude / Grok タブで表示し、初期は Codex、他社への切替はありません；第三者中継や他社の同名モデルは除外し、models.dev は集約公開カタログであり各社の公式リアルタイムAPIではありません；タブ右の同期設定ボタンからネイティブダイアログを開き、自動同期、間隔、詳細状態は3社共通で即時保存、ダイアログを閉じても同期は継続します
- API換算コストは過去のローカル使用量を現在の料金で再推定；合計、日別、プロジェクト、モデル、セッションの一部費用を明示し、不足する料金・ルールを全件表示、価格不明は —、明示的な無料は $0 とし、サブスクリプション枠と reset credits の意味は維持します
- 料金同期は保存済みのプロキシ設定で公開カタログを取得するだけで、セッションログを送信しません；オフラインでは内蔵料金と最後に保存成功したカタログを使い、同期エラーを表示します
- 常時表示の `Codex | Claude | Grok` タブがあります。macOS の設定「メニューバー項目」では、クォータバー、カウントダウン、バッテリー（横または縦）、クォータリング（値をリング内に表示可能）、進捗テキストの 5 つのスタイルをプレビューできます。メニューバーはメインパネルに従うか、3 つのプロバイダーのうち任意のものを個別に表示し、クォータバー・カウントダウン・バッテリーでは 2 つのプロバイダーを 1 項目に上下に重ねることもできます。保存済みの「Codex + Grok」はその組のままです。
- 任意で表示できる最前面のデスクトップフローティングウィンドウは、最大 3 つのプロバイダーの残り使用量をカプセル、カード、リング、バッテリーの 4 スタイルで表示します。クリックでメインパネルを開き、ドラッグで移動し、右クリックでその他のオプションを表示します。不透明度、位置の固定、フルスクリーンアプリの上に表示するかを設定できます。
- 今日、現在のサイクル、前のサイクル、今月、またはカスタム範囲の API 換算コストと token 使用量を表示します。
- 年初来の token 使用量チャートを、ヒートマップ、折れ線、棒グラフ（日次 / 週次 / 累計）で表示します。
- Skills ページでは各プロバイダーのグローバル skill について、出どころ、CLI に読み込まれるか、ほかのどの CLI と共有しているかを表示します。検索、`SKILL.md` の閲覧、フォルダを開く、CLI 自身の設定でのオン / オフ（書き込み前にバックアップ）、ほかの CLI へのコピーや移動（同名の skill は上書きしません）、Quosparr の Skill ゴミ箱への移動、リンク切れ、入れ子のフォルダ、名前の重複の診断ができ、多くはワンクリックで修復でき、複数の skill を選んでまとめて削除・修復することもできます。いずれもリスクを示して確認を取ってから実行します。設定の「Skill 管理」はプロバイダーごとのタブで状態や種類を絞り込め、詳細ページにインストール元、ファイル、frontmatter、GitHub リポジトリやホームページへのリンクを表示し、ゴミ箱では削除した skill を復元するかシステムのゴミ箱へ移せます。設定の「Skill 管理」では、公開 GitHub リポジトリ（URL または `owner/repo`）、ローカルフォルダ、zip から skill を追加することもできます。深い階層にある skill フォルダも探し、バージョン管理、ビルド成果物、システムファイルは取り込まず、インストールする skill と対象の CLI または共有フォルダ `~/.agents/skills` を選べます。同名のフォルダは上書きしません。GitHub から追加した skill は出どころを記録して更新を確認でき、更新するとリポジトリの最新版に置き換え、古い版は Skill ゴミ箱へ移します。
- 設定の「セッション管理」では、プロバイダーごとにローカルのすべてのセッションを一覧できます。タイトル、プロジェクト、ID で検索し、プロジェクトと期間で絞り込めます。詳細ページにはモデル、トークンの内訳、API 換算コスト、関連ファイル、会話のプレビューを表示します。選んだセッションは確認後に 1 件ずつシステムのゴミ箱へ移し、完全には削除しません。公式のインデックスやデータベースはそのままで、派生統計は残ったログから再計算します。
- 設定の「MCP 管理」では、プロバイダーごとにユーザーレベルのグローバル MCP サーバー（Codex と Grok は `config.toml`、Claude は `.claude.json`）を一覧できます。コマンド、サイト、構造は秘匿化した要約だけを表示し、env や headers の値は表示せず、サーバーへの接続や起動も行いません。複数選んでまとめて削除でき、削除前に `~/.quosparr/backups` へバックアップし、ほかの設定は保持します。変更は該当する CLI を再起動すると反映されます。Codex と Grok のサーバーには有効 / 無効スイッチもあり、同じく先にバックアップしてから保存します。Codex は項目の `enabled` を書き、Grok は公式 CLI と同じく `disabled_mcp_servers` も更新します。Claude Code は MCP の有効・無効をプロジェクトごとに記録するため、Claude のサーバーにはスイッチがありません。プロジェクト、プラグイン、企業管理の MCP は表示しません。
- ライト、ダーク、システムの外観に対応し、言語は English、简体中文、繁體中文、한국어、日本語、Русский、Français から選べます。
- アップデート機能を内蔵し、macOS 14+ ではデスクトップウィジェットも利用できます。

## スクリーンショット

<table>
  <tr>
    <th colspan="2">メインパネル</th>
    <th colspan="2">API 換算コスト</th>
    <th colspan="2">サブスク枠の推定</th>
  </tr>
  <tr>
    <td colspan="2" align="center" valign="top"><img src="docs/images/1-en.webp?v2" alt="Quosparr メインパネル" width="240"></td>
    <td colspan="2" align="center" valign="top"><img src="docs/images/3.webp?v1" alt="Quosparr API 換算コスト" width="240"></td>
    <td colspan="2" align="center" valign="top"><img src="docs/images/6.webp?v1" alt="Quosparr サブスク枠の推定" width="240"></td>
  </tr>
  <tr>
    <th colspan="3">プロバイダーアカウント</th>
    <th colspan="3">表示設定</th>
  </tr>
  <tr>
    <td colspan="3" align="center" valign="top"><img src="docs/images/4.webp?v1" alt="Quosparr プロバイダーアカウント" width="370"></td>
    <td colspan="3" align="center" valign="top"><img src="docs/images/5.webp?v1" alt="Quosparr 表示設定" width="370"></td>
  </tr>
  <tr>
    <th colspan="6">Windows の設定と Grok パネル</th>
  </tr>
  <tr>
    <td colspan="6" align="center" valign="top"><img src="docs/images/2.webp?v1" alt="Quosparr Windows の設定と Grok パネル" width="780"></td>
  </tr>
</table>

スクリーンショットは macOS と Windows 上の v1 の画面です。

## インストール

[GitHub Releases](https://github.com/Licoy/Quosparr/releases) からダウンロードし、お使いの OS と CPU に合うパッケージを選んでください。

| システム | パッケージ | インストール方法 |
| --- | --- | --- |
| macOS 12+（Apple Silicon / Intel） | `Quosparr-macos-arm64.dmg` / `Quosparr-macos-x86_64.dmg`、または同名の `.zip` | `Quosparr.app` を `Applications` にドラッグします |
| Windows 10 / 11（amd64 / arm64） | `Quosparr.Setup.<version>.<arch>.exe` | `%LOCALAPPDATA%\Programs\Quosparr` にユーザー単位でインストールします。管理者権限は不要です |
| Linux（amd64 / arm64） | `quosparr-<version>-linux-<arch>.tar.gz` と `install.sh`、または `quosparr_<version>_<arch>.deb` | tar.gz と同じ場所で `sh install.sh` を実行して `~/.local` にインストールするか、パッケージマネージャーで deb をインストールします |

各パッケージには対応するアーキテクチャ向けの SQLite ランタイムが同梱されているため、初回起動時に何かをダウンロードすることはありません。初回起動時、Quosparr はログイン時の起動を既定で有効にします。変更するには **コントロールパネル → 一般 → ログイン時に起動** を使います。

Homebrew でのインストールはまだ提供されていません。当面は DMG または ZIP を使用してください。

以前 CodexRunway を使っていた場合は、[アップグレードガイド](docs/upgrade-from-codexrunway_JA.md)を参照してください。

### macOS のセキュリティ警告

macOS 版は現在 ad-hoc signed で、公証されていません。開発元を確認できないと表示された場合は、`Quosparr.app` を右クリックして「開く」を選ぶか、システム設定 → プライバシーとセキュリティで「このまま開く」をクリックしてください。

アプリが壊れていると表示された場合は、多くの場合ダウンロード時の隔離属性が原因です。アプリを `Applications` に移動したあと、次を実行します。

```bash
xattr -dr com.apple.quarantine /Applications/Quosparr.app
```

### Linux の要件

Linux には GTK3、AppIndicator、トレイに対応したデスクトップ、およびプロキシ資格情報を保存するための動作する Secret Service セッション（GNOME Keyring や KWallet など）が必要です。deb はこれらの依存関係を宣言しており、パッケージマネージャー経由で更新されます。

## 動作要件

- この端末に Codex がインストールされ、使用されていることを推奨します。Codex のデータは `CODEX_HOME`（既定は `~/.codex`、Windows では `%USERPROFILE%\.codex`）から読み込まれます。
- サポートされるのは `file` 資格情報ストアのみです（Codex の `config.toml` で `cli_auth_credentials_store` が未設定の場合の既定）。`keyring`、`auto`、`ephemeral` は非対応として報告され、アカウントの取り込みと切り替えができなくなります。
- Claude パネルは `CLAUDE_CONFIG_DIR`、未設定なら `~/.claude` の Claude Code 設定を読みます。サブスクリプション使用量には Claude.ai ログインが必要で、API キーやゲートウェイでは取得できません。ブラウザログインは CLI のインストールや事前実行が不要で、明示的に切り替えるまではライブラリだけに保存します。既存の CLI ログインは自動で追加され、資格情報 JSON の貼り付けやファイルからの取り込みも可能です。実際のブラウザ認証と切り替えは、隔離アカウントでの受け入れ確認が必要です。
- Grok パネルには[公式 Grok CLI](https://docs.x.ai/build/overview) が必要です。

  ```bash
  curl -fsSL https://x.ai/cli/install.sh | bash
  ```

- Grok のアカウントとクォータは、OAuth / SuperGrok と互換のレガシーセッションのみに対応します。`grok login --oauth` を実行するか、アプリ内でサインインするか、現在のログインを取り込むか、資格情報 JSON を貼り付けてください。API キーのみのログインは管理アカウントとして追加されません。`GROK_HOME` が設定されていればそれを使い、設定されていなければ `~/.grok` を使います。
- Grok Bot アシスタントのパネルに Grok CLI は不要です。読むのは `~/.quosparr/grokbot-usage.json`（Windows では `%USERPROFILE%\.quosparr\grokbot-usage.json`）だけです。[`examples/grokbot-usage.json`](examples/grokbot-usage.json) をそのパスへコピーするか、そのパネルの **JSON を取り込む** を使ってください。
- Cursor 使用量パネルに Grok CLI は不要で、このアプリの中で Cursor にログインもしません。読むのは `~/.quosparr/cursor-usage.json`（Windows では `%USERPROFILE%\.quosparr\cursor-usage.json`）だけです。[`examples/cursor-usage.json`](examples/cursor-usage.json) をそのパスへコピーするか、そのパネルの **JSON を取り込む** を使ってください。必須項目だけのファイルは [`examples/cursor-usage.min.json`](examples/cursor-usage.min.json) です。

## デスクトップウィジェット

ウィジェットには macOS 14 以降が必要で、Windows と Linux には用意されていません。クォータ概要、token 推移、主要指標、Codex リセット更新の 4 種類があります。各ウィジェットは「ウィジェットを編集」で Codex、Claude、Grok、または一組（Codex + Claude、Codex + Grok、Grok + Claude）を選べます。すでに Codex + Grok のウィジェットはその組のままです。タップすると対応するページが開きます。

スナップショットは `~/.quosparr/widget-snapshot.json` にあり、権限は `0600` で、秘密情報を含まない派生データのみを保持します。メールアドレス、アカウント ID、token、認証 JSON、イベントの原文は含まれません。書き込むのはメインアプリだけで、ウィジェットは読み取りのみです。

## ネットワークプロキシ

**コントロールパネル → 一般 → ネットワーク** で、**システム設定に従う**、HTTP、SOCKS5 のいずれかを選択します。カスタムプロキシではユーザー名とパスワードを指定できます。この設定は、アプリが発行する Codex、Claude、Grok のリクエスト、CLI 操作、公開ステータスと価格のフィード、ソフトウェアアップデートに適用されます。システムのプロキシは変更されず、ブラウザのサインインページや実行中の外部 CLI はそれぞれの設定をそのまま使います。

変更は **プロキシ設定を保存** したあとにのみ反映されます。**接続をテスト** は、現在の下書きを使ってテストサイト（既定は `https://github.com`、変更可能）にアクセスするだけで、アカウント資格情報は送信せず、下書きも保存しません。プロキシのパスワードは macOS キーチェーン、Windows 資格情報マネージャー、Linux Secret Service に保存されます。設定が無効、資格情報が利用できない、またはシステムプロキシを解決できない場合は、直接接続へ切り替えません。設定済みの HTTP、SOCKS、またはシステムプロキシへの接続が拒否されたかタイムアウトしたときは、そのリクエストを一度だけ直接接続で再試行し、画面に示します。

## プライバシー

- Quosparr は独自のデータを `~/.quosparr` に保存し、アプリのキャッシュを `~/.codex` に書き込みません。既定で保存されるのは、秘密情報を含まない派生データのみです。
- プロバイダーアカウントの資格情報と上記の明示確認済み Desktop プロファイルは管理された例外です、資格情報は`~/.quosparr/accounts/<id>/auth.json` に保存されます（ディレクトリは `0700`、ファイルは `0600`、原子的に書き込み。Windows では現在のユーザーのみがアクセスできる ACL を使用）。アカウントインデックスには token は含まれません。
- 公式の `~/.codex/auth.json` に書き込まれるのは、OAuth のリフレッシュ時か、アカウント切り替えを確認したときだけです。現在のアカウント以外をリフレッシュした場合は、ライブラリ内のコピーだけが更新されます。
- Claude の資格情報は、macOS では `/usr/bin/security` でキーチェーンを読み、Windows と Linux では `<config>/.credentials.json` を読みます。管理コピーは `~/.quosparr/accounts/claude-<stable-id>/auth.json` にあり、`claudeAiOauth` だけを含みます。索引 `claude-index.json` にトークンはありません。切り替えは MCP トークンを残し、`.claude.json` の `oauthAccount` だけを置き換えます。Quosparr は現在のアカウントを更新しません。
- API アカウントのキーは管理コピー（`~/.quosparr/accounts/<id>/auth.json` または `claude-api-<id>/auth.json`）にのみ保存され、インデックスには Base URL、モデル、残高キャッシュだけが残ります。残高とモデルの取得は入力したサイトにのみ、プロキシ設定に従って送られます。変更する `config.toml` と `settings.json` は先に `~/.quosparr/backups` にバックアップされ、`settings.json` への初回書き込みではキーが名前順に並び替えられます。
- Grok の資格情報は `$GROK_HOME/auth.json` から読み込まれ、管理コピーは `~/.quosparr/accounts/grok-<stable-id>/auth.json` にあります。Grok の切り替えでは OAuth / 互換レガシーのスコープだけが置き換えられ、API キーと未知のスコープは残ります。保証されるのは新しいセッションのみで、実行中の Grok プロセスは終了されず、以前のアカウントを書き戻す可能性があるため、切り替え前にアプリが警告を表示します。
- セッションの修復が対象とするのは `session_index.jsonl` だけで、書き込み前にバックアップを作成します。セッションログをシステムのゴミ箱へ移すのは「セッション管理」で選択して確認したときだけで、移す前にファイルが変わっていないか 1 件ずつ再確認します。完全に削除することはなく、公式のインデックスやデータベースも書き換えません。
- Skill と MCP のオン / オフ、修復、削除は確認を取ってから実行します。設定ファイルは書き込み前に `~/.quosparr/backups` へバックアップし、削除した skill はまず復元できる Quosparr の Skill ゴミ箱へ移します。MCP の一覧にキーや env、headers の値は表示しません。
- skill の追加では、まず取り込み元のファイル一覧を確認し、skill のフォルダだけをメモリに読み込んでから、確認後に書き込みます。大きなリポジトリでもダウンロード時間が増えるだけでメモリは増えず、含まれるスクリプトは実行しません。GitHub からの取得はプロキシ設定に従い HTTPS で GitHub にのみ接続します。追加した skill の出どころ（リポジトリ、ref、コミット、リポジトリ内のフォルダ、またはローカルパス）は `~/.quosparr/skill-sources.json`（`0600`）に記録し、秘密情報は含みません。
- API 換算コストはローカルのセッションログから計算されます。セッションの内容はアップロードされません。
- token や API キーが、ログ、セルフチェックの出力、issue テンプレートに書き込まれることはありません。
- Grok Bot アシスタントのスナップショットは、秘密ではない使用量の統計だけを持ちます。取り込んだ文書にあるトークン、Cookie、API キーは捨てられ、ログにも書かれません。
- Cursor 使用量のスナップショットは別ファイルで、秘密ではない使用量の統計だけを持ちます。取り込んだ文書にあるトークン、Cookie、API キーは捨てられ、ログにも書かれません。
- Codex リセット更新とアップデートの確認は公開データのみを取得し、アカウント、token、セッションの内容は送信されません。

## データソース

- **Codex のクォータ / reset credits / 公式使用量**: ローカルの資格情報を使って、公式の ChatGPT / Codex バックエンド API から取得します。公式の使用量は現在のアカウントに属し、遅延や改訂が生じることがあります。
- **週次クォータの推定**: 非公式で、週次使用率と日次 Credits から外挿します（1000 Credits ≈ $40、バージョン `credits-usd-2026-08-26`）。`quota-estimate-history.json` には、派生した Credits と使用率のみが保存されます。
- **Codex の API 換算コスト**: ローカルのセッションログを、組み込みの価格表 `openai-builtin-2026-09-27` で計算します。対応する現在の料金は models.dev から取得し、検証済みの内蔵ルールを保持します。未知のモデルは正確なコストとして報告されません。
- **Claude Code のサブスクリプション使用量**: 文書化されていない OAuth usage エンドポイントで、現在の Claude Code サブスクリプションだけを表します。公式の課金 API とは書きません。ブラウザの Cookie や Console の API 請求は読みません、Desktop プロファイルは上記の独立した明示確認付き機能でのみ管理します。
- **Claude の API 換算コスト**: Claude Code のローカルセッションログと価格表 `anthropic-builtin-2026-10-14` で計算します。同じ応答はメッセージ ID で一度だけ数えます。キャッシュ書き込みは 5 分と 1 時間に分け、Fast 料金は別です。未知のモデルに正確な費用は出しません。日界はローカルタイムゾーンです。
- **Grok のクォータ**: 公式 CLI の chat-proxy `/v1/billing?format=credits` エンドポイントのみを使います。ブラウザの Cookie は使わず、ローカルログからクォータを推測せず、API 課金も混ぜません。
- **Grok Bot アシスタントの週次クォータ**: ローカル JSON スナップショットのみ、スキーマバージョン 1、パスは `~/.quosparr/grokbot-usage.json`（Windows は `%USERPROFILE%\.quosparr\grokbot-usage.json`）。公開されたサードパーティ課金 API はありません。例とフィールドは [`examples/grokbot-usage.json`](examples/grokbot-usage.json) にあります。`schemaVersion`、`fetchedAt`、`plan`、`usedPercent`、`resetAt`、`onDemand.status`、`periodLabel`、`breakdown[].bot`、`breakdown[].sharePercent`、`breakdown[].activities[].name`、`breakdown[].activities[].sharePercent`、`notes`。任意項目も schema 1 のままです。`periodStart`、`history[]`（`at`、`usedPercent`）、`cache`（`hitPercent`、`hits`、`misses`、`notes`）、`usage`（`inputTokens`、`outputTokens`、`cacheReadTokens`、`model`、`costUsd`）、`costLines[]`。パネルはこれらで消費速度、リセット、枠の見積もり、キャッシュヒット、API 相当コストを表示します。アシスタントには修復できるローカルセッション索引がありません。必須項目だけのファイルは [`examples/grokbot-usage.min.json`](examples/grokbot-usage.min.json) です。`schemaVersion` は `1` でなければなりません。`usedPercent` は 0 から 100 です。時刻は RFC3339 です。ファイルは秘密ではない使用量統計で、権限は `0600` です。パネルから取り込むか、例をそのパスへコピーしてください。置き換えはアトミックに行います。同じディレクトリに一時ファイルを書いてから名前を変更します。このアプリは Cookie、Grok Bot の設定ページ、Cursor ダッシュボードを読みません。
- **Cursor 製品の使用量**: ローカル JSON スナップショットのみ、スキーマバージョン 1、パスは `~/.quosparr/cursor-usage.json`（Windows は `%USERPROFILE%\.quosparr\cursor-usage.json`）。公開されたサードパーティ課金 API はありません。例とフィールドは [`examples/cursor-usage.json`](examples/cursor-usage.json) にあります。`schemaVersion`、`fetchedAt`、`plan`、`included.usedPercent`、`included.resetAt`、`included.periodLabel`、`includedPools.cursorModels.usedPercent`、`includedPools.cursorModels.label`、`includedPools.cursorModels.note`、`includedPools.otherModels.usedPercent`、`includedPools.otherModels.label`、`includedPools.otherModels.note`、`onDemand.status`、`onDemand.usedUsd`、`onDemand.limitUsd`、`breakdown[].name`、`breakdown[].sharePercent`、`grokBotWeeklyRef`、`notes`。任意項目も schema 1 のままです。`periodStart`、`history[]`（`at`、`usedPercent`）。`includedPools` が無いときは、全体の含まれる使用率と任意の内訳を表示します。どちらかのプールがあるときは、そのプールが独立した進捗行になります。`label` を省くとローカライズされた既定名（Cursor モデル / その他のモデル）を使います。Cursor モデルの `note` を省くと、ローカライズされた「Cursor Grok と Composer を含みます」を使います。`grokBotWeeklyRef` は Grok Bot アシスタントのパネルを指すだけで、週次の数値は保存しません。`included.periodLabel` は `this_week`、`this_month`、または `Included in Pro` のような短い表示名です。短いファイルは [`examples/cursor-usage.min.json`](examples/cursor-usage.min.json) で、必須項目とラベル無しの 2 つのプールです。`includedPools` を省いたスナップショットも有効です。`schemaVersion` は `1` でなければなりません。`included` は必須です。`usedPercent` と `sharePercent` は 0 から 100 です。存在する各プールの `usedPercent` も 0 から 100 です。`usedUsd` と `limitUsd` は null か 0 以上の数です。時刻は RFC3339 です。ファイルは秘密ではない使用量統計で、権限は `0600` です。パネルから取り込むか、例をそのパスへコピーしてください。置き換えはアトミックに行います。同じディレクトリに一時ファイルを書いてから名前を変更します。このアプリは Cookie、ダッシュボード HTML、Cursor コンソールを読みません。Grok Bot アシスタントの使用量は `grokbot-usage.json` のままです。
- **Grok の API 換算コスト**: ローカルの `turn_completed` 使用量を `xai-builtin-2026-09-27` で計算し、入力、キャッシュ、出力、長いコンテキストの料金を区別します。`costUsdTicks` はサブスクリプションの会計値であり、API 換算コストには使いません。
- **Codex リセット更新**: 公開フィード [didcodexreset.com/api/status.json](https://didcodexreset.com/api/status.json) から取得します。非公式で参考情報にすぎず、遅れることがあります。
- チャートの「公式統計（全デバイス）」は現在のアカウントに属し、「ローカルログ（全セッション）」はこの端末をスキャンするため、複数のアカウントにまたがることがあります。両者は一方がもう一方の部分集合というわけではなく、引き算してはいけません。どちらの Codex の日次値も UTC の日で集計されます。

## 開発

メインアプリは Go と [mygo](https://github.com/egoist/mygo) のネイティブ UI で書かれており、Go と mygo のバージョンは `go.mod` で固定されています。macOS のウィジェット拡張は引き続き Swift です。

```bash
go test ./...
go vet ./...
go build ./...
```

実際のアカウントや設定に影響しないよう、一時的な HOME と別の設定ドメインを使って実行します。

```bash
QUOSPARR_TEST_HOME="$(mktemp -d)"
HOME="$QUOSPARR_TEST_HOME" QUOSPARR_PREFS_DOMAIN=com.quosparr.app.v1-test go run .
HOME="$QUOSPARR_TEST_HOME" QUOSPARR_PREFS_DOMAIN=com.quosparr.app.v1-test go run . --self-check
```

セルフチェックはネットワークリクエストを行わず、資格情報も読み込みません。現時点では、バージョン、プラットフォーム、独自ディレクトリの権限、ロック、Codex の auth とセッションインデックスの有無を確認します。それ以外の診断は `not ported yet` と表示されます。

macOS の開発用パッケージをビルドします（アップデーターは無効。`ARCH=arm64` または `ARCH=x86_64` を任意で指定可能）。

```bash
V1_DEVELOPMENT_PACKAGE=1 bash Scripts/package-v1-app.sh
```

Go、Python 3、NSIS の `makensis` があるマシンで Windows amd64 の開発用パッケージをビルドします（アップデーターは無効）。出力ディレクトリは `build/v1-windows-amd64/windows-amd64/` で、`Quosparr Setup <version> amd64.exe` を含みます。

```bash
V1_DEVELOPMENT_PACKAGE=1 bash Scripts/package-v1-platform.sh windows/amd64
```

Windows / Linux のパッケージ化とリリース確認については [docs/v1-release.md](docs/v1-release.md)、貢献に関する注意は [CONTRIBUTORS.md](CONTRIBUTORS.md) を参照してください。

## コミュニティ

- [LinuxDO](https://linux.do/)

## ライセンス

このフォークは GNU Affero General Public License v3.0（[LICENSE](LICENSE)）に従い、[Licoy/Quosparr](https://github.com/Licoy/Quosparr) のフォークです。
