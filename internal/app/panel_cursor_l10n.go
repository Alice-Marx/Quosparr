package app

// Cursor product usage strings. English, Simplified Chinese, Traditional
// Chinese, Korean, Japanese, Russian, French — the same order as
// l10n.Languages. Shared burn, reset, and error sentences stay on the
// Grok Bot keys; these strings name Cursor so the tab cannot be read as
// the assistant weekly quota.

func init() {
	extraTexts[xCursorTab] = [7]string{
		"Cursor", "Cursor 额度", "Cursor 額度", "Cursor", "Cursor", "Cursor", "Cursor",
	}
	extraTexts[xCursorTitle] = [7]string{
		"Cursor usage", "Cursor 用量", "Cursor 用量", "Cursor 사용량", "Cursor 使用量", "Использование Cursor", "Usage Cursor",
	}
	extraTexts[xCursorLocal] = [7]string{
		"Local snapshot", "本地快照", "本機快照", "로컬 스냅샷", "ローカルスナップショット", "Локальный снимок", "Instantané local",
	}
	extraTexts[xCursorEmptyTitle] = [7]string{
		"No Cursor usage snapshot yet", "还没有 Cursor 额度快照", "還沒有 Cursor 額度快照", "Cursor 사용량 스냅샷이 없습니다", "Cursor の使用量スナップショットがありません", "Снимка использования Cursor ещё нет", "Pas encore d’instantané d’usage Cursor",
	}
	extraTexts[xCursorEmptyDetail] = [7]string{
		"Put a JSON file at %@, or import one below. Cursor product usage has no public third-party billing API. This panel does not scrape cookies, dashboard HTML, or the Cursor console. Grok Bot assistant usage stays in grokbot-usage.json.",
		"将 JSON 放到 %@，或使用下方的导入。Cursor 产品用量没有公开的第三方计费接口。此面板不会抓取 Cookie、仪表盘 HTML 或 Cursor 控制台。Grok Bot 助手用量仍在 grokbot-usage.json。",
		"將 JSON 放到 %@，或使用下方的匯入。Cursor 產品用量沒有公開的第三方計費介面。此面板不會抓取 Cookie、儀表板 HTML 或 Cursor 主控台。Grok Bot 助手用量仍在 grokbot-usage.json。",
		"JSON을 %@에 두거나 아래에서 가져오세요. Cursor 제품 사용량에는 공개된 제3자 결제 API가 없습니다. 이 패널은 쿠키, 대시보드 HTML, Cursor 콘솔을 긁어오지 않습니다. Grok Bot 도우미 사용량은 grokbot-usage.json에 있습니다.",
		"JSON を %@ に置くか、下の読み込みを使ってください。Cursor 製品の使用量に公開の第三者課金 API はありません。このパネルは Cookie、ダッシュボード HTML、Cursor コンソールを取得しません。Grok Bot アシスタントの使用量は grokbot-usage.json のままです。",
		"Положите JSON в %@ или импортируйте его ниже. У использования продукта Cursor нет публичного стороннего API биллинга. Эта панель не снимает cookie, HTML панели и консоль Cursor. Использование помощника Grok Bot остаётся в grokbot-usage.json.",
		"Placez un JSON dans %@, ou importez-le ci-dessous. L’usage du produit Cursor n’a pas d’API de facturation tierce publique. Ce panneau ne récupère ni cookies, ni HTML de tableau de bord, ni la console Cursor. L’usage de l’assistant Grok Bot reste dans grokbot-usage.json.",
	}
	extraTexts[xCursorImport] = [7]string{
		"Import JSON", "导入 JSON", "匯入 JSON", "JSON 가져오기", "JSON を読み込む", "Импорт JSON", "Importer le JSON",
	}
	extraTexts[xCursorSourceNote] = [7]string{
		"Local snapshot, not a public third-party billing API. Separate from the Grok Bot assistant file.",
		"本地快照，不是公开的第三方计费接口。与 Grok Bot 助手的快照分开。",
		"本機快照，不是公開的第三方計費介面。與 Grok Bot 助手的快照分開。",
		"로컬 스냅샷이며 공개된 제3자 결제 API가 아닙니다. Grok Bot 도우미 파일과 별개입니다.",
		"ローカルスナップショットであり、公開の第三者課金 API ではありません。Grok Bot アシスタントのファイルとは別です。",
		"Локальный снимок, не публичный сторонний API биллинга. Отдельно от файла помощника Grok Bot.",
		"Instantané local, pas une API de facturation tierce publique. Distinct du fichier de l’assistant Grok Bot.",
	}
	extraTexts[xCursorOnDemandUsed] = [7]string{
		"On-demand used %@", "按需已用 %@", "隨需已用 %@", "온디맨드 사용 %@", "オンデマンド使用 %@", "По запросу использовано %@", "À la demande utilisé %@",
	}
	extraTexts[xCursorOnDemandCap] = [7]string{
		"On-demand limit %@", "按需上限 %@", "隨需上限 %@", "온디맨드 한도 %@", "オンデマンド上限 %@", "Лимит по запросу %@", "Plafond à la demande %@",
	}
	extraTexts[xCursorOnDemandRange] = [7]string{
		"On-demand %@ / %@", "按需 %@ / %@", "隨需 %@ / %@", "온디맨드 %@ / %@", "オンデマンド %@ / %@", "По запросу %@ / %@", "À la demande %@ / %@",
	}
	extraTexts[xCursorPeriodMonth] = [7]string{
		"This month", "本月", "本月", "이번 달", "今月", "Этот месяц", "Ce mois-ci",
	}
	extraTexts[xCursorPaceTitle] = [7]string{
		"Usage pace", "用量推算", "用量推算", "사용량 추정", "使用量の見積もり", "Оценка использования", "Projection d’usage",
	}
}

const (
	xCursorTab           textKey = "cursorTab"
	xCursorTitle         textKey = "cursorTitle"
	xCursorLocal         textKey = "cursorLocal"
	xCursorEmptyTitle    textKey = "cursorEmptyTitle"
	xCursorEmptyDetail   textKey = "cursorEmptyDetail"
	xCursorImport        textKey = "cursorImport"
	xCursorSourceNote    textKey = "cursorSourceNote"
	xCursorOnDemandUsed  textKey = "cursorOnDemandUsed"
	xCursorOnDemandCap   textKey = "cursorOnDemandCap"
	xCursorOnDemandRange textKey = "cursorOnDemandRange"
	xCursorPeriodMonth   textKey = "cursorPeriodMonth"
	xCursorPaceTitle     textKey = "cursorPaceTitle"
)
