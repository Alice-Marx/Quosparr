package app

func init() {
	extraTexts[xProxyDialFallback] = [7]string{
		"A refused or timed-out proxy dial is retried once directly.",
		"代理拨号被拒绝或超时后，会自动改试一次直连。",
		"代理撥號被拒絕或逾時後，會自動改試一次直接連線。",
		"프록시 연결이 거부되거나 시간 초과되면 한 번 직접 연결합니다.",
		"プロキシ接続が拒否またはタイムアウトしたら、一度だけ直接接続します。",
		"Если прокси отклоняет соединение или не отвечает, запрос один раз идёт напрямую.",
		"Un proxy qui refuse ou dépasse le délai est réessayé une fois en direct.",
	}
}

const xProxyDialFallback textKey = "proxyDialFallback"
