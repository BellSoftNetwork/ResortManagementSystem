/**
 * 네트워크 상태 관리 서비스
 * Pinia store를 래핑하여 기존 API와의 호환성 유지
 */
import { useNetworkStore } from "stores/network";

class NetworkStatusService {
  /**
   * Pinia store 인스턴스 가져오기
   * Note: store는 Pinia가 초기화된 후에만 사용 가능하므로 lazy하게 접근
   */
  private getStore() {
    return useNetworkStore();
  }

  /**
   * 현재 네트워크 상태 반환
   */
  get isOffline(): boolean {
    return this.getStore().isOffline;
  }

  /**
   * 네트워크 상태를 오프라인으로 설정
   */
  setOffline(): void {
    this.getStore().setOffline();
  }

  /**
   * 네트워크 오류 설정 (서버 응답 없음)
   */
  setNetworkError(error?: string): void {
    this.getStore().setNetworkError(error);
  }

  /**
   * 서버 오류 설정 (5xx 응답)
   */
  setServerError(statusCode: number, error?: string): void {
    this.getStore().setServerError(statusCode, error);
  }

  /**
   * 네트워크 상태를 온라인으로 설정
   */
  setOnline(): void {
    this.getStore().setOnline();
  }
}

// 싱글톤 인스턴스 생성
export const networkStatusService = new NetworkStatusService();
