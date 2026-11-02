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
)
