package l10n

// CodexNotice is a short runtime sentence that is not in the generated tables.
// Languages follow Languages() order.
func (l Localizer) CodexNotice(key string) string {
	row, ok := codexNotices[key]
	if !ok {
		return key
	}
	for i, language := range Languages() {
		if language == l.Language() {
			return row[i]
		}
	}
	return row[0]
}

// ProxyDial explains a refused or timed-out proxy after one direct retry.
func (l Localizer) ProxyDial(endpoint string, triedDirect bool) string {
	if endpoint == "" {
		endpoint = "—"
	}
	key := "proxy_dial"
	if !triedDirect {
		key = "proxy_down"
	}
	return Format(l.CodexNotice(key), endpoint)
}

// ProxyDirectUsed is the quiet note after a dial failure succeeded directly.
func (l Localizer) ProxyDirectUsed(endpoint string) string {
	if endpoint == "" {
		endpoint = "—"
	}
	return Format(l.CodexNotice("proxy_direct"), endpoint)
}

var codexNotices = map[string][7]string{
	"proxy_dial": {
		"Local proxy %@ is unreachable. A direct retry also failed. Check the proxy, or switch Settings › Network to Follow System or Direct.",
		"本机代理不可用（%@）。已改试直连，仍然失败。请检查代理，或在设置 › 网络中改为跟随系统或直连。",
		"本機代理無法使用（%@）。已改試直接連線，仍然失敗。請檢查代理，或在設定 › 網路中改為跟隨系統或直接連線。",
		"로컬 프록시 %@ 에 연결할 수 없습니다. 직접 연결도 실패했습니다. 프록시를 확인하거나 설정 › 네트워크에서 시스템 따름 또는 직접 연결로 바꾸세요.",
		"ローカルプロキシ %@ が使えません。直接接続でも失敗しました。プロキシを確認するか、設定 › ネットワークでシステムに従うか直接接続に変えてください。",
		"Локальный прокси %@ недоступен. Прямая попытка тоже не удалась. Проверьте прокси или выберите в Настройках › Сеть «как в системе» или прямое соединение.",
		"Le proxy local %@ est injoignable. Un nouvel essai direct a aussi échoué. Vérifiez le proxy, ou choisissez Suivre le système ou Direct dans Réglages › Réseau.",
	},
	"proxy_down": {
		"Local proxy %@ is unreachable. Check the proxy, or switch Settings › Network to Follow System or Direct.",
		"本机代理不可用（%@）。请检查代理，或在设置 › 网络中改为跟随系统或直连。",
		"本機代理無法使用（%@）。請檢查代理，或在設定 › 網路中改為跟隨系統或直接連線。",
		"로컬 프록시 %@ 에 연결할 수 없습니다. 프록시를 확인하거나 설정 › 네트워크에서 시스템 따름 또는 직접 연결로 바꾸세요.",
		"ローカルプロキシ %@ が使えません。プロキシを確認するか、設定 › ネットワークでシステムに従うか直接接続に変えてください。",
		"Локальный прокси %@ недоступен. Проверьте прокси или выберите в Настройках › Сеть «как в системе» или прямое соединение.",
		"Le proxy local %@ est injoignable. Vérifiez le proxy, ou choisissez Suivre le système ou Direct dans Réglages › Réseau.",
	},
	"proxy_direct": {
		"Local proxy %@ could not be reached, so this request used a direct connection.",
		"本机代理 %@ 连不上，这次已改用直连。",
		"本機代理 %@ 連不上，這次已改用直接連線。",
		"로컬 프록시 %@ 에 연결할 수 없어 이번 요청은 직접 연결했습니다.",
		"ローカルプロキシ %@ に接続できなかったため、今回は直接接続しました。",
		"Локальный прокси %@ недоступен, этот запрос ушёл напрямую.",
		"Le proxy local %@ était injoignable : cette requête est passée en direct.",
	},
	"reset_feed": {
		"The public reset feed is temporarily unavailable. Local quota is unaffected. Refresh later.",
		"公开的重置状态暂时取不到。这不影响本机额度，请稍后刷新。",
		"公開的重置狀態暫時取不到。這不影響本機額度，請稍後重新整理。",
		"공개 재설정 상태를 잠시 가져올 수 없습니다. 로컬 할당량에는 영향이 없습니다. 나중에 새로고침하세요.",
		"公開のリセット状況を今は取得できません。ローカルのクォータには影響しません。あとで更新してください。",
		"Публичный статус сброса временно недоступен. Локальная квота не затронута. Обновите позже.",
		"Le flux public de réinitialisation est temporairement indisponible. Le quota local n’est pas affecté. Actualisez plus tard.",
	},
	"cost_local": {
		"Official usage was not fetched. Figures below are this computer’s Codex sessions for the calendar month, not the subscription cycle.",
		"官方用量未取到。以下为本机 Codex 会话，按自然月统计，不是订阅周期。",
		"官方用量未取到。以下為本機 Codex 工作階段，按日曆月統計，不是訂閱週期。",
		"공식 사용량을 가져오지 못했습니다. 아래는 이 컴퓨터의 Codex 세션이며 달력 월 기준입니다. 구독 주기가 아닙니다.",
		"公式の使用量を取得できませんでした。以下はこのコンピューターの Codex セッションで、暦月の集計です。サブスクリプション周期ではありません。",
		"Официальное потребление не получено. Ниже сессии Codex этого компьютера за календарный месяц, не цикл подписки.",
		"La consommation officielle n’a pas été obtenue. Les chiffres ci-dessous viennent des sessions Codex de cet ordinateur pour le mois civil, pas du cycle d’abonnement.",
	},
	"cost_none": {
		"Official usage was not fetched, and this computer has no priced sessions.",
		"官方用量未取到，本机也没有可计价的会话。",
		"官方用量未取到，本機也沒有可計價的工作階段。",
		"공식 사용량을 가져오지 못했고, 이 컴퓨터에도 가격을 매길 세션이 없습니다.",
		"公式の使用量を取得できず、このコンピューターにも料金換算できるセッションがありません。",
		"Официальное потребление не получено, и на этом компьютере нет сессий с ценой.",
		"La consommation officielle n’a pas été obtenue, et cet ordinateur n’a pas de session tarifée.",
	},
	"heat_local": {
		"Official token usage was not fetched. The chart uses sessions on this computer.",
		"官方 Token 用量未取到。图表使用本机会话。",
		"官方 Token 用量未取到。圖表使用本機工作階段。",
		"공식 토큰 사용량을 가져오지 못했습니다. 차트는 이 컴퓨터의 세션을 사용합니다.",
		"公式のトークン使用量を取得できませんでした。グラフはこのコンピューターのセッションです。",
		"Официальный расход токенов не получен. График строится по сессиям этого компьютера.",
		"La consommation officielle de tokens n’a pas été obtenue. Le graphique utilise les sessions de cet ordinateur.",
	},
	"heat_empty": {
		"Official usage and local sessions both have no token records.",
		"官方用量与本机会话都没有 Token 记录。",
		"官方用量與本機工作階段都沒有 Token 記錄。",
		"공식 사용량과 로컬 세션 모두 토큰 기록이 없습니다.",
		"公式の使用量にもローカルセッションにもトークン記録がありません。",
		"Нет записей токенов ни в официальном потреблении, ни в локальных сессиях.",
		"Ni la consommation officielle ni les sessions locales n’ont d’enregistrement de tokens.",
	},
	"repair_queued": {
		"Index repair will run when the full refresh finishes.",
		"全量刷新结束后将修复索引。",
		"全量重新整理結束後將修復索引。",
		"전체 새로고침이 끝나면 인덱스를 고칩니다.",
		"全体更新が終わってからインデックスを修復します。",
		"Индекс будет исправлен после полного обновления.",
		"L’index sera réparé à la fin de l’actualisation complète.",
	},
	"repair_left": {
		"Some index rows still do not match local sessions.",
		"修复后仍有未与本机会话对齐的索引项。",
		"修復後仍有未與本機工作階段對齊的索引項。",
		"복구 후에도 로컬 세션과 맞지 않는 인덱스 항목이 있습니다.",
		"修復後もローカルセッションと一致しない索引が残っています。",
		"После восстановления часть строк индекса всё ещё не совпадает с локальными сессиями.",
		"Après la réparation, certaines lignes d’index ne correspondent toujours pas aux sessions locales.",
	},
}
