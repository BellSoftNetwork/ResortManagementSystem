import { describe, it, expect, beforeEach } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useNetworkStore, RETRY_INTERVAL_MS } from "../network";

describe("useNetworkStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  describe("초기 상태", () => {
    it("기본 상태값이 올바르게 설정된다", () => {
      // given: 새로운 스토어 인스턴스

      // when: 스토어 생성
      const store = useNetworkStore();

      // then: 초기 상태 확인
      expect(store.isOffline).toBe(false);
      expect(store.errorType).toBe("NONE");
      expect(store.lastError).toBeNull();
      expect(store.statusCode).toBeNull();
      expect(store.retryCount).toBe(0);
      expect(store.autoRetryEnabled).toBe(false);
      expect(store.autoRetryCount).toBe(0);
    });
  });

  describe("getters", () => {
    it("isOnline은 isOffline의 반대값을 반환한다", () => {
      // given: 스토어 생성
      const store = useNetworkStore();

      // when/then: 온라인 상태
      expect(store.isOnline).toBe(true);

      // when: 오프라인으로 변경
      store.setNetworkError();

      // then: 오프라인 상태
      expect(store.isOnline).toBe(false);
    });

    it("isNetworkError는 errorType이 NETWORK_ERROR일 때 true를 반환한다", () => {
      // given: 스토어 생성
      const store = useNetworkStore();

      // when/then: 초기 상태
      expect(store.isNetworkError).toBe(false);

      // when: 네트워크 에러 설정
      store.setNetworkError();

      // then: 네트워크 에러 상태
      expect(store.isNetworkError).toBe(true);
    });

    it("isServerError는 errorType이 SERVER_ERROR일 때 true를 반환한다", () => {
      // given: 스토어 생성
      const store = useNetworkStore();

      // when/then: 초기 상태
      expect(store.isServerError).toBe(false);

      // when: 서버 에러 설정
      store.setServerError(500);

      // then: 서버 에러 상태
      expect(store.isServerError).toBe(true);
    });

    it("errorMessage는 커스텀 메시지가 있으면 반환하고, 없으면 기본 메시지를 반환한다", () => {
      // given: 스토어 생성
      const store = useNetworkStore();

      // when: 기본 네트워크 에러
      store.setNetworkError();

      // then: 기본 메시지
      expect(store.errorMessage).toBe("서버에 연결할 수 없습니다");

      // when: 커스텀 에러 메시지
      store.setNetworkError("커스텀 에러 메시지");

      // then: 커스텀 메시지
      expect(store.errorMessage).toBe("커스텀 에러 메시지");
    });

    it("canAutoRetry는 autoRetryEnabled가 true이고 최대 횟수 미만일 때 true를 반환한다", () => {
      // given: 스토어 생성
      const store = useNetworkStore();

      // when/then: 초기 상태 (비활성화)
      expect(store.canAutoRetry).toBe(false);

      // when: 자동 재시도 시작
      store.startAutoRetry();

      // then: 재시도 가능
      expect(store.canAutoRetry).toBe(true);

      // when: 최대 횟수 도달 (10회)
      for (let i = 0; i < 10; i++) {
        store.incrementAutoRetry();
      }

      // then: 재시도 불가
      expect(store.canAutoRetry).toBe(false);
    });

    it("retryProgress는 현재 재시도 횟수와 최대 횟수를 문자열로 반환한다", () => {
      // given: 스토어 생성
      const store = useNetworkStore();

      // when/then: 초기 상태
      expect(store.retryProgress).toBe("0/10");

      // when: 재시도 횟수 증가
      store.incrementAutoRetry();
      store.incrementAutoRetry();
      store.incrementAutoRetry();

      // then: 진행 상태
      expect(store.retryProgress).toBe("3/10");
    });
  });

  describe("actions", () => {
    describe("setNetworkError", () => {
      it("네트워크 에러 상태를 설정한다", () => {
        // given: 스토어 생성
        const store = useNetworkStore();

        // when: 네트워크 에러 설정
        store.setNetworkError();

        // then: 상태 확인
        expect(store.isOffline).toBe(true);
        expect(store.errorType).toBe("NETWORK_ERROR");
        expect(store.lastError).toBeNull();
        expect(store.statusCode).toBeNull();
      });

      it("커스텀 에러 메시지와 함께 네트워크 에러를 설정한다", () => {
        // given: 스토어 생성
        const store = useNetworkStore();

        // when: 커스텀 에러 메시지와 함께 설정
        store.setNetworkError("네트워크 연결이 끊어졌습니다");

        // then: 상태 확인
        expect(store.isOffline).toBe(true);
        expect(store.errorType).toBe("NETWORK_ERROR");
        expect(store.lastError).toBe("네트워크 연결이 끊어졌습니다");
      });
    });

    describe("setServerError", () => {
      it("서버 에러 상태를 설정한다", () => {
        // given: 스토어 생성
        const store = useNetworkStore();

        // when: 서버 에러 설정
        store.setServerError(500);

        // then: 상태 확인
        expect(store.isOffline).toBe(true);
        expect(store.errorType).toBe("SERVER_ERROR");
        expect(store.statusCode).toBe(500);
        expect(store.lastError).toBeNull();
      });

      it("커스텀 에러 메시지와 함께 서버 에러를 설정한다", () => {
        // given: 스토어 생성
        const store = useNetworkStore();

        // when: 커스텀 에러 메시지와 함께 설정
        store.setServerError(503, "서비스 점검 중입니다");

        // then: 상태 확인
        expect(store.isOffline).toBe(true);
        expect(store.errorType).toBe("SERVER_ERROR");
        expect(store.statusCode).toBe(503);
        expect(store.lastError).toBe("서비스 점검 중입니다");
      });
    });

    describe("setOffline (deprecated)", () => {
      it("setNetworkError와 동일하게 동작한다", () => {
        // given: 스토어 생성
        const store = useNetworkStore();

        // when: deprecated 메서드 호출
        store.setOffline("연결 끊김");

        // then: setNetworkError와 동일한 결과
        expect(store.isOffline).toBe(true);
        expect(store.errorType).toBe("NETWORK_ERROR");
        expect(store.lastError).toBe("연결 끊김");
      });
    });

    describe("setOnline", () => {
      it("모든 상태를 초기화하고 온라인 상태로 변경한다", () => {
        // given: 오프라인 상태의 스토어
        const store = useNetworkStore();
        store.setServerError(500, "서버 에러");
        store.incrementRetryCount();
        store.startAutoRetry();
        store.incrementAutoRetry();

        // when: 온라인으로 변경
        store.setOnline();

        // then: 모든 상태 초기화
        expect(store.isOffline).toBe(false);
        expect(store.errorType).toBe("NONE");
        expect(store.lastError).toBeNull();
        expect(store.statusCode).toBeNull();
        expect(store.retryCount).toBe(0);
        expect(store.autoRetryEnabled).toBe(false);
        expect(store.autoRetryCount).toBe(0);
      });
    });

    describe("incrementRetryCount", () => {
      it("재시도 횟수를 증가시킨다", () => {
        // given: 스토어 생성
        const store = useNetworkStore();

        // when: 재시도 횟수 증가
        store.incrementRetryCount();
        store.incrementRetryCount();

        // then: 횟수 확인
        expect(store.retryCount).toBe(2);
      });
    });

    describe("자동 재시도 관리", () => {
      it("startAutoRetry는 자동 재시도를 활성화하고 카운터를 초기화한다", () => {
        // given: 스토어 생성
        const store = useNetworkStore();
        store.incrementAutoRetry();
        store.incrementAutoRetry();

        // when: 자동 재시도 시작
        store.startAutoRetry();

        // then: 활성화 및 카운터 초기화
        expect(store.autoRetryEnabled).toBe(true);
        expect(store.autoRetryCount).toBe(0);
      });

      it("incrementAutoRetry는 자동 재시도 횟수를 증가시킨다", () => {
        // given: 스토어 생성
        const store = useNetworkStore();

        // when: 자동 재시도 횟수 증가
        store.incrementAutoRetry();
        store.incrementAutoRetry();
        store.incrementAutoRetry();

        // then: 횟수 확인
        expect(store.autoRetryCount).toBe(3);
      });

      it("stopAutoRetry는 자동 재시도를 비활성화한다", () => {
        // given: 자동 재시도 활성화 상태
        const store = useNetworkStore();
        store.startAutoRetry();

        // when: 자동 재시도 중단
        store.stopAutoRetry();

        // then: 비활성화 확인
        expect(store.autoRetryEnabled).toBe(false);
      });

      it("resetRetryState는 자동 재시도 상태를 초기화한다", () => {
        // given: 자동 재시도 진행 중 상태
        const store = useNetworkStore();
        store.startAutoRetry();
        store.incrementAutoRetry();
        store.incrementAutoRetry();

        // when: 재시도 상태 초기화
        store.resetRetryState();

        // then: 초기화 확인
        expect(store.autoRetryEnabled).toBe(false);
        expect(store.autoRetryCount).toBe(0);
      });
    });
  });

  describe("상수 값", () => {
    it("RETRY_INTERVAL_MS는 10초(10000ms)이다", () => {
      expect(RETRY_INTERVAL_MS).toBe(10000);
    });
  });
});
