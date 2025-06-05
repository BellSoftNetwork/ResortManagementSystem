import { Notify } from "quasar";
import { MAX_RETRY_COUNT, NOTIFICATION_TIMEOUT, RETRY_DELAY_BASE } from "../constants";

/**
 * 알림 관련 기능을 관리하는 서비스
 */
class NotificationService {
  /**
   * 오프라인 상태 알림 표시
   * @param retryCount 재시도 횟수
   */
  showOfflineNotification(retryCount?: number): void {
    let message = "오프라인 상태입니다.";

    if (retryCount !== undefined && retryCount > 0) {
      // 재시도 중인 경우 (retryCount가 1, 2, 3일 때)
      message = `오프라인 상태입니다. 재연결 시도 중 (${retryCount}/${MAX_RETRY_COUNT})`;
    } else if (retryCount === -1) {
      // 모든 재시도 실패 후
      message = "오프라인 상태입니다. 잠시 후 다시 시도해주세요.";
    } else {
      // 기본 메시지
      message = "오프라인 상태입니다. 재연결을 시도합니다.";
    }

    // 타임아웃 설정: 재시도 중이면 다음 재시도 직전에 사라지도록, 아니면 기본 타임아웃 적용
    const timeout =
      retryCount !== undefined && retryCount > 0 ? retryCount * RETRY_DELAY_BASE - 100 : NOTIFICATION_TIMEOUT;

    Notify.create({
      type: "negative",
      message: message,
      position: "bottom",
      timeout: timeout,
      progress: true, // 진행 표시기 추가
      actions: [
        {
          label: "닫기",
          color: "white",
        },
      ],
    });
  }

  /**
   * 온라인 상태 복귀 알림 표시
   */
  showOnlineNotification(): void {
    Notify.create({
      type: "positive",
      message: "온라인 모드로 복귀되었습니다.",
      position: "bottom",
      timeout: NOTIFICATION_TIMEOUT,
    });
  }
}

// 싱글톤 인스턴스 생성
export const notificationService = new NotificationService();
