package app

// Grok Bot assistant strings. English, Simplified Chinese, Traditional
// Chinese, Korean, Japanese, Russian, French — the same order as
// l10n.Languages.

func init() {
	extraTexts[xGrokBotTab] = [7]string{
		"Grok Bot", "Grok Bot 助手", "Grok Bot 助手", "Grok Bot 도우미", "Grok Bot アシスタント", "Grok Bot", "Grok Bot",
	}
	extraTexts[xGrokBotTitle] = [7]string{
		"Weekly quota", "本周额度", "本週額度", "이번 주 할당량", "今週のクォータ", "Квота за неделю", "Quota de la semaine",
	}
	extraTexts[xGrokBotLocal] = [7]string{
		"Local snapshot", "本地快照", "本機快照", "로컬 스냅샷", "ローカルスナップショット", "Локальный снимок", "Instantané local",
	}
	extraTexts[xGrokBotEmptyTitle] = [7]string{
		"No usage snapshot yet", "还没有用量快照", "還沒有用量快照", "사용량 스냅샷이 없습니다", "使用量スナップショットがありません", "Снимка использования ещё нет", "Pas encore d’instantané d’usage",
	}
	extraTexts[xGrokBotEmptyDetail] = [7]string{
		"Put a JSON file at %@, or import one below. Grok Bot assistant has no public third-party billing API. This panel does not scrape cookies, Grok Bot settings pages, or Cursor dashboards.",
		"将 JSON 放到 %@，或使用下方的导入。Grok Bot 助手没有公开的第三方计费接口。此面板不会抓取 Cookie、Grok Bot 设置页或 Cursor 控制台。",
		"將 JSON 放到 %@，或使用下方的匯入。Grok Bot 助手沒有公開的第三方計費介面。此面板不會抓取 Cookie、Grok Bot 設定頁或 Cursor 控制台。",
		"JSON을 %@에 두거나 아래에서 가져오세요. Grok Bot 도우미에는 공개된 제3자 결제 API가 없습니다. 이 패널은 쿠키, Grok Bot 설정 페이지, Cursor 대시보드를 긁어오지 않습니다.",
		"JSON を %@ に置くか、下の読み込みを使ってください。Grok Bot アシスタントに公開の第三者課金 API はありません。このパネルは Cookie、Grok Bot の設定ページ、Cursor ダッシュボードを取得しません。",
		"Положите JSON в %@ или импортируйте его ниже. У ассистента Grok Bot нет публичного стороннего API биллинга. Эта панель не снимает cookie, страницы настроек Grok Bot и панели Cursor.",
		"Placez un JSON dans %@, ou importez-le ci-dessous. L’assistant Grok Bot n’a pas d’API de facturation tierce publique. Ce panneau ne récupère ni cookies, ni pages de réglages Grok Bot, ni tableaux de bord Cursor.",
	}
	extraTexts[xGrokBotImport] = [7]string{
		"Import JSON", "导入 JSON", "匯入 JSON", "JSON 가져오기", "JSON を読み込む", "Импорт JSON", "Importer le JSON",
	}
	extraTexts[xGrokBotUsed] = [7]string{
		"%@% used", "已用 %@%", "已用 %@%", "%@% 사용", "%@% 使用", "использовано %@%", "%@ % utilisé",
	}
	extraTexts[xGrokBotLeft] = [7]string{
		"%@% left", "剩余 %@%", "剩餘 %@%", "%@% 남음", "残り %@%", "осталось %@%", "%@ % restant",
	}
	extraTexts[xGrokBotOnDemand] = [7]string{
		"On-demand: %@", "按需：%@", "隨需：%@", "온디맨드: %@", "オンデマンド: %@", "По запросу: %@", "À la demande : %@",
	}
	extraTexts[xGrokBotOnDemandUnlimited] = [7]string{
		"Unlimited", "不限", "不限", "무제한", "無制限", "без ограничений", "Illimité",
	}
	extraTexts[xGrokBotOnDemandLimited] = [7]string{
		"Limited", "有限额", "有限額", "제한", "制限あり", "с лимитом", "Limité",
	}
	extraTexts[xGrokBotOnDemandOff] = [7]string{
		"Off", "关闭", "關閉", "꺼짐", "オフ", "выкл.", "Désactivé",
	}
	extraTexts[xGrokBotFetched] = [7]string{
		"Snapshot time %@", "快照时间 %@", "快照時間 %@", "스냅샷 시각 %@", "スナップショット時刻 %@", "Время снимка %@", "Heure de l’instantané %@",
	}
	extraTexts[xGrokBotBreakdown] = [7]string{
		"Breakdown", "用量构成", "用量構成", "구성", "内訳", "Разбивка", "Répartition",
	}
	extraTexts[xGrokBotInvalid] = [7]string{
		"Snapshot could not be read", "快照无法读取", "快照無法讀取", "스냅샷을 읽을 수 없습니다", "スナップショットを読めません", "Снимок не прочитан", "Instantané illisible",
	}
	extraTexts[xGrokBotImportFailed] = [7]string{
		"Import failed", "导入失败", "匯入失敗", "가져오기 실패", "読み込みに失敗しました", "Импорт не удался", "Échec de l’importation",
	}
	extraTexts[xGrokBotPeriodWeek] = [7]string{
		"This week", "本周", "本週", "이번 주", "今週", "Эта неделя", "Cette semaine",
	}
	extraTexts[xGrokBotSourceNote] = [7]string{
		"Local snapshot, not a public third-party billing API.",
		"本地快照，不是公开的第三方计费接口。",
		"本機快照，不是公開的第三方計費介面。",
		"로컬 스냅샷이며 공개된 제3자 결제 API가 아닙니다.",
		"ローカルスナップショットであり、公開の第三者課金 API ではありません。",
		"Локальный снимок, не публичный сторонний API биллинга.",
		"Instantané local, pas une API de facturation tierce publique.",
	}
	extraTexts[xGrokBotErrRead] = [7]string{
		"The snapshot file could not be read.", "无法读取快照文件。", "無法讀取快照檔案。", "스냅샷 파일을 읽을 수 없습니다.", "スナップショットファイルを読めません。", "Файл снимка не прочитан.", "Le fichier d’instantané est illisible.",
	}
	extraTexts[xGrokBotErrJSON] = [7]string{
		"The snapshot is not valid JSON.", "快照不是有效的 JSON。", "快照不是有效的 JSON。", "스냅샷이 올바른 JSON이 아닙니다.", "スナップショットは有効な JSON ではありません。", "Снимок не является корректным JSON.", "L’instantané n’est pas un JSON valide.",
	}
	extraTexts[xGrokBotErrSchema] = [7]string{
		"The snapshot schema version is not supported.", "不支持该快照的 schema 版本。", "不支援此快照的 schema 版本。", "지원하지 않는 스냅샷 스키마 버전입니다.", "このスナップショットの schema バージョンは未対応です。", "Версия схемы снимка не поддерживается.", "La version de schéma de l’instantané n’est pas prise en charge.",
	}
	extraTexts[xGrokBotErrPercent] = [7]string{
		"usedPercent must be from 0 to 100.", "usedPercent 必须在 0 到 100 之间。", "usedPercent 必須介於 0 到 100。", "usedPercent는 0에서 100 사이여야 합니다.", "usedPercent は 0 から 100 である必要があります。", "usedPercent должен быть от 0 до 100.", "usedPercent doit être entre 0 et 100.",
	}
	extraTexts[xGrokBotErrTime] = [7]string{
		"A timestamp must be RFC3339.", "时间必须是 RFC3339。", "時間必須是 RFC3339。", "시각은 RFC3339여야 합니다.", "時刻は RFC3339 である必要があります。", "Время должно быть в RFC3339.", "L’horodatage doit être en RFC3339.",
	}
	extraTexts[xGrokBotErrBreakdown] = [7]string{
		"The breakdown is invalid.", "用量构成无效。", "用量構成無效。", "구성이 올바르지 않습니다.", "内訳が無効です。", "Разбивка некорректна.", "La répartition est invalide.",
	}
	extraTexts[xGrokBotErrField] = [7]string{
		"A snapshot field is invalid.", "快照字段无效。", "快照欄位無效。", "스냅샷 필드가 올바르지 않습니다.", "スナップショットの項目が無効です。", "Поле снимка некорректно.", "Un champ de l’instantané est invalide.",
	}
	extraTexts[xGrokBotBurnTitle] = [7]string{"Burn rate", "速率分析", "速率分析", "소모 속도", "消費速度", "Скорость расхода", "Vitesse de consommation"}
	extraTexts[xGrokBotBurnEmpty] = [7]string{
		"Need periodStart, or at least two history samples, before a rate can be shown.",
		"需要 periodStart，或至少两个 history 采样，才能计算速率。",
		"需要 periodStart，或至少兩個 history 取樣，才能計算速率。",
		"속도를 보려면 periodStart 또는 history 샘플이 두 개 이상 필요합니다.",
		"速度を出すには periodStart か、history のサンプルが 2 点以上必要です。",
		"Для скорости нужны periodStart или хотя бы две точки history.",
		"Il faut periodStart, ou au moins deux échantillons history, pour afficher une vitesse.",
	}
	extraTexts[xGrokBotBurnRate] = [7]string{"%@% / day", "每天 %@%", "每天 %@%", "하루 %@%", "1 日 %@%", "%@% в день", "%@ % / jour"}
	extraTexts[xGrokBotBurnHistory] = [7]string{"From history samples.", "来自 history 采样。", "來自 history 取樣。", "history 샘플 기준입니다.", "history のサンプルから計算しています。", "По точкам history.", "D’après les échantillons history."}
	extraTexts[xGrokBotBurnCoarse] = [7]string{"Coarse rate from periodStart to the snapshot time.", "由 periodStart 到快照时间的粗算速率。", "由 periodStart 到快照時間的粗算速率。", "periodStart부터 스냅샷 시각까지의 대략적인 속도입니다.", "periodStart からスナップショット時刻までの概算です。", "Грубая оценка от periodStart до времени снимка.", "Estimation grossière de periodStart jusqu’à l’heure de l’instantané."}
	extraTexts[xGrokBotBurnETA] = [7]string{"Linear ETA to 100%: %@", "按当前速率到 100%：%@", "按目前速率到 100%：%@", "현재 속도면 100% 시각: %@", "この速度で 100% になる時刻: %@", "Линейный срок до 100%: %@", "ETA linéaire jusqu’à 100 % : %@"}
	extraTexts[xGrokBotBurnFull] = [7]string{"The snapshot is already at 100%.", "快照已到 100%。", "快照已到 100%。", "스냅샷이 이미 100%입니다.", "スナップショットはすでに 100% です。", "Снимок уже на 100%.", "L’instantané est déjà à 100 %."}
	extraTexts[xGrokBotBurnFlat] = [7]string{"Usage is not rising, so there is no ETA.", "用量没有上升，因此没有预计时间。", "用量沒有上升，因此沒有預計時間。", "사용량이 늘지 않아 예상 시각이 없습니다.", "使用量は増えていないため、到達時刻はありません。", "Расход не растёт, поэтому срока нет.", "L’usage n’augmente pas, donc pas d’échéance."}
	extraTexts[xGrokBotResetTitle] = [7]string{"Reset status", "重置状态", "重置狀態", "재설정 상태", "リセット状態", "Статус сброса", "État de réinitialisation"}
	extraTexts[xGrokBotResetMissing] = [7]string{"The snapshot has no reset time.", "快照未提供重置时间。", "快照未提供重置時間。", "스냅샷에 재설정 시각이 없습니다.", "スナップショットにリセット時刻がありません。", "В снимке нет времени сброса.", "L’instantané n’a pas d’heure de réinitialisation."}
	extraTexts[xGrokBotPeriodStart] = [7]string{"Window start %@", "窗口开始 %@", "視窗開始 %@", "구간 시작 %@", "期間の開始 %@", "Начало окна %@", "Début de la fenêtre %@"}
	extraTexts[xGrokBotPaceTitle] = [7]string{"Quota pace", "订阅额度推算", "訂閱額度推算", "구독 할당량 추정", "サブスクリプション枠の見積もり", "Оценка квоты подписки", "Projection du quota"}
	extraTexts[xGrokBotPaceMissing] = [7]string{
		"Need %@ before this projection.",
		"推算还需要：%@。",
		"推算還需要：%@。",
		"추정에 필요한 항목: %@.",
		"見積もりに必要な項目: %@。",
		"Для оценки не хватает: %@.",
		"Il manque pour la projection : %@.",
	}
	extraTexts[xGrokBotPaceProjected] = [7]string{"Projected %@% by reset", "按当前进度到重置时约 %@%", "按目前進度到重置時約 %@%", "현재 속도면 재설정 시 약 %@%", "このペースだとリセット時に約 %@%", "К сбросу около %@%", "Environ %@ % au moment de la réinitialisation"}
	extraTexts[xGrokBotPaceElapsed] = [7]string{"%@% of the window has elapsed.", "窗口已过 %@%。", "視窗已過 %@%。", "구간의 %@%가 지났습니다.", "期間の %@% が経過しています。", "Прошло %@% окна.", "%@ % de la fenêtre se sont écoulés."}
	extraTexts[xGrokBotPaceOver] = [7]string{"This pace would exhaust the quota early.", "按这个速度会提前用完额度。", "按這個速度會提前用完額度。", "이 속도면 할당량이 일찍 소진됩니다.", "このペースでは枠を早く使い切ります。", "В этом темпе квота закончится раньше.", "À ce rythme le quota sera épuisé trop tôt."}
	extraTexts[xGrokBotPaceUnder] = [7]string{"This pace stays under the window.", "按这个速度，重置前用不完。", "按這個速度，重置前用不完。", "이 속도면 재설정 전에 다 쓰지 않습니다.", "このペースならリセット前に使い切りません。", "В этом темпе квота сохранится до сброса.", "À ce rythme le quota tient jusqu’à la réinitialisation."}
	extraTexts[xGrokBotPaceEven] = [7]string{"This pace is close to the window.", "速度与窗口大致相当。", "速度與視窗大致相當。", "속도가 구간과 비슷합니다.", "ペースは期間とほぼ同じです。", "Темп близок к окну.", "Le rythme est proche de la fenêtre."}
	extraTexts[xGrokBotPaceAssumed] = [7]string{
		"No periodStart; the window is the seven days ending at resetAt.",
		"没有 periodStart，窗口按 resetAt 往前 7 天估算。",
		"沒有 periodStart，視窗按 resetAt 往前 7 天估算。",
		"periodStart가 없어 resetAt 기준 7일로 잡았습니다.",
		"periodStart がないため、resetAt から 7 日前を開始にしています。",
		"Нет periodStart: окно — семь дней до resetAt.",
		"Pas de periodStart : la fenêtre est les sept jours qui finissent à resetAt.",
	}
	extraTexts[xGrokBotCacheTitle] = [7]string{"Cache hits", "缓存命中分析", "快取命中分析", "캐시 적중", "キャッシュヒット", "Попадания в кэш", "Succès de cache"}
	extraTexts[xGrokBotCacheNA] = [7]string{
		"The assistant weekly quota has no cache-hit metric, unlike Claude Code and Grok CLI local sessions.",
		"助手周额度没有缓存命中指标（不同于 Claude Code / Grok CLI 本地会话）。",
		"助手週額度沒有快取命中指標（不同於 Claude Code / Grok CLI 本機工作階段）。",
		"도우미 주간 할당량에는 캐시 적중 지표가 없습니다. Claude Code / Grok CLI 로컬 세션과는 다릅니다.",
		"アシスタントの週次クォータにはキャッシュヒット指標がありません。Claude Code / Grok CLI のローカルセッションとは異なります。",
		"У недельной квоты ассистента нет метрики попаданий в кэш, в отличие от локальных сессий Claude Code и Grok CLI.",
		"Le quota hebdomadaire de l’assistant n’a pas d’indicateur de cache, contrairement aux sessions locales de Claude Code et de Grok CLI.",
	}
	extraTexts[xGrokBotCacheHit] = [7]string{"Hit rate %@%", "命中率 %@%", "命中率 %@%", "적중률 %@%", "ヒット率 %@%", "Доля попаданий %@%", "Taux de succès %@ %"}
	extraTexts[xGrokBotCacheCounts] = [7]string{"Hits %@ · Misses %@", "命中 %@ · 未命中 %@", "命中 %@ · 未命中 %@", "적중 %@ · 실패 %@", "ヒット %@ · ミス %@", "Попадания %@ · Промахи %@", "Succès %@ · Échecs %@"}
	extraTexts[xGrokBotCacheOptional] = [7]string{
		"Shown only because the snapshot included cache.",
		"因为快照带了 cache 字段才显示。",
		"因為快照帶了 cache 欄位才顯示。",
		"스냅샷에 cache 필드가 있어서 표시합니다.",
		"スナップショットに cache があるため表示しています。",
		"Показано, потому что в снимке есть cache.",
		"Affiché parce que l’instantané contient cache.",
	}
	extraTexts[xGrokBotCostTitle] = [7]string{"API-equivalent cost", "API等价成本", "API等價成本", "API 환산 비용", "API相当コスト", "Эквивалентная стоимость API", "Coût équivalent API"}
	extraTexts[xGrokBotCostEmpty] = [7]string{
		"This snapshot has no token or cost fields. The weekly quota is a subscription share, not an API bill.",
		"快照没有 token 或成本字段。每周额度是订阅份额，不是 API 账单。",
		"快照沒有 token 或成本欄位。每週額度是訂閱份額，不是 API 帳單。",
		"스냅샷에 토큰이나 비용 필드가 없습니다. 주간 할당량은 구독 몫이며 API 청구서가 아닙니다.",
		"スナップショットにトークンやコストの項目がありません。週次クォータはサブスクリプションの割合であり、API の請求ではありません。",
		"В снимке нет полей токенов или стоимости. Недельная квота — доля подписки, а не счёт API.",
		"Cet instantané n’a pas de champs de tokens ou de coût. Le quota hebdomadaire est une part d’abonnement, pas une facture d’API.",
	}
	extraTexts[xGrokBotCostUntitled] = [7]string{"Usage", "用量", "用量", "사용량", "使用量", "Использование", "Usage"}
	extraTexts[xGrokBotCostTokens] = [7]string{"In %@ · Out %@ · Cache read %@", "输入 %@ · 输出 %@ · 缓存读 %@", "輸入 %@ · 輸出 %@ · 快取讀 %@", "입력 %@ · 출력 %@ · 캐시 읽기 %@", "入力 %@ · 出力 %@ · キャッシュ読取り %@", "Вход %@ · Выход %@ · Чтение кэша %@", "Entrée %@ · Sortie %@ · Lecture cache %@"}
	extraTexts[xGrokBotCostMissingPrice] = [7]string{"No price", "缺价格", "缺價格", "가격 없음", "価格なし", "нет цены", "prix manquant"}
	extraTexts[xGrokBotRepairTitle] = [7]string{"Session repair", "会话修复", "工作階段修復", "세션 복구", "セッション修復", "Восстановление сессий", "Réparation des sessions"}
	extraTexts[xGrokBotRepairNA] = [7]string{
		"The assistant has no local session log to repair. Session repair applies to Codex, Claude, and Grok CLI.",
		"助手无本机会话日志可修复。会话修复仅适用于 Codex、Claude 与 Grok CLI。",
		"助手無本機工作階段日誌可修復。工作階段修復僅適用於 Codex、Claude 與 Grok CLI。",
		"도우미에는 고칠 로컬 세션 로그가 없습니다. 세션 복구는 Codex, Claude, Grok CLI에만 적용됩니다.",
		"アシスタントには修復できるローカルセッションログがありません。セッション修復は Codex、Claude、Grok CLI 向けです。",
		"У ассистента нет локального журнала сессий для восстановления. Оно относится только к Codex, Claude и Grok CLI.",
		"L’assistant n’a pas de journal de session local à réparer. La réparation concerne Codex, Claude et Grok CLI.",
	}
}

