import { boot } from "quasar/wrappers";
import { useNetworkStore } from "stores/network";

/**
 * 브라우저 네트워크 이벤트 리스너 등록
 * - online: 브라우저가 네트워크 연결 감지 시 헬스체크 시도
 * - offline: 브라우저가 네트워크 연결 끊김 감지 시 로깅만 (실제 오프라인 처리는 API 호출 실패 시)
 */
export default boot(() => {
  // Pinia가 초기화된 후에 실행되도록 setTimeout 사용
  setTimeout(() => {
    const networkStore = useNetworkStore();

    // 브라우저 online 이벤트 감지
    window.addEventListener("online", async () => {
      console.log("[NetworkEvents] Browser online detected");

      // 현재 오프라인 상태라면 즉시 헬스체크
      if (networkStore.isOffline) {
        try {
          const response = await fetch("/api/v1/config");
          if (response.ok) {
            console.log("[NetworkEvents] Server is reachable, recovering...");
            networkStore.setOnline();
            window.location.reload();
          }
        } catch {
          console.log("[NetworkEvents] Server still unreachable");
        }
      }
    });

    // 브라우저 offline 이벤트 감지 (로깅만)
    window.addEventListener("offline", () => {
      console.log("[NetworkEvents] Browser offline detected");
    });
  }, 100);
});
