import axios, { AxiosError } from "axios";

/**
 * API 에러에서 사용자 친화적 메시지를 추출합니다.
 * - 네트워크 오류: "서버에 연결할 수 없습니다"
 * - 5xx 서버 오류: "서버에 문제가 발생했습니다 (오류 코드: {status})"
 * - 기타: API 응답 메시지 또는 "요청에 실패했습니다"
 *
 * @param error - catch 블록에서 받은 에러 객체
 * @returns 사용자에게 표시할 에러 메시지
 */
export function getErrorMessage(error: unknown): string {
  // Axios 에러인지 확인
  if (!axios.isAxiosError(error)) {
    return "요청에 실패했습니다";
  }

  const axiosError = error as AxiosError<{ message?: string }>;

  // 서버 오류 (5xx) 또는 네트워크 오류 감지
  const status = axiosError.response?.status;
  const isServerError = status && status >= 500;
  const isNetworkError = !axiosError.response;

  if (isServerError) {
    return `서버에 문제가 발생했습니다 (오류 코드: ${status})`;
  } else if (isNetworkError) {
    return "서버에 연결할 수 없습니다";
  } else {
    // 일반 인증 오류 (400, 401, 403 등)
    // data가 null이거나 message가 빈 문자열인 경우 처리
    const message = axiosError.response?.data?.message;
    return message && message.trim() !== "" ? message : "요청에 실패했습니다";
  }
}
