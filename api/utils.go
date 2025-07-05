package api

import "net/http"

// 확장에서 접근을 위해 CORS를 허용
// TODO: Chainalysis, TRM, Trevse?? 로  URL 제ㅈ
func enableCORS(w *http.ResponseWriter) {
	// 모든 Origin을 허용: 개발 단계에서만 사용하고, 운영 환경에서는 특정 Origin으로 제한해야 합니다.
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
