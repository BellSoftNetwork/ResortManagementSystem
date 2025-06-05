/**
 * 네트워크 상태 관리 서비스
 */
class NetworkStatusService {
  private _isOffline = false;

  /**
   * 현재 네트워크 상태 반환
   */
  get isOffline(): boolean {
    return this._isOffline;
  }

  /**
   * 네트워크 상태를 오프라인으로 설정
   */
  setOffline(): void {
    this._isOffline = true;
  }

  /**
   * 네트워크 상태를 온라인으로 설정
   */
  setOnline(): void {
    this._isOffline = false;
  }
}

// 싱글톤 인스턴스 생성
export const networkStatusService = new NetworkStatusService();