const (
	xGrokBotTab               textKey = "grokBotTab"
	xGrokBotTitle             textKey = "grokBotTitle"
	xGrokBotLocal             textKey = "grokBotLocal"
	xGrokBotEmptyTitle        textKey = "grokBotEmptyTitle"
	xGrokBotEmptyDetail       textKey = "grokBotEmptyDetail"
	xGrokBotImport            textKey = "grokBotImport"
	xGrokBotUsed              textKey = "grokBotUsed"
	xGrokBotLeft              textKey = "grokBotLeft"
	xGrokBotOnDemand          textKey = "grokBotOnDemand"
	xGrokBotOnDemandUnlimited textKey = "grokBotOnDemandUnlimited"
	xGrokBotOnDemandLimited   textKey = "grokBotOnDemandLimited"
	xGrokBotOnDemandOff       textKey = "grokBotOnDemandOff"
	xGrokBotFetched           textKey = "grokBotFetched"
	xGrokBotBreakdown         textKey = "grokBotBreakdown"
	xGrokBotInvalid           textKey = "grokBotInvalid"
	xGrokBotImportFailed      textKey = "grokBotImportFailed"
	xGrokBotPeriodWeek        textKey = "grokBotPeriodWeek"
	xGrokBotSourceNote        textKey = "grokBotSourceNote"
	xGrokBotErrRead           textKey = "grokBotErrRead"
	xGrokBotErrJSON           textKey = "grokBotErrJSON"
	xGrokBotErrSchema         textKey = "grokBotErrSchema"
	xGrokBotErrPercent        textKey = "grokBotErrPercent"
	xGrokBotErrTime           textKey = "grokBotErrTime"
	xGrokBotErrBreakdown      textKey = "grokBotErrBreakdown"
	xGrokBotErrField          textKey = "grokBotErrField"
	xGrokBotBurnTitle         textKey = "grokBotBurnTitle"
	xGrokBotBurnEmpty         textKey = "grokBotBurnEmpty"
	xGrokBotBurnRate          textKey = "grokBotBurnRate"
	xGrokBotBurnHistory       textKey = "grokBotBurnHistory"
	xGrokBotBurnCoarse        textKey = "grokBotBurnCoarse"
	xGrokBotBurnETA           textKey = "grokBotBurnETA"
	xGrokBotBurnFull          textKey = "grokBotBurnFull"
	xGrokBotBurnFlat          textKey = "grokBotBurnFlat"
	xGrokBotResetTitle        textKey = "grokBotResetTitle"
	xGrokBotResetMissing      textKey = "grokBotResetMissing"
	xGrokBotPeriodStart       textKey = "grokBotPeriodStart"
	xGrokBotPaceTitle         textKey = "grokBotPaceTitle"
	xGrokBotPaceMissing       textKey = "grokBotPaceMissing"
	xGrokBotPaceProjected     textKey = "grokBotPaceProjected"
	xGrokBotPaceElapsed       textKey = "grokBotPaceElapsed"
	xGrokBotPaceOver          textKey = "grokBotPaceOver"
	xGrokBotPaceUnder         textKey = "grokBotPaceUnder"
	xGrokBotPaceEven          textKey = "grokBotPaceEven"
	xGrokBotPaceAssumed       textKey = "grokBotPaceAssumed"
	xGrokBotCacheTitle        textKey = "grokBotCacheTitle"
	xGrokBotCacheNA           textKey = "grokBotCacheNA"
	xGrokBotCacheHit          textKey = "grokBotCacheHit"
	xGrokBotCacheCounts       textKey = "grokBotCacheCounts"
	xGrokBotCacheOptional     textKey = "grokBotCacheOptional"
	xGrokBotCostTitle         textKey = "grokBotCostTitle"
	xGrokBotCostEmpty         textKey = "grokBotCostEmpty"
	xGrokBotCostUntitled      textKey = "grokBotCostUntitled"
	xGrokBotCostTokens        textKey = "grokBotCostTokens"
	xGrokBotCostMissingPrice  textKey = "grokBotCostMissingPrice"
	xGrokBotRepairTitle       textKey = "grokBotRepairTitle"
	xGrokBotRepairNA          textKey = "grokBotRepairNA"
)
