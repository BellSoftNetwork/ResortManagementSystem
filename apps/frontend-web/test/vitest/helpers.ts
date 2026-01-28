/**
 * API 응답 모킹을 위한 헬퍼 함수들
 */

/**
 * 목 API 응답 생성 (단일 값)
 * @param value - 응답 데이터
 * @returns 단일 값 응답 객체
 */
export function createMockApiResponse<T>(value: T): { value: T };

/**
 * 목 API 응답 생성 (목록)
 * @param value - 응답 데이터 배열
 * @param totalCount - 전체 항목 수
 * @returns 목록 응답 객체
 */
export function createMockApiResponse<T>(value: T[], totalCount: number): { value: T[]; totalCount: number };

/**
 * 목 API 응답 생성 (구현)
 */
export function createMockApiResponse<T>(value: T | T[], totalCount?: number) {
  if (totalCount !== undefined) {
    return { value, totalCount };
  }
  return { value };
}

/**
 * 목 API 에러 생성
 * @param message - 에러 메시지
 * @returns Error 객체
 */
export function createMockApiError(message: string): Error {
  return new Error(message);
}
